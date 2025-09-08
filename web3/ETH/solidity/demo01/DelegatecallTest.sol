// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

contract Implementation {
    uint256 public data;

    event DataStored(address indexed sender, uint256 data);

    function storeData(uint256 _data) public {
        data = _data;
        emit DataStored(msg.sender, _data);
    }
}

contract Proxy {
    uint256 public data;
    address public implementation;

    constructor(address _implementation) {
        implementation = _implementation;
    }

    fallback() external payable {
        (bool success, ) = implementation.delegatecall(msg.data);
        require(success, "Delegatecall failed");
    }
    receive() external payable { }
}