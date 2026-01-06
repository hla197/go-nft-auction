import { useEffect, useState } from 'react';
import { Contract, ethers  } from 'ethers';
import CreateAuctionModal from './CreateAuctionModal';
import { ERC721_ABI } from '../abi/ERC721';
import { NFT_ADDRESS } from '../config/config';

interface Props {
    signer: any;
    account: string;
    onSelect: (tokenId: number) => void;
}

const TRANSFER_TOPIC = ethers.id(
    'Transfer(address,address,uint256)'
);

export async function getMyERC721Tokens(
    provider: ethers.Provider,
    nftAddress: string,
    owner: string,
    fromBlock = 0
): Promise<number[]> {
    const ownerTopic = ethers.zeroPadValue(owner, 32);

    // 1️⃣ 收到的 NFT
    const receivedLogs = await provider.getLogs({
        address: nftAddress,
        fromBlock,
        topics: [
            TRANSFER_TOPIC,
            null,
            ownerTopic
        ]
    });
    console.log('receivedLogs', receivedLogs);

    // 2️⃣ 转出的 NFT
    const sentLogs = await provider.getLogs({
        address: nftAddress,
        fromBlock,
        topics: [
            TRANSFER_TOPIC,
            ownerTopic,
            null
        ]
    });

    console.log('sentLogs', sentLogs);
    const owned = new Set<number>();

    // 收到 = 加
    for (const log of receivedLogs) {
        const tokenId = Number(log.topics[3]);
        owned.add(tokenId);
    }

    // 转出 = 删
    for (const log of sentLogs) {
        const tokenId = Number(log.topics[3]);
        owned.delete(tokenId);
    }

    return [...owned];
}

export default function MyNFTs({ signer, account }: Props) {
    const [nfts, setNfts] = useState<number[]>([]);
    const [selectedTokenId, setSelectedTokenId] = useState<number | null>(null);

    async function loadMyNFTs() {
        console.log('loadMyNFTs');
        if (!signer || !account) return;

        const nft = new Contract(NFT_ADDRESS, ERC721_ABI, signer);
        const balance = await nft.balanceOf(account);
        console.log('balance', balance);
        const list: number[] = [];
        for (let i = 0; i < balance; i++) {
            const tokenId = await nft.tokenOfOwnerByIndex(account, i);
            list.push(Number(tokenId));
        }
        setNfts(list);
    }

    useEffect(() => {
        loadMyNFTs();
    }, [signer, account]);

    return (
        <div>
            <h2>我的 NFT</h2>

            <div style={{ display: 'flex', gap: 16, flexWrap: 'wrap' }}>
                {nfts.map((id) => (
                    <div
                        key={id}
                        style={{
                            border: '1px solid #ccc',
                            padding: 12,
                            cursor: 'pointer'
                        }}
                        onClick={() => setSelectedTokenId(id)}
                    >
                        <p>NFT #{id}</p>
                    </div>
                ))}
            </div>

            {selectedTokenId !== null && (
                <CreateAuctionModal
                    signer={signer}
                    nftAddress={NFT_ADDRESS}
                    tokenId={selectedTokenId}
                    onClose={() => setSelectedTokenId(null)}
                />
            )}
        </div>
    );
}
