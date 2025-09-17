# ShibStyleMemeToken - SHIB风格Meme代币合约

##  项目概述

ShibStyleMemeToken是一个功能完整的SHIB风格Meme代币智能合约，基于OpenZeppelin 5.x构建，具备现代DeFi代币的所有核心功能。

### 核心特性

- **动态税费系统** - 买入/卖出/转账不同税率
- **自动Swap机制** - 税费自动转换为ETH并分配
- **反机器人保护** - 交易限制和冷却时间
- **权限管理系统** - 基于位标志的高效权限控制
- **代币销毁机制** - 自动通缩模型
- **紧急控制功能** - 暂停、救援等安全机制

##  快速部署指南

### 前置要求

- **Solidity版本**: ^0.8.20 或更高
- **开发环境**: Remix IDE 或 Hardhat
- **网络**: 以太坊主网或测试网
- **钱包**: MetaMask 或其他Web3钱包

### 部署参数

#### 基础参数
```
name: "ShibToken"
symbol: "SHIB"
totalSupply: 1000000000000000000000000000  // 10亿代币(18位小数)
```

#### 网络地址
```javascript
// 主网
routerAddress: "0x7a250d5630B4cF539739dF2C5dAcb4c659F2488D"

// 测试网 (Sepolia)
routerAddress: "0xC532a74256D3Db42D0Bf7a0400fEFDbad7694008"
```

#### 钱包配置
```
marketingAddress: "你的营销钱包地址"
liquidityAddress: "你的流动性钱包地址"
```

#### 税率配置 (推荐)
```javascript
initialTaxRates: [300, 500, 100]
// 300 = 3% 买入税
// 500 = 5% 卖出税  
// 100 = 1% 转账税
```

#### 费用分配 (推荐)
```javascript
initialFeeDistribution: [40, 30, 20, 10]
// 40 = 40% 营销
// 30 = 30% 流动性
// 20 = 20% 代币销毁
// 10 = 10% 持币奖励
```

### 部署步骤

#### 方法1: Remix IDE部署

1. **打开Remix**: https://remix.ethereum.org
2. **创建新文件**: `ShibStyleMemeToken.sol`
3. **粘贴合约代码**
4. **编译合约**:
    - 选择Solidity版本: 0.8.20+
    - 点击"Compile"
5. **部署合约**:
    - 切换到"Deploy & Run"
    - 选择环境(Injected Web3)
    - 填写构造函数参数
    - 点击"Deploy"

#### 方法2: Hardhat部署

1. **初始化项目**:
```bash
npm init -y
npm install --save-dev hardhat @openzeppelin/contracts
npx hardhat init
```

2. **创建部署脚本** (`scripts/deploy.js`):
```javascript
const { ethers } = require("hardhat");

async function main() {
    const [deployer] = await ethers.getSigners();
    
    console.log("Deploying contracts with account:", deployer.address);
    
    const ShibStyleMemeToken = await ethers.getContractFactory("ShibStyleMemeToken");
    
    const token = await ShibStyleMemeToken.deploy(
        "ShibToken",                                    // name
        "SHIB",                                        // symbol
        ethers.parseEther("1000000000"),               // totalSupply
        "0x7a250d5630B4cF539739dF2C5dAcb4c659F2488D", // routerAddress
        "YOUR_MARKETING_WALLET",                       // marketingAddress
        "YOUR_LIQUIDITY_WALLET",                       // liquidityAddress
        [300, 500, 100],                              // initialTaxRates
        [40, 30, 20, 10]                              // initialFeeDistribution
    );
    
    await token.waitForDeployment();
    console.log("Token deployed to:", await token.getAddress());
}

main()
    .then(() => process.exit(0))
    .catch((error) => {
        console.error(error);
        process.exit(1);
    });
```

3. **执行部署**:
```bash
npx hardhat run scripts/deploy.js --network mainnet
```

##  合约管理

### 启用交易

部署后合约默认禁用交易，需要手动启用：

```javascript
// 调用enableTrading函数
await contract.enableTrading();
```

### 权限管理

#### 设置地址权限
```javascript
// 设置免税地址
await contract.setAddressFlags(address, 0x01);

// 设置黑名单
await contract.setAddressFlags(address, 0x04);

// 批量设置
await contract.batchSetAddressFlags([addr1, addr2], [0x01, 0x04]);
```

#### 权限标志说明
- `0x01` - 免税地址
- `0x02` - 免限制地址
- `0x04` - 黑名单地址
- `0x08` - AMM交易对
- `0x10` - 机器人标记(扩展)

### 参数调整

#### 更新税率
```javascript
await contract.updateTaxRates(
    200,  // 2% 买入税
    400,  // 4% 卖出税
    50    // 0.5% 转账税
);
```

