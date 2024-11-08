package config

import "github.com/spf13/viper"

// Configuration is the application configuration
type Configuration struct {
	// Port is the port on which we will listen for requests
	Port int `yaml:"port"`

	// ConnectionString is the database connection string
	ConnectionString string `yaml:"connectionString"`
}

var configuration *Configuration

// Get returns the application configuration
func Get() (*Configuration, error) {
	if configuration != nil {
		return configuration, nil
	}

	viper.SetConfigFile("/etc/call-manager/configuration.yml")

	err := viper.ReadInConfig()
	if err != nil {
		return nil, err
	}

	var config Configuration

	err = viper.Unmarshal(&config)
	if err != nil {
		return nil, err
	}

	configuration = &config

	return configuration, nil
}
