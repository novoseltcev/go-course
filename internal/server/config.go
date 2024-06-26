package server

type Config struct {
	Address string			`env:"ADDRESS"`
	StoreInterval int 		`env:"STORE_INTERVAL"`
	FileStoragePath string	`env:"FILE_STORAGE_PATH"`
	Restore bool 			`env:"RESTORE"`
	DatabaseDsn string		`env:"DATABASE_DSN"`
	SecretKey string		`env:"KEY"`
}
