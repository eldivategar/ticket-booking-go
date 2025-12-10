package configs

import "github.com/spf13/viper"

type Config struct {
	// Server configurations
	ServerHost    string `mapstructure:"SERVER_HOST"`
	ServerPort    int    `mapstructure:"SERVER_PORT"`
	ServerMode    string `mapstructure:"SERVER_MODE"`
	ServerVersion string `mapstructure:"SERVER_VERSION"`

	// Database configurations
	DBHost     string `mapstructure:"DB_HOST"`
	DBPort     int    `mapstructure:"DB_PORT"`
	DBUser     string `mapstructure:"DB_USER"`
	DBPassword string `mapstructure:"DB_PASSWORD"`
	DBName     string `mapstructure:"DB_NAME"`
	DBSSLMode  string `mapstructure:"DB_SSL_MODE"`

	// Redis configurations
	RedisHost     string `mapstructure:"REDIS_HOST"`
	RedisPort     int    `mapstructure:"REDIS_PORT"`
	RedisPassword string `mapstructure:"REDIS_PASSWORD"`
	RedisDB       int    `mapstructure:"REDIS_DB"`

	// S3 configurations (minio for default)
	MinioEndpoint       string `mapstructure:"MINIO_ENDPOINT"`
	MinioPublicEndpoint string `mapstructure:"MINIO_PUBLIC_ENDPOINT"`
	MinioAccessKey      string `mapstructure:"MINIO_ACCESS_KEY"`
	MinioSecretKey      string `mapstructure:"MINIO_SECRET_KEY"`
	MinioRegion         string `mapstructure:"MINIO_REGION"`
	MinioBucket         string `mapstructure:"MINIO_BUCKET"`
	MinioUseSSL         bool   `mapstructure:"MINIO_USE_SSL"`

	// JWT configurations
	JWTAccessSecret       string `mapstructure:"JWT_ACCESS_SECRET"`
	JWTRefreshSecret      string `mapstructure:"JWT_REFRESH_SECRET"`
	JWTAccessTokenExpire  string `mapstructure:"JWT_ACCESS_TOKEN_EXPIRE"`
	JWTRefreshTokenExpire string `mapstructure:"JWT_REFRESH_TOKEN_EXPIRE"`

	// Message Broker configurations
	RabbitMQURL string `mapstructure:"RABBITMQ_URL"`
}

func LoadConfig() (cfg Config, err error) {
	viper.AddConfigPath(".")
	viper.SetConfigFile(".env")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&cfg)
	return
}
