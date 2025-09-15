package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

type Config struct {
	BucketName      string `mapstructure:"bucket_name"`
	AWSRegion       string `mapstructure:"aws_region"`
	AWSAccessKeyID  string `mapstructure:"aws_access_key_id"`
	AWSSecretKey    string `mapstructure:"aws_secret_access_key"`
}

func Load() (*Config, error) {
	// Set default values
	viper.SetDefault("aws_region", "us-east-1")

	// Config file name (without extension)
	viper.SetConfigName("tincan")
	viper.SetConfigType("yaml")

	// Look for config in home directory and current directory
	if home, err := os.UserHomeDir(); err == nil {
		viper.AddConfigPath(home)
		viper.AddConfigPath(filepath.Join(home, ".config"))
	}
	viper.AddConfigPath(".")

	// Environment variable prefix
	viper.SetEnvPrefix("TINCAN")
	viper.AutomaticEnv()

	// Read config file (optional)
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("error reading config file: %w", err)
		}
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("unable to decode config: %w", err)
	}

	// Validate required fields
	if config.BucketName == "" {
		return nil, fmt.Errorf("bucket_name is required (set TINCAN_BUCKET_NAME environment variable or add to config file)")
	}

	return &config, nil
}