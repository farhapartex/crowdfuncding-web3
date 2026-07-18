// SPDX-License-Identifier: MIT
pragma solidity ^0.8.24;

import {Test} from "forge-std/Test.sol";
import {CrowdFunding} from "../src/CrowdFunding.sol";
import {MockERC20} from "../src/MockERC20.sol";
import {Campaign, CampaignStatus, CurrencyMode} from "../src/CrowdFundingTypes.sol";

contract CrowdFundingTest is Test {
    event CampaignCreated(
        uint256 indexed campaignId,
        address indexed owner,
        CurrencyMode currencyMode,
        address token,
        uint256 goalEth,
        uint256 goalToken,
        uint256 deadline
    );
    event ContributionMade(uint256 indexed campaignId, address indexed contributor, address token, uint256 amount);
    event FundsWithdrawn(uint256 indexed campaignId, address indexed owner, uint256 ethAmount, uint256 tokenAmount);
    event ContributionRefunded(
        uint256 indexed campaignId, address indexed contributor, uint256 ethAmount, uint256 tokenAmount
    );
    event CampaignClosed(uint256 indexed campaignId, address indexed owner);

    string constant TITLE = "Save the Turtles";
    string constant DESCRIPTION = "Help us protect sea turtles";
    uint256 constant GOAL = 10 ether;
    uint256 constant DURATION = 7 days;
    uint256 constant TOKEN_GOAL = 10_000e6;

    CrowdFunding crowdFunding;
    MockERC20 token;
    address campaignOwner = makeAddr("campaignOwner");
    address alice = makeAddr("alice");
    address bob = makeAddr("bob");

    function setUp() public {
        crowdFunding = new CrowdFunding();
        token = new MockERC20("Test USD Coin", "tUSDC", 6, 0);

        vm.deal(alice, 100 ether);
        vm.deal(bob, 100 ether);

        token.mint(alice, 1_000_000e6);
        token.mint(bob, 1_000_000e6);
    }

    function _createCampaign(uint256 goal, uint256 duration) internal returns (uint256 campaignId) {
        vm.prank(campaignOwner);
        campaignId = crowdFunding.createCampaign(TITLE, DESCRIPTION, CurrencyMode.EthOnly, address(0), goal, 0, duration);
    }

    function _createTokenCampaign(uint256 goalToken, uint256 duration) internal returns (uint256 campaignId) {
        vm.prank(campaignOwner);
        campaignId =
            crowdFunding.createCampaign(TITLE, DESCRIPTION, CurrencyMode.TokenOnly, address(token), 0, goalToken, duration);
    }

    function _createBothCampaign(uint256 goalEth, uint256 goalToken, uint256 duration)
        internal
        returns (uint256 campaignId)
    {
        vm.prank(campaignOwner);
        campaignId = crowdFunding.createCampaign(
            TITLE, DESCRIPTION, CurrencyMode.Both, address(token), goalEth, goalToken, duration
        );
    }

    // ---------------------------------------------------------------------
    // createCampaign
    // ---------------------------------------------------------------------

    function test_CreateCampaign_StoresCampaignData() public {
        uint256 expectedDeadline = block.timestamp + DURATION;
        uint256 campaignId = _createCampaign(GOAL, DURATION);

        Campaign memory campaign = crowdFunding.getCampaign(campaignId);
        assertEq(campaign.owner, campaignOwner);
        assertEq(campaign.title, TITLE);
        assertEq(campaign.description, DESCRIPTION);
        assertEq(uint256(campaign.currencyMode), uint256(CurrencyMode.EthOnly));
        assertEq(campaign.token, address(0));
        assertEq(campaign.goalEth, GOAL);
        assertEq(campaign.goalToken, 0);
        assertEq(campaign.deadline, expectedDeadline);
        assertEq(campaign.amountRaisedEth, 0);
        assertEq(campaign.amountRaisedToken, 0);
        assertFalse(campaign.withdrawn);
    }

    function test_CreateCampaign_EmitsCampaignCreated() public {
        uint256 expectedDeadline = block.timestamp + DURATION;

        vm.expectEmit(true, true, false, true);
        emit CampaignCreated(0, campaignOwner, CurrencyMode.EthOnly, address(0), GOAL, 0, expectedDeadline);

        _createCampaign(GOAL, DURATION);
    }

    function test_CreateCampaign_IncreasesCampaignCount() public {
        _createCampaign(GOAL, DURATION);
        _createCampaign(GOAL, DURATION);

        assertEq(crowdFunding.campaignCount(), 2);
    }

    function test_CreateCampaign_RevertsWhenGoalIsZero() public {
        vm.expectRevert(CrowdFunding.GoalMustBeGreaterThanZero.selector);
        crowdFunding.createCampaign(TITLE, DESCRIPTION, CurrencyMode.EthOnly, address(0), 0, 0, DURATION);
    }

    function test_CreateCampaign_RevertsWhenDurationIsZero() public {
        vm.expectRevert(CrowdFunding.DurationMustBeGreaterThanZero.selector);
        crowdFunding.createCampaign(TITLE, DESCRIPTION, CurrencyMode.EthOnly, address(0), GOAL, 0, 0);
    }

    function test_CreateCampaign_RevertsWhenTokenOnlyMissingTokenAddress() public {
        vm.expectRevert(CrowdFunding.TokenAddressRequired.selector);
        crowdFunding.createCampaign(TITLE, DESCRIPTION, CurrencyMode.TokenOnly, address(0), 0, TOKEN_GOAL, DURATION);
    }

    function test_CreateCampaign_RevertsWhenTokenOnlyGoalIsZero() public {
        vm.expectRevert(CrowdFunding.GoalMustBeGreaterThanZero.selector);
        crowdFunding.createCampaign(TITLE, DESCRIPTION, CurrencyMode.TokenOnly, address(token), 0, 0, DURATION);
    }

    function test_CreateCampaign_RevertsWhenBothMissingEitherGoal() public {
        vm.expectRevert(CrowdFunding.GoalMustBeGreaterThanZero.selector);
        crowdFunding.createCampaign(TITLE, DESCRIPTION, CurrencyMode.Both, address(token), GOAL, 0, DURATION);
    }

    // ---------------------------------------------------------------------
    // contributeEth
    // ---------------------------------------------------------------------

    function test_ContributeEth_UpdatesAmountRaisedAndContribution() public {
        uint256 campaignId = _createCampaign(GOAL, DURATION);

        vm.prank(alice);
        crowdFunding.contributeEth{value: 3 ether}(campaignId);

        assertEq(crowdFunding.getCampaign(campaignId).amountRaisedEth, 3 ether);
        assertEq(crowdFunding.getContributionEth(campaignId, alice), 3 ether);
    }

    function test_ContributeEth_AccumulatesMultipleContributionsFromSameAddress() public {
        uint256 campaignId = _createCampaign(GOAL, DURATION);

        vm.startPrank(alice);
        crowdFunding.contributeEth{value: 2 ether}(campaignId);
        crowdFunding.contributeEth{value: 1 ether}(campaignId);
        vm.stopPrank();

        assertEq(crowdFunding.getContributionEth(campaignId, alice), 3 ether);
    }

    function test_ContributeEth_EmitsContributionMade() public {
        uint256 campaignId = _createCampaign(GOAL, DURATION);

        vm.expectEmit(true, true, false, true);
        emit ContributionMade(campaignId, alice, address(0), 3 ether);

        vm.prank(alice);
        crowdFunding.contributeEth{value: 3 ether}(campaignId);
    }

    function test_ContributeEth_RevertsWhenCampaignDoesNotExist() public {
        vm.expectRevert(CrowdFunding.CampaignDoesNotExist.selector);
        vm.prank(alice);
        crowdFunding.contributeEth{value: 1 ether}(0);
    }

    function test_ContributeEth_RevertsWhenCampaignHasEnded() public {
        uint256 campaignId = _createCampaign(GOAL, DURATION);
        vm.warp(block.timestamp + DURATION + 1);

        vm.expectRevert(CrowdFunding.CampaignHasEnded.selector);
        vm.prank(alice);
        crowdFunding.contributeEth{value: 1 ether}(campaignId);
    }

    function test_ContributeEth_RevertsWhenValueIsZero() public {
        uint256 campaignId = _createCampaign(GOAL, DURATION);

        vm.expectRevert(CrowdFunding.ContributionMustBeGreaterThanZero.selector);
        vm.prank(alice);
        crowdFunding.contributeEth{value: 0}(campaignId);
    }

    function test_ContributeEth_RevertsOnTokenOnlyCampaign() public {
        uint256 campaignId = _createTokenCampaign(TOKEN_GOAL, DURATION);

        vm.expectRevert(CrowdFunding.EthNotAccepted.selector);
        vm.prank(alice);
        crowdFunding.contributeEth{value: 1 ether}(campaignId);
    }

    // ---------------------------------------------------------------------
    // contributeToken
    // ---------------------------------------------------------------------

    function test_ContributeToken_RevertsOnEthOnlyCampaign() public {
        uint256 campaignId = _createCampaign(GOAL, DURATION);

        vm.prank(alice);
        token.approve(address(crowdFunding), 100e6);

        vm.expectRevert(CrowdFunding.TokenNotAccepted.selector);
        vm.prank(alice);
        crowdFunding.contributeToken(campaignId, 100e6);
    }

    function test_ContributeToken_RevertsWithoutApproval() public {
        uint256 campaignId = _createTokenCampaign(TOKEN_GOAL, DURATION);

        vm.prank(alice);
        vm.expectRevert();
        crowdFunding.contributeToken(campaignId, 100e6);
    }

    function test_ContributeToken_PullsApprovedTokensAndUpdatesState() public {
        uint256 campaignId = _createTokenCampaign(TOKEN_GOAL, DURATION);

        vm.startPrank(alice);
        token.approve(address(crowdFunding), 500e6);
        crowdFunding.contributeToken(campaignId, 500e6);
        vm.stopPrank();

        assertEq(crowdFunding.getCampaign(campaignId).amountRaisedToken, 500e6);
        assertEq(crowdFunding.getContributionToken(campaignId, alice), 500e6);
        assertEq(token.balanceOf(address(crowdFunding)), 500e6);
    }

    function test_ContributeToken_EmitsContributionMadeWithTokenAddress() public {
        uint256 campaignId = _createTokenCampaign(TOKEN_GOAL, DURATION);

        vm.startPrank(alice);
        token.approve(address(crowdFunding), 500e6);

        vm.expectEmit(true, true, false, true);
        emit ContributionMade(campaignId, alice, address(token), 500e6);

        crowdFunding.contributeToken(campaignId, 500e6);
        vm.stopPrank();
    }

    function test_ContributeToken_RevertsWhenAmountIsZero() public {
        uint256 campaignId = _createTokenCampaign(TOKEN_GOAL, DURATION);

        vm.expectRevert(CrowdFunding.ContributionMustBeGreaterThanZero.selector);
        vm.prank(alice);
        crowdFunding.contributeToken(campaignId, 0);
    }

    // ---------------------------------------------------------------------
    // closeCampaign (ETH-only path retained; behavior is currency-agnostic)
    // ---------------------------------------------------------------------

    function test_CloseCampaign_SetsDeadlineToCurrentTimestamp() public {
        uint256 campaignId = _createCampaign(GOAL, DURATION);

        vm.prank(campaignOwner);
        crowdFunding.closeCampaign(campaignId);

        assertEq(crowdFunding.getCampaign(campaignId).deadline, block.timestamp);
    }

    function test_CloseCampaign_EmitsCampaignClosed() public {
        uint256 campaignId = _createCampaign(GOAL, DURATION);

        vm.expectEmit(true, true, false, true);
        emit CampaignClosed(campaignId, campaignOwner);

        vm.prank(campaignOwner);
        crowdFunding.closeCampaign(campaignId);
    }

    function test_CloseCampaign_RevertsWhenCallerIsNotOwner() public {
        uint256 campaignId = _createCampaign(GOAL, DURATION);

        vm.expectRevert(CrowdFunding.NotCampaignOwner.selector);
        vm.prank(alice);
        crowdFunding.closeCampaign(campaignId);
    }

    function test_CloseCampaign_RevertsWhenAlreadyEnded() public {
        uint256 campaignId = _createCampaign(GOAL, DURATION);
        vm.warp(block.timestamp + DURATION + 1);

        vm.expectRevert(CrowdFunding.CampaignHasEnded.selector);
        vm.prank(campaignOwner);
        crowdFunding.closeCampaign(campaignId);
    }

    function test_CloseCampaign_BlocksFurtherContributions() public {
        uint256 campaignId = _createCampaign(GOAL, DURATION);

        vm.prank(campaignOwner);
        crowdFunding.closeCampaign(campaignId);

        vm.expectRevert(CrowdFunding.CampaignHasEnded.selector);
        vm.prank(alice);
        crowdFunding.contributeEth{value: 1 ether}(campaignId);
    }

    function test_CloseCampaign_AllowsImmediateRefundWhenGoalNotReached() public {
        uint256 campaignId = _createCampaign(GOAL, DURATION);
        vm.prank(alice);
        crowdFunding.contributeEth{value: 3 ether}(campaignId);

        vm.prank(campaignOwner);
        crowdFunding.closeCampaign(campaignId);

        uint256 aliceBalanceBefore = alice.balance;

        vm.prank(alice);
        crowdFunding.refund(campaignId);

        assertEq(alice.balance, aliceBalanceBefore + 3 ether);
    }

    function test_CloseCampaign_DoesNotAllowRefundWhenGoalAlreadyReached() public {
        uint256 campaignId = _createCampaign(GOAL, DURATION);
        vm.prank(alice);
        crowdFunding.contributeEth{value: GOAL}(campaignId);

        vm.prank(campaignOwner);
        crowdFunding.closeCampaign(campaignId);

        vm.expectRevert(CrowdFunding.GoalAlreadyReached.selector);
        vm.prank(alice);
        crowdFunding.refund(campaignId);
    }

    // ---------------------------------------------------------------------
    // withdraw
    // ---------------------------------------------------------------------

    function test_Withdraw_TransfersFundsToOwnerWhenGoalReached() public {
        uint256 campaignId = _createCampaign(GOAL, DURATION);
        vm.prank(alice);
        crowdFunding.contributeEth{value: GOAL}(campaignId);

        uint256 ownerBalanceBefore = campaignOwner.balance;

        vm.prank(campaignOwner);
        crowdFunding.withdraw(campaignId);

        assertEq(campaignOwner.balance, ownerBalanceBefore + GOAL);
        assertEq(address(crowdFunding).balance, 0);
        assertTrue(crowdFunding.getCampaign(campaignId).withdrawn);
    }

    function test_Withdraw_EmitsFundsWithdrawn() public {
        uint256 campaignId = _createCampaign(GOAL, DURATION);
        vm.prank(alice);
        crowdFunding.contributeEth{value: GOAL}(campaignId);

        vm.expectEmit(true, true, false, true);
        emit FundsWithdrawn(campaignId, campaignOwner, GOAL, 0);

        vm.prank(campaignOwner);
        crowdFunding.withdraw(campaignId);
    }

    function test_Withdraw_RevertsWhenCallerIsNotOwner() public {
        uint256 campaignId = _createCampaign(GOAL, DURATION);
        vm.prank(alice);
        crowdFunding.contributeEth{value: GOAL}(campaignId);

        vm.expectRevert(CrowdFunding.NotCampaignOwner.selector);
        vm.prank(bob);
        crowdFunding.withdraw(campaignId);
    }

    function test_Withdraw_RevertsWhenGoalNotReached() public {
        uint256 campaignId = _createCampaign(GOAL, DURATION);
        vm.prank(alice);
        crowdFunding.contributeEth{value: GOAL - 1}(campaignId);

        vm.expectRevert(CrowdFunding.GoalNotReached.selector);
        vm.prank(campaignOwner);
        crowdFunding.withdraw(campaignId);
    }

    function test_Withdraw_RevertsWhenAlreadyWithdrawn() public {
        uint256 campaignId = _createCampaign(GOAL, DURATION);
        vm.prank(alice);
        crowdFunding.contributeEth{value: GOAL}(campaignId);

        vm.startPrank(campaignOwner);
        crowdFunding.withdraw(campaignId);

        vm.expectRevert(CrowdFunding.FundsAlreadyWithdrawn.selector);
        crowdFunding.withdraw(campaignId);
        vm.stopPrank();
    }

    function test_Withdraw_TokenOnly_TransfersTokensToOwnerWhenGoalReached() public {
        uint256 campaignId = _createTokenCampaign(TOKEN_GOAL, DURATION);

        vm.startPrank(alice);
        token.approve(address(crowdFunding), TOKEN_GOAL);
        crowdFunding.contributeToken(campaignId, TOKEN_GOAL);
        vm.stopPrank();

        vm.prank(campaignOwner);
        crowdFunding.withdraw(campaignId);

        assertEq(token.balanceOf(campaignOwner), TOKEN_GOAL);
        assertEq(token.balanceOf(address(crowdFunding)), 0);
    }

    function test_Withdraw_Both_TransfersBothCurrenciesOnceEitherGoalReached() public {
        // Only the ETH goal is reached; token goal is not — withdraw should still
        // pay out whatever balance exists in both currencies.
        uint256 campaignId = _createBothCampaign(GOAL, TOKEN_GOAL, DURATION);

        vm.prank(alice);
        crowdFunding.contributeEth{value: GOAL}(campaignId);

        vm.startPrank(bob);
        token.approve(address(crowdFunding), 1_000e6);
        crowdFunding.contributeToken(campaignId, 1_000e6);
        vm.stopPrank();

        uint256 ownerEthBefore = campaignOwner.balance;

        vm.prank(campaignOwner);
        crowdFunding.withdraw(campaignId);

        assertEq(campaignOwner.balance, ownerEthBefore + GOAL);
        assertEq(token.balanceOf(campaignOwner), 1_000e6);
    }

    // ---------------------------------------------------------------------
    // refund
    // ---------------------------------------------------------------------

    function test_Refund_ReturnsContributionWhenGoalNotReachedAfterDeadline() public {
        uint256 campaignId = _createCampaign(GOAL, DURATION);
        vm.prank(alice);
        crowdFunding.contributeEth{value: 3 ether}(campaignId);

        vm.warp(block.timestamp + DURATION + 1);

        uint256 aliceBalanceBefore = alice.balance;

        vm.prank(alice);
        crowdFunding.refund(campaignId);

        assertEq(alice.balance, aliceBalanceBefore + 3 ether);
        assertEq(crowdFunding.getContributionEth(campaignId, alice), 0);
    }

    function test_Refund_EmitsContributionRefunded() public {
        uint256 campaignId = _createCampaign(GOAL, DURATION);
        vm.prank(alice);
        crowdFunding.contributeEth{value: 3 ether}(campaignId);

        vm.warp(block.timestamp + DURATION + 1);

        vm.expectEmit(true, true, false, true);
        emit ContributionRefunded(campaignId, alice, 3 ether, 0);

        vm.prank(alice);
        crowdFunding.refund(campaignId);
    }

    function test_Refund_RevertsWhenCampaignStillActive() public {
        uint256 campaignId = _createCampaign(GOAL, DURATION);
        vm.prank(alice);
        crowdFunding.contributeEth{value: 3 ether}(campaignId);

        vm.expectRevert(CrowdFunding.CampaignStillActive.selector);
        vm.prank(alice);
        crowdFunding.refund(campaignId);
    }

    function test_Refund_RevertsWhenGoalAlreadyReached() public {
        uint256 campaignId = _createCampaign(GOAL, DURATION);
        vm.prank(alice);
        crowdFunding.contributeEth{value: GOAL}(campaignId);

        vm.warp(block.timestamp + DURATION + 1);

        vm.expectRevert(CrowdFunding.GoalAlreadyReached.selector);
        vm.prank(alice);
        crowdFunding.refund(campaignId);
    }

    function test_Refund_RevertsWhenCallerHasNoContribution() public {
        uint256 campaignId = _createCampaign(GOAL, DURATION);
        vm.prank(alice);
        crowdFunding.contributeEth{value: 3 ether}(campaignId);

        vm.warp(block.timestamp + DURATION + 1);

        vm.expectRevert(CrowdFunding.NoContributionToRefund.selector);
        vm.prank(bob);
        crowdFunding.refund(campaignId);
    }

    function test_Refund_TokenOnly_ReturnsTokensWhenGoalNotReachedAfterDeadline() public {
        uint256 campaignId = _createTokenCampaign(TOKEN_GOAL, DURATION);

        vm.startPrank(alice);
        token.approve(address(crowdFunding), 500e6);
        crowdFunding.contributeToken(campaignId, 500e6);
        vm.stopPrank();

        vm.warp(block.timestamp + DURATION + 1);

        vm.prank(alice);
        crowdFunding.refund(campaignId);

        assertEq(token.balanceOf(alice), 1_000_000e6);
        assertEq(crowdFunding.getContributionToken(campaignId, alice), 0);
    }

    function test_Refund_Both_RefundsEachContributorInTheirOwnCurrency() public {
        uint256 campaignId = _createBothCampaign(GOAL, TOKEN_GOAL, DURATION);

        vm.prank(alice);
        crowdFunding.contributeEth{value: 1 ether}(campaignId);

        vm.startPrank(bob);
        token.approve(address(crowdFunding), 500e6);
        crowdFunding.contributeToken(campaignId, 500e6);
        vm.stopPrank();

        vm.warp(block.timestamp + DURATION + 1);

        uint256 aliceBalanceBefore = alice.balance;
        vm.prank(alice);
        crowdFunding.refund(campaignId);
        assertEq(alice.balance, aliceBalanceBefore + 1 ether);

        vm.prank(bob);
        crowdFunding.refund(campaignId);
        assertEq(token.balanceOf(bob), 1_000_000e6);
    }

    // ---------------------------------------------------------------------
    // getCampaignStatus
    // ---------------------------------------------------------------------

    function test_GetCampaignStatus_ReturnsActiveBeforeDeadlineWhenGoalNotReached() public {
        uint256 campaignId = _createCampaign(GOAL, DURATION);
        vm.prank(alice);
        crowdFunding.contributeEth{value: 1 ether}(campaignId);

        assertEq(uint256(crowdFunding.getCampaignStatus(campaignId)), uint256(CampaignStatus.Active));
    }

    function test_GetCampaignStatus_ReturnsSuccessfulWhenGoalReachedBeforeDeadline() public {
        uint256 campaignId = _createCampaign(GOAL, DURATION);
        vm.prank(alice);
        crowdFunding.contributeEth{value: GOAL}(campaignId);

        assertEq(uint256(crowdFunding.getCampaignStatus(campaignId)), uint256(CampaignStatus.Successful));
    }

    function test_GetCampaignStatus_ReturnsFailedAfterDeadlineWhenGoalNotReached() public {
        uint256 campaignId = _createCampaign(GOAL, DURATION);
        vm.prank(alice);
        crowdFunding.contributeEth{value: 1 ether}(campaignId);

        vm.warp(block.timestamp + DURATION + 1);

        assertEq(uint256(crowdFunding.getCampaignStatus(campaignId)), uint256(CampaignStatus.Failed));
    }

    function test_GetCampaignStatus_Both_SuccessfulWhenOnlyTokenGoalReached() public {
        uint256 campaignId = _createBothCampaign(GOAL, TOKEN_GOAL, DURATION);

        vm.startPrank(alice);
        token.approve(address(crowdFunding), TOKEN_GOAL);
        crowdFunding.contributeToken(campaignId, TOKEN_GOAL);
        vm.stopPrank();

        assertEq(uint256(crowdFunding.getCampaignStatus(campaignId)), uint256(CampaignStatus.Successful));
    }

    // ---------------------------------------------------------------------
    // getCampaigns / pagination
    // ---------------------------------------------------------------------

    function test_GetCampaigns_ReturnsAllCreatedCampaignsWithinOnePage() public {
        _createCampaign(GOAL, DURATION);
        _createCampaign(GOAL * 2, DURATION * 2);

        Campaign[] memory page = crowdFunding.getCampaigns(0, 10);

        assertEq(page.length, 2);
        assertEq(page[0].goalEth, GOAL);
        assertEq(page[1].goalEth, GOAL * 2);
    }

    function test_GetCampaigns_ReturnsRequestedSlice() public {
        _createCampaign(GOAL, DURATION);
        _createCampaign(GOAL * 2, DURATION * 2);
        _createCampaign(GOAL * 3, DURATION * 3);

        Campaign[] memory page = crowdFunding.getCampaigns(1, 1);

        assertEq(page.length, 1);
        assertEq(page[0].goalEth, GOAL * 2);
    }

    function test_GetCampaigns_ClampsLimitToMaxPageSize() public {
        for (uint256 i = 0; i < crowdFunding.MAX_PAGE_SIZE() + 5; i++) {
            _createCampaign(GOAL, DURATION);
        }

        Campaign[] memory page = crowdFunding.getCampaigns(0, crowdFunding.MAX_PAGE_SIZE() + 5);

        assertEq(page.length, crowdFunding.MAX_PAGE_SIZE());
    }

    function test_GetCampaigns_ReturnsEmptyWhenOffsetBeyondLength() public {
        _createCampaign(GOAL, DURATION);

        Campaign[] memory page = crowdFunding.getCampaigns(5, 10);

        assertEq(page.length, 0);
    }

    function test_GetCampaign_RevertsWhenCampaignDoesNotExist() public {
        vm.expectRevert(CrowdFunding.CampaignDoesNotExist.selector);
        crowdFunding.getCampaign(0);
    }
}
