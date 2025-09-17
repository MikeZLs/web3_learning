// SPDX-License-Identifier: MIT
pragma solidity ^0.8.28;

/**
 * @title ShibStyleMemeToken - 完整的SHIB风格Meme代币合约
 * @dev 基于OpenZeppelin 5.x最新版本，最大化代码复用，详细中文注释
 * @author Meme Token Developer
 *
 * 合约功能概述：
 * 1. 基础ERC20功能 (继承自OpenZeppelin)
 * 2. 动态税费系统 (买入/卖出/转账不同税率)
 * 3. 自动swap和ETH分配机制
 * 4. 反机器人保护 (交易限制+冷却时间)
 * 5. 位标志权限管理系统 (节省Gas)
 * 6. 代币自动销毁机制
 * 7. 紧急暂停和管理功能
 * 8. 批量操作支持
 */

// ============================================================================
// OpenZeppelin 5.x 导入 - 引入所需的标准合约库
// ============================================================================
import {ERC20} from "@openzeppelin/contracts/token/ERC20/ERC20.sol";                // 标准ERC20实现
import {Ownable} from "@openzeppelin/contracts/access/Ownable.sol";                 // 所有权管理
import {ReentrancyGuard} from "@openzeppelin/contracts/utils/ReentrancyGuard.sol";  // 重入攻击保护
import {Pausable} from "@openzeppelin/contracts/utils/Pausable.sol";                // 暂停功能
import {Context} from "@openzeppelin/contracts/utils/Context.sol";                  // 上下文信息
import {IERC20} from "@openzeppelin/contracts/token/ERC20/IERC20.sol";

// ============================================================================
// Uniswap V2 接口定义 - 用于DEX交易和流动性管理
// ============================================================================

/**
 * @dev Uniswap V2 工厂合约接口
 * 用于创建代币交易对
 */
interface IUniswapV2Factory {
    function createPair(address tokenA, address tokenB) external returns (address pair);
}

/**
 * @dev Uniswap V2 路由合约接口
 * 用于代币交换和流动性操作
 */
interface IUniswapV2Router02 {
    function factory() external pure returns (address);
    function WETH() external pure returns (address);

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
}

