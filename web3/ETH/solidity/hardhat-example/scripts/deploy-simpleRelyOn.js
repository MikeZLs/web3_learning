const hre = require("hardhat");

// 案例3：要部署多个合约，合约部署有顺序上的依赖
async function main() {
  console.log("Starting multi-contract deployment...");
  
  // 部署合约A
  const myContractA_f = await hre.ethers.getContractFactory("MyContractA");
  const myContractA = await myContractA_f.deploy(); // 合约A 无构造参数
  await myContractA.waitForDeployment();
  const contractA_address = await myContractA.getAddress();
  console.log("ContractA deployed to:", contractA_address);
  
  // 部署合约B
  const myContractB_f = await hre.ethers.getContractFactory("MyContractB");
  const myContractB = await myContractB_f.deploy(); // 合约B 无构造参数
  await myContractB.waitForDeployment();
  const contractB_address = await myContractB.getAddress();
  console.log("ContractB deployed to:", contractB_address);
  
  // 部署合约C，使用A和B的地址作为构造参数
  const myContractC_f = await hre.ethers.getContractFactory("MyContractC");
  const myContractC = await myContractC_f.deploy(contractA_address, contractB_address);
  await myContractC.waitForDeployment();
  const contractC_address = await myContractC.getAddress();
  console.log("ContractC deployed to:", contractC_address);
  
  // 部署摘要
  console.log("\n=== Deployment Summary ===");
  console.log("ContractA:", contractA_address);
  console.log("ContractB:", contractB_address);
  console.log("ContractC:", contractC_address);
}

main().catch((error) => {
  console.error(error);
  process.exitCode = 1;
});


// 执行命令
// npx hardhat run .\scripts\deploy-simpleRelyOn.js
// npx hardhat run .\scripts\deploy-simpleRelyOn.js  --network eth_testnet