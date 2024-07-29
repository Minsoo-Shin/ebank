package config

import (
	"flag"
	"log"
	"time"
)

type Config struct {
	DB     DBConfig
	Jwt    JwtConfig
	Server ServerConfig
}

type DBConfig struct {
	UserTablePath        string
	AccountTablePath     string
	TransactionTablePath string
}

type JwtConfig struct {
	SecretKey string
	Duration  time.Duration
}

type ServerConfig struct {
	Port string
}

func New() Config {
	userFilePathPtr := flag.String("user_file_path", "data/user.json", "user_file_path")
	accountFilePathPtr := flag.String("account_file_path", "data/account.json", "account_file_path")
	transactionFilePathPtr := flag.String("transaction_file_path", "data/transaction.json", "transaction_file_path")

	secretPtr := flag.String("secret", "happy_coding", "secret key")
	durationPtr := flag.Duration("duration", 15*time.Minute, "token duration")

	portPtr := flag.String("port", ":50051", "port number")

	flag.Parse()

	config := Config{
		DB: DBConfig{
			UserTablePath:        *userFilePathPtr,
			AccountTablePath:     *accountFilePathPtr,
			TransactionTablePath: *transactionFilePathPtr,
		},
		Jwt: JwtConfig{
			SecretKey: *secretPtr,
			Duration:  *durationPtr,
		},
		Server: ServerConfig{
			Port: *portPtr,
		},
	}

	config.Validate()

	return config
}

func (r Config) Validate() {
	if r.DB.UserTablePath == "" || r.DB.AccountTablePath == "" || r.DB.TransactionTablePath == "" {
		log.Fatal("File paths cannot be empty")
	}
	if r.Jwt.SecretKey == "" {
		log.Fatal("Secret key cannot be empty")
	}
	if r.Jwt.Duration == 0 {
		log.Fatal("Duration cannot be 0")
	}
	if r.Server.Port == "" {
		log.Fatal("Port number cannot be empty")
	}
}
