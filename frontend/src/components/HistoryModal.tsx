import React, { useState, useEffect } from 'react';
import { getNftHistory } from '../api/api'; // 导入 getNftHistory
import { ethers } from 'ethers';

interface HistoryModalProps {
  isOpen: boolean;
  onClose: () => void;
  nftContract: string;
  nftTokenId: string;
  chainId: number;
}

const HistoryModal: React.FC<HistoryModalProps> = ({ isOpen, onClose, nftContract, nftTokenId, chainId }) => {
  const [historicalAuctions, setHistoricalAuctions] = useState<any[]>([]);
  const [loadingHistory, setLoadingHistory] = useState(false);

  // 获取历史拍卖数据
  useEffect(() => {
    if (!isOpen) return;

    async function fetchHistory() {
      setLoadingHistory(true);
      try {
        const response = await getNftHistory(nftContract, nftTokenId, chainId.toString());
        setHistoricalAuctions(response.data);
      } catch (error) {
        console.error('Error fetching historical auctions:', error);
        alert('获取历史拍卖失败');
      } finally {
        setLoadingHistory(false);
      }
    }

    fetchHistory();
  }, [isOpen, nftContract, nftTokenId, chainId]);

  // 处理 bid_address，将中间部分用 * 代替
  const obfuscateAddress = (address: string) => {
    if (!address) return '';
    const start = address.slice(0, 6);
    const end = address.slice(-4);
    return `${start}***${end}`;
  };

  return (
    isOpen && (
      <div style={{ position: 'fixed', inset: 0, background: 'rgba(0,0,0,0.4)' }}>
        <div style={{ background: '#fff', padding: 20, width: 700, margin: '100px auto', borderRadius: '8px' }}>
          <h3 style={{ textAlign: 'center', marginBottom: '20px' }}>历史拍卖记录</h3>

          {loadingHistory ? (
            <p>加载历史拍卖中...</p>
          ) : (
            <table style={{ width: '100%', borderCollapse: 'collapse', borderRadius: '8px' }}>
              <thead style={{ backgroundColor: '#f1f1f1' }}>
                <tr>
                  <th style={{ padding: '10px', textAlign: 'left', fontWeight: 'bold' }}>Auction ID</th>
                  <th style={{ padding: '10px', textAlign: 'left', fontWeight: 'bold' }}>出价者</th>
                  <th style={{ padding: '10px', textAlign: 'left', fontWeight: 'bold' }}>出价金额 (USD)</th>
                  <th style={{ padding: '10px', textAlign: 'left', fontWeight: 'bold' }}>拍卖结束时间</th>
                </tr>
              </thead>
              <tbody>
                {historicalAuctions.map((auction: any) => (
                  <tr key={auction.auction_id} style={{ borderBottom: '1px solid #ddd' }}>
                    <td style={{ padding: '10px' }}>{auction.auction_id}</td>
                    <td style={{ padding: '10px' }}>{obfuscateAddress(auction.bid_address)}</td>
                    <td style={{ padding: '10px' }}>{ethers.formatUnits(auction.bid_price_usd, 8)}</td>
                    <td style={{ padding: '10px' }}>{new Date(auction.end_time * 1000).toLocaleString()}</td>
                  </tr>
                ))}
              </tbody>
            </table>
          )}

          <div style={{ marginTop: 16, textAlign: 'center' }}>
            <button
              onClick={onClose}
              style={{
                padding: '10px 20px',
                backgroundColor: '#007bff',
                color: '#fff',
                border: 'none',
                borderRadius: '5px',
                cursor: 'pointer',
              }}
            >
              关闭
            </button>
          </div>
        </div>
      </div>
    )
  );
};

export default HistoryModal;
