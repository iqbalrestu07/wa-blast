package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"time"

	"wa-blast/configs"
	"wa-blast/crypto"
	"wa-blast/flags"
	"wa-blast/loggers"
	"wa-blast/repositories"
	"wa-blast/router"
	"wa-blast/services"
	"wa-blast/util"

	"github.com/rs/cors"
)

var loc *time.Location

var log = loggers.Get()

func main() {
	os.Setenv("TZ", "Asia/Jakarta")

	showVersion := flag.Bool("version", false, "Version")

	// Get starting time
	start := time.Now()

	flag.Parse()

	if *showVersion {
		printVersion()
		os.Exit(0)
	}

	configFile := flag.String("c", "config.yml", "config.yml")
	if *configFile == "" {
		fmt.Println("Configuration file is not specified. Use -c config.yml to specify one")
		os.Exit(1)
	}

	configs.Init(*configFile)

	r := router.New(start)
	// Start server

	repositories.Init()

	util.Init()

	services.Init()

	// Initiate crypto
	crypto.Init()

	// Where ORIGIN_ALLOWED is like `scheme://dns[:port]`, or `*` (insecure)
	c := cors.New(cors.Options{
		AllowCredentials: true,
		AllowedHeaders:   []string{"*"},
		AllowedOrigins:   []string{"http://infobankdata.octarine.id", "http://localhost:8080", "http://127.0.0.1:8080"}, // All origins
		AllowedMethods:   []string{"GET", "HEAD", "POST", "PUT", "DELETE"},                                              // Allowing only get, just an example
	})

	log.Infof("Boot time: %s", time.Since(start))
	log.Fatal(http.ListenAndServe(getPort(), c.Handler(r)))
}

// getPort retrieve port from config
func getPort() string {
	port := configs.MustGetString("server.port")
	return ":" + port
}

func printVersion() {
	fmt.Printf("%s version %s, build %s\n", flags.AppName, flags.AppVersion, flags.AppCommitHash)
}

func setTimezone(tz string) error {
	location, err := time.LoadLocation("Asia/Jakarta")
	if err != nil {
		return err
	}
	loc = location
	return nil
}
