// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

contract SimpleStorage {
    // 定义一个名为 DataStored 的事件
    event DataStored(address indexed sender, uint256 data);

    uint256 private data;

    // 存储数据并触发事件
    function storeData(uint256 _data) public {
        data = _data;
        emit DataStored(msg.sender, _data); // 触发事件
    }

    // 读取数据
    function getData() view public returns (uint256) {
        return data;
    }
}