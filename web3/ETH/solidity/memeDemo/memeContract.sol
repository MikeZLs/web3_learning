// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

/**
 * @title ShibStyleMemeToken - SHIB风格Meme代币合约
 * @dev 实现代币税、流动性池集成、交易限制等核心功能
 * @author Lance
 */

interface IERC20 {
    function totalSupply() external view returns (uint256);
    function balanceOf(address account) external view returns (uint256);
    function transfer(address recipient, uint256 amount) external returns (bool);
    function allowance(address owner, address spender) external view returns (uint256);
    function approve(address spender, uint256 amount) external returns (bool);
    function transferFrom(address sender, address recipient, uint256 amount) external returns (bool);

    event Transfer(address indexed from, address indexed to, uint256 value);
    event Approval(address indexed owner, address indexed spender, uint256 value);
}

// 简化的Uniswap接口
interface IUniswapV2Router {
    function addLiquidityETH(
        address token,
        uint amountTokenDesired,
        uint amountTokenMin,
        uint amountETHMin,
        address to,
        uint deadline
    ) external payable returns (uint amountToken, uint amountETH, uint liquidity);

    function swapExactTokensForETHSupportingFeeOnTransferTokens(
        uint amountIn,
        uint amountOutMin,
        address[] calldata path,
        address to,
        uint deadline
    ) external;

    function factory() external pure returns (address);
    function WETH() external pure returns (address);
}

interface IUniswapV2Factory {
    function createPair(address tokenA, address tokenB) external returns (address pair);
}

