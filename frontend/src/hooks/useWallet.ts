import { useEffect, useState, useCallback } from 'react';
import { BrowserProvider, JsonRpcSigner } from 'ethers';
import { CHAINS } from '../config/config';

declare global {
    interface Window {
        ethereum?: any;
    }
}

export async function switchNetwork(chainKey: keyof typeof CHAINS) {
    const chain = CHAINS[chainKey];

    try {
        await window.ethereum.request({
            method: 'wallet_switchEthereumChain',
            params: [{ chainId: chain.chainId }]
        });
    } catch (err: any) {
        if (err.code === 4902) {
            await window.ethereum.request({
                method: 'wallet_addEthereumChain',
                params: [chain]
            });
        } else {
            throw err;
        }
    }
}


async function addSepolia() {
    await window.ethereum.request({
        method: 'wallet_addEthereumChain',
        params: [
            {
                chainId: '0xaa36a7',
                chainName: 'Sepolia Test Network',
                nativeCurrency: {
                    name: 'Sepolia ETH',
                    symbol: 'ETH',
                    decimals: 18
                },
                rpcUrls: ['https://rpc.sepolia.org'],
                blockExplorerUrls: ['https://sepolia.etherscan.io']
            }
        ]
    });
}



export function useWallet() {
    const [provider, setProvider] = useState<BrowserProvider | null>(null);
    const [signer, setSigner] = useState<JsonRpcSigner | null>(null);
    const [account, setAccount] = useState<string | null>(null);
    const [networkName, setNetworkName] = useState<string | null>(null);
    const [chainId, setChainID] = useState<bigint | null>(null);

    // 👉 只在用户主动点击时调用
    const connect = useCallback(async () => {
        if (!window.ethereum) {
            alert('请安装 MetaMask');
            return;
        }

        const accounts = await window.ethereum.request({
            method: 'eth_requestAccounts'
        });

        const provider = new BrowserProvider(window.ethereum);
        const signer = await provider.getSigner(accounts[0]);
        console.log('Network:', (await provider.getNetwork()).chainId);

        setProvider(provider);
        setSigner(signer);
        setAccount(accounts[0]);
        setNetworkName(await provider.getNetwork().then(network => network.name));
        setChainID(await provider.getNetwork().then(network => network.chainId));
    }, []);

    const disconnect = () => {
        setProvider(null);
        setSigner(null);
        setAccount(null);
    };

    // ✅ 正确监听账户切换
    useEffect(() => {
        if (!window.ethereum) return;

        const handleAccountsChanged = async (accounts: string[]) => {
            console.log('accountsChanged:', accounts);

            if (accounts.length === 0) {
                disconnect();
                return;
            }

            const provider = new BrowserProvider(window.ethereum);
            const signer = await provider.getSigner(accounts[0]);

            setProvider(provider);
            setSigner(signer);
            setAccount(accounts[0]);
            setNetworkName(await provider.getNetwork().then(network => network.name));
            setChainID(await provider.getNetwork().then(network => network.chainId));
        };

        window.ethereum.on('accountsChanged', handleAccountsChanged);

        return () => {
            window.ethereum.removeListener('accountsChanged', handleAccountsChanged);
        };
    }, []);

    return {
        provider,
        signer,
        account,
        networkName,
        chainId,
        connect,
        disconnect
    };
}
