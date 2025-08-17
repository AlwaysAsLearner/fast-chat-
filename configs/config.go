package config

import (
	"fmt"

	"github.com/AlwaysAsLearner/fast-chat/backend/internal/utils"
	"github.com/spf13/viper"
)

type JWTConfig struct {
	Secret     string `mapstructure:"secret"`
	Expiration int    `mapstructure:"expiration_minutes"`
}

type LoggerConfig struct {
	Level  string `mapstructure:"level"`
	Format string `mapstructure:"format"`
	Output string `mapstructure:"output"`
}

type WebSocketConfig struct {
	ReadTimeout    int `mapstructure:"read_timeout"`
	WriteTimeout   int `mapstructure:"write_timeout"`
	MaxConnections int `mapstructure:"max_connections"`
	BufferSize     int `mapstructure:"buffer_size"`
}

/*
database:
  driver: "postgres"
  host: "localhost"
  port: 5432
  user: "chat_user"
  password: "secure_pass"
  dbname: "chat_db"

*/

type DBConfig struct {
	Driver   string `mapstructure:"driver"`
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	Dbname   string `mapstructure:"dbname"`
	SSLMode  string `mapstructure:"ssl_mode"`
	TimeZone string `mapstructure:"time_zone"`
}

// Our config accesors availabe for other packages and main program
var JWT JWTConfig
var Logger LoggerConfig
var WebSocket WebSocketConfig
var DB DBConfig

func Load() {
	viper.SetConfigType("yaml")

	// Load jwt config
	viper.SetConfigName("jwt")     // targets jwt.yaml file
	viper.AddConfigPath("configs") // the directory where to search this file
	viper.ReadInConfig()
	viper.UnmarshalKey("jwt", &JWT)

	viper.SetConfigName("logger")
	viper.AddConfigPath("configs")
	viper.ReadInConfig()
	viper.UnmarshalKey("logger", &Logger)

	viper.SetConfigName("websocket")
	viper.AddConfigPath("configs")
	viper.ReadInConfig()
	viper.UnmarshalKey("websocket", &WebSocket)

	utils.ReadConfig("database", &DB)

}

func (db *DBConfig) ToDsn() string {
	return fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%d sslmode=%s TimeZone=%s",
		db.Host, db.User, db.Password, db.Dbname, db.Port, db.SSLMode, db.TimeZone,
	)
}
