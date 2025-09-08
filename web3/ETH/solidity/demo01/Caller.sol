// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

contract Caller {
    function callStoreData(address _callee, uint256 _data) public returns (bool, bytes memory) {
        // 使用 abi.encodeWithSignature
        bytes memory payload = abi.encodeWithSignature("storeData(uint256)", _data);
        (bool success, bytes memory returnData) = _callee.call(payload);
        return (success, returnData);
    }

    function callStoreDataWithSelector(address _callee, uint256 _data) public returns (bool, bytes memory) {
        // 手动构建选择器
        bytes4 selector = bytes4(keccak256("storeData(uint256)"));
        bytes memory payload = abi.encodeWithSelector(selector, _data);
        (bool success, bytes memory returnData) = _callee.call(payload);
        return (success, returnData);
    }
}