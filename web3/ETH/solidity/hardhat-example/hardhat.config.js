require("@nomicfoundation/hardhat-toolbox");
// require("@nomiclabs/hardhat-waffle");
// require("./tasks/claimtoken");
// require("./tasks/eth_foxe_getmerkle");
// require("./tasks/generateMerkle")
// require("./tasks/taskdemo");
// require("./tasks/account_utils");
// 配置gas估算插件
require("hardhat-gas-reporter");
// require("dotenv").config();
/** @type import('hardhat/config').HardhatUserConfig */
const {projectId, mnemonic1,ALCHEMY_API_KEY} = process.env
module.exports = {
  networks: {
    // npx hardhat run scripts/deploy.js --network eth_testnet
    // 领水地址 : https://goerlifaucet.com/
    eth_testnet: {
      url: `https://sepolia.infura.io/v3/a8d3d794739b4366b3ed18393c5ee524`,  
      accounts: {
        mnemonic:" oisfhoak awrjaoilj adsaj ", // 私钥，切勿泄露
      }
    },
    ganache: {
      url: "http://127.0.0.1:7545",
      accounts: [
        // 从 Ganache 的 ACCOUNTS 页面复制私钥
        "0xd5686d988fe5506ceb63f546929898b5bf57370c6acbe0322f0f0eb640e1b033",
        // "0x2e26b48c71a98f118cce82fad43171d6f0aec8e55eb246c4a742fd5fae8f7779",
        // 可以添加多个账户
      ]
    },
    bsc: {
      url: `https://bsc-dataseed.binance.org/`,
      accounts: {
        mnemonic: "obscure satoshi lecture culture lady pattern fog shoe emerge step wonder sword"
      }
    },
    bsc_testnet: {
      url: `https://data-seed-prebsc-1-s1.binance.org:8545/`,
      accounts: {
        mnemonic: "obscure satoshi lecture culture lady pattern fog shoe emerge step wonder sword"
      }
    },
    hardhat: {
      forking: {
        url: `https://eth-mainnet.alchemyapi.io/v2/${ALCHEMY_API_KEY}`,
        blockNumber: 17342219
      }
    }
  },
  solidity: {
    // 可修改编译器版本
    compilers: [
      {
        version: "0.8.28",
      },
      {
        version: "0.8.18",
      },
      {
        version: "0.6.12",
      },
      {
        version: "0.8.7",
      },
      {
        version: "0.8.0",
      },
    ],
  },
  gasReporter: {
    // 开启gas估算插件 设置人民币和gas费换算，也可以设置为美元 usd
  	enabled: false,
    currency: 'USD',
    token: "ETH",
    coinmarketcap:process.env.COINMARKETCAP_API_KEY,
  }
};