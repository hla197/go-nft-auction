import { ethers, Contract, Log } from 'ethers';
import { ERC721_ABI } from '../abi/ERC721';
import { nftUtil, NFTItem } from './common';
import { NFT_START_BLOCK } from '../config/config';

/* ================== 常量 ================== */

const TRANSFER_TOPIC = ethers.id('Transfer(address,address,uint256)');

const common = nftUtil();

/* ================== 类型 ================== */

interface TokenMetadata {
  tokenURI: string;
  metadata: any;
}

/**
 * 刷新用户 ERC721 NFT（增量扫描 + 缓存）
 */
export async function refreshMyERC721Tokens(provider: ethers.Provider, nftAddress: string, owner: string, chainId: number, onNFTUpdate: (nft: NFTItem) => void) {
  const ownerTopic = ethers.zeroPadValue(ethers.getAddress(owner), 32);

  const latestBlock = await provider.getBlockNumber();
  const fromBlock = NFT_START_BLOCK;

  /* ---------- 扫描 Transfer 事件 ---------- */

  const owned = await scanOwnedERC721({
    provider,
    nftAddress,
    owner,
    fromBlock,
    toBlock: latestBlock
  });
  // // 缓存拥有的 NFT ID 列表
  const ownedStrList = Array.from(owned).map((id) => id.toString());
  localStorage.setItem(common.ownerNftList(chainId, nftAddress, owner), JSON.stringify(ownedStrList));
  console.log('refreshMyERC721Tokens owned', ownedStrList);

  /* ---------- 加载 NFT 元数据 ---------- */

  const contract = new Contract(nftAddress, ERC721_ABI, provider);

  for (const tokenId of owned) {
    (async () => {
      const tokenMetadata = await getTokenURI(chainId, nftAddress, contract, tokenId);
      onNFTUpdate({
        tokenId,
        tokenURI: tokenMetadata.tokenURI,
        rawImage: tokenMetadata.metadata.image,
        image: common.ipfsToHttp(tokenMetadata.metadata.image),
        name: tokenMetadata.metadata.name,
        contract: nftAddress
      });
    })();
  }
}

/**
 * 扫描某地址当前持有的 ERC721 TokenIds
 */
export async function scanOwnedERC721({
  provider,
  nftAddress,
  owner,
  fromBlock = 0,
  toBlock
}: {
  provider: ethers.Provider;
  nftAddress: string;
  owner: string;
  fromBlock?: number;
  toBlock?: number;
}): Promise<number[]> {
  const latestBlock = toBlock ?? (await provider.getBlockNumber());

  const TRANSFER_TOPIC = ethers.id('Transfer(address,address,uint256)');

  /* ---------- 1️⃣ 拉取所有 Transfer ---------- */
  const logs: Log[] = await provider.getLogs({
    address: nftAddress,
    fromBlock,
    toBlock: latestBlock,
    topics: [TRANSFER_TOPIC]
  });

  const iface = new ethers.Interface(ERC721_ABI);

  /* ---------- 2️⃣ 解析日志 ---------- */
  const transfers = logs.map((log) => {
    const parsed = iface.parseLog(log);
    if (!parsed) return null; // ⭐ 关键

    return {
      tokenId: parsed.args.tokenId,
      from: parsed.args.from.toLowerCase(),
      to: parsed.args.to.toLowerCase(),
      blockNumber: log.blockNumber,
      logIndex: log.index
    };
  });
  /* ---------- 3️⃣ 按时间排序（关键） ---------- */
  transfers.sort((a, b) => (a.blockNumber !== b.blockNumber ? a.blockNumber - b.blockNumber : a.logIndex - b.logIndex));

  /* ---------- 4️⃣ 回放 ownership ---------- */
  const ownerMap = new Map<number, string>();

  for (const t of transfers) {
    ownerMap.set(t.tokenId, t.to);
  }

  /* ---------- 5️⃣ 过滤出 owner 当前持有 ---------- */
  return [...ownerMap.entries()].filter(([_, currentOwner]) => currentOwner === owner.toLowerCase()).map(([tokenId]) => tokenId);
}

export async function getTokenURI(chainId: number, nftAddress: string, contract: Contract, tokenId: number): Promise<TokenMetadata> {
  const key = common.cacheKey(chainId, nftAddress, tokenId);
  const cached = common.getCache(key);
  let tokenMetadata = new Promise<TokenMetadata>((resolve) => resolve({ tokenURI: '', metadata: null }));

  if (cached && cached.metadata) {
    tokenMetadata = new Promise<TokenMetadata>((resolve) => resolve({ tokenURI: cached.tokenURI, metadata: cached.metadata }));
    return tokenMetadata;
  }

  try {
    const tokenURI = await contract.tokenURI(tokenId);
    const metadataUrl = common.ipfsToHttp(tokenURI);
    const metadata = await fetch(metadataUrl).then((r) => r.json());

    common.setCache(key, {
      tokenId: tokenId.toString(),
      tokenURI: tokenURI,
      metadata: metadata,
      updatedAt: Date.now(),
      contract: nftAddress
    });

    tokenMetadata = new Promise<TokenMetadata>((resolve) => resolve({ tokenURI: tokenURI, metadata: metadata }));
  } catch (err) {
    // 防止反复重试
    common.setCache(key, {
      tokenId: tokenId.toString(),
      tokenURI: '',
      metadata: null,
      updatedAt: Date.now(),
      contract: nftAddress
    });
  }

  return tokenMetadata;
}

/* ================== 直接从缓存读取（页面初始化用） ================== */

export function loadCachedNFTs(chainId: number, nftAddress: string, account: string): NFTItem[] {
  let owned = localStorage.getItem(common.ownerNftList(chainId, nftAddress, account));
  owned = owned ? JSON.parse(owned) : [];
  console.log('loadCachedNFTs owned', owned);
  if (!owned) return [];

  const result: NFTItem[] = [];
  for (const tokenId of owned) {
    const key = common.cacheKey(chainId, nftAddress, Number(tokenId));
    const cached = common.getCache(key);
    if (!cached || !cached.metadata) {
      // 骨架空数据
      result.push({
        tokenId: Number(tokenId),
        tokenURI: '',
        rawImage: '',
        image: '',
        contract: nftAddress
      });
    } else {
      result.push({
        tokenId: Number(tokenId),
        tokenURI: cached.tokenURI,
        rawImage: cached.metadata.image,
        image: common.ipfsToHttp(cached.metadata.image),
        contract: nftAddress,
        name: cached.metadata.name
      });
    }
  }

  return result;
}
