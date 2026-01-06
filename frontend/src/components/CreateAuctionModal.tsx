import { useState } from 'react';
import { Contract, ethers, parseEther } from 'ethers';
import { ERC721_ABI } from '../abi/ERC721';
import { AUCTION_ABI } from '../abi/NFTAuction';
import { AUCTION_ADDRESS } from '../config/config';

export default function CreateAuctionModal({
  provider,
  account,
  nft,
  onClose
}: any) {
  const [duration, setDuration] = useState('3600');
  const [price, setPrice] = useState('0.01');
  const [loading, setLoading] = useState(false);
  const [step, setStep] = useState('');
  const [error, setError] = useState('');

  async function handleCreate() {
    setError('');

    if (Number(duration) <= 0) {
      setError('拍卖时长必须大于 0');
      return;
    }

    if (Number(price) <= 0) {
      setError('起拍价必须大于 0');
      return;
    }

    try {
      setLoading(true);

      const signer = await provider.getSigner();
      const erc721 = new Contract(nft.contract, ERC721_ABI, signer);
      const auction = new Contract(AUCTION_ADDRESS, AUCTION_ABI, signer);

      // 1️⃣ 校验 owner
      setStep('检查 NFT 所有权');
      const owner = await erc721.ownerOf(nft.tokenId);
      if (owner.toLowerCase() !== account.toLowerCase()) {
        throw new Error('你不是该 NFT 的拥有者');
      }

      // 2️⃣ 授权
      setStep('检查授权状态');
      const approvedForAll = await erc721.isApprovedForAll(account, AUCTION_ADDRESS);
      if (!approvedForAll) {
        setStep('请求 NFT 授权');
        const tx = await erc721.setApprovalForAll(AUCTION_ADDRESS, true);
        await tx.wait();
      }

      // 3️⃣ 创建拍卖
      setStep('创建拍卖中');
      console.log('创建拍卖中', nft.contract, nft.tokenId);


      const tx = await auction.startAuction(
        BigInt(duration),
        parseEther(price),
        nft.contract,
        nft.tokenId,
        ethers.ZeroAddress
      );
      await tx.wait();

      onClose();
    } catch (e: any) {
      setError(e.message ?? '创建拍卖失败');
    } finally {
      setLoading(false);
      setStep('');
    }
  }

  return (
    <div style={mask}>
      <div style={modal}>
        {/* 头部 */}
        <div style={header}>
          <h3 style={{ margin: 0 }}>创建拍卖</h3>
          <button onClick={onClose} style={closeBtn}>✕</button>
        </div>

        {/* NFT 信息 */}
        <div style={nftPreview}>
          {nft.image && (
            <img
              src={nft.image}
              alt={nft.name}
              style={nftImage}
            />
          )}
          <div>
            <div style={{ fontWeight: 600 }}>
              {nft.name ?? `NFT #${nft.tokenId.toString()}`}
            </div>
            <div style={{ fontSize: 12, color: '#666' }}>
              Token ID: {nft.tokenId.toString()}
            </div>
          </div>
        </div>

        {/* 表单 */}
        <div style={form}>
          <label style={label}>拍卖时长（秒）</label>
          <input
            style={input}
            value={duration}
            onChange={(e) => setDuration(e.target.value)}
          />

          <label style={label}>起拍价（ETH）</label>
          <input
            style={input}
            value={price}
            onChange={(e) => setPrice(e.target.value)}
          />
        </div>

        {/* 步骤 */}
        {step && (
          <div style={stepBox}>
            <span style={spinner} />
            <span>{step}</span>
          </div>
        )}

        {/* 错误 */}
        {error && (
          <div style={errorBox}>
            ⚠ {error}
          </div>
        )}

        {/* 操作 */}
        <div style={actions}>
          <button onClick={onClose} style={cancelBtn} disabled={loading}>
            取消
          </button>
          <button
            onClick={handleCreate}
            style={{ ...confirmBtn, opacity: loading ? 0.6 : 1 }}
            disabled={loading}
          >
            {loading ? '处理中...' : '确认创建'}
          </button>
        </div>
      </div>
    </div>
  );
}

/* ================= 样式 ================= */

const mask = {
  position: 'fixed' as const,
  inset: 0,
  background: 'rgba(0,0,0,0.5)',
  display: 'flex',
  justifyContent: 'center',
  alignItems: 'center'
};

const modal = {
  width: 400,
  background: '#fff',
  borderRadius: 10,
  padding: 20
};

const header = {
  display: 'flex',
  justifyContent: 'space-between',
  alignItems: 'center',
  marginBottom: 12
};

const closeBtn = {
  border: 'none',
  background: 'transparent',
  cursor: 'pointer',
  fontSize: 16
};

const nftPreview = {
  display: 'flex',
  gap: 12,
  marginBottom: 16
};

const nftImage = {
  width: 64,
  height: 64,
  borderRadius: 6,
  objectFit: 'cover' as const
};

const form = {
  display: 'flex',
  flexDirection: 'column' as const,
  gap: 6
};

const label = {
  fontSize: 12,
  color: '#555'
};

const input = {
  height: 34,
  padding: '0 8px',
  borderRadius: 6,
  border: '1px solid #ddd'
};

const stepBox = {
  display: 'flex',
  alignItems: 'center',
  gap: 8,
  fontSize: 12,
  color: '#2563eb',
  marginTop: 10
};

const spinner = {
  width: 12,
  height: 12,
  border: '2px solid #ccc',
  borderTopColor: '#2563eb',
  borderRadius: '50%'
};

const errorBox = {
  marginTop: 8,
  fontSize: 12,
  color: '#dc2626'
};

const actions = {
  display: 'flex',
  justifyContent: 'flex-end',
  gap: 8,
  marginTop: 16
};

const cancelBtn = {
  padding: '6px 12px'
};

const confirmBtn = {
  padding: '6px 12px',
  background: '#2563eb',
  color: '#fff',
  border: 'none',
  borderRadius: 6
};
