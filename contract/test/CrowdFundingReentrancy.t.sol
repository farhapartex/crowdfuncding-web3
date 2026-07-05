// SPDX-License-Identifier: MIT
pragma solidity ^0.8.24;

import {Test} from "forge-std/Test.sol";
import {CrowdFunding} from "../src/CrowdFunding.sol";
import {Campaign} from "../src/CrowdFundingTypes.sol";
import {ReentrantRefunder} from "./attackers/ReentrantRefunder.sol";
import {ReentrantWithdrawer} from "./attackers/ReentrantWithdrawer.sol";

contract CrowdFundingReentrancyTest is Test {
    string constant TITLE = "Attack Campaign";
    string constant DESCRIPTION = "Attempting to drain funds via reentrancy";
    uint256 constant GOAL = 10 ether;
    uint256 constant DURATION = 7 days;

    CrowdFunding crowdFunding;
    address alice = makeAddr("alice");

    function setUp() public {
        crowdFunding = new CrowdFunding();
        vm.deal(alice, 100 ether);
    }

    function test_RefundReentrancy_RevertsEntireTransaction() public {
        ReentrantRefunder attacker = new ReentrantRefunder(crowdFunding);
        uint256 campaignId = crowdFunding.createCampaign(TITLE, DESCRIPTION, GOAL, DURATION);

        attacker.contribute{value: 3 ether}(campaignId);

        vm.warp(block.timestamp + DURATION + 1);

        vm.expectRevert(CrowdFunding.TransferFailed.selector);
        attacker.attack();

        assertEq(crowdFunding.getContribution(campaignId, address(attacker)), 3 ether);
        assertEq(address(crowdFunding).balance, 3 ether);
    }

    function test_WithdrawReentrancy_RevertsEntireTransaction() public {
        ReentrantWithdrawer attacker = new ReentrantWithdrawer(crowdFunding);
        uint256 campaignId = attacker.createCampaign(TITLE, DESCRIPTION, GOAL, DURATION);

        vm.prank(alice);
        crowdFunding.contribute{value: GOAL}(campaignId);

        vm.expectRevert(CrowdFunding.TransferFailed.selector);
        attacker.attack();

        Campaign memory campaign = crowdFunding.getCampaign(campaignId);
        assertFalse(campaign.withdrawn);
        assertEq(address(crowdFunding).balance, GOAL);
    }
}