#### 更新交易配置
```javascript
await contract.updateTradingConfig({
    maxTxAmount: totalSupply * 3 / 100,     // 3% 最大交易
    maxWalletAmount: totalSupply * 8 / 100, // 8% 最大持有
    swapThreshold: totalSupply / 2000,      // 0.05% swap阈值
    cooldownSeconds: 60,                    // 60秒冷却
    limitsEnabled: true
});
```

#### 更新费用分配
```javascript
await contract.updateFeeDistribution({
    marketing: 50,   // 50% 营销
    liquidity: 25,   // 25% 流动性
    burn: 15,        // 15% 销毁
    reflection: 10   // 10% 奖励
});
```

## 维护和监控

### 状态查询

#### 获取合约信息
```javascript
const info = await contract.getContractInfo();
console.log("交易启用:", info._tradingEnabled);
console.log("税率配置:", info._taxRates);
console.log("合约余额:", info._contractTokenBalance);
```

#### 查询地址权限
```javascript
const isExcluded = await contract.isExcludedFromFees(address);
const isBlacklisted = await contract.isBlacklisted(address);
const cooldown = await contract.getRemainingCooldown(address);
```

### 紧急功能

#### 暂停合约
```javascript
await contract.pause();    // 暂停
await contract.unpause();  // 恢复
```

#### 手动执行Swap
```javascript
await contract.manualSwap();
```

#### 资产救援
```javascript
// 救援ETH
await contract.rescueETH();

// 救援其他代币
await contract.rescueTokens(tokenAddress, amount);
```

## 监控和分析

### 重要事件监听

```javascript
// 监听交易启用
contract.on("TradingEnabled", (timestamp) => {
    console.log("交易已启用:", new Date(timestamp * 1000));
});

// 监听Swap事件
contract.on("TokensSwappedForETH", (tokensSwapped, ethReceived) => {
    console.log(`Swap: ${tokensSwapped} tokens -> ${ethReceived} ETH`);
});

// 监听税率更新
contract.on("TaxRatesUpdated", (buyTax, sellTax, transferTax) => {
    console.log(`新税率: ${buyTax/100}%/${sellTax/100}%/${transferTax/100}%`);
});
```

### 关键指标监控

```javascript
// 定期检查合约状态
async function monitorContract() {
    const info = await contract.getContractInfo();
    
    console.log("=== 合约状态 ===");
    console.log("交易启用:", info._tradingEnabled);
    console.log("合约代币余额:", ethers.formatEther(info._contractTokenBalance));
    console.log("合约ETH余额:", ethers.formatEther(info._contractETHBalance));
    
    // 检查是否需要手动干预
    if (info._contractTokenBalance > tradingConfig.swapThreshold * 10) {
        console.log("⚠️  合约代币余额过高，考虑手动Swap");
    }
}

setInterval(monitorContract, 60000); // 每分钟检查一次
```

##  安全注意事项

### 权限安全

1. **所有者权限**: 妥善保管所有者私钥
2. **多签钱包**: 考虑使用多签钱包作为所有者
3. **权限分散**: 避免单点故障

### 参数设置

1. **税率限制**: 不要设置过高的税率(建议<10%)
2. **交易限制**: 适度的限制保护项目
3. **冷却时间**: 不要设置过长的冷却时间(<5分钟)

### 合约验证

1. **源码验证**: 在区块链浏览器验证合约源码
2. **审计报告**: 考虑进行专业安全审计
3. **测试覆盖**: 充分的测试网测试

## 故障排除

### 常见问题

#### 1. 合约部署失败
- 检查构造函数参数格式
- 确认账户有足够ETH支付gas
- 验证网络配置正确

#### 2. 交易失败
```
"Trading not enabled" - 需要调用enableTrading()
"Blacklisted address" - 地址在黑名单中
"Cooldown period active" - 冷却时间未结束
"Exceeds max tx amount" - 超过单笔交易限制
```

#### 3. Swap不执行
- 检查合约代币余额是否达到阈值
- 确认交易已启用
- 验证Uniswap流动性是否充足

### 调试工具

#### 使用Remix调试
1. 在Remix中连接已部署合约
2. 使用"At Address"功能
3. 调用查询函数检查状态

#### 使用ethers.js脚本
```javascript
// 连接到合约
const contract = new ethers.Contract(
    contractAddress,
    contractABI,
    provider
);

// 查询状态
async function debug() {
    console.log("Owner:", await contract.owner());
    console.log("Trading enabled:", await contract.tradingEnabled());
    console.log("Contract balance:", await contract.balanceOf(contractAddress));
}
```

## 性能优化

### Gas优化建议

1. **批量操作**: 使用`batchSetAddressFlags`批量设置权限
2. **合理阈值**: 设置适当的swap阈值避免频繁swap
3. **事件精简**: 只记录必要的事件参数

### 维护最佳实践

1. **定期监控**: 设置自动化监控脚本
2. **参数调优**: 根据市场情况调整参数
3. **社区沟通**: 重要变更前与社区沟通

##  许可证

MIT License - 详见LICENSE文件
---

**免责声明**: 这仅仅是一个教学性质的测试代币。