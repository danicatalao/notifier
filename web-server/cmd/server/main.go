package main

import (
	"fmt"
	"log"

	config "github.com/danicatalao/notifier/web-server/configs"
)

func main() {
	// Loading .env variables into config
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatalf("Config error: %s", err)
	}
	fmt.Printf("%+v\n", cfg)

	fmt.Println("Hello, World!")
}
