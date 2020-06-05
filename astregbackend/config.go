package main

import (
	"gopkg.in/ini.v1"
	"log"
	"os"
)

var (
	Config Configuration
)

type Configuration struct {
	ListenAddress   string
	Digits          int
	RedisDB         int
	RedisExpiration int
	RedisHost       string
	RedisPW         string
	Verbose         bool
}

// LoadConfig loads the configuration from given file path
func LoadConfig(filename string) {
	ini, err := ini.Load(filename)
	if err != nil {
		log.Printf("Fail to read file: %v", err)
		os.Exit(1)
	}

	Config.ListenAddress = ini.Section("").Key("ListenAddress").String()
	Config.Digits = ini.Section("").Key("Digits").MustInt()
	Config.RedisDB = ini.Section("").Key("RedisDB").MustInt()
	Config.RedisExpiration = ini.Section("").Key("RedisExpiration").MustInt()
	Config.RedisHost = ini.Section("").Key("RedisHost").String()
	Config.RedisPW = ini.Section("").Key("RedisPW").String()
	Config.Verbose = ini.Section("").Key("Verbose").MustBool()
}
