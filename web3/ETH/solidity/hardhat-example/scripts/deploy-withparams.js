const hre = require("hardhat");

// 案例2：合约有构造参数，需要进行传参
async function main() {
  // 获取部署者账户
  const [deployer] = await hre.ethers.getSigners();
  
  // 指定需要部署的合约名称
  const simple_with_params_f = await hre.ethers.getContractFactory("SimpleWithParams");
  
  // 部署合约，传入构造函数参数
  const simple_with_params_f_contract = await simple_with_params_f.deploy(
    await deployer.getAddress(),
    '0x356faDD245d35ff8F1a207aC83BE7EEa911abeEE',
    100,
    97,
    {
      value: hre.ethers.parseEther("0.01") // 发送 0.01 ether
    }
  );
  
  // 等待部署完成
  await simple_with_params_f_contract.waitForDeployment();
  
  // 获取合约地址
  const contractAddress = await simple_with_params_f_contract.getAddress();
  
  console.log("Contract deployed to:", contractAddress);
  console.log("Deployer address:", await deployer.getAddress());
  console.log("Transaction hash:", simple_with_params_f_contract.deploymentTransaction().hash);
}

main().catch((error) => {
  console.error(error);
  process.exitCode = 1;
});
// 执行命令
// npx hardhat run .\scripts\deploy-withparams.js
// npx hardhat run .\scripts\deploy-withparams.js --network eth_testnet