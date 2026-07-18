// SPDX-License-Identifier: MIT
pragma solidity ^0.8.24;

import {IERC20} from "@openzeppelin/contracts/token/ERC20/IERC20.sol";
import {SafeERC20} from "@openzeppelin/contracts/token/ERC20/utils/SafeERC20.sol";

import {Campaign, CampaignStatus, CurrencyMode} from "./CrowdFundingTypes.sol";

contract CrowdFunding {
    using SafeERC20 for IERC20;

    uint256 public constant MAX_PAGE_SIZE = 50;

    Campaign[] public campaigns;
    mapping(uint256 => mapping(address => uint256)) public contributionsEth;
    mapping(uint256 => mapping(address => uint256)) public contributionsToken;

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

    error GoalMustBeGreaterThanZero();
    error DurationMustBeGreaterThanZero();
    error TokenAddressRequired();
    error EthNotAccepted();
    error TokenNotAccepted();
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

    function createCampaign(
        string calldata title,
        string calldata description,
        CurrencyMode currencyMode,
        address token,
        uint256 goalEth,
        uint256 goalToken,
        uint256 durationInSeconds
    ) external returns (uint256 campaignId) {
        if (durationInSeconds == 0) revert DurationMustBeGreaterThanZero();

        if (currencyMode == CurrencyMode.EthOnly) {
            if (goalEth == 0) revert GoalMustBeGreaterThanZero();
        } else if (currencyMode == CurrencyMode.TokenOnly) {
            if (goalToken == 0) revert GoalMustBeGreaterThanZero();
            if (token == address(0)) revert TokenAddressRequired();
        } else {
            if (goalEth == 0 || goalToken == 0) revert GoalMustBeGreaterThanZero();
            if (token == address(0)) revert TokenAddressRequired();
        }

        campaignId = campaigns.length;
        campaigns.push(
            Campaign({
                owner: msg.sender,
                title: title,
                description: description,
                currencyMode: currencyMode,
                token: token,
                goalEth: goalEth,
                goalToken: goalToken,
                deadline: block.timestamp + durationInSeconds,
                amountRaisedEth: 0,
                amountRaisedToken: 0,
                withdrawn: false
            })
        );

        emit CampaignCreated(campaignId, msg.sender, currencyMode, token, goalEth, goalToken, campaigns[campaignId].deadline);
    }

    function contributeEth(uint256 campaignId) external payable campaignExists(campaignId) {
        Campaign storage campaign = campaigns[campaignId];

        if (campaign.currencyMode == CurrencyMode.TokenOnly) revert EthNotAccepted();
        // forge-lint: disable-next-line(block-timestamp)
        if (block.timestamp >= campaign.deadline) revert CampaignHasEnded();
        if (msg.value == 0) revert ContributionMustBeGreaterThanZero();

        campaign.amountRaisedEth += msg.value;
        contributionsEth[campaignId][msg.sender] += msg.value;

        emit ContributionMade(campaignId, msg.sender, address(0), msg.value);
    }

    function contributeToken(uint256 campaignId, uint256 amount) external campaignExists(campaignId) {
        Campaign storage campaign = campaigns[campaignId];

        if (campaign.currencyMode == CurrencyMode.EthOnly) revert TokenNotAccepted();
        // forge-lint: disable-next-line(block-timestamp)
        if (block.timestamp >= campaign.deadline) revert CampaignHasEnded();
        if (amount == 0) revert ContributionMustBeGreaterThanZero();

        IERC20(campaign.token).safeTransferFrom(msg.sender, address(this), amount);

        campaign.amountRaisedToken += amount;
        contributionsToken[campaignId][msg.sender] += amount;

        emit ContributionMade(campaignId, msg.sender, campaign.token, amount);
    }

    function closeCampaign(uint256 campaignId) external campaignExists(campaignId) {
        Campaign storage campaign = campaigns[campaignId];

        if (msg.sender != campaign.owner) revert NotCampaignOwner();
        // forge-lint: disable-next-line(block-timestamp)
        if (block.timestamp >= campaign.deadline) revert CampaignHasEnded();

        campaign.deadline = block.timestamp;

        emit CampaignClosed(campaignId, msg.sender);
    }

    function withdraw(uint256 campaignId) external campaignExists(campaignId) {
        Campaign storage campaign = campaigns[campaignId];

        if (msg.sender != campaign.owner) revert NotCampaignOwner();
        if (campaign.withdrawn) revert FundsAlreadyWithdrawn();
        if (!_goalReached(campaign)) revert GoalNotReached();

        campaign.withdrawn = true;

        uint256 ethAmount = campaign.amountRaisedEth;
        uint256 tokenAmount = campaign.amountRaisedToken;

        if (ethAmount > 0) {
            (bool success,) = payable(campaign.owner).call{value: ethAmount}("");
            if (!success) revert TransferFailed();
        }
        if (tokenAmount > 0) {
            IERC20(campaign.token).safeTransfer(campaign.owner, tokenAmount);
        }

        emit FundsWithdrawn(campaignId, campaign.owner, ethAmount, tokenAmount);
    }

    function refund(uint256 campaignId) external campaignExists(campaignId) {
        Campaign storage campaign = campaigns[campaignId];

        // forge-lint: disable-next-line(block-timestamp)
        if (block.timestamp < campaign.deadline) revert CampaignStillActive();
        if (_goalReached(campaign)) revert GoalAlreadyReached();

        uint256 ethContributed = contributionsEth[campaignId][msg.sender];
        uint256 tokenContributed = contributionsToken[campaignId][msg.sender];
        if (ethContributed == 0 && tokenContributed == 0) revert NoContributionToRefund();

        if (ethContributed > 0) {
            contributionsEth[campaignId][msg.sender] = 0;
            (bool success,) = payable(msg.sender).call{value: ethContributed}("");
            if (!success) revert TransferFailed();
        }
        if (tokenContributed > 0) {
            contributionsToken[campaignId][msg.sender] = 0;
            IERC20(campaign.token).safeTransfer(msg.sender, tokenContributed);
        }

        emit ContributionRefunded(campaignId, msg.sender, ethContributed, tokenContributed);
    }

    function getCampaign(uint256 campaignId) external view campaignExists(campaignId) returns (Campaign memory) {
        return campaigns[campaignId];
    }

    function getCampaigns(uint256 offset, uint256 limit) external view returns (Campaign[] memory) {
        uint256 totalCampaigns = campaigns.length;
        if (offset >= totalCampaigns) {
            return new Campaign[](0);
        }

        uint256 pageSize = limit < MAX_PAGE_SIZE ? limit : MAX_PAGE_SIZE;
        uint256 remaining = totalCampaigns - offset;
        if (pageSize > remaining) {
            pageSize = remaining;
        }

        Campaign[] memory page = new Campaign[](pageSize);
        for (uint256 i = 0; i < pageSize; i++) {
            page[i] = campaigns[offset + i];
        }

        return page;
    }

    function getCampaignStatus(uint256 campaignId) external view campaignExists(campaignId) returns (CampaignStatus) {
        Campaign storage campaign = campaigns[campaignId];

        if (_goalReached(campaign)) {
            return CampaignStatus.Successful;
        }
        // forge-lint: disable-next-line(block-timestamp)
        if (block.timestamp >= campaign.deadline) {
            return CampaignStatus.Failed;
        }
        return CampaignStatus.Active;
    }

    function getContributionEth(uint256 campaignId, address contributor)
        external
        view
        campaignExists(campaignId)
        returns (uint256)
    {
        return contributionsEth[campaignId][contributor];
    }

    function getContributionToken(uint256 campaignId, address contributor)
        external
        view
        campaignExists(campaignId)
        returns (uint256)
    {
        return contributionsToken[campaignId][contributor];
    }

    function campaignCount() external view returns (uint256) {
        return campaigns.length;
    }

    function _goalReached(Campaign storage campaign) private view returns (bool) {
        return (campaign.goalEth > 0 && campaign.amountRaisedEth >= campaign.goalEth)
            || (campaign.goalToken > 0 && campaign.amountRaisedToken >= campaign.goalToken);
    }
}
