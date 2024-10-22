package main

import (
	"github.com/Ruclo/apitesto/internal/config"
	"fmt"
)

func main() {
	config, err := config.LoadConfigFromYAML("./input.yaml")
	if err != nil {
		return
	}
	fmt.Println(config)
}