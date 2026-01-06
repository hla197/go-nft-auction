// src/config/contracts.ts
export const NFT_ADDRESS = '0xb834404C86987E297b074C7b4C148368E64F8d7C';
export const AUCTION_ADDRESS = '0x9a22E86A2Db8AbD5bF3aD302d73735F87510ACF9';
export const ERC20_ADDRESS = '0x9F9557a99E38C1A0b98e90f1dD0141C307fe3B13';
export const BASE_URL = "http://127.0.0.1:8088/";

export const NFT_START_BLOCK = 9975260; // NFT合约创建事件

export const CHAINS = {
    sepolia: {
        chainId: '0xaa36a7',
        chainName: 'Sepolia Test Network',
        rpcUrls: ['https://rpc.sepolia.org'],
        blockExplorerUrls: ['https://sepolia.etherscan.io'],
        nativeCurrency: {
            name: 'Sepolia ETH',
            symbol: 'ETH',
            decimals: 18
        }
    },
    hardhat: {
        chainId: '0x7a69',
        chainName: 'Hardhat Local',
        rpcUrls: ['http://127.0.0.1:8545'],
        nativeCurrency: {
            name: 'ETH',
            symbol: 'ETH',
            decimals: 18
        }
    }
};