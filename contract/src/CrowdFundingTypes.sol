// SPDX-License-Identifier: MIT
pragma solidity ^0.8.24;

enum CampaignStatus {
    Active,
    Successful,
    Failed
}

enum CurrencyMode {
    EthOnly,
    TokenOnly,
    Both
}

struct Campaign {
    address owner;
    string title;
    string description;
    CurrencyMode currencyMode;
    address token;
    uint256 goalEth;
    uint256 goalToken;
    uint256 deadline;
    uint256 amountRaisedEth;
    uint256 amountRaisedToken;
    bool withdrawn;
}
