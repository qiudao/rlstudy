package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/qiudao/rlstudy"
	"github.com/qiudao/rlstudy/pkg/env"
)

func main() {
	cfgPath := flag.String("config", "config.json", "config file path")
	seed := flag.Int64("seed", 42, "random seed")
	flag.Parse()

	cfg, err := rlstudy.LoadConfig(*cfgPath)
	if err != nil {
		log.Printf("using default config: %v", err)
		cfg = rlstudy.DefaultConfig()
	}

	s := env.NewServer(cfg.Arms, *seed)
	addr := fmt.Sprintf(":%d", cfg.Port)
	log.Printf("bandit-env starting on %s (arms=%d, seed=%d)", addr, cfg.Arms, *seed)
	log.Fatal(http.ListenAndServe(addr, s))
}
