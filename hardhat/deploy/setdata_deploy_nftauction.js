const { ethers, deployments, getNamedAccounts } = require('hardhat');
const fs = require('fs');
const path = require('path');

module.exports = async function () {
  const { deployer } = await getNamedAccounts();

  console.log('Call method on existing NftAuction proxy, deployer:', deployer);

  // 1. 读取 proxy 地址
  const storePath = './deploy/cache/proxyNftAuction.json';
  const deployedData = JSON.parse(fs.readFileSync(storePath, 'utf8'));
  const proxyAddress = deployedData.proxyAddress;

  console.log('Proxy address:', proxyAddress);

  // 2. 获取合约实例（⚠️ 用 V2 ABI）
  const nftAuction = await ethers.getContractAt('NftAuctionV2', proxyAddress);

  // 测试网 USDC和ETH的预言机
  // 3. 调用方法
  //   const tx1 = await nftAuction.setPriceFeed(ethers.ZeroAddress, '0x694AA1769357215DE4FAC081bf1f309aDC325306');
  //   await tx1.wait();

    const tx2 = await nftAuction.setPriceFeed("0x9F9557a99E38C1A0b98e90f1dD0141C307fe3B13", '0xA2F78ab2355fe2f984D808B5CeE7FD0A93D5270E');
    await tx2.wait();

  //   const price = await nftAuction.getChainlinkDataFeedLatestAnswer(ethers.ZeroAddress);
  //   console.log('price:', price.toString());

// {
//         address seller; // 拍卖人
//         uint256 startPrice; // 原始数量（ETH wei 或 ERC20 token）
//         uint256 startPriceUsd; // USD 统一价格，8 decimals
//         uint256 highestBid; // 最高投标价
//         address highestBidder; // 最高投标人
//         uint256 highestBidUsd; // USD 8 decimals
//         address highestBidderToken; // 最高投标人使用的代币
//         uint256 endTime; // 结束时间
//         uint256 startTime; // 开始时间
//         address nftAddress; // NFT合约地址
//         uint256 tokenId; // NFT的tokenI
//         address token; // 代币地址
//         bool active; // 状态
//     }

    // const auction = await nftAuction.auctions(1);
    // console.log("Auction Details:", auction);

  console.log('end');
};

module.exports.tags = ['CallNftAuction'];