contract ShibStyleMemeToken is IERC20 {
    // ============ 基础代币信息 ============
    string public constant name = "ShibStyle";
    string public constant symbol = "SHBS";
    uint8 public constant decimals = 18;
    uint256 private _totalSupply = 1000000000 * 10**decimals; // 10亿代币

    // ============ 余额和授权映射 ============
    mapping(address => uint256) private _balances;
    mapping(address => mapping(address => uint256)) private _allowances;

    // ============ 税费和限制参数 ============
    uint256 public buyTax = 300;  // 3% 买入税
    uint256 public sellTax = 500; // 5% 卖出税
    uint256 public maxTxAmount = _totalSupply * 2 / 100; // 最大交易额：2%
    uint256 public maxWalletAmount = _totalSupply * 3 / 100; // 最大钱包持有量：3%

    // ============ 地址管理 ============
    address public owner;
    address public marketingWallet; // 营销钱包
    address public liquidityWallet; // 流动性钱包

    // ============ Uniswap集成 ============
    IUniswapV2Router public immutable uniswapV2Router;
    address public immutable uniswapV2Pair;

    // ============ 状态管理 ============
    bool public tradingEnabled = false; // 交易开关
    bool private inSwap = false; // 防重入标志
    uint256 public swapTokensAtAmount = _totalSupply / 1000; // 0.1% 触发swap

    // ============ 白名单和黑名单 ============
    mapping(address => bool) public isExcludedFromFees; // 免税名单
    mapping(address => bool) public isExcludedFromLimits; // 免限制名单
    mapping(address => bool) public blacklist; // 黑名单

    // ============ 交易冷却 ============
    mapping(address => uint256) public lastTransactionTime; // 最后交易时间
    uint256 public cooldownTime = 30; // 冷却时间30秒

    // ============ 事件定义 ============
    event TradingEnabled();
    event TaxesUpdated(uint256 buyTax, uint256 sellTax);
    event LimitsUpdated(uint256 maxTxAmount, uint256 maxWalletAmount);
    event TokensSwapped(uint256 tokensSwapped, uint256 ethReceived);

    // ============ 修饰符 ============
    modifier onlyOwner() {
        require(msg.sender == owner, "Not owner");
        _;
    }

    modifier lockTheSwap() {
        inSwap = true;
        _;
        inSwap = false;
    }

    /**
     * @dev 构造函数：初始化代币合约
     * @param _router Uniswap路由器地址
     * @param _marketingWallet 营销钱包地址
     */
    constructor(address _router, address _marketingWallet) {
        owner = msg.sender;
        marketingWallet = _marketingWallet;
        liquidityWallet = msg.sender;

        // 初始化Uniswap
        uniswapV2Router = IUniswapV2Router(_router);
        uniswapV2Pair = IUniswapV2Factory(uniswapV2Router.factory())
            .createPair(address(this), uniswapV2Router.WETH());

        // 设置免税和免限制地址
        isExcludedFromFees[owner] = true;
        isExcludedFromFees[address(this)] = true;
        isExcludedFromFees[marketingWallet] = true;

        isExcludedFromLimits[owner] = true;
        isExcludedFromLimits[address(this)] = true;
        isExcludedFromLimits[marketingWallet] = true;
        isExcludedFromLimits[uniswapV2Pair] = true;

        // 将所有代币分配给创建者
        _balances[owner] = _totalSupply;
        emit Transfer(address(0), owner, _totalSupply);
    }

    // ============ ERC20基础功能 ============

    function totalSupply() public view override returns (uint256) {
        return _totalSupply;
    }

    function balanceOf(address account) public view override returns (uint256) {
        return _balances[account];
    }

    function transfer(address recipient, uint256 amount) public override returns (bool) {
        _transfer(msg.sender, recipient, amount);
        return true;
    }

    function allowance(address owner, address spender) public view override returns (uint256) {
        return _allowances[owner][spender];
    }

    function approve(address spender, uint256 amount) public override returns (bool) {
        _approve(msg.sender, spender, amount);
        return true;
    }

    function transferFrom(address sender, address recipient, uint256 amount) public override returns (bool) {
        uint256 currentAllowance = _allowances[sender][msg.sender];
        require(currentAllowance >= amount, "Transfer amount exceeds allowance");

        _transfer(sender, recipient, amount);
        _approve(sender, msg.sender, currentAllowance - amount);

        return true;
    }

    /**
     * @dev 内部转账函数：实现税费和限制逻辑
     */
    function _transfer(address from, address to, uint256 amount) internal {
        require(from != address(0), "Transfer from zero address");
        require(to != address(0), "Transfer to zero address");
        require(!blacklist[from] && !blacklist[to], "Address blacklisted");
        require(amount > 0, "Transfer amount must be greater than zero");

        // 检查交易是否开启（免税地址可以在未开启时交易）
        if (!tradingEnabled) {
            require(isExcludedFromFees[from] || isExcludedFromFees[to], "Trading not enabled");
        }

        // 检查冷却时间
        if (!isExcludedFromLimits[from]) {
            require(block.timestamp >= lastTransactionTime[from] + cooldownTime, "Cooldown period active");
            lastTransactionTime[from] = block.timestamp;
        }

        // 检查交易限制
        if (!isExcludedFromLimits[from] && !isExcludedFromLimits[to]) {
            require(amount <= maxTxAmount, "Transfer amount exceeds max tx amount");

            // 检查钱包最大持有量（买入时）
            if (from == uniswapV2Pair) {
                require(_balances[to] + amount <= maxWalletAmount, "Max wallet amount exceeded");
            }
        }

        // 检查是否需要swap合约内的代币
        uint256 contractTokenBalance = _balances[address(this)];
        bool canSwap = contractTokenBalance >= swapTokensAtAmount;

        if (canSwap && !inSwap && from != uniswapV2Pair && !isExcludedFromFees[from]) {
            _swapTokensForEth(contractTokenBalance);
        }

        // 计算是否收税
        bool takeFee = !inSwap && !isExcludedFromFees[from] && !isExcludedFromFees[to];

        if (takeFee) {
            uint256 fees = _calculateFees(from, to, amount);
            if (fees > 0) {
                amount = amount - fees;
                _balances[from] -= fees;
                _balances[address(this)] += fees;
                emit Transfer(from, address(this), fees);
            }
        }

        // 执行转账
        _balances[from] -= amount;
        _balances[to] += amount;

        emit Transfer(from, to, amount);
    }

    /**
     * @dev 计算交易手续费
     */
    function _calculateFees(address from, address to, uint256 amount) internal view returns (uint256) {
        uint256 feeRate = 0;

        // 买入（从Uniswap pair转入）
        if (from == uniswapV2Pair) {
            feeRate = buyTax;
        }
            // 卖出（转出到Uniswap pair）
        else if (to == uniswapV2Pair) {
            feeRate = sellTax;
        }

        return amount * feeRate / 10000; // 基点计算（10000 = 100%）
    }

    /**
     * @dev 将合约内的代币换成ETH并分配
     */
    function _swapTokensForEth(uint256 tokenAmount) private lockTheSwap {
        address[] memory path = new address[](2);
        path[0] = address(this);
        path[1] = uniswapV2Router.WETH();

        _approve(address(this), address(uniswapV2Router), tokenAmount);

        try uniswapV2Router.swapExactTokensForETHSupportingFeeOnTransferTokens(
            tokenAmount,
            0,
            path,
            address(this),
            block.timestamp
        ) {
            uint256 ethBalance = address(this).balance;
            if (ethBalance > 0) {
                // 70%给营销，30%添加流动性
                uint256 marketingEth = ethBalance * 70 / 100;
                uint256 liquidityEth = ethBalance - marketingEth;

                if (marketingEth > 0) {
                    payable(marketingWallet).transfer(marketingEth);
                }
                if (liquidityEth > 0) {
                    payable(liquidityWallet).transfer(liquidityEth);
                }
            }

            emit TokensSwapped(tokenAmount, ethBalance);
        } catch {
            // Swap失败时静默处理
        }
    }

    function _approve(address owner, address spender, uint256 amount) internal {
        require(owner != address(0), "Approve from zero address");
        require(spender != address(0), "Approve to zero address");

        _allowances[owner][spender] = amount;
        emit Approval(owner, spender, amount);
    }

    // ============ 管理员功能 ============

    /**
     * @dev 开启交易
     */
    function enableTrading() external onlyOwner {
        require(!tradingEnabled, "Trading already enabled");
        tradingEnabled = true;
        emit TradingEnabled();
    }

    /**
     * @dev 更新税率
     */
    function updateTaxes(uint256 _buyTax, uint256 _sellTax) external onlyOwner {
        require(_buyTax <= 1000 && _sellTax <= 1000, "Tax too high"); // 最大10%
        buyTax = _buyTax;
        sellTax = _sellTax;
        emit TaxesUpdated(_buyTax, _sellTax);
    }

    /**
     * @dev 更新交易限制
     */
    function updateLimits(uint256 _maxTxAmount, uint256 _maxWalletAmount) external onlyOwner {
        require(_maxTxAmount >= _totalSupply / 1000, "Max tx amount too low"); // 至少0.1%
        require(_maxWalletAmount >= _totalSupply / 1000, "Max wallet amount too low");
        maxTxAmount = _maxTxAmount;
        maxWalletAmount = _maxWalletAmount;
        emit LimitsUpdated(_maxTxAmount, _maxWalletAmount);
    }

    /**
     * @dev 设置免税地址
     */
    function setExcludedFromFees(address account, bool excluded) external onlyOwner {
        isExcludedFromFees[account] = excluded;
    }

    /**
     * @dev 设置免限制地址
     */
    function setExcludedFromLimits(address account, bool excluded) external onlyOwner {
        isExcludedFromLimits[account] = excluded;
    }

    /**
     * @dev 设置黑名单
     */
    function setBlacklist(address account, bool blacklisted) external onlyOwner {
        blacklist[account] = blacklisted;
    }

    /**
     * @dev 更新冷却时间
     */
    function setCooldownTime(uint256 _cooldownTime) external onlyOwner {
        require(_cooldownTime <= 300, "Cooldown time too long"); // 最大5分钟
        cooldownTime = _cooldownTime;
    }

    /**
     * @dev 更新营销钱包
     */
    function setMarketingWallet(address _marketingWallet) external onlyOwner {
        marketingWallet = _marketingWallet;
    }

    /**
     * @dev 手动触发swap
     */
    function manualSwap() external onlyOwner {
        uint256 contractBalance = _balances[address(this)];
        require(contractBalance > 0, "No tokens to swap");
        _swapTokensForEth(contractBalance);
    }

    /**
     * @dev 紧急提取ETH
     */
    function emergencyWithdrawETH() external onlyOwner {
        uint256 balance = address(this).balance;
        require(balance > 0, "No ETH to withdraw");
        payable(owner).transfer(balance);
    }

    /**
     * @dev 转移所有权
     */
    function transferOwnership(address newOwner) external onlyOwner {
        require(newOwner != address(0), "New owner cannot be zero address");
        owner = newOwner;
    }

    // ============ 接收ETH ============
    receive() external payable {}

    // ============ 查询功能 ============

    /**
     * @dev 获取合约ETH余额
     */
    function getContractETHBalance() external view returns (uint256) {
        return address(this).balance;
    }

    /**
     * @dev 获取合约代币余额
     */
    function getContractTokenBalance() external view returns (uint256) {
        return _balances[address(this)];
    }

    /**
     * @dev 查询用户剩余冷却时间
     */
    function getRemainingCooldown(address user) external view returns (uint256) {
        if (isExcludedFromLimits[user]) return 0;
        uint256 timePassed = block.timestamp - lastTransactionTime[user];
        if (timePassed >= cooldownTime) return 0;
        return cooldownTime - timePassed;
    }
}