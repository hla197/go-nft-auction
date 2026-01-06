package config

type AppConfig struct {
	Env  string
	Port string
}

type DatabaseConfig struct {
	Host     string
	Port     int
	Username string
	Password string
	Name     string
}

type RedisConfig struct {
	Host     string
	Port     int
	Password string
	DB       int
	PoolSize int
}

type BlockChainConfig struct {
	RpcUrl               string
	WssUrl               string
	NftAuctionContract   string
	NftAuctionStartBlock uint64
	NftContract          string
	NftStartBlock        uint64
}

type JwtConfig struct {
	Secret string
}
