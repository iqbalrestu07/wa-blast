package configs

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/viper"
)

func Init(fileName string) {
	splits := strings.Split(filepath.Base(fileName), ".")
	viper.SetConfigType("yaml")
	viper.SetConfigName(filepath.Base(splits[0]))
	viper.AddConfigPath(filepath.Dir(fileName))
	err := viper.ReadInConfig()
	// If an error occurred while reading config file, exit app
	if err != nil {
		fmt.Printf("Unable to load configuration file. Error: %s\n", err.Error())
		os.Exit(2)
	}
}

// MustGetString retrieve configuration by key. If an error occurred, it will panic and exit application
func MustGetString(key string) string {
	validateKey(key)
	return viper.GetString(key)
}

func MustGetInt(key string) int {
	validateKey(key)
	return viper.GetInt(key)
}

func MustGetFloat32(key string) float32 {
	validateKey(key)
	return float32(viper.GetFloat64(key))
}

// validateKey exit app if a configuration is not set
func validateKey(key string) {
	if !viper.IsSet(key) {
		fmt.Printf("Configuration %s is not set\n", key)
		os.Exit(3)
	}
}
