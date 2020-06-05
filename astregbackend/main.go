package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/go-redis/redis/v8"
	"log"
	"net/http"
	"os"
	"runtime"
)

const (
	productName = "astregbackend"
	logFlags    = log.Lshortfile | log.Ldate | log.Ltime
)

var (
	showVersion = flag.Bool("v", false, "Show version information")
	configfile  = flag.String("c", "/etc/astregbackend.conf", "Path to configfile")

	usage = fmt.Sprintf("Usage: %s -v | -h | OPTIONS", os.Args[0])

	// Version can be set at build time using:
	//    -ldflags "-X main.version=0.42"
	version, versionInfo string

	ctx = context.Background()
)

func init() {
	if version == "" {
		version = "snapshot"
	}
	versionInfo = fmt.Sprintf("%s %s (%s)", productName, version, runtime.Version())

	flag.Usage = func() {
		fmt.Println(versionInfo)
		fmt.Println(usage)
		flag.PrintDefaults()
	}

	log.SetFlags(logFlags)
}

func main() {
	flag.Parse()

	if *showVersion {
		fmt.Println(version)
		os.Exit(0)
	}

	log.Println("Going to read settings from ", *configfile)
	LoadConfig(*configfile)

	log.Println("Connecting to Redis on ", Config.RedisHost)
	rediscon := redis.NewClient(&redis.Options{
		Addr:     Config.RedisHost,
		Password: Config.RedisPW,
		DB:       Config.RedisDB,
	})

	NewRealtimeHandler("/ps_aors/", NewDummyHandler(rediscon, Config.Digits, "aors"))
	NewRealtimeHandler("/ps_auths/", NewDummyHandler(rediscon, Config.Digits, "auth"))
	NewRealtimeHandler("/ps_endpoints/", NewDummyHandler(rediscon, Config.Digits, "endpoint"))

	log.Println("Starting HTTP server on", Config.ListenAddress)
	log.Fatal(http.ListenAndServe(Config.ListenAddress, nil))
}
