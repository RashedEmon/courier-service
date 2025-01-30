package config

import (
	"fmt"
	"log"
	"path/filepath"
	"runtime"
	"sync"

	"github.com/spf13/viper"
)

type Config struct {
	Env   string `mapstructure:"ENV"`
	Debug bool   `mapstructure:"DEBUG"`
	ProjectRoot string `mapstructure:"PROJECTROOT"`
	DB    struct {
		Host     string `mapstructure:"HOST"`
		Port     int    `mapstructure:"PORT"`
		Name     string `mapstructure:"NAME"`
		Username string `mapstructure:"USERNAME"`
		Password string `mapstructure:"PASSWORD"`
	} `mapstructure:"DB"`
}

var (
	ConfigInstance *Config
	once           sync.Once
)

// It will return project root
func getProjectRoot() string {
	_, filename, _, _ := runtime.Caller(0)
	configDir := filepath.Dir(filename)
	projectRoot := filepath.Dir(configDir)
	return projectRoot
}

// Load config with singleton behave
func LoadConfig() *Config {
	projectRoot := getProjectRoot()

	once.Do(func() {
		viper.SetConfigName("config")
		viper.SetConfigType("yaml")
		viper.AddConfigPath(projectRoot)

		if err := viper.ReadInConfig(); err != nil {
			log.Fatalf("Error reading config file: %v", err)
		}

		ConfigInstance = &Config{ProjectRoot: projectRoot}
		if err := viper.Unmarshal(ConfigInstance); err != nil {
			log.Fatalf("Unable to decode into struct: %v", err)
		}
		fmt.Println("Configuration loaded successfully!")
	})

	return ConfigInstance
}
