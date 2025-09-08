// We require the Hardhat Runtime Environment explicitly here. This is optional
// but useful for running the script in a standalone fashion through `node <script>`.
//
// You can also run a script with `npx hardhat run <script>`. If you do that, Hardhat
// will compile your contracts, add the Hardhat Runtime Environment's members to the
// global scope, and execute the script.
const hre = require("hardhat");

async function main() {
  // 计算解锁时间（当前时间 + 60秒）
  const currentTimestampInSeconds = Math.round(Date.now() / 1000);
  const unlockTime = currentTimestampInSeconds + 60;
  
  // 锁定金额
  const lockedAmount = hre.ethers.parseEther("0.001");
  
  console.log("Deploying Lock contract...");
  console.log("Unlock time:", unlockTime);
  console.log("Locked amount:", hre.ethers.formatEther(lockedAmount), "ETH");
  
  // 获取合约工厂并部署
  const Lock = await hre.ethers.getContractFactory("Lock");
  const lock = await Lock.deploy(unlockTime, { value: lockedAmount });
  
  // 等待部署完成
  await lock.waitForDeployment();
  
  // 获取合约地址
  const contractAddress = await lock.getAddress();
  
  console.log(
    `Lock with ${hre.ethers.formatEther(
      lockedAmount
    )} ETH and unlock timestamp ${unlockTime} deployed to ${contractAddress}`
  );
  
  console.log("Transaction hash:", lock.deploymentTransaction().hash);
}

main().catch((error) => {
  console.error(error);
  process.exitCode = 1;
});