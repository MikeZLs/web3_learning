// SPDX-License-Identifier: MIT
pragma solidity ^0.8.18;

contract TipJar {
    address public owner;

    constructor(){
        owner = msg.sender;  // 谁发起部署，owner就是谁
    }

    modifier onlyOwner(){
        require(msg.sender == owner,"you are not owner");
        _;
    }

    function tip() public payable {
        require(msg.value > 0,"you should send a tip to use this function");
    }

    function withdraw() public onlyOwner{
        uint256 contractBalance = address(this).balance;

        require(contractBalance > 0,"there are no tips to withdraw");

        payable(owner).transfer(contractBalance);
    }

    function getBalance() view public onlyOwner returns (uint256) {
        return address(this).balance;
    }

    // // 添加一个查看当前调用者的函数
    // function getCurrentSender() public view returns (address) {
    //     return msg.sender;
    // }

}