const { ethers, upgrades } = require('hardhat');
const path = require('path');
const fs = require('fs');

module.exports = async function () {
  const storePath = path.resolve(__dirname, './cache/MyERC20.json');

  this.MyErc20 = await ethers.getContractFactory('MyErc20');
  deploy = await this.MyErc20.deploy('MyErc20', 'MEC20');
  await deploy.waitForDeployment();

  const deployAddress = await deploy.getAddress();

  console.log('deploy address:', deployAddress);

  console.log('Saved deployed addresses to', storePath);

  fs.writeFileSync(storePath, JSON.stringify({deployAddress}));
};

module.exports.tags = ['MyErc20'];
