### 项目地址 https://github.com/ethereum/go-ethereum/blob/master/cmd/geth/main.go
## Geth 主要命令分类

### 1. **核心命令**

- **`geth`** - 启动以太坊节点的主命令
- **`account`** - 管理账户相关操作
- **`attach`** - 连接到正在运行的节点，启动交互式 JavaScript 环境
- **`console`** - 启动带有交互式 JavaScript 环境的节点
- **`init`** - 初始化新的创世区块

### 2. **数据库和区块链操作**

- **`db`** - 低级数据库操作
- **`dump`** - 导出特定区块的存储内容
- **`dumpconfig`** - 以 TOML 格式导出配置值
- **`dumpgenesis`** - 将创世区块的 JSON 配置导出到标准输出
- **`export`** - 将区块链导出到文件
- **`export-history`** - 将区块链历史导出为 Era 归档
- **`import`** - 导入区块链文件
- **`import-history`** - 导入 Era 归档
- **`import-preimages`** - 从 RLP 流导入预映像数据库

### 3. **开发和调试工具**

- **`js`** - 执行指定的 JavaScript 文件（已弃用）
- **`license`** - 显示许可证信息

## 重要参数分类

### 1. **网络配置参数**

- **`--networkid`** - 明确设置网络 ID（整数）
- **`--sepolia`** - 连接到 Sepolia 测试网
- **`--holesky`** - 连接到 Holesky 测试网
- **`--syncmode`** - 同步模式（snap/full/light）
    - `snap`（默认）- 快速同步，下载更多数据但避免处理整个历史
    - `full` - 完全同步
    - `light` - 轻节点模式

### 2. **缓存和性能参数**

- **`--cache`** - 分配给内部缓存的内存（MB）
    - 默认：主网全节点 4096 MB，轻模式 128 MB
- **`--cache.blocklogs`** - 用于过滤的日志缓存大小（默认：32 个区块）
- **`--cache.database`** - 用于数据库 IO 的缓存内存百分比（默认：50%）
- **`--cache.gc`** - 用于 trie 修剪的缓存内存百分比（默认：25%）
- **`--cache.noprefetch`** - 禁用启发式状态预取
- **`--cache.preimages`** - 启用记录 trie 键的 SHA3/keccak 预映像
- **`--cache.snapshot`** - 用于快照缓存的内存百分比（默认：10%）
- **`--cache.trie`** - 用于 trie 缓存的内存百分比（默认：15%）

### 3. **Gas 价格预言机参数**

- **`--gpo.blocks`** - 检查 gas 价格的最近区块数（默认：20）
- **`--gpo.ignoreprice`** - GPO 忽略交易的最低 gas 价格（默认：2）
- **`--gpo.maxprice`** - GPO 推荐的最大交易优先费（默认：500 gwei）
- **`--gpo.percentile`** - GPO 建议的 gas 价格百分位（默认：60）

### 4. **RPC 和 API 参数**

- **`--http`** - 启用 HTTP-RPC 服务器
- **`--http.addr`** - HTTP-RPC 服务器监听地址（默认：localhost）
- **`--http.port`** - HTTP-RPC 服务器端口（默认：8545）
- **`--ws`** - 启用 WebSocket-RPC 服务器
- **`--ws.port`** - WebSocket-RPC 服务器端口（默认：8546）

### 5. **数据目录和存储**

- **`--datadir`** - 数据库和密钥库的数据目录
- **`--keystore`** - 密钥库目录

### 6. **挖矿相关参数**

- **`--mine`** - 启用挖矿
- **`--miner.threads`** - 用于挖矿的 CPU 线程数
- **`--miner.etherbase`** - 挖矿奖励接收地址

## 常用命令组合示例

1. **启动主网节点**：
   ```bash
   geth --syncmode snap
   ```

2. **启动测试网节点**：
   ```bash
   geth --sepolia --syncmode snap
   ```

3. **启动带控制台的节点**：
   ```bash
   geth --sepolia console
   ```

4. **连接到运行中的节点**：
   ```bash
   geth attach
   ```

5. **创建新账户**：
   ```bash
   geth account new
   ```

6. **启动轻节点**：
   ```bash
   geth --syncmode light
   ```

7. **启用 RPC 服务**：
   ```bash
   geth --http --http.addr 0.0.0.0
   ```

## JavaScript 控制台命令

在 Geth 控制台中，可以使用以下主要对象：

- **`eth`** - 以太坊相关操作
    - `eth.accounts` - 列出账户
    - `eth.getBalance()` - 获取余额
    - `eth.sendTransaction()` - 发送交易
    - `eth.syncing` - 同步状态

- **`admin`** - 节点管理
    - `admin.nodeInfo` - 节点信息
    - `admin.peers` - 连接的对等节点

- **`personal`** - 账户管理
    - `personal.newAccount()` - 创建新账户
    - `personal.unlockAccount()` - 解锁账户

- **`miner`** - 挖矿控制
    - `miner.start()` - 开始挖矿
    - `miner.stop()` - 停止挖矿

