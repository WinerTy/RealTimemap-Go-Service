package main

import (
	"fmt"
	"realtimemap-service/internal/config"
)

func main() {
	cfg := config.MustLoad()

	fmt.Print(cfg)
}
