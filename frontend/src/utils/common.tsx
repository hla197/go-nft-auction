const CACHE_PREFIX = 'nft-cache';
const LAST_BLOCK_PREFIX = 'nft-last-block';
const CACHE_TTL = 1000 * 60 * 60 * 24 * 7; // 7 天

export type NFTItem = {
  tokenId: number;
  tokenURI: string;
  image: string;
  rawImage: string;
  name?: string;
  contract: string;
};

export type CachedNFT = {
  tokenId: string;
  tokenURI: string;
  metadata: any | null;
  updatedAt: number;
  contract: string;
};

export type AuctionStruct = {
  seller: string;
  startPrice: bigint;
  startPriceUsd: bigint;
  highestBid: bigint;
  highestBidder: string;
  highestBidUsd: bigint;
  highestBidderToken: string;
  endTime: bigint;
  startTime: bigint;
  nftAddress: string;
  tokenId: bigint;
  token: string;
  active: boolean;
  auctionId: number;

  // 前端附加
  nftMeta?: NFTItem | null;
};

export function nftUtil() {
  /* ================== key ================== */

  function cacheKey(chainId: number, nftAddress: string, tokenId: number) {
    return `${CACHE_PREFIX}:${chainId}:${nftAddress.toLowerCase()}:${tokenId.toString()}`;
  }

  function lastBlockKey(chainId: number, nftAddress: string, owner: string) {
    return `${LAST_BLOCK_PREFIX}:${chainId}:${nftAddress.toLowerCase()}:${owner.toLowerCase()}`;
  }

  function ownerNftList(chainId: number, nftAddress: string, owner: string) {
    return `${chainId}:${nftAddress.toLowerCase()}:${owner.toLowerCase()}`;
  }

  /* ================== localStorage ================== */

  function getCache(key: string): CachedNFT | null {
    const raw = localStorage.getItem(key);
    if (!raw) return null;

    const cached = JSON.parse(raw) as CachedNFT;
    if (Date.now() - cached.updatedAt > CACHE_TTL) {
      localStorage.removeItem(key);
      return null;
    }
    return cached;
  }

  function setCache(key: string, data: CachedNFT) {
    localStorage.setItem(key, JSON.stringify(data));
  }

  function getCachedOwnerNFTs(chainId: number, nftAddress: string, owner: string): string[] {
    const raw = localStorage.getItem(ownerNftList(chainId, nftAddress, owner));
    return raw ? JSON.parse(raw) : [];
  }

  function setCachedOwnerNFTs(chainId: number, nftAddress: string, owner: string, tokenIds: string[]) {
    localStorage.setItem(ownerNftList(chainId, nftAddress, owner), JSON.stringify(tokenIds));
  }

  function getLastBlock(chainId: number, nftAddress: string, owner: string): number {
    const raw = localStorage.getItem(lastBlockKey(chainId, nftAddress, owner));
    return raw ? Number(raw) : 0;
  }

  function setLastBlock(chainId: number, nftAddress: string, owner: string, block: number) {
    localStorage.setItem(lastBlockKey(chainId, nftAddress, owner), block.toString());
  }

  function ipfsToHttp(uri: string): string {
    if (!uri) return '';
    return uri.startsWith('ipfs://') ? uri.replace('ipfs://', 'https://ipfs.io/ipfs/') : uri;
  }

  /* ================== NFT ================== */

  async function getNFTMetadata(chainId: number, contract: any, nftAddress: string, tokenId: number): Promise<NFTItem | null> {
    const key = cacheKey(chainId, nftAddress, tokenId);
    const cached = getCache(key);

    if (cached?.metadata) {
      return {
        tokenId,
        tokenURI: cached.tokenURI,
        image: ipfsToHttp(cached.metadata.image),
        rawImage: cached.metadata.image,
        name: cached.metadata.name,
        contract: nftAddress
      };
    }

    const tokenURI = await contract.tokenURI(tokenId);
    const metadata = await fetch(ipfsToHttp(tokenURI)).then((r) => r.json());

    setCache(key, {
      tokenId: tokenId.toString(),
      tokenURI,
      metadata,
      updatedAt: Date.now(),
      contract: nftAddress
    });

    return {
      tokenId,
      tokenURI,
      image: ipfsToHttp(metadata.image),
      rawImage: metadata.image,
      name: metadata.name,
      contract: nftAddress
    };
  }

  // 通用交易处理器，捕获 revert message
  async function handleTx(txFn: () => Promise<any>, successMsg?: string, fallbackMsg: string = '交易失败') {
    try {
      // callStatic 先模拟，捕获 revert 原因
      if (txFn.callStatic) {
        await txFn.callStatic();
      }

      // 发送真实交易
      const tx = await txFn();
      await tx.wait();

      if (successMsg) alert(successMsg);
    } catch (err: any) {
      console.error(err);

      // ethers v6 missing revert data fallback
      let message =
        err?.reason || // revert reason
        err?.data?.message || // rpc message
        err?.message || // ethers error
        fallbackMsg; // 默认提示

      alert(message);
    }
  }

  /* ================== ✅ EXPORT ================== */

  return {
    ownerNftList,
    cacheKey,
    lastBlockKey,
    getCache,
    setCache,
    getNFTMetadata,
    getCachedOwnerNFTs,
    setCachedOwnerNFTs,
    getLastBlock,
    setLastBlock,
    ipfsToHttp,
    handleTx,
    clearAll(chainId?: number) {
      Object.keys(localStorage).forEach((key) => {
        if (!chainId || key.includes(`:${chainId}:`)) {
          localStorage.removeItem(key);
        }
      });
    }
  };
}