这些命令和参数构成了 Geth 的核心功能，允许用户运行以太坊节点、管理账户、同步区块链数据以及与以太坊网络进行交互。

## Geth 主要命令分类

### 1. **核心命令**

- **`geth`** - 启动以太坊节点的主命令
- **`account`** - 管理账户相关操作
- **`attach`** - 连接到正在运行的节点，启动交互式 JavaScript 环境
- **`console`** - 启动带有交互式 JavaScript 环境的节点
- **`init`** - 初始化新的创世区块

### 2. **数据库和区块链操作**

- **`db`** - 低级数据库操作
- **`dump`** - 导出特定区块的存储内容
- **`dumpconfig`** - 以 TOML 格式导出配置值
- **`dumpgenesis`** - 将创世区块的 JSON 配置导出到标准输出
- **`export`** - 将区块链导出到文件
- **`export-history`** - 将区块链历史导出为 Era 归档
- **`import`** - 导入区块链文件
- **`import-history`** - 导入 Era 归档
- **`import-preimages`** - 从 RLP 流导入预映像数据库

### 3. **开发和调试工具**

- **`js`** - 执行指定的 JavaScript 文件（已弃用）
- **`license`** - 显示许可证信息

## 重要参数分类

### 1. **网络配置参数**

- **`--networkid`** - 明确设置网络 ID（整数）
- **`--sepolia`** - 连接到 Sepolia 测试网
- **`--holesky`** - 连接到 Holesky 测试网
- **`--syncmode`** - 同步模式（snap/full/light）
    - `snap`（默认）- 快速同步，下载更多数据但避免处理整个历史
    - `full` - 完全同步
    - `light` - 轻节点模式

### 2. **缓存和性能参数**

- **`--cache`** - 分配给内部缓存的内存（MB）
    - 默认：主网全节点 4096 MB，轻模式 128 MB
- **`--cache.blocklogs`** - 用于过滤的日志缓存大小（默认：32 个区块）
- **`--cache.database`** - 用于数据库 IO 的缓存内存百分比（默认：50%）
- **`--cache.gc`** - 用于 trie 修剪的缓存内存百分比（默认：25%）
- **`--cache.noprefetch`** - 禁用启发式状态预取
- **`--cache.preimages`** - 启用记录 trie 键的 SHA3/keccak 预映像
- **`--cache.snapshot`** - 用于快照缓存的内存百分比（默认：10%）
- **`--cache.trie`** - 用于 trie 缓存的内存百分比（默认：15%）

### 3. **Gas 价格预言机参数**

- **`--gpo.blocks`** - 检查 gas 价格的最近区块数（默认：20）
- **`--gpo.ignoreprice`** - GPO 忽略交易的最低 gas 价格（默认：2）
- **`--gpo.maxprice`** - GPO 推荐的最大交易优先费（默认：500 gwei）
- **`--gpo.percentile`** - GPO 建议的 gas 价格百分位（默认：60）

### 4. **RPC 和 API 参数**

- **`--http`** - 启用 HTTP-RPC 服务器
- **`--http.addr`** - HTTP-RPC 服务器监听地址（默认：localhost）
- **`--http.port`** - HTTP-RPC 服务器端口（默认：8545）
- **`--ws`** - 启用 WebSocket-RPC 服务器
- **`--ws.port`** - WebSocket-RPC 服务器端口（默认：8546）

### 5. **数据目录和存储**

- **`--datadir`** - 数据库和密钥库的数据目录
- **`--keystore`** - 密钥库目录

### 6. **挖矿相关参数**

- **`--mine`** - 启用挖矿
- **`--miner.threads`** - 用于挖矿的 CPU 线程数
- **`--miner.etherbase`** - 挖矿奖励接收地址

## 常用命令组合示例

1. **启动主网节点**：
   ```bash
   geth --syncmode snap
   ```

2. **启动测试网节点**：
   ```bash
   geth --sepolia --syncmode snap
   ```

3. **启动带控制台的节点**：
   ```bash
   geth --sepolia console
   ```

4. **连接到运行中的节点**：
   ```bash
   geth attach
   ```

5. **创建新账户**：
   ```bash
   geth account new
   ```

6. **启动轻节点**：
   ```bash
   geth --syncmode light
   ```

7. **启用 RPC 服务**：
   ```bash
   geth --http --http.addr 0.0.0.0
   ```

## JavaScript 控制台命令

在 Geth 控制台中，可以使用以下主要对象：

- **`eth`** - 以太坊相关操作
    - `eth.accounts` - 列出账户
    - `eth.getBalance()` - 获取余额
    - `eth.sendTransaction()` - 发送交易
    - `eth.syncing` - 同步状态

- **`admin`** - 节点管理
    - `admin.nodeInfo` - 节点信息
    - `admin.peers` - 连接的对等节点

- **`personal`** - 账户管理
    - `personal.newAccount()` - 创建新账户
    - `personal.unlockAccount()` - 解锁账户

- **`miner`** - 挖矿控制
    - `miner.start()` - 开始挖矿
    - `miner.stop()` - 停止挖矿
