// SPDX-License-Identifier: MIT
pragma solidity ^0.8.24;

enum CampaignStatus {
    Active,
    Successful,
    Failed
}

struct Campaign {
    address owner;
    string title;
    string description;
    uint256 goal;
    uint256 deadline;
    uint256 amountRaised;
    bool withdrawn;
}
