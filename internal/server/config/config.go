package config

import (
	"fmt"
	"net"
	"os"
	"time"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

type Config struct {
	LogLevel    string
	DSN         string
	HashCost    int
	TokenSecret []byte
	TokenTTL    time.Duration
	TCPAddr     *net.TCPAddr
}

func MustLoad() *Config {
	initConfigPath()

	viper.SetConfigType("yaml")
	err := viper.ReadInConfig()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	c := &Config{
		LogLevel:    viper.GetString("LogLevel"),
		DSN:         viper.GetString("DSN"),
		HashCost:    viper.GetInt("HashCost"),
		TokenSecret: []byte(viper.GetString("TokenSecret")),
		TokenTTL:    time.Duration(viper.GetInt("TokenTTL")) * time.Hour,
		TCPAddr:     mustResolveTCPAddr(viper.GetString("TCPAddr")),
	}

	return c
}

func initConfigPath() {
	const (
		configEnv  = "GOPHKEEPER_CONFIG"
		configFlag = "config"
	)

	viper.BindEnv(configEnv)

	pflag.StringP(configFlag,
		"c", "server.config.yaml", "path to config file")
	pflag.Parse()
	viper.BindPFlags(pflag.CommandLine)

	if path := viper.GetString(configEnv); path != "" {
		viper.SetConfigFile(path)
		return
	}

	viper.SetConfigFile(viper.GetString(configFlag))
}

func mustResolveTCPAddr(addr string) *net.TCPAddr {
	a, err := net.ResolveTCPAddr("tcp", addr)
	if err != nil {
		fmt.Printf("incorrect 'TCPAddr' config: %q\n", err.Error())
		os.Exit(1)
	}
	return a
}
