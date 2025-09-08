// SPDX-License-Identifier: MIT
pragma solidity ^0.8.18;

contract ProfileStatus {
    struct Status{
        string name;
        string message;
    }

    mapping (address => Status) private userStatus;

    function createOrUpdateStatus(string memory _name,string memory _message) public {
        userStatus[msg.sender].name = _name;
        userStatus[msg.sender].message = _message;
    }

    function getStatus() view public returns(string memory name,string memory message){
        return (userStatus[msg.sender].name,userStatus[msg.sender].message);
    }


}