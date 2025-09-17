以下是Hardhat的常用命令清单：

## 基础命令

### 项目初始化和编译

# 初始化项目
npx hardhat

# 指定使用v2最新版初始化solidity项目
npx hardhat@^2.0.0

# 编译合约
npx hardhat compile

# 清理编译文件
npx hardhat clean

# 查看版本
npx hardhat --version

# 查看帮助
npx hardhat help


## 测试相关命令

# 运行所有测试
npx hardhat test

# 运行指定测试文件
npx hardhat test test/MyContract.test.js

# 运行测试并显示gas使用情况
npx hardhat test --gas-report

# 运行测试时显示详细日志
npx hardhat test --verbose


## 网络和节点命令


# 启动本地节点（默认端口8545）
npx hardhat node

# 指定端口启动节点
npx hardhat node --port 8546

# 启动节点时预设账户
npx hardhat node --accounts 20

## 部署和脚本执行


# 运行部署脚本（默认本地网络）
npx hardhat run scripts/deploy.js

# 部署到指定网络
npx hardhat run scripts/deploy.js --network goerli

# 运行其他脚本
npx hardhat run scripts/myScript.js --network localhost


## 控制台命令


# 启动Hardhat控制台（本地网络）
npx hardhat console

# 连接到指定网络的控制台
npx hardhat console --network goerli

# 在控制台中可以直接执行JavaScript代码与合约交互


## 验证合约


# 在区块链浏览器上验证合约
npx hardhat verify --network mainnet DEPLOYED_CONTRACT_ADDRESS "Constructor argument 1"

# 验证时指定合约文件
npx hardhat verify --contract contracts/MyContract.sol:MyContract --network goerli ADDRESS


## 账户和余额


# 查看账户列表
npx hardhat accounts

# 查看指定网络的账户
npx hardhat accounts --network localhost


## Flatten 合约


# 将合约及其依赖合并到一个文件（需要安装插件）
npx hardhat flatten contracts/MyContract.sol > flattened.sol


## 大小检查

# 检查合约大小（需要安装 hardhat-contract-sizer 插件）
npx hardhat size-contracts


## 覆盖率测试


# 运行代码覆盖率测试（需要安装 solidity-coverage 插件）
npx hardhat coverage


## 自定义任务

# 运行自定义任务
npx hardhat myCustomTask

# 带参数的自定义任务
npx hardhat myTask --param1 value1 --param2 value2


## 网络相关


# 列出所有配置的网络
npx hardhat --help

# Fork 主网进行本地测试
npx hardhat node --fork https://eth-mainnet.alchemyapi.io/v2/YOUR-API-KEY

# 重置fork的状态
npx hardhat node --fork URL --fork-block-number 12345678


## Gas 相关


# 查看gas使用情况（需要hardhat-gas-reporter插件）
npx hardhat test --gas-report

# 设置gas价格
npx hardhat run scripts/deploy.js --network goerli --gas-price 20000000000


## 常用插件命令

### Etherscan 验证

npm install --save-dev @nomiclabs/hardhat-etherscan
npx hardhat verify --network mainnet CONTRACT_ADDRESS "constructor args"


### Gas Reporter

npm install --save-dev hardhat-gas-reporter
# 配置后运行测试时自动显示gas报告


### Contract Sizer

npm install --save-dev hardhat-contract-sizer
npx hardhat size-contracts


## 环境变量使用


# 使用.env文件中的变量
npx hardhat run scripts/deploy.js --network mainnet