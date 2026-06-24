package main

import (
	"fmt"

	"github.com/invaderjt/blog-aggregator/internal/config"
)

func main() {
	cfg, err := config.ReadConfig()
	if err != nil {
		fmt.Println("Could not read config file")
	}

	cfg.SetUser("Jack")

	cfg, err = config.ReadConfig()
	if err != nil {
		fmt.Println("Could not read config file")
	}

	fmt.Println(cfg)

}
