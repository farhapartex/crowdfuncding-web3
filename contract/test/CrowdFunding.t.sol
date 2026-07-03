// SPDX-License-Identifier: MIT
pragma solidity ^0.8.24;

import {Test} from "forge-std/Test.sol";
import {CrowdFunding} from "../src/CrowdFunding.sol";
import {Campaign, CampaignStatus} from "../src/CrowdFundingTypes.sol";

contract CrowdFundingTest is Test {
    event CampaignCreated(uint256 indexed campaignId, address indexed owner, uint256 goal, uint256 deadline);
    event ContributionMade(uint256 indexed campaignId, address indexed contributor, uint256 amount);
    event FundsWithdrawn(uint256 indexed campaignId, address indexed owner, uint256 amount);
    event ContributionRefunded(uint256 indexed campaignId, address indexed contributor, uint256 amount);

    string constant TITLE = "Save the Turtles";
    string constant DESCRIPTION = "Help us protect sea turtles";
    uint256 constant GOAL = 10 ether;
    uint256 constant DURATION = 7 days;

    CrowdFunding crowdFunding;
    address campaignOwner = makeAddr("campaignOwner");
    address alice = makeAddr("alice");
    address bob = makeAddr("bob");

    function setUp() public {
        crowdFunding = new CrowdFunding();
        vm.deal(alice, 100 ether);
        vm.deal(bob, 100 ether);
    }

    function _createCampaign(uint256 goal, uint256 duration) internal returns (uint256 campaignId) {
        vm.prank(campaignOwner);
        campaignId = crowdFunding.createCampaign(TITLE, DESCRIPTION, goal, duration);
    }

    function test_CreateCampaign_StoresCampaignData() public {
        uint256 expectedDeadline = block.timestamp + DURATION;
        uint256 campaignId = _createCampaign(GOAL, DURATION);

        Campaign memory campaign = crowdFunding.getCampaign(campaignId);
        assertEq(campaign.owner, campaignOwner);
        assertEq(campaign.title, TITLE);
        assertEq(campaign.description, DESCRIPTION);
        assertEq(campaign.goal, GOAL);
        assertEq(campaign.deadline, expectedDeadline);
        assertEq(campaign.amountRaised, 0);
        assertFalse(campaign.withdrawn);
    }

    function test_CreateCampaign_EmitsCampaignCreated() public {
        uint256 expectedDeadline = block.timestamp + DURATION;

        vm.expectEmit(true, true, false, true);
        emit CampaignCreated(0, campaignOwner, GOAL, expectedDeadline);

        _createCampaign(GOAL, DURATION);
    }

    function test_CreateCampaign_IncreasesCampaignCount() public {
        _createCampaign(GOAL, DURATION);
        _createCampaign(GOAL, DURATION);

        assertEq(crowdFunding.campaignCount(), 2);
    }

    function test_CreateCampaign_RevertsWhenGoalIsZero() public {
        vm.expectRevert(CrowdFunding.GoalMustBeGreaterThanZero.selector);
        crowdFunding.createCampaign(TITLE, DESCRIPTION, 0, DURATION);
    }

    function test_CreateCampaign_RevertsWhenDurationIsZero() public {
        vm.expectRevert(CrowdFunding.DurationMustBeGreaterThanZero.selector);
        crowdFunding.createCampaign(TITLE, DESCRIPTION, GOAL, 0);
    }

    function test_Contribute_UpdatesAmountRaisedAndContribution() public {
        uint256 campaignId = _createCampaign(GOAL, DURATION);

        vm.prank(alice);
        crowdFunding.contribute{value: 3 ether}(campaignId);

        assertEq(crowdFunding.getCampaign(campaignId).amountRaised, 3 ether);
        assertEq(crowdFunding.getContribution(campaignId, alice), 3 ether);
    }

    function test_Contribute_AccumulatesMultipleContributionsFromSameAddress() public {
        uint256 campaignId = _createCampaign(GOAL, DURATION);

        vm.startPrank(alice);
        crowdFunding.contribute{value: 2 ether}(campaignId);
        crowdFunding.contribute{value: 1 ether}(campaignId);
        vm.stopPrank();

        assertEq(crowdFunding.getContribution(campaignId, alice), 3 ether);
    }

    function test_Contribute_EmitsContributionMade() public {
        uint256 campaignId = _createCampaign(GOAL, DURATION);

        vm.expectEmit(true, true, false, true);
        emit ContributionMade(campaignId, alice, 3 ether);

        vm.prank(alice);
        crowdFunding.contribute{value: 3 ether}(campaignId);
    }

    function test_Contribute_RevertsWhenCampaignDoesNotExist() public {
        vm.expectRevert(CrowdFunding.CampaignDoesNotExist.selector);
        vm.prank(alice);
        crowdFunding.contribute{value: 1 ether}(0);
    }

    function test_Contribute_RevertsWhenCampaignHasEnded() public {
        uint256 campaignId = _createCampaign(GOAL, DURATION);
        vm.warp(block.timestamp + DURATION + 1);

        vm.expectRevert(CrowdFunding.CampaignHasEnded.selector);
        vm.prank(alice);
        crowdFunding.contribute{value: 1 ether}(campaignId);
    }

    function test_Contribute_RevertsWhenValueIsZero() public {
        uint256 campaignId = _createCampaign(GOAL, DURATION);

        vm.expectRevert(CrowdFunding.ContributionMustBeGreaterThanZero.selector);
        vm.prank(alice);
        crowdFunding.contribute{value: 0}(campaignId);
    }

    function test_Withdraw_TransfersFundsToOwnerWhenGoalReached() public {
        uint256 campaignId = _createCampaign(GOAL, DURATION);
        vm.prank(alice);
        crowdFunding.contribute{value: GOAL}(campaignId);

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
        crowdFunding.contribute{value: GOAL}(campaignId);

        vm.expectEmit(true, true, false, true);
        emit FundsWithdrawn(campaignId, campaignOwner, GOAL);

        vm.prank(campaignOwner);
        crowdFunding.withdraw(campaignId);
    }

    function test_Withdraw_RevertsWhenCallerIsNotOwner() public {
        uint256 campaignId = _createCampaign(GOAL, DURATION);
        vm.prank(alice);
        crowdFunding.contribute{value: GOAL}(campaignId);

        vm.expectRevert(CrowdFunding.NotCampaignOwner.selector);
        vm.prank(bob);
        crowdFunding.withdraw(campaignId);
    }

    function test_Withdraw_RevertsWhenGoalNotReached() public {
        uint256 campaignId = _createCampaign(GOAL, DURATION);
        vm.prank(alice);
        crowdFunding.contribute{value: GOAL - 1}(campaignId);

        vm.expectRevert(CrowdFunding.GoalNotReached.selector);
        vm.prank(campaignOwner);
        crowdFunding.withdraw(campaignId);
    }

    function test_Withdraw_RevertsWhenAlreadyWithdrawn() public {
        uint256 campaignId = _createCampaign(GOAL, DURATION);
        vm.prank(alice);
        crowdFunding.contribute{value: GOAL}(campaignId);

        vm.startPrank(campaignOwner);
        crowdFunding.withdraw(campaignId);

        vm.expectRevert(CrowdFunding.FundsAlreadyWithdrawn.selector);
        crowdFunding.withdraw(campaignId);
        vm.stopPrank();
    }

    function test_Refund_ReturnsContributionWhenGoalNotReachedAfterDeadline() public {
        uint256 campaignId = _createCampaign(GOAL, DURATION);
        vm.prank(alice);
        crowdFunding.contribute{value: 3 ether}(campaignId);

        vm.warp(block.timestamp + DURATION + 1);

        uint256 aliceBalanceBefore = alice.balance;

        vm.prank(alice);
        crowdFunding.refund(campaignId);

        assertEq(alice.balance, aliceBalanceBefore + 3 ether);
        assertEq(crowdFunding.getContribution(campaignId, alice), 0);
    }

    function test_Refund_EmitsContributionRefunded() public {
        uint256 campaignId = _createCampaign(GOAL, DURATION);
        vm.prank(alice);
        crowdFunding.contribute{value: 3 ether}(campaignId);

        vm.warp(block.timestamp + DURATION + 1);

        vm.expectEmit(true, true, false, true);
        emit ContributionRefunded(campaignId, alice, 3 ether);

        vm.prank(alice);
        crowdFunding.refund(campaignId);
    }

    function test_Refund_RevertsWhenCampaignStillActive() public {
        uint256 campaignId = _createCampaign(GOAL, DURATION);
        vm.prank(alice);
        crowdFunding.contribute{value: 3 ether}(campaignId);

        vm.expectRevert(CrowdFunding.CampaignStillActive.selector);
        vm.prank(alice);
        crowdFunding.refund(campaignId);
    }

    function test_Refund_RevertsWhenGoalAlreadyReached() public {
        uint256 campaignId = _createCampaign(GOAL, DURATION);
        vm.prank(alice);
        crowdFunding.contribute{value: GOAL}(campaignId);

        vm.warp(block.timestamp + DURATION + 1);

        vm.expectRevert(CrowdFunding.GoalAlreadyReached.selector);
        vm.prank(alice);
        crowdFunding.refund(campaignId);
    }

    function test_Refund_RevertsWhenCallerHasNoContribution() public {
        uint256 campaignId = _createCampaign(GOAL, DURATION);
        vm.prank(alice);
        crowdFunding.contribute{value: 3 ether}(campaignId);

        vm.warp(block.timestamp + DURATION + 1);

        vm.expectRevert(CrowdFunding.NoContributionToRefund.selector);
        vm.prank(bob);
        crowdFunding.refund(campaignId);
    }

    function test_GetCampaignStatus_ReturnsActiveBeforeDeadlineWhenGoalNotReached() public {
        uint256 campaignId = _createCampaign(GOAL, DURATION);
        vm.prank(alice);
        crowdFunding.contribute{value: 1 ether}(campaignId);

        assertEq(uint256(crowdFunding.getCampaignStatus(campaignId)), uint256(CampaignStatus.Active));
    }

    function test_GetCampaignStatus_ReturnsSuccessfulWhenGoalReachedBeforeDeadline() public {
        uint256 campaignId = _createCampaign(GOAL, DURATION);
        vm.prank(alice);
        crowdFunding.contribute{value: GOAL}(campaignId);

        assertEq(uint256(crowdFunding.getCampaignStatus(campaignId)), uint256(CampaignStatus.Successful));
    }

    function test_GetCampaignStatus_ReturnsFailedAfterDeadlineWhenGoalNotReached() public {
        uint256 campaignId = _createCampaign(GOAL, DURATION);
        vm.prank(alice);
        crowdFunding.contribute{value: 1 ether}(campaignId);

        vm.warp(block.timestamp + DURATION + 1);

        assertEq(uint256(crowdFunding.getCampaignStatus(campaignId)), uint256(CampaignStatus.Failed));
    }

    function test_GetCampaigns_ReturnsAllCreatedCampaigns() public {
        _createCampaign(GOAL, DURATION);
        _createCampaign(GOAL * 2, DURATION * 2);

        Campaign[] memory allCampaigns = crowdFunding.getCampaigns();

        assertEq(allCampaigns.length, 2);
        assertEq(allCampaigns[0].goal, GOAL);
        assertEq(allCampaigns[1].goal, GOAL * 2);
    }

    function test_GetCampaign_RevertsWhenCampaignDoesNotExist() public {
        vm.expectRevert(CrowdFunding.CampaignDoesNotExist.selector);
        crowdFunding.getCampaign(0);
    }
}
