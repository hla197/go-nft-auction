# 项目结构说明文档
## 项目概述

这是一个包含 Go 后端和 前端 React/Vite 应用的全栈项目，旨在实现一个拍卖平台或相关的去中心化应用（DApp）。本项目包含后端服务、前端展示和配置文件，以及相应的模块化代码结构。以下是各个文件夹和文件的详细说明。

internal: 为go语言的模块

hardhat: 这是一个包含智能合约的目录，包含合约的源代码、编译后的字节码和 ABI 文件。

frontend: 为前端 React 应用的目录，包含 HTML、CSS、JavaScript 文件等。

web: frontend build出来的文件

```
├── cmd                        # 项目的主应用入口
│   └── main.go                # Go 应用的启动入口文件，包含了应用的主要逻辑
├── config.dev.yaml            # 开发环境的配置文件，通常包括数据库、API密钥等开发所需配置
├── config.prod.yaml           # 生产环境的配置文件，包含用于生产环境的配置项
├── config.smple.yaml          # 示例配置文件，通常用于说明如何配置应用
├── frontend                   # 前端相关的文件夹，通常包含 HTML、CSS、JavaScript 文件等
│   ├── index.html             # 前端应用的入口 HTML 文件
│   ├── node_modules           # 项目的依赖包，自动生成并管理，包含所有的 npm 包
│   ├── package-lock.json      # 记录具体的版本依赖，确保项目依赖的版本一致
│   ├── package.json           # 项目的 npm 配置文件，包含依赖项、脚本等
│   ├── src                    # 前端源代码文件夹，包含 React 组件、样式等
│   └── vite.config.ts         # Vite 配置文件，用于配置开发环境和构建过程
├── go.mod                     # Go 语言的模块依赖管理文件
├── go.sum                     # Go 语言的依赖校验文件，确保项目依赖的一致性
├── hardhat                    # 项目根目录，包含与 Hardhat 相关的文件和文件夹
│   ├── README.md              # 项目文档文件，通常用于描述项目的功能、安装步骤、配置和使用方法
│   ├── artifacts              # 存储编译后的智能合约的文件夹。包括合约的 ABI 和字节码等信息
│   ├── cache                  # 存储 Hardhat 编译过程中的缓存文件，加速后续编译操作
│   ├── contracts              # 存放 Solidity 智能合约源代码的文件夹，所有的 .sol 文件都在这里
│   ├── deploy                 # 存放智能合约部署脚本的文件夹，用于将合约部署到不同的网络
│   ├── deployments            # 存储已部署合约的地址和状态的文件夹，包含不同环境下的部署记录
│   ├── hardhat.config.js      # Hardhat 项目的配置文件，包含编译器配置、网络配置、插件等设置
│   ├── package-lock.json      # 由 npm 自动生成，锁定项目所有依赖的版本，确保一致性
│   ├── package.json           # 项目的配置文件，定义了项目的依赖、开发脚本以及项目的元数据
│   └── test                   # 存放与智能合约相关的测试文件夹，通常使用 JavaScript 或 TypeScript 编写
├── internal                   # 内部包，包含应用的核心代码
│   ├── abi                    # 包含与区块链交互的 ABI 文件
│   ├── config                 # 配置文件夹，包含应用的配置信息（如数据库、API接口等）
│   ├── global                 # 全局设置和常量文件
│   ├── handlers               # 包含 HTTP 请求的处理逻辑，例如控制器（controllers）
│   ├── infra                  # 基础设施相关代码，如数据库连接、外部 API 调用等
│   ├── logger                 # 日志相关的工具和配置
│   ├── middleware             # 中间件相关代码，例如跨域设置、请求验证等
│   ├── models                 # 数据库模型定义，包含表结构和 ORM 相关代码
│   ├── routers                # 路由相关代码，定义 HTTP 路径和请求方法的映射
│   ├── scanner                # 扫描器相关功能，如文件扫描、事件监听等
│   └── utils                  # 工具函数库，包含各种辅助工具函数
├── logs                       # 日志文件夹，存储应用的运行日志
│   └── gin.log                # Gin 框架的运行日志文件
└── web                        # 静态资源和前端文件
    ├── assets                 # 静态资源文件夹
    └── index.html             # 网页应用的首页

```

## 目录和文件说明
### cmd

main.go: 项目的入口文件，负责初始化和启动服务。通常包含应用程序的主要启动逻辑，如连接数据库、启动 HTTP 服务等。

### config

config.dev.yaml: 开发环境配置文件。存储与开发环境相关的配置，如数据库连接、API 密钥等。

config.prod.yaml: 生产环境配置文件。用于生产环境的配置，通常包括安全相关的密钥、数据库的生产配置等。

config.smple.yaml: 示例配置文件。用来展示如何配置该应用程序的环境参数，用户可以参考此文件来配置其他环境的设置。

### frontend

前端部分使用 React 和 Vite 打包工具开发。

index.html: 作为前端应用程序的入口 HTML 文件，包含前端应用的基础结构和加载资源。

node_modules: 包含所有由 npm 安装的依赖包。

package-lock.json: 自动生成的文件，锁定项目依赖的版本，以保证每次安装时的依赖一致性。

package.json: 存储前端项目的依赖、脚本命令等配置信息。

src: 前端源代码文件夹，包含前端组件、样式、页面逻辑等。

vite.config.ts: Vite 配置文件，设置了开发服务器、打包设置等。

### go.mod & go.sum

这些文件用于 Go 项目的模块管理，记录项目依赖包及其版本。

internal
这是一个 Go 项目的内部文件夹，通常存放业务逻辑的实现，包括：

abi 存放与区块链交互的 ABI 文件，定义合约接口。

config 存放应用的配置文件。

handlers 处理具体的 HTTP 请求，通常定义控制器（Controller）。

infra 包含应用的基础设施代码，如数据库、外部 API 的连接和服务。

logger 负责应用的日志管理。

middleware 存放应用的中间件，例如请求验证、跨域等功能。

models 包含数据库模型定义。

routers 负责应用的路由配置，定义不同的路由和 HTTP 方法。

scanner 可能涉及事件监听和文件扫描功能。

utils 包含一些常用的工具函数。

### logs/gin.log

存放 Gin 框架运行时的日志文件，记录系统运行的相关信息和错误。

### web
存放静态资源和前端文件

主要frontend将前端打包到build后的文件放在web下，然后通过gin启动服务，将web下的文件作为静态资源返回给前端。