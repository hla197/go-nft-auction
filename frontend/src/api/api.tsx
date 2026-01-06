import { get, post } from "../utils/https";

// 获取某个 NFT 的详细信息
export const getNftHistory = async (nft_contract: string, nft_token_id: string, chain_id: string) => {
  return await post(`/auction/getNftHistory`, {
    nft_contract,
    nft_token_id,
    chain_id
  });
};

// 创建拍卖
