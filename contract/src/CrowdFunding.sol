// SPDX-License-Identifier: MIT
pragma solidity ^0.8.24;

import {Campaign, CampaignStatus} from "./CrowdFundingTypes.sol";

contract CrowdFunding {
    Campaign[] public campaigns;
    mapping(uint256 => mapping(address => uint256)) public contributions;

    event CampaignCreated(uint256 indexed campaignId, address indexed owner, uint256 goal, uint256 deadline);
    event ContributionMade(uint256 indexed campaignId, address indexed contributor, uint256 amount);
    event FundsWithdrawn(uint256 indexed campaignId, address indexed owner, uint256 amount);
    event ContributionRefunded(uint256 indexed campaignId, address indexed contributor, uint256 amount);

    error GoalMustBeGreaterThanZero();
    error DurationMustBeGreaterThanZero();
    error CampaignDoesNotExist();
    error CampaignHasEnded();
    error CampaignStillActive();
    error ContributionMustBeGreaterThanZero();
    error GoalNotReached();
    error GoalAlreadyReached();
    error FundsAlreadyWithdrawn();
    error NotCampaignOwner();
    error NoContributionToRefund();
    error TransferFailed();

    modifier campaignExists(uint256 campaignId) {
        if (campaignId >= campaigns.length) revert CampaignDoesNotExist();
        _;
    }

    function createCampaign(string calldata title, string calldata description, uint256 goal, uint256 durationInSeconds)
        external
        returns (uint256 campaignId)
    {
        if (goal == 0) revert GoalMustBeGreaterThanZero();
        if (durationInSeconds == 0) revert DurationMustBeGreaterThanZero();

        campaignId = campaigns.length;
        campaigns.push(
            Campaign({
                owner: msg.sender,
                title: title,
                description: description,
                goal: goal,
                deadline: block.timestamp + durationInSeconds,
                amountRaised: 0,
                withdrawn: false
            })
        );

        emit CampaignCreated(campaignId, msg.sender, goal, campaigns[campaignId].deadline);
    }

    function contribute(uint256 campaignId) external payable campaignExists(campaignId) {
        Campaign storage campaign = campaigns[campaignId];

        // forge-lint: disable-next-line(block-timestamp)
        if (block.timestamp >= campaign.deadline) revert CampaignHasEnded();
        if (msg.value == 0) revert ContributionMustBeGreaterThanZero();

        campaign.amountRaised += msg.value;
        contributions[campaignId][msg.sender] += msg.value;

        emit ContributionMade(campaignId, msg.sender, msg.value);
    }

    function withdraw(uint256 campaignId) external campaignExists(campaignId) {
        Campaign storage campaign = campaigns[campaignId];

        if (msg.sender != campaign.owner) revert NotCampaignOwner();
        if (campaign.amountRaised < campaign.goal) revert GoalNotReached();
        if (campaign.withdrawn) revert FundsAlreadyWithdrawn();

        campaign.withdrawn = true;
        uint256 amountToWithdraw = campaign.amountRaised;

        (bool success,) = payable(campaign.owner).call{value: amountToWithdraw}("");
        if (!success) revert TransferFailed();

        emit FundsWithdrawn(campaignId, campaign.owner, amountToWithdraw);
    }

    function refund(uint256 campaignId) external campaignExists(campaignId) {
        Campaign storage campaign = campaigns[campaignId];

        // forge-lint: disable-next-line(block-timestamp)
        if (block.timestamp < campaign.deadline) revert CampaignStillActive();
        if (campaign.amountRaised >= campaign.goal) revert GoalAlreadyReached();

        uint256 contributedAmount = contributions[campaignId][msg.sender];
        if (contributedAmount == 0) revert NoContributionToRefund();

        contributions[campaignId][msg.sender] = 0;

        (bool success,) = payable(msg.sender).call{value: contributedAmount}("");
        if (!success) revert TransferFailed();

        emit ContributionRefunded(campaignId, msg.sender, contributedAmount);
    }

    function getCampaign(uint256 campaignId) external view campaignExists(campaignId) returns (Campaign memory) {
        return campaigns[campaignId];
    }

    function getCampaigns() external view returns (Campaign[] memory) {
        return campaigns;
    }

    function getCampaignStatus(uint256 campaignId) external view campaignExists(campaignId) returns (CampaignStatus) {
        Campaign storage campaign = campaigns[campaignId];

        if (campaign.amountRaised >= campaign.goal) {
            return CampaignStatus.Successful;
        }
        // forge-lint: disable-next-line(block-timestamp)
        if (block.timestamp >= campaign.deadline) {
            return CampaignStatus.Failed;
        }
        return CampaignStatus.Active;
    }

    function getContribution(uint256 campaignId, address contributor)
        external
        view
        campaignExists(campaignId)
        returns (uint256)
    {
        return contributions[campaignId][contributor];
    }

    function campaignCount() external view returns (uint256) {
        return campaigns.length;
    }
}
