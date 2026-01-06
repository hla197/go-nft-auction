import { useEffect, useState } from 'react';
import { BrowserProvider } from 'ethers';
import { loadCachedNFTs, refreshMyERC721Tokens } from '../utils/erc721';
import { NFT_ADDRESS, NFT_START_BLOCK } from '../config/config';
import CreateAuctionModal from '../components/CreateAuctionModal';
import { NFTItem } from '../utils/common';

export default function MyNFTs({ provider, account, chainId }: any) {
  const [tokens, setTokens] = useState<NFTItem[]>([]);
  const [loading, setLoading] = useState(false);
  const [selectedNFT, setSelectedNFT] = useState<NFTItem | null>(null);
  const [showAuctionModal, setShowAuctionModal] = useState(false);

  useEffect(() => {
    // 🚫 钱包没准备好，直接清空
    if (!provider || !account || !chainId) {
      setTokens([]);
      return;
    }
    console.log('MyNFTs', provider, account, chainId);
    let cancelled = false;

    /**
     * 加载用户的NFT收藏
     *
     * 该函数采用两阶段加载策略：
     * 1. 首先从本地缓存中加载已保存的NFT数据，提供即时响应
     * 2. 然后在后台刷新最新的NFT数据，确保信息准确性
     *
     * 使用cancelled标志位来防止组件卸载后的状态更新
     */
    async function load() {
      // 先显示缓存，保证 UI 快速响应
      const cached = loadCachedNFTs(chainId, NFT_ADDRESS, account);
      setTokens(cached);

      // 后台刷新，刷新数据覆盖 setTokens
      const refreshedNFTs: typeof cached = []; // 临时存储刷新结果

      await refreshMyERC721Tokens(provider, NFT_ADDRESS, account, chainId, (nft) => {
        // 每次回调刷新都覆盖之前的 token
        const idx = refreshedNFTs.findIndex((t) => t.tokenId === nft.tokenId);
        if (idx === -1) {
          refreshedNFTs.push(nft);
        } else {
          refreshedNFTs[idx] = nft; // 覆盖
        }

        // 将当前刷新结果直接 setTokens
        setTokens([...refreshedNFTs]);
      });
    }

    load();

    // 🧹 清理（切账号 / 切链 时防止脏更新）
    return () => {
      cancelled = true;
    };
  }, [provider, NFT_ADDRESS, account, chainId]);

  return (
    <div>
      <h2>我的 NFT（标准 ERC721）</h2>

      <div style={{ display: 'grid', gridTemplateColumns: 'repeat(4, 1fr)', gap: 16 }}>
        {tokens.map((nft) => (
          <div
            onClick={() => {
              setSelectedNFT(nft);
              setShowAuctionModal(true);
            }}
            key={nft.tokenId.toString()}
            style={{
              border: '1px solid #ddd',
              padding: 12,
              borderRadius: 8
            }}
          >
            {nft.image && <img src={nft.image} alt={nft.name ?? `NFT #${nft.tokenId}`} style={{ width: '100%', borderRadius: 4 }} />}

            <p style={{ marginTop: 8 }}>{nft.name ?? `NFT #${nft.tokenId.toString()}`}</p>

            <small>Token ID: {nft.tokenId.toString()}</small>
          </div>
        ))}
      </div>
      {showAuctionModal && selectedNFT && <CreateAuctionModal provider={provider} account={account} nft={selectedNFT} onClose={() => setShowAuctionModal(false)} />}
    </div>
  );
}
