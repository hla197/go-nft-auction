export const AUCTION_ABI = [
  // -------- Read --------
  'function auctionIdCounter() view returns (uint256)',
  'function auctions(uint256) view returns (tuple(' +
    'address seller,' +
    'uint256 startPrice,' +
    'uint256 startPriceUsd,' +
    'uint256 highestBid,' +
    'address highestBidder,' +
    'uint256 highestBidUsd,' +
    'address highestBidderToken,' +
    'uint256 endTime,' +
    'uint256 startTime,' +
    'address nftAddress,' +
    'uint256 tokenId,' +
    'address token,' +
    'bool active' +
  '))',
  'function getChainlinkDataFeedLatestAnswer(address token) public view returns (uint256)',

  // -------- Write --------
  'function startAuction(uint256,uint256,address,uint256,address)',
  'function placeBid(uint256,address,uint256) payable',
  'function endAuction(uint256)',

  // -------- Events --------
  'event AuctionCreated(address indexed seller,uint256 indexed auctionId,uint256 minPrice,uint256 endTime)',
  'event BidPlaced(address indexed bidder,uint256 indexed auctionId,uint256 bidAmount)',
  'event AuctionEnded(address indexed winner,uint256 indexed auctionId,uint256 bidAmount)',
] as const;
