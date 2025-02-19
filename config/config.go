package config

import (
	"fmt"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	AppName         string
	AppEnv          string
	AppHost         string
	AppPort         string
	DBHost          string
	DBUser          string
	DBPass          string
	DBName          string
	DBPort          string
	LogLevel        string
	LogFormat       string
	SecretKey       string
	DBSslMode       string
	ShutdownTimeout time.Duration
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime time.Duration
}

func LoadConfig(configPath string) (*Config, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(configPath)

	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("error reading config file, %s", err)
	}

	viper.AutomaticEnv()
	env := viper.GetString("ENV")

	if env == "" {
		env = "local"
	}

	viper.SetConfigName("config." + env)
	if err := viper.MergeInConfig(); err != nil {
		return nil, fmt.Errorf("config file for env: %s not found: %e", env, err)
	}

	AppName := viper.GetString("app.name")
	AppHost := viper.GetString("app.host")
	AppPort := viper.GetString("app.port")
	DBHost := viper.GetString("db.host")
	DBUser := viper.GetString("db.user")
	DBPass := viper.GetString("db.password")
	DBName := viper.GetString("db.name")
	DBPort := viper.GetString("db.port")
	DBSslMode := viper.GetString("db.ssl-mode")
	LogLevel := viper.GetString("log.level")
	LogFormat := viper.GetString("log.format")
	SecretKey := viper.GetString("auth.secret-key")
	ShutdownTimeoutStr := viper.GetString("app.shutdown-timeout")
	MaxOpenConns := viper.GetInt("db.max-open-conns")
	MaxIdleConns := viper.GetInt("db.max-idle-conns")
	ConnMaxLifetimeStr := viper.GetString("db.conn-max-lifetime")

	if ShutdownTimeoutStr == "" {
		ShutdownTimeoutStr = "10"
	}

	if AppName == "" || AppHost == "" || AppPort == "" {
		return nil, fmt.Errorf("one or more app configuration fields are empty")
	}

	if SecretKey == "" {
		return nil, fmt.Errorf("jwt secret key is empty")
	}

	ShutdownTimeout, err := time.ParseDuration(ShutdownTimeoutStr + "s")
	if err != nil {
		return nil, fmt.Errorf("error parsing shutdown timeout: %w", err)
	}

	ConnMaxLifetime, err := time.ParseDuration(ConnMaxLifetimeStr + "m")
	if err != nil {
		return nil, fmt.Errorf("error parsing shutdown timeout: %w", err)
	}

	if DBHost == "" || DBUser == "" || DBPass == "" || DBName == "" || DBPort == "" || DBSslMode == "" {
		return nil, fmt.Errorf("one or more database configuration fields are empty")
	}

	if LogLevel == "" || LogFormat == "" {
		return nil, fmt.Errorf("one or more logger configuration fields are empty")
	}

	config := Config{
		AppName:         AppName,
		AppHost:         AppHost,
		AppPort:         AppPort,
		AppEnv:          env,
		DBHost:          DBHost,
		DBUser:          DBUser,
		DBPass:          DBPass,
		DBName:          DBName,
		DBPort:          DBPort,
		DBSslMode:       DBSslMode,
		LogLevel:        LogLevel,
		LogFormat:       LogFormat,
		SecretKey:       SecretKey,
		ShutdownTimeout: ShutdownTimeout,
		MaxOpenConns:    MaxOpenConns,
		MaxIdleConns:    MaxIdleConns,
		ConnMaxLifetime: ConnMaxLifetime,
	}

	return &config, nil
}
