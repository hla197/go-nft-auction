import React, { useState, useEffect } from 'react';
import { Contract, ethers, formatEther, parseEther, ZeroAddress } from 'ethers';
import { AUCTION_ABI } from '../abi/NFTAuction';
import { ERC20_ABI } from '../abi/ERC20';
import { ERC721_ABI } from '../abi/ERC721';
import { AUCTION_ADDRESS, ERC20_ADDRESS, NFT_ADDRESS } from '../config/config';
import { getTokenURI } from '../utils/erc721';
import { AuctionStruct, nftUtil } from '../utils/common';
import BidModal from '../components/BidModal'; // 出价弹窗
import HistoryModal from '../components/HistoryModal'; // 历史拍卖弹窗

const common = nftUtil();

export default function AuctionList({ signer, account, provider, chainId }: any) {
  const [auctionContract, setAuctionContract] = useState<Contract | null>(null);
  const [auctions, setAuctions] = useState<AuctionStruct[]>([]);

  /* ---------- 出价弹窗 ---------- */
  const [bidOpen, setBidOpen] = useState(false);
  const [bidAuctionId, setBidAuctionId] = useState<number | null>(null);
  const [bidAmount, setBidAmount] = useState('');
  const [payToken, setPayToken] = useState<'ETH' | 'ERC20'>('ERC20');
  const [ethPrice, setEthPrice] = useState<BigInt>(0n);
  const [erc20Price, setErc20Price] = useState<BigInt>(0n);

  /* ---------- 历史拍卖弹窗 ---------- */
  const [historyModalOpen, setHistoryModalOpen] = useState(false);
  const [selectedNftContract, setSelectedNftContract] = useState('');
  const [selectedNftTokenId, setSelectedNftTokenId] = useState('');
  const [selectedChainId, setSelectedChainId] = useState('');

  /* ---------- 初始化合约 ---------- */
  useEffect(() => {
    if (!signer) return;
    setAuctionContract(new Contract(AUCTION_ADDRESS, AUCTION_ABI, signer));
  }, [signer]);

  /* ---------- 加载拍卖 ---------- */
  async function loadAuctions() {
    if (!auctionContract) return;

    const count: bigint = await auctionContract.auctionIdCounter();
    const nftContract = new Contract(NFT_ADDRESS, ERC721_ABI, provider);
    const list: AuctionStruct[] = [];

    const ethPrice = await auctionContract.getChainlinkDataFeedLatestAnswer(ZeroAddress);
    const erc20Price = await auctionContract.getChainlinkDataFeedLatestAnswer(ERC20_ADDRESS);
    setErc20Price(erc20Price);
    setEthPrice(ethPrice);

    for (let id = 1; id <= Number(count); id++) {
      const a = await auctionContract.auctions(id);
      if (!a.active) continue;
      const tokenMetadata = await getTokenURI(chainId, NFT_ADDRESS, nftContract, a.tokenId);
      if (tokenMetadata?.metadata?.image) {
        tokenMetadata.metadata.image = common.ipfsToHttp(tokenMetadata.metadata.image);
      }
      const auction: AuctionStruct = {
        auctionId: id,
        seller: a.seller,
        nftAddress: a.nftAddress,
        tokenId: a.tokenId,
        startPrice: a.startPrice,
        highestBid: a.highestBid,
        startPriceUsd: a.startPriceUsd,
        highestBidUsd: a.highestBidUsd,
        highestBidder: a.highestBidder,
        highestBidderToken: a.highestBidderToken,
        endTime: a.endTime,
        startTime: a.startTime,
        token: a.token,
        active: a.active,
        nftMeta: tokenMetadata?.metadata
      };

      list.push(auction);
    }

    setAuctions(list);
  }

  useEffect(() => {
    if (auctionContract) loadAuctions();
  }, [auctionContract]);

  /* ---------- 工具函数 ---------- */
  function formatBid(a: AuctionStruct) {
    if (a.highestBidder === ethers.ZeroAddress) return '暂无';

    if (a.highestBidderToken === ethers.ZeroAddress) {
      return `${formatEther(a.highestBid)} ETH`;
    }
    return `${ethers.formatUnits(a.highestBid, 6)} USDC`;
  }

  function isBidTooLow(a: AuctionStruct, amount: number, token: 'ETH' | 'ERC20') {
    // 转换为 USD
    const bidUsd = getUsdPrice(token, amount.toString());
    const highestUsd =
      a.highestBidderToken === ZeroAddress
        ? Number(a.highestBidUsd) / 1e8 // ETH 出价对应 USD
        : Number(a.highestBidUsd) / 1e8; // ERC20 出价对应 USD

    return bidUsd <= highestUsd;
  }

  function getUsdPrice(payToken: 'ETH' | 'ERC20', value: string) {
    if (!value) return 0;
    const usdprice = payToken === 'ETH' ? Number(ethPrice) * parseFloat(value) : Number(erc20Price) * parseFloat(value);
    return usdprice;
  }

  /* ---------- 提交出价 ---------- */
  async function submitBid() {
    if (!auctionContract || bidAuctionId === null) return;
    if (!bidAmount || Number(bidAmount) <= 0) return alert('请输入有效出价');

    const auction = auctions.find((a) => a.auctionId === bidAuctionId);
    if (!auction) return;

    if (isBidTooLow(auction, Number(bidAmount), payToken)) {
      return alert('出价必须高于当前最高 USD 价');
    }

    try {
      if (payToken === 'ETH') {
        const value = parseEther(bidAmount);
        const tx = await auctionContract.placeBid(bidAuctionId, ZeroAddress, value, { value });
        await tx.wait();
      } else {
        const erc20 = new Contract(ERC20_ADDRESS, ERC20_ABI, signer);
        const decimals = await erc20.decimals();
        const amount = ethers.parseUnits(bidAmount, decimals);

        const allowance = await erc20.allowance(account, AUCTION_ADDRESS);
        if (allowance < amount) {
          const approveTx = await erc20.approve(AUCTION_ADDRESS, amount);
          await approveTx.wait();
        }

        const tx = await auctionContract.placeBid(bidAuctionId, ERC20_ADDRESS, amount);
        await tx.wait();
      }

      setBidOpen(false);
      setBidAmount('');
      loadAuctions();
      const event = new CustomEvent('balance-refresh', {});
      window.dispatchEvent(event);
    } catch (err: any) {
      console.error(err);
      alert(err?.reason || err?.data?.message || err?.message || '出价失败');
    }
  }

  /* ---------- 结束拍卖 ---------- */
  async function endAuction(id: number) {
    if (!auctionContract) return;
    await common.handleTx(() => auctionContract.endAuction(id), '拍卖已成功结束', '拍卖未结束或交易失败');
    loadAuctions();
  }

  return (
    <div>
      <h2>拍卖列表</h2>

      {auctions.map((a) => (
        <div
          key={a.auctionId}
          style={{
            border: '1px solid #e5e7eb',
            borderRadius: 12,
            padding: 16,
            marginBottom: 16,
            display: 'flex',
            gap: 16
          }}
        >
          {/* NFT 图片 */}
          <div
            style={{
              flex: 1,
              borderRadius: 10,
              overflow: 'hidden',
              background: '#f3f4f6'
            }}
          >
            <img
              src={a.nftMeta?.image}
              style={{
                width: '100%',
                height: '100%',
                objectFit: 'cover'
              }}
            />
          </div>
          {/* 信息 */}
          <div style={{ flex: 2 }}>
            <h3>{a.nftMeta?.name ?? 'NFT'}</h3>

            <p>
              <b>Auction ID：</b>
              {a.auctionId}
            </p>
            <p>
              <b>NFT 合约：</b>
              {a.nftAddress}
            </p>
            <p>
              <b>Token ID：</b>
              {a.tokenId !== undefined ? a.tokenId.toString() : '-'}
            </p>

            <p>
              <b>起拍价：</b>
              {formatEther(a.startPrice)} ETH
              {a.startPriceUsd && `（≈ ${ethers.formatUnits(a.startPriceUsd, 8)} USD）`}
            </p>

            <p>
              <b>最高出价：</b>
              {formatBid(a)}
              {a.highestBidUsd && `（≈ ${ethers.formatUnits(a.highestBidUsd, 8)} USD）`}
            </p>

            <p>
              <b>最高出价者：</b>
              {a.highestBidder === ethers.ZeroAddress ? '暂无' : a.highestBidder}
            </p>

            <p>
              <b>出价币种：</b>
              {a.highestBidderToken === ethers.ZeroAddress ? 'ETH' : 'USDC'}
            </p>

            <p>
              <b>开始时间：</b>
              {new Date(Number(a.startTime) * 1000).toLocaleString()}
            </p>

            <p>
              <b>结束时间：</b>
              {new Date(Number(a.endTime) * 1000).toLocaleString()}
            </p>

            <p>
              <b>状态：</b>
              {a.active ? '进行中' : '已结束'}
            </p>

            {/* 操作 */}
            {a.active && (
              <div style={{ marginTop: 12 }}>
                <button
                  onClick={() => {
                    setBidAuctionId(a.auctionId);
                    setPayToken('ERC20');
                    setBidAmount('');
                    setBidOpen(true);
                  }}
                >
                  出价
                </button>

                <button style={{ marginLeft: 10 }} onClick={() => endAuction(a.auctionId)}>
                  结束拍卖
                </button>

                {/* 获取历史拍卖价格按钮 */}
                <button
                  style={{ marginLeft: 10 }}
                  onClick={() => {
                    setSelectedNftContract(a.nftAddress);
                    setSelectedNftTokenId(a.tokenId.toString());
                    setSelectedChainId(chainId.toString());
                    setHistoryModalOpen(true);
                  }}
                >
                  获取历史拍卖价格
                </button>
              </div>
            )}
          </div>
        </div>
      ))}

      {/* ---------- 出价弹窗 ---------- */}
      <BidModal
        isOpen={bidOpen}
        onClose={() => setBidOpen(false)}
        auctionId={bidAuctionId || 0}
        signer={signer}
        account={account}
        payToken={payToken}
        ethPrice={ethPrice}
        erc20Price={erc20Price}
        onBidPlaced={loadAuctions}
      />

      {/* ---------- 历史拍卖弹窗 ---------- */}
      <HistoryModal isOpen={historyModalOpen} onClose={() => setHistoryModalOpen(false)} nftContract={selectedNftContract} nftTokenId={selectedNftTokenId} chainId={Number(selectedChainId)} />
    </div>
  );
}
