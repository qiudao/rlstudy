package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"strconv"
	"syscall"

	"github.com/qiudao/rlstudy"
	"github.com/qiudao/rlstudy/pkg/env"
)

func main() {
	cfgPath := flag.String("config", "config.json", "config file path")
	seed := flag.Int64("seed", 42, "random seed")
	daemon := flag.Bool("daemon", false, "run as daemon (background mode, write PID file)")
	pidFile := flag.String("pidfile", ".bandit-env.pid", "PID file path (daemon mode)")
	flag.Parse()

	cfg, err := rlstudy.LoadConfig(*cfgPath)
	if err != nil {
		log.Printf("using default config: %v", err)
		cfg = rlstudy.DefaultConfig()
	}

	s := env.NewServer(cfg.Arms, *seed)
	addr := fmt.Sprintf(":%d", cfg.Port)

	ln, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("listen %s: %v", addr, err)
	}

	if *daemon {
		absPid, _ := filepath.Abs(*pidFile)
		os.WriteFile(absPid, []byte(strconv.Itoa(os.Getpid())), 0644)
		log.Printf("bandit-env daemon started on %s (pid=%d, pidfile=%s)", addr, os.Getpid(), absPid)

		srv := &http.Server{Handler: s}
		go func() {
			sigCh := make(chan os.Signal, 1)
			signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
			<-sigCh
			log.Println("shutting down...")
			srv.Shutdown(context.Background())
			os.Remove(absPid)
		}()
		log.Fatal(srv.Serve(ln))
	} else {
		log.Printf("bandit-env starting on %s (arms=%d, seed=%d)", addr, cfg.Arms, *seed)
		log.Fatal(http.Serve(ln, s))
	}
}
