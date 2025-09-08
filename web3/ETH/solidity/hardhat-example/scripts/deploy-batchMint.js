const hre = require("hardhat");

// 案例1：合约没有构造函数 不需要传参时，直接部署
async function main() {
  console.log("Starting deployment...");
  
  // 指定需要部署的合约名称 （非合约文件名）
  const batchMint_factory = await hre.ethers.getContractFactory("BatchMintClips");
  
  // 部署合约
  console.log("Deploying BatchMintClips...");
  const BatchMintClips_contract = await batchMint_factory.deploy();
  
  // 等待合约部署完成
  await BatchMintClips_contract.waitForDeployment();
  
  // 获取合约地址
  const contractAddress = await BatchMintClips_contract.getAddress();
  
  // 打印合约地址和交易信息
  console.log("BatchMintClips deployed to:", contractAddress);
  console.log("Deployment transaction hash:", BatchMintClips_contract.deploymentTransaction().hash);
}

main().catch((error) => {
  console.error(error);
  process.exitCode = 1;
});

// 执行命令  // 0.01 ether = 10000000 GWEI
// npx hardhat run .\scripts\deploy-batchMint.js