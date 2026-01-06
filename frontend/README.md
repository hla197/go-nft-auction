```
├── App.tsx                   # 应用的根组件，包含了整个应用的结构和路由配置。
├── abi                       # 存放智能合约的 ABI 文件，用于与区块链交互。
│   ├── ERC20.ts              # ERC20 代币合约的 ABI 文件，用于与 ERC20 代币合约进行交互。
│   ├── ERC721.ts             # ERC721 标准 NFT 合约的 ABI 文件，用于与 ERC721 合约进行交互。
│   └── NFTAuction.ts         # NFT 拍卖合约的 ABI 文件，用于与 NFT 拍卖合约进行交互。
├── api                       # 存放与后端交互的 API 调用函数。
│   └── api.tsx               # 主要的 API 请求文件，包含发送请求、处理响应等功能。
├── components                # 存放应用的 UI 组件。
│   ├── BidModal.tsx          # 出价弹窗组件，允许用户在拍卖中出价。
│   ├── CreateAuctionModal.tsx# 创建拍卖的弹窗组件，允许用户创建新的拍卖。
│   ├── HistoryModal.tsx      # 显示历史拍卖记录的弹窗组件。
│   └── MyNFTs.tsx            # 展示用户 NFT 列表的组件，展示用户拥有的 NFT。
├── config                    # 存放应用的配置文件。
│   └── config.ts             # 配置文件，包含一些常量、API 密钥、智能合约地址等。
├── hooks                     # 存放自定义 hook 的目录。
│   └── useWallet.ts          # 自定义 hook，用于管理钱包连接、账户等与钱包相关的逻辑。
├── main.tsx                  # 应用的入口文件，通常在这里渲染根组件 App.tsx。
├── pages                     # 存放应用的页面组件。
│   ├── AuctionList.tsx       # 拍卖列表页面，显示所有正在进行的拍卖。
│   ├── Home.tsx              # 首页，展示应用的欢迎信息和导航。
│   └── MyNFTPage.tsx         # 用户的 NFT 页面，展示当前用户的 NFT 资产。
└── utils                     # 存放一些工具函数的目录。
    ├── common.tsx            # 通用工具函数，可能包含一些常见的格式化、转换函数。
    ├── erc721.tsx            # 与 ERC721 合约相关的工具函数，可能包括 NFT 的获取、解析等。
    └── https.tsx             # 用于处理 HTTPS 请求的工具函数，封装一些请求逻辑。

```

详细说明：
根目录文件

App.tsx:
这是 React 应用的根组件，通常负责应用的总体结构，包含路由配置（例如使用 React Router）和全局状态管理。

main.tsx:
这是应用的入口文件，通常在这里渲染根组件 App.tsx。这个文件通常会设置一些全局配置（如 React Router、Redux 等）并将应用渲染到 DOM 上。

文件夹：abi

ERC20.ts:
该文件包含 ERC20 代币合约的 ABI（应用二进制接口）。它允许前端与 ERC20 合约进行交互，查询代币余额、转账等操作。

ERC721.ts:
包含 ERC721 标准（NFT 合约）的 ABI 文件，允许与 NFT 合约进行交互，查询、转移 NFT 等操作。

NFTAuction.ts:
该文件包含 NFT 拍卖合约的 ABI，允许与 NFT 拍卖合约进行交互，包括拍卖创建、出价、结束拍卖等功能。

文件夹：api

api.tsx:
该文件负责处理前端与后端 API 的交互。它可能会封装 GET、POST 请求，处理响应数据等。所有与后端的 API 调用都可能在这里定义。

文件夹：components

BidModal.tsx:
这是一个弹窗组件，允许用户在拍卖中进行出价。该组件可能包括输入框、确认按钮等元素，允许用户提交出价。

CreateAuctionModal.tsx:
这是一个弹窗组件，允许用户创建新的拍卖。通常包括输入拍卖的 NFT 信息、起拍价、拍卖时间等。

HistoryModal.tsx:
显示历史拍卖记录的弹窗组件，可以查看过去的拍卖详情，例如出价者、出价金额、拍卖时间等。

MyNFTs.tsx:
展示用户拥有的所有 NFT 组件。通过与智能合约交互，查询并展示用户的 NFT 列表。

文件夹：config

config.ts:
该文件存放一些配置项，如智能合约地址、API 密钥、其他常量等。这些配置在整个应用中可能都会使用到。

文件夹：hooks

useWallet.ts:
这是一个自定义 Hook，用于管理与钱包的交互。它可能包含连接钱包、获取账户地址、切换网络等功能。它使得与钱包的交互更加模块化和可复用。

文件夹：pages

AuctionList.tsx:
显示所有正在进行的拍卖列表页面，用户可以查看拍卖的详细信息，并参与出价。

Home.tsx:
主页组件，通常是应用的欢迎页面，展示一些基本的信息或导航功能。

MyNFTPage.tsx:
显示当前用户持有的 NFT 资产的页面，允许用户查看他们的 NFT 和参与拍卖。

文件夹：utils

common.tsx:
通用的工具函数，包含常见的格式化、数据处理、状态管理等辅助函数。

erc721.tsx:
这个文件可能包含与 ERC721 合约相关的工具函数，例如获取 NFT 元数据、解析 NFT 信息等。

https.tsx:
用于封装 HTTP 请求逻辑，可能包含用于与后端进行 HTTPS 通信的函数（如 axios 请求封装、错误处理等）。