// SPDX-License-Identifier: MIT
pragma solidity ^0.8.18;

contract FavoriteNumber {
    mapping(string => uint256) private nameToFavoriteNumber;

    function createOrUpdateFavoriteNumber(string memory name, uint256 number) public {
        nameToFavoriteNumber[name] = number;
    }

    function getNumber(string memory name) public view returns(uint256){
        return nameToFavoriteNumber[name];
    }

}