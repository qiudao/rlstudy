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
	useTCP := flag.Bool("tcp", false, "use TCP instead of Unix socket")
	flag.Parse()

	cfg, err := rlstudy.LoadConfig(*cfgPath)
	if err != nil {
		log.Printf("using default config: %v", err)
		cfg = rlstudy.DefaultConfig()
	}

	s := env.NewServer(cfg.Arms, *seed)

	var ln net.Listener
	var listenAddr string

	if *useTCP {
		listenAddr = fmt.Sprintf(":%d", cfg.Port)
		ln, err = net.Listen("tcp", listenAddr)
	} else {
		listenAddr = cfg.Socket
		os.Remove(listenAddr) // clean up stale socket
		ln, err = net.Listen("unix", listenAddr)
	}
	if err != nil {
		log.Fatalf("listen %s: %v", listenAddr, err)
	}

	cleanup := func() {
		if !*useTCP {
			os.Remove(listenAddr)
		}
	}

	if *daemon {
		absPid, _ := filepath.Abs(*pidFile)
		os.WriteFile(absPid, []byte(strconv.Itoa(os.Getpid())), 0644)
		log.Printf("bandit-env daemon started on %s (pid=%d, pidfile=%s)", listenAddr, os.Getpid(), absPid)

		srv := &http.Server{Handler: s}
		go func() {
			sigCh := make(chan os.Signal, 1)
			signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
			<-sigCh
			log.Println("shutting down...")
			srv.Shutdown(context.Background())
			os.Remove(absPid)
			cleanup()
		}()
		log.Fatal(srv.Serve(ln))
	} else {
		log.Printf("bandit-env starting on %s (arms=%d, seed=%d)", listenAddr, cfg.Arms, *seed)
		defer cleanup()
		log.Fatal(http.Serve(ln, s))
	}
}
