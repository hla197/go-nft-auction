import React, { useState, useEffect } from 'react';
import { Contract, ethers } from 'ethers';
import { AUCTION_ADDRESS, ERC20_ADDRESS } from '../config/config';
import { ERC20_ABI } from '../abi/ERC20';
import { AUCTION_ABI } from '../abi/NFTAuction';

interface BidModalProps {
  isOpen: boolean;
  onClose: () => void;
  auctionId: number;
  signer: any;
  account: string;
  payToken: 'ETH' | 'ERC20';
  ethPrice: BigInt;
  erc20Price: BigInt;
  onBidPlaced: () => void;
}

const BidModal: React.FC<BidModalProps> = ({ isOpen, onClose, auctionId, signer, account, payToken, ethPrice, erc20Price, onBidPlaced }) => {
  const [bidAmount, setBidAmount] = useState<string>('');
  const [loading, setLoading] = useState<boolean>(false);
  const [loadingMessage, setLoadingMessage] = useState<string>('');
  const [selectedPayToken, setSelectedPayToken] = useState<'ETH' | 'ERC20'>(payToken);

  const submitBid = async () => {
    if (!bidAmount || parseFloat(bidAmount) <= 0) return alert('请输入有效的出价');

    try {
      setLoading(true);
      setLoadingMessage('正在提交出价...');
      const auctionContract = new Contract(AUCTION_ADDRESS, AUCTION_ABI, signer);

      if (selectedPayToken === 'ETH') {
        const value = ethers.parseEther(bidAmount);
        const tx = await auctionContract.placeBid(auctionId, ethers.ZeroAddress, value, { value });
        await tx.wait();
      } else {
        const erc20Contract = new Contract(ERC20_ADDRESS, ERC20_ABI, signer);
        const decimals = await erc20Contract.decimals();
        const amount = ethers.parseUnits(bidAmount, decimals);

        const allowance = await erc20Contract.allowance(account, AUCTION_ADDRESS);
        if (allowance < amount) {
          setLoadingMessage('正在授权ERC20代币');
          const approveTx = await erc20Contract.approve(AUCTION_ADDRESS, amount);
          await approveTx.wait();
        }
        setLoadingMessage('正在提交出价');
        const tx = await auctionContract.placeBid(auctionId, ERC20_ADDRESS, amount);
        await tx.wait();
      }

      setLoadingMessage('出价成功');
      // 关闭模态框
      onClose();
      onBidPlaced(); // 通知主组件刷新拍卖列表
    } catch (err) {
      console.error(err);
      alert('出价失败');
      setLoading(false);
      setLoadingMessage('');
    }
  };

  useEffect(() => {
    if (!isOpen) {
      setBidAmount('');
      setLoading(false);
      setLoadingMessage('');
    }
  }, [isOpen]);

  // 更新选中的币种
  const handlePayTokenChange = (token: 'ETH' | 'ERC20') => {
    setSelectedPayToken(token);
    setBidAmount(''); // 重置出价金额
  };

  if (!isOpen) return null;

  return (
    <div style={styles.overlay}>
      <div style={styles.modal}>
        <h3>出价</h3>

        <div style={styles.radioGroup}>
          <label>
            <input
              type="radio"
              checked={selectedPayToken === 'ERC20'}
              onChange={() => handlePayTokenChange('ERC20')}
            />{' '}
            USDC
          </label>
          <label style={styles.radioLabel}>
            <input
              type="radio"
              checked={selectedPayToken === 'ETH'}
              onChange={() => handlePayTokenChange('ETH')}
            />{' '}
            ETH
          </label>
        </div>

        <p style={styles.rateText}>
          {selectedPayToken} To USD Rate: {ethers.formatUnits(selectedPayToken === 'ETH' ? ethPrice : erc20Price, 8)}
        </p>

        <input
          type="number"
          value={bidAmount}
          placeholder={selectedPayToken === 'ETH' ? 'ETH 数量' : 'USDC 数量'}
          onChange={(e) => setBidAmount(e.target.value)}
          style={styles.inputField}
        />

        <p style={styles.estimateText}>
          估算 USD: ${parseFloat(ethers.formatUnits(selectedPayToken === 'ETH' ? ethPrice : erc20Price, 8)) * parseFloat(bidAmount).toFixed(2)}
        </p>

        {loading && <p>{loadingMessage}</p>}

        <div style={styles.buttonGroup}>
          <button onClick={submitBid} disabled={loading} style={styles.submitButton}>
            {loading ? '提交中...' : '确认'}
          </button>
          <button onClick={onClose} style={styles.cancelButton}>
            取消
          </button>
        </div>
      </div>
    </div>
  );
};

const styles = {
  overlay: {
    position: 'fixed',
    inset: 0,
    background: 'rgba(0, 0, 0, 0.4)',
  },
  modal: {
    background: '#fff',
    padding: '20px',
    width: '350px',
    margin: '100px auto',
    borderRadius: '8px',
    textAlign: 'center',
  },
  radioGroup: {
    display: 'flex',
    justifyContent: 'center',
    marginBottom: '15px',
  },
  radioLabel: {
    marginLeft: '20px',
  },
  rateText: {
    marginTop: '10px',
    fontSize: '14px',
    color: '#555',
  },
  inputField: {
    marginTop: '10px',
    padding: '10px',
    borderRadius: '5px',
    border: '1px solid #ccc',
    fontSize: '16px',
    textAlign: 'center',
  },
  estimateText: {
    marginTop: '15px',
    fontSize: '16px',
    fontWeight: 'bold',
  },
  buttonGroup: {
    marginTop: '20px',
    display: 'flex',
    justifyContent: 'space-around',
  },
  submitButton: {
    padding: '10px 20px',
    backgroundColor: '#28a745',
    color: '#fff',
    border: 'none',
    borderRadius: '5px',
    cursor: 'pointer',
    width: '45%',
  },
  cancelButton: {
    padding: '10px 20px',
    backgroundColor: '#dc3545',
    color: '#fff',
    border: 'none',
    borderRadius: '5px',
    cursor: 'pointer',
    width: '45%',
  },
};

export default BidModal;
