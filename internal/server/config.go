package server

type Config struct {
	Address         string `env:"ADDRESS"`
	FileStoragePath string `env:"FILE_STORAGE_PATH"`
	DatabaseDsn     string `env:"DATABASE_DSN"`
	SecretKey       string `env:"KEY"`
	StoreInterval   int8   `env:"STORE_INTERVAL"`
	Restore         bool   `env:"RESTORE"`
	CryptoKey       string `env:"CRYPTO_KEY,file"`
}
