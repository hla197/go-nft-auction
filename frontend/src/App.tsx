import { Routes, Route, Link } from 'react-router-dom';
import { useWallet, switchNetwork } from './hooks/useWallet';
import { useEffect, useState } from 'react';
import { Contract, ethers, formatEther, parseEther } from 'ethers';

import Home from './pages/Home';
import AuctionList from './pages/AuctionList';
import MyNFTPage from './pages/MyNFTPage';
import { ERC20_ADDRESS } from './config/config';
import { ERC20_ABI } from './abi/ERC20';

export default function App() {
  const wallet = useWallet();
  const [ethBalance, setEthBalance] = useState('0');
  const [usdcBalance, setUsdcBalance] = useState('0');
  async function getBalance() {
    if (!wallet.provider || !wallet.account) return;

    const eth = await wallet.provider.getBalance(wallet.account);
    setEthBalance(formatEther(eth));

    const erc20 = new Contract(ERC20_ADDRESS, ERC20_ABI, wallet.provider);
    const bal = await erc20.balanceOf(wallet.account);
    setUsdcBalance(ethers.formatUnits(bal, 6));

    const decimals = await erc20.decimals(); // 6
    console.log(decimals, ethers.formatUnits(bal, decimals));
  }

  useEffect(() => {
    wallet.connect();
    getBalance();
    const handler = (e: CustomEvent) => {
        console.log('收到刷新事件');
        // 执行刷新逻辑
        getBalance();
    };

    window.addEventListener('balance-refresh', handler as EventListener);

    return () => {
        window.removeEventListener('balance-refresh', handler as EventListener);
    };

  }, []);

  useEffect(() => {
    getBalance();
  }, [wallet]);

  function refresh() {
    getBalance();
  }

  return (
    <div style={{ padding: 20 }}>
      <header style={{ marginBottom: 20 }}>
        <Link to="/">首页</Link> | <Link to="/market">拍卖市场</Link> | <Link to="/my-nfts">我的 NFT</Link>
      </header>

      {wallet.account && (
        <div>
          <div>
            <button onClick={wallet.disconnect}>断开连接</button>
            <p>
              当前链：{wallet.networkName} {wallet.chainId}
            </p>
            {wallet.chainId !== BigInt(11155111) && <button onClick={() => switchNetwork('sepolia')}>切换到 Sepolia 测试网</button>}
          </div>
          <div>
            <p style={{ fontSize: 13, color: '#6b7280' }}>
              ETH 余额：{ethBalance} | USDC 余额：{usdcBalance}
            </p>
          </div>
        </div>
      )}

      {!wallet.account ? <button onClick={wallet.connect}>连接钱包</button> : <p>当前账户：{wallet.account}</p>}

      <Routes>
        <Route path="/" element={<Home />} />
        <Route path="/market" element={<AuctionList signer={wallet.signer} account={wallet.account} chainId={wallet.chainId} provider={wallet.provider} />} />
        <Route path="/my-nfts" element={<MyNFTPage signer={wallet.signer} account={wallet.account} chainId={wallet.chainId} provider={wallet.provider} />} />
      </Routes>
    </div>
  );
}
