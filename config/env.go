package config

import (
	"log"
	"time"

	"github.com/spf13/viper"
)

type Env struct {
	AppEnv                 string        `mapstructure:"APP_ENV"`
	ServerAddress          string        `mapstructure:"SERVER_ADDRESS"`
	ContextTimeout         int           `mapstructure:"CONTEXT_TIMEOUT"`
	MongoAtlas             bool          `mapstructure:"MONGO_ATLAS"`
	DBUrl                  string        `mapstructure:"DB_URL"`
	DBHost                 string        `mapstructure:"DB_HOST"`
	DBPort                 string        `mapstructure:"DB_PORT"`
	DBUser                 string        `mapstructure:"MONGO_INITDB_ROOT_USERNAME"`
	DBPass                 string        `mapstructure:"MONGO_INITDB_ROOT_PASSWORD"`
	DBName                 string        `mapstructure:"MONGO_INITDB_DATABASE"`
	AccessTokenExpiryHour  int           `mapstructure:"ACCESS_TOKEN_EXPIRY_HOUR"`
	RefreshTokenExpiryHour int           `mapstructure:"REFRESH_TOKEN_EXPIRY_HOUR"`
	AccessTokenSecret      string        `mapstructure:"ACCESS_TOKEN_SECRET"`
	RefreshTokenSecret     string        `mapstructure:"REFRESH_TOKEN_SECRET"`
	ClientID               string        `mapstructure:"CLIENT_ID"`
	ClientSecret           string        `mapstructure:"CLIENT_SECRET"`
	RedirectURL            string        `mapstructure:"REDIRECT_URL"`
	IssuerURL              string        `mapstructure:"ISSUER_URL"`
	RedisHost              string        `mapstructure:"REDIS_HOST"`
	RedisPort              int           `mapstructure:"REDIS_PORT"`
	RedisPassword          string        `mapstructure:"REDIS_PASSWORD"`
	RedisContextTime       time.Duration `mapstructure:"REDIS_TIMEOUT"`
}

func NewEnv() *Env {
	env := Env{}
	viper.SetConfigFile(".env")

	err := viper.ReadInConfig()
	if err != nil {
		log.Fatal("Can't find the file .env : ", err)
	}

	err = viper.Unmarshal(&env)
	if err != nil {
		log.Fatal("Environment can't be loaded: ", err)
	}

	if env.AppEnv == "development" {
		log.Println("The App is running in development env")
	}

	return &env
}
