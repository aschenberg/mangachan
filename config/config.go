package config

import (
	"fmt"

	"os"
	"path/filepath"
	"runtime"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Server   ServerConfig
	JWT      JWTConfig
	Arango   ArangoConfig
	Redis    RedisConfig
	Password PasswordConfig
	Cors     CorsConfig
	Logger   LoggerConfig
	Otp      OtpConfig
	Oidc     OIDC
}

type ServerConfig struct {
	InternalPort string `env:"INTERNAL_PORT"`
	ExternalPort string `env:"EXTERNAL_PORT"`
	RunMode      string `env:"APP_MODE"`
}

type JWTConfig struct {
	AccessTokenExpireHour  int    `env:"JWT_ACCESS_TOKEN_EXPIRY_HOUR"`
	RefreshTokenExpireHour int    `env:"JWT_REFRESH_TOKEN_EXPIRY_HOUR"`
	AccessTokenSecret      string `env:"JWT_ACCESS_TOKEN_SECRET"`
	RefreshTokenSecret     string `env:"JWT_REFRESH_TOKEN_SECRET"`
}

type LoggerConfig struct {
	FilePath string `env:"JWT_REFRESH_SECRET"`
	Encoding string `env:"JWT_REFRESH_SECRET"`
	Level    string `env:"LOG_LEVEL"`
	Logger   string `env:"LOGGER"`
}

//	type PostgresConfig struct {
//		Host            string
//		Port            string
//		User            string
//		Password        string
//		DbName          string
//		SSLMode         string
//		MaxIdleConns    int
//		MaxOpenConns    int
//		ConnMaxLifetime time.Duration
//	}
type ArangoConfig struct {
	Host     string        `env:"ARANGO_HOST"`
	Port     string        `env:"ARANGO_PORT"`
	User     string        `env:"ARANGO_USERNAME"`
	Password string        `env:"ARANGO_ROOT_PASSWORD"`
	DbName   string        `env:"ARANGO_DBNAME"`
	Timeout  time.Duration `env:"ARANGO_TIMEOUT"`
}

type RedisConfig struct {
	Host               string        `env:"REDIS_HOST"`
	Port               string        `env:"REDIS_PORT"`
	Password           string        `env:"REDIS_PASSWORD"`
	Db                 string        `env:"REDIS_DB"`
	DialTimeout        time.Duration `env:"REDIS_DIAL_TIMEOUT"`
	ReadTimeout        time.Duration `env:"REDIS_READ_TIMEOUT"`
	WriteTimeout       time.Duration `env:"REDIS_WRITE_TIMEOUT"`
	IdleCheckFrequency time.Duration `env:"REDIS_IDLE_CHECK_FREQ"`
	PoolSize           int           `env:"REDIS_POOLSIZE"`
	PoolTimeout        time.Duration `env:"REDIS_POOL_TIMEOUT"`
}

type PasswordConfig struct {
	IncludeChars     bool
	IncludeDigits    bool
	MinLength        int
	MaxLength        int
	IncludeUppercase bool
	IncludeLowercase bool
}

type CorsConfig struct {
	AllowOrigins string
}

type OtpConfig struct {
	ExpireTime time.Duration
	Digits     int
	Limiter    time.Duration
}

type OIDC struct {
	ClientId     string `env:"OIDC_CLIENT_ID"`
	ClientSecret string `env:"OIDC_CLIENT_SECRET"`
	RedirectUrl  string `env:"OIDC_REDIRECT_URL"`
	IssuerUrl    string `env:"OIDC_ISSUER_URL"`
}

func NewConfig() *Config {
	cfg := &Config{}
	cwd := projectRoot()
	envFilePath := cwd + ".env"

	err := readEnv(envFilePath, cfg)
	if err != nil {
		panic(err)
	}

	return cfg
}

func readEnv(envFilePath string, cfg *Config) error {
	envFileExists := checkFileExists(envFilePath)

	if envFileExists {
		err := cleanenv.ReadConfig(envFilePath, cfg)
		if err != nil {
			return fmt.Errorf("config error: %w", err)
		}
	} else {
		err := cleanenv.ReadEnv(cfg)
		if err != nil {

			if _, statErr := os.Stat(envFilePath + ".example"); statErr == nil {
				return fmt.Errorf("missing environmentvariables: %w\n\nprovide all required environment variables or rename and update .env.example to .env for convinience", err)
			}

			return err
		}
	}
	return nil
}

func checkFileExists(fileName string) bool {
	envFileExists := false
	if _, err := os.Stat(fileName); err == nil {
		envFileExists = true
	}
	return envFileExists
}

func projectRoot() string {
	_, b, _, _ := runtime.Caller(0)
	projectRoot := filepath.Dir(b)

	return projectRoot + "/../"
}
