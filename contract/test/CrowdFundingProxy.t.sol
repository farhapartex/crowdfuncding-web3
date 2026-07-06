// SPDX-License-Identifier: MIT
pragma solidity ^0.8.24;

import {Test} from "forge-std/Test.sol";
import {CrowdFunding} from "../src/CrowdFunding.sol";
import {CrowdFundingProxy} from "../src/CrowdFundingProxy.sol";
import {CrowdFundingV2} from "./mocks/CrowdFundingV2.sol";

contract CrowdFundingProxyTest is Test {
    string constant TITLE = "Save the Turtles";
    string constant DESCRIPTION = "Help us protect sea turtles";
    uint256 constant GOAL = 10 ether;
    uint256 constant DURATION = 7 days;

    CrowdFunding implementationV1;
    CrowdFundingProxy proxy;
    CrowdFunding crowdFunding;

    address admin = makeAddr("admin");
    address alice = makeAddr("alice");

    function setUp() public {
        implementationV1 = new CrowdFunding();
        proxy = new CrowdFundingProxy(address(implementationV1), admin);
        crowdFunding = CrowdFunding(address(proxy));

        vm.deal(alice, 10 ether);
    }

    function test_ProxyReturnsCorrectAdminAndImplementation() public view {
        assertEq(proxy.admin(), admin);
        assertEq(proxy.implementation(), address(implementationV1));
    }

    function test_CallsThroughProxyReachImplementationLogic() public {
        uint256 campaignId = crowdFunding.createCampaign(TITLE, DESCRIPTION, GOAL, DURATION);

        vm.prank(alice);
        crowdFunding.contribute{value: 3 ether}(campaignId);

        assertEq(crowdFunding.getContribution(campaignId, alice), 3 ether);
        assertEq(crowdFunding.getCampaign(campaignId).amountRaised, 3 ether);
    }

    function test_UpgradeTo_RevertsWhenCallerIsNotAdmin() public {
        CrowdFundingV2 implementationV2 = new CrowdFundingV2();

        vm.expectRevert(CrowdFundingProxy.NotAdmin.selector);
        proxy.upgradeTo(address(implementationV2));
    }

    function test_UpgradeTo_PreservesExistingDataAndAddsNewLogic() public {
        uint256 campaignId = crowdFunding.createCampaign(TITLE, DESCRIPTION, GOAL, DURATION);

        vm.prank(alice);
        crowdFunding.contribute{value: 3 ether}(campaignId);

        CrowdFundingV2 implementationV2 = new CrowdFundingV2();

        vm.prank(admin);
        proxy.upgradeTo(address(implementationV2));

        assertEq(proxy.implementation(), address(implementationV2));

        assertEq(crowdFunding.getContribution(campaignId, alice), 3 ether);
        assertEq(crowdFunding.getCampaign(campaignId).amountRaised, 3 ether);

        CrowdFundingV2 crowdFundingV2 = CrowdFundingV2(address(proxy));
        assertEq(crowdFundingV2.totalRaised(), 3 ether);
    }
}