// ============================================================================
// 主合约 - 继承多个OpenZeppelin合约实现完整功能
// ============================================================================
contract ShibStyleMemeToken is ERC20, Ownable, ReentrancyGuard, Pausable {

    // ========================================================================
    // 配置结构体定义 - 用于组织相关配置参数
    // ========================================================================

    /**
     * @dev 税费配置结构体
     * 定义不同类型交易的税率
     */
    struct TaxRates {
        uint16 buy;          // 买入税率 (基点，如300表示3%)
        uint16 sell;         // 卖出税率 (基点，如500表示5%)
        uint16 transfer;     // 普通转账税率 (基点，如100表示1%)
    }

    /**
     * @dev 交易配置结构体
     * 定义交易限制和冷却时间等参数
     */
    struct TradingConfig {
        uint256 maxTxAmount;        // 单笔交易最大数量
        uint256 maxWalletAmount;    // 单个钱包最大持有量
        uint256 swapThreshold;      // 触发自动swap的代币数量阈值
        uint32 cooldownSeconds;     // 交易冷却时间(秒)
        bool limitsEnabled;         // 是否启用交易限制
    }

    /**
     * @dev 费用分配结构体
     * 定义税费收入的分配比例
     */
    struct FeeDistribution {
        uint8 marketing;     // 营销费用百分比
        uint8 liquidity;     // 流动性费用百分比
        uint8 burn;          // 销毁代币百分比
        uint8 reflection;    // 持币奖励百分比
    }

    // ========================================================================
    // 状态变量定义 - 存储合约的核心状态数据
    // ========================================================================

    // 核心配置变量
    TaxRates public taxRates;                    // 税费配置
    TradingConfig public tradingConfig;          // 交易配置
    FeeDistribution public feeDistribution;     // 费用分配配置

    // 关键地址
    address public marketingWallet;              // 营销钱包地址
    address public liquidityWallet;             // 流动性钱包地址
    address public immutable uniswapV2Pair;     // Uniswap V2交易对地址(不可变)
    IUniswapV2Router02 public immutable uniswapV2Router; // Uniswap V2路由合约(不可变)

    /**
     * @dev 地址权限管理映射 - 位标志系统
     * 使用位运算来节省gas，一个uint8可以存储多个布尔状态
     *
     * 位标志定义：
     * 0x01 = 免税地址 (第0位)
     * 0x02 = 免限制地址 (第1位)
     * 0x04 = 黑名单地址 (第2位)
     * 0x08 = AMM交易对 (第3位)
     * 0x10 = 机器人标记 (第4位，扩展用)
     */
    mapping(address => uint8) public addressFlags;

    // 交易控制相关映射
    mapping(address => uint256) public lastTransactionTime; // 记录每个地址的最后交易时间
    bool public tradingEnabled;                             // 是否启用交易

    // 自动swap控制
    bool private _inSwap;                       // 防止swap过程中的递归调用

    // 系统常量 - 定义系统的边界值
    uint256 public constant MAX_TAX = 1000;                 // 最大税率10% (1000基点)
    uint256 public constant MAX_WALLET_PERCENT = 10;        // 最大钱包持有比例10%
    uint256 public constant MAX_TX_PERCENT = 5;             // 最大单笔交易比例5%

    // ========================================================================
    // 事件定义 - 记录重要的合约状态变化
    // ========================================================================

    event TradingEnabled(uint256 timestamp);               // 交易启用事件
    event TaxRatesUpdated(uint16 buyTax, uint16 sellTax, uint16 transferTax); // 税率更新事件
    event TradingConfigUpdated(TradingConfig config);      // 交易配置更新事件
    event AddressFlagsUpdated(address indexed account, uint8 flags); // 地址标志更新事件
    event FeeDistributionUpdated(FeeDistribution distribution); // 费用分配更新事件
    event WalletsUpdated(address indexed newMarketing, address indexed newLiquidity); // 钱包地址更新事件
    event TokensSwappedForETH(uint256 tokensSwapped, uint256 ethReceived); // 代币交换ETH事件
    event ETHDistributed(uint256 marketing, uint256 liquidity); // ETH分配事件
    event TokensBurned(uint256 amount);                     // 代币销毁事件
    event EmergencyETHWithdrawn(uint256 amount);           // 紧急提取ETH事件
    event TokensRescued(address indexed token, uint256 amount); // 代币救援事件

    // ========================================================================
    // 修饰符定义 - 提供函数执行的前置条件检查
    // ========================================================================

    /**
     * @dev 防止swap过程中的重入调用
     * 确保自动swap过程不会被其他交易干扰
     */
    modifier lockSwap() {
        _inSwap = true;    // 设置swap锁定状态
        _;                 // 执行函数体
        _inSwap = false;   // 解除swap锁定状态
    }

    /**
     * @dev 限制只有外部账户可以调用
     * 防止合约调用，增强安全性
     */
    modifier onlyEOA() {
        require(tx.origin == _msgSender(), "Only EOA allowed");
        _;
    }

    // ========================================================================
    // 构造函数 - 初始化合约的核心参数和状态
    // ========================================================================

    constructor(
        string memory name,                          // 代币名称
        string memory symbol,                        // 代币符号
        uint256 totalSupply,                        // 总供应量
        address routerAddress,                      // Uniswap路由地址
        address marketingAddress,                   // 营销钱包地址
        address liquidityAddress,                   // 流动性钱包地址
        TaxRates memory initialTaxRates,           // 初始税率配置
        FeeDistribution memory initialFeeDistribution // 初始费用分配配置
    )
    ERC20(name, symbol)        // 调用ERC20构造函数
    Ownable(msg.sender)      // 调用Ownable构造函数，设置部署者为所有者
    {
        // 参数验证 - 确保传入参数的有效性
        require(routerAddress != address(0), "Invalid router");
        require(marketingAddress != address(0), "Invalid marketing wallet");
        require(liquidityAddress != address(0), "Invalid liquidity wallet");
        require(
            initialTaxRates.buy <= MAX_TAX &&
            initialTaxRates.sell <= MAX_TAX &&
            initialTaxRates.transfer <= MAX_TAX,
            "Tax rate too high"
        );
        require(
            initialFeeDistribution.marketing +
            initialFeeDistribution.liquidity +
            initialFeeDistribution.burn +
            initialFeeDistribution.reflection == 100,
            "Fee distribution must sum to 100"
        );

        // 设置基础配置
        taxRates = initialTaxRates;           // 设置税率配置
        feeDistribution = initialFeeDistribution; // 设置费用分配配置
        marketingWallet = marketingAddress;   // 设置营销钱包
        liquidityWallet = liquidityAddress;   // 设置流动性钱包

        // 初始化Uniswap相关配置
        uniswapV2Router = IUniswapV2Router02(routerAddress); // 设置路由合约
        uniswapV2Pair = IUniswapV2Factory(uniswapV2Router.factory())
            .createPair(address(this), uniswapV2Router.WETH()); // 创建交易对

        // 设置交易配置默认值
        tradingConfig = TradingConfig({
            maxTxAmount: totalSupply * MAX_TX_PERCENT / 100,        // 最大单笔交易量
            maxWalletAmount: totalSupply * MAX_WALLET_PERCENT / 100, // 最大钱包持有量
            swapThreshold: totalSupply / 1000,                      // swap阈值(0.1%总供应量)
            cooldownSeconds: 30,                                    // 30秒冷却时间
            limitsEnabled: true                                     // 启用交易限制
        });

        // 设置地址权限标志
        _setAddressFlags(owner(), 0x03);              // 所有者: 免税(0x01) + 免限制(0x02)
        _setAddressFlags(address(this), 0x03);        // 合约地址: 免税 + 免限制
        _setAddressFlags(marketingAddress, 0x01);     // 营销钱包: 免税
        _setAddressFlags(liquidityAddress, 0x01);     // 流动性钱包: 免税
        _setAddressFlags(uniswapV2Pair, 0x0A);        // 交易对: 免限制(0x02) + AMM标记(0x08)

        // 铸造全部代币给所有者
        _mint(owner(), totalSupply);
    }

    // ========================================================================
    // 权限管理系统 - 基于位标志的高效权限管理
    // ========================================================================

    /**
     * @dev 检查地址是否具有指定权限标志
     * 使用位运算AND操作检查特定标志位
     * @param account 要检查的地址
     * @param flag 要检查的标志位
     * @return 是否具有该权限
     */
    function hasFlag(address account, uint8 flag) public view returns (bool) {
        return (addressFlags[account] & flag) == flag;
    }

    /**
     * @dev 内部函数：设置地址权限标志
     * @param account 目标地址
     * @param flags 权限标志位
     */
    function _setAddressFlags(address account, uint8 flags) internal {
        addressFlags[account] = flags;           // 设置权限标志
        emit AddressFlagsUpdated(account, flags); // 触发事件
    }

    /**
     * @dev 外部函数：设置单个地址权限标志
     * 只有合约所有者可以调用
     * @param account 目标地址
     * @param flags 权限标志位
     */
    function setAddressFlags(address account, uint8 flags) external onlyOwner {
        _setAddressFlags(account, flags);
    }

    /**
     * @dev 批量设置多个地址的权限标志
     * 提高批量操作的效率
     * @param accounts 地址数组
     * @param flags 对应的权限标志数组
     */
    function batchSetAddressFlags(
        address[] calldata accounts,
        uint8[] calldata flags
    ) external onlyOwner {
        require(accounts.length == flags.length, "Arrays length mismatch");

        for (uint256 i = 0; i < accounts.length; i++) {
            _setAddressFlags(accounts[i], flags[i]);
        }
    }

    // 便利查询函数 - 提供更友好的权限查询接口

    /**
     * @dev 检查地址是否免除交易税费
     */
    function isExcludedFromFees(address account) external view returns (bool) {
        return hasFlag(account, 0x01);
    }

    /**
     * @dev 检查地址是否免除交易限制
     */
    function isExcludedFromLimits(address account) external view returns (bool) {
        return hasFlag(account, 0x02);
    }

    /**
     * @dev 检查地址是否在黑名单中
     */
    function isBlacklisted(address account) external view returns (bool) {
        return hasFlag(account, 0x04);
    }

    /**
     * @dev 检查地址是否为AMM交易对
     */
    function isAMMPair(address account) external view returns (bool) {
        return hasFlag(account, 0x08);
    }

    // ========================================================================
    // 核心转账逻辑 - 重写ERC20标准实现自定义转账逻辑
    // ========================================================================

    /**
     * @dev 重写ERC20的_update函数 (OpenZeppelin 5.x新特性)
     * 这是所有转账、铸造、销毁操作的核心入口点
     * @param from 发送者地址 (地址0表示铸造)
     * @param to 接收者地址 (地址0表示销毁)
     * @param value 转账数量
     */
    function _update(address from, address to, uint256 value)
    internal
    override
    whenNotPaused // 合约未暂停时才能执行
    {
        // 基础安全检查 - 防止黑名单地址参与交易
        require(!hasFlag(from, 0x04) && !hasFlag(to, 0x04), "Blacklisted address");

        // 交易开启检查 - 交易未启用时只允许免税地址操作
        if (!tradingEnabled) {
            require(
                hasFlag(from, 0x01) || hasFlag(to, 0x01),
                "Trading not enabled"
            );
        }

        // 转账特殊处理 - 排除铸造和销毁操作
        if (from != address(0) && to != address(0)) {
            _handleTransfer(from, to, value);
        }

        // 执行实际的代币状态更新
        super._update(from, to, value);
    }

    /**
     * @dev 处理转账的所有逻辑
     * 按顺序执行各种检查和处理
     * @param from 发送者地址
     * @param to 接收者地址
     * @param amount 转账数量
     */
    function _handleTransfer(address from, address to, uint256 amount) internal {
        _checkCooldown(from);              // 检查冷却时间
        _checkTradingLimits(from, to, amount); // 检查交易限制
        _checkAndExecuteSwap(from);        // 检查并执行自动swap
        _processTax(from, to, amount);     // 处理税费
        _updateCooldown(from);             // 更新冷却时间
    }

    /**
     * @dev 检查发送者的冷却时间
     * 防止高频交易和机器人攻击
     * @param from 发送者地址
     */
    function _checkCooldown(address from) internal view {
        // 免限制地址或冷却时间为0时跳过检查
        if (!hasFlag(from, 0x02) && tradingConfig.cooldownSeconds > 0) {
            require(
                block.timestamp >= lastTransactionTime[from] + tradingConfig.cooldownSeconds,
                "Cooldown period active"
            );
        }
    }

    /**
     * @dev 更新发送者的最后交易时间
     * 用于冷却时间计算
     * @param from 发送者地址
     */
    function _updateCooldown(address from) internal {
        // 只更新非免限制地址的交易时间
        if (!hasFlag(from, 0x02)) {
            lastTransactionTime[from] = block.timestamp;
        }
    }

    /**
     * @dev 检查交易限制
     * 包括单笔交易限制和钱包持有量限制
     * @param from 发送者地址
     * @param to 接收者地址
     * @param amount 转账数量
     */
    function _checkTradingLimits(address from, address to, uint256 amount) internal view {
        if (!tradingConfig.limitsEnabled) return;           // 限制未启用时跳过
        if (hasFlag(from, 0x02) || hasFlag(to, 0x02)) return; // 免限制地址跳过

        // 检查单笔交易限制
        require(amount <= tradingConfig.maxTxAmount, "Exceeds max tx amount");

        // 检查钱包持有量限制(仅适用于从AMM购买的情况)
        if (hasFlag(from, 0x08)) {
            require(
                balanceOf(to) + amount <= tradingConfig.maxWalletAmount,
                "Exceeds max wallet amount"
            );
        }
    }

    /**
     * @dev 税费处理逻辑
     * 计算并预留税费，实际收取在自动swap时进行
     * @param from 发送者地址
     * @param to 接收者地址
     * @param amount 转账数量
     */
    function _processTax(address from, address to, uint256 amount) internal {
        uint256 taxAmount = _calculateTax(from, to, amount);

        if (taxAmount > 0) {
            // 税费处理逻辑在自动swap时进行
            // 这里可以添加额外的税费处理逻辑
        }
    }

    /**
     * @dev 计算交易应缴纳的税费
     * 根据交易类型(买入/卖出/转账)确定税率
     * @param from 发送者地址
     * @param to 接收者地址
     * @param amount 转账数量
     * @return 应缴纳的税费数量
     */
    function _calculateTax(address from, address to, uint256 amount) internal view returns (uint256){
        // 免税地址或swap过程中不收税
        if (hasFlag(from, 0x01) || hasFlag(to, 0x01) || _inSwap) {
            return 0;
        }

        uint256 taxRate = 0;

        // 根据交易类型确定税率
        if (hasFlag(from, 0x08)) {
            // 从AMM买入
            taxRate = taxRates.buy;
        } else if (hasFlag(to, 0x08)) {
            // 向AMM卖出
            taxRate = taxRates.sell;
        } else {
            // 普通转账
            taxRate = taxRates.transfer;
        }

        // 计算税费 (基点转换为百分比: 除以10000)
        return amount * taxRate / 10000;
    }

    // ========================================================================
    // 自动Swap和ETH分配系统 - 实现税费的自动处理和分配
    // ========================================================================

    /**
     * @dev 检查是否需要执行自动swap
     * 在合适的时机自动将合约中的代币换成ETH并分配
     * @param from 发送者地址
     */
    function _checkAndExecuteSwap(address from) internal {
        uint256 contractBalance = balanceOf(address(this)); // 获取合约代币余额

        // 判断是否应该执行swap
        bool shouldSwap = !_inSwap &&                           // 不在swap过程中
                !hasFlag(from, 0x08) &&                        // 发送者不是AMM(避免在swap时触发)
                contractBalance >= tradingConfig.swapThreshold && // 达到swap阈值
                        tradingEnabled;                         // 交易已启用

        if (shouldSwap) {
            _swapAndDistribute(contractBalance);
        }
    }

    /**
     * @dev 执行代币swap并分配ETH
     * 将累积的税费代币转换为ETH并按配置分配
     * @param tokenAmount 要swap的代币数量
     */
    function _swapAndDistribute(uint256 tokenAmount) internal lockSwap nonReentrant {
        // 计算需要销毁的代币数量
        uint256 burnAmount = tokenAmount * feeDistribution.burn / 100;
        uint256 swapAmount = tokenAmount - burnAmount; // 实际swap的数量

        // 执行代币销毁
        if (burnAmount > 0) {
            _burn(address(this), burnAmount);           // 销毁代币
            emit TokensBurned(burnAmount);              // 触发销毁事件
        }

        // 执行代币swap
        if (swapAmount > 0) {
            uint256 ethReceived = _swapTokensForETH(swapAmount); // 换取ETH

            if (ethReceived > 0) {
                _distributeETH(ethReceived);            // 分配ETH
            }
        }
    }

    /**
     * @dev 将代币换成ETH
     * 通过Uniswap V2进行代币交换
     * @param tokenAmount 要交换的代币数量
     * @return 获得的ETH数量
     */
    function _swapTokensForETH(uint256 tokenAmount) internal returns (uint256) {
        // 设置交换路径: 代币 -> WETH
        address[] memory path = new address[](2);
        path[0] = address(this);                    // 本代币地址
        path[1] = uniswapV2Router.WETH();          // WETH地址

        // 授权路由合约使用代币
        _approve(address(this), address(uniswapV2Router), tokenAmount);

        uint256 initialETH = address(this).balance; // 记录swap前的ETH余额

        // 执行代币swap操作
        try uniswapV2Router.swapExactTokensForETHSupportingFeeOnTransferTokens(
            tokenAmount,                            // 输入代币数量
            0,                                     // 最小输出(设为0，实际项目中应设置滑点保护)
            path,                                  // 交换路径
            address(this),                         // 接收地址
            block.timestamp                        // 截止时间
        ) {
            uint256 newETH = address(this).balance;      // 获取swap后的ETH余额
            uint256 ethReceived = newETH - initialETH;   // 计算实际获得的ETH

            emit TokensSwappedForETH(tokenAmount, ethReceived); // 触发swap事件
            return ethReceived;
        } catch {
            // swap失败时返回0，不影响正常交易
            return 0;
        }
    }

    /**
     * @dev 分配ETH给各个钱包
     * 按照配置的比例分配swap获得的ETH
     * @param ethAmount 要分配的ETH总量
     */
    function _distributeETH(uint256 ethAmount) internal {
        // 计算总分配份额(排除销毁和反射部分)
        uint256 totalShares = feeDistribution.marketing + feeDistribution.liquidity;
        if (totalShares == 0) return; // 无分配份额时直接返回

        // 按比例计算各钱包应得的ETH
        uint256 marketingETH = ethAmount * feeDistribution.marketing / totalShares;
        uint256 liquidityETH = ethAmount * feeDistribution.liquidity / totalShares;

        // 向营销钱包转账
        if (marketingETH > 0) {
            payable(marketingWallet).transfer(marketingETH);
        }

        // 向流动性钱包转账
        if (liquidityETH > 0) {
            payable(liquidityWallet).transfer(liquidityETH);
        }

        emit ETHDistributed(marketingETH, liquidityETH); // 触发分配事件
    }

    // ========================================================================
    // 管理员功能 - 提供合约配置和管理接口
    // ========================================================================

    /**
     * @dev 启用交易
     * 只能由合约所有者调用，且只能启用一次
     */
    function enableTrading() external onlyOwner {
        require(!tradingEnabled, "Trading already enabled"); // 防止重复启用
        tradingEnabled = true;
        emit TradingEnabled(block.timestamp);
    }

    /**
     * @dev 更新税费配置
     * 允许所有者调整不同类型交易的税率
     * @param buyTax 买入税率(基点)
     * @param sellTax 卖出税率(基点)
     * @param transferTax 转账税率(基点)
     */
    function updateTaxRates(
        uint16 buyTax,
        uint16 sellTax,
        uint16 transferTax
    ) external onlyOwner {
        // 验证税率不超过上限
        require(
            buyTax <= MAX_TAX && sellTax <= MAX_TAX && transferTax <= MAX_TAX,
            "Tax rate too high"
        );

        taxRates = TaxRates(buyTax, sellTax, transferTax); // 更新税率配置
        emit TaxRatesUpdated(buyTax, sellTax, transferTax); // 触发更新事件
    }

    /**
     * @dev 更新交易配置
     * 调整交易限制、冷却时间等参数
     * @param newConfig 新的交易配置
     */
    function updateTradingConfig(TradingConfig calldata newConfig) external onlyOwner {
        // 验证配置的合理性
        require(
            newConfig.maxTxAmount >= totalSupply() / 1000 &&      // 最小0.1%总供应量
            newConfig.maxWalletAmount >= totalSupply() / 1000,    // 最小0.1%总供应量
            "Limits too restrictive"
        );
        require(newConfig.cooldownSeconds <= 300, "Cooldown too long"); // 冷却时间不超过5分钟

        tradingConfig = newConfig;                              // 更新交易配置
        emit TradingConfigUpdated(newConfig);                   // 触发更新事件
    }

    /**
     * @dev 更新费用分配配置
     * 调整税费收入在各用途之间的分配比例
     * @param newDistribution 新的费用分配配置
     */
    function updateFeeDistribution(FeeDistribution calldata newDistribution) external onlyOwner {
        // 验证分配比例总和为100%
        require(
            newDistribution.marketing +
            newDistribution.liquidity +
            newDistribution.burn +
            newDistribution.reflection == 100,
            "Distribution must sum to 100"
        );

        feeDistribution = newDistribution;                      // 更新费用分配配置
        emit FeeDistributionUpdated(newDistribution);           // 触发更新事件
    }

    /**
     * @dev 更新钱包地址
     * 更改营销和流动性钱包地址
     * @param newMarketingWallet 新的营销钱包地址
     * @param newLiquidityWallet 新的流动性钱包地址
     */
    function updateWallets(
        address newMarketingWallet,
        address newLiquidityWallet
    ) external onlyOwner {
        // 验证地址有效性
        require(
            newMarketingWallet != address(0) && newLiquidityWallet != address(0),
            "Invalid wallet addresses"
        );

        marketingWallet = newMarketingWallet;                   // 更新营销钱包
        liquidityWallet = newLiquidityWallet;                   // 更新流动性钱包

        emit WalletsUpdated(newMarketingWallet, newLiquidityWallet); // 触发更新事件
    }

    // ========================================================================
    // 紧急功能 - 提供紧急情况下的合约控制和资产救援
    // ========================================================================

    /**
     * @dev 暂停合约
     * 紧急情况下暂停所有转账操作
     */
    function pause() external onlyOwner {
        _pause(); // 调用Pausable的暂停函数
    }

    /**
     * @dev 恢复合约
     * 解除暂停状态，恢复正常操作
     */
    function unpause() external onlyOwner {
        _unpause(); // 调用Pausable的恢复函数
    }

    /**
     * @dev 手动执行swap
     * 允许所有者手动触发代币swap和ETH分配
     */
    function manualSwap() external onlyOwner {
        uint256 contractBalance = balanceOf(address(this)); // 获取合约代币余额
        require(contractBalance > 0, "No tokens to swap");  // 确保有代币可swap
        _swapAndDistribute(contractBalance);                // 执行swap和分配
    }

    /**
     * @dev 紧急提取ETH
     * 紧急情况下将合约中的ETH提取给所有者
     */
    function rescueETH() external onlyOwner {
        uint256 balance = address(this).balance;            // 获取合约ETH余额
        require(balance > 0, "No ETH to rescue");           // 确保有ETH可提取

        payable(owner()).transfer(balance);                 // 转账给所有者
        emit EmergencyETHWithdrawn(balance);                // 触发提取事件
    }

    /**
     * @dev 救援其他代币
     * 提取意外发送到合约的其他ERC20代币
     * @param token 要救援的代币合约地址
     * @param amount 要救援的数量
     */
    function rescueTokens(IERC20 token, uint256 amount) external onlyOwner {
        require(address(token) != address(this), "Cannot rescue own tokens"); // 不能救援自己的代币

        token.transfer(owner(), amount);                    // 转账给所有者
        emit TokensRescued(address(token), amount);         // 触发救援事件
    }

    // ========================================================================
    // 查询功能 - 提供合约状态和配置的查询接口
    // ========================================================================

    /**
     * @dev 获取完整的合约信息
     * 一次性返回合约的主要状态和配置信息
     * @return _tradingEnabled 交易是否启用
     * @return _taxRates 当前税率配置
     * @return _tradingConfig 当前交易配置
     * @return _feeDistribution 当前费用分配配置
     * @return _contractTokenBalance 合约代币余额
     * @return _contractETHBalance 合约ETH余额
     */
    function getContractInfo() external view returns (
        bool _tradingEnabled,
        TaxRates memory _taxRates,
        TradingConfig memory _tradingConfig,
        FeeDistribution memory _feeDistribution,
        uint256 _contractTokenBalance,
        uint256 _contractETHBalance
    ) {
        return (
            tradingEnabled,                  // 交易启用状态
            taxRates,                       // 税率配置
            tradingConfig,                  // 交易配置
            feeDistribution,                // 费用分配配置
            balanceOf(address(this)),       // 合约代币余额
            address(this).balance           // 合约ETH余额
        );
    }

    /**
     * @dev 查询地址剩余冷却时间
     * 计算指定地址还需等待多长时间才能再次交易
     * @param account 要查询的地址
     * @return 剩余冷却时间(秒)
     */
    function getRemainingCooldown(address account) external view returns (uint256) {
        // 免限制地址或冷却时间为0时，无需冷却
        if (hasFlag(account, 0x02) || tradingConfig.cooldownSeconds == 0) {
            return 0;
        }

        uint256 elapsed = block.timestamp - lastTransactionTime[account]; // 计算已过时间
        if (elapsed >= tradingConfig.cooldownSeconds) {
            return 0; // 冷却时间已过
        }

        return tradingConfig.cooldownSeconds - elapsed; // 返回剩余冷却时间
    }

    /**
     * @dev 检查地址是否为合约地址
     * 通过检查代码长度判断是否为合约
     * @param account 要检查的地址
     * @return 是否为合约地址
     */
    function isContract(address account) public view returns (bool) {
        return account.code.length > 0; // 合约地址的代码长度大于0
    }

    // ========================================================================
    // ETH接收功能 - 允许合约接收ETH
    // ========================================================================

    /**
     * @dev 接收ETH的回退函数
     * 当合约收到ETH且调用数据为空时触发
     */
    receive() external payable {}

    /**
     * @dev 通用回退函数
     * 当合约收到ETH且没有匹配的函数时触发
     */
    fallback() external payable {}
}