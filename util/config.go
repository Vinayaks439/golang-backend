package util

import "github.com/spf13/viper"

// Config This config struct will be used to pass configuration to the application.
// The values are read by viper from a config file or environment variables.
type Config struct {
	DbSource      string `mapstructure:"DB_SOURCE"`
	ServerAddress string `mapstructure:"HTTP_SERVER_ADDRESS"`
}

// LoadConfig loads the config from the given path or env variables.
func LoadConfig(path string) (Config, error) {
	var config Config
	viper.AddConfigPath(path)
	viper.SetConfigName("app")
	viper.SetConfigType("env") // json, toml, yaml, yml, properties, props, prop, env, dotenv
	viper.AutomaticEnv()
	err := viper.ReadInConfig()
	if err != nil {
		return config, err
	}
	err = viper.Unmarshal(&config)
	return config, err
}
