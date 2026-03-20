package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/qiudao/rlstudy"
	"github.com/qiudao/rlstudy/pkg/client"
	"github.com/qiudao/rlstudy/pkg/runner"

)

func main() {
	cfgPath := flag.String("config", "config.json", "config file path")
	envURL := flag.String("env", "", "environment URL (overrides config port)")
	epsilonsStr := flag.String("epsilons", "0,0.01,0.1", "comma-separated epsilon values")
	runs := flag.Int("runs", 2000, "number of runs")
	steps := flag.Int("steps", 1000, "steps per run")
	seed := flag.Int64("seed", 42, "random seed")
	csvPath := flag.String("csv", "", "output CSV file path")
	flag.Parse()

	cfg, err := rlstudy.LoadConfig(*cfgPath)
	if err != nil {
		log.Printf("using default config: %v", err)
		cfg = rlstudy.DefaultConfig()
	}

	epsilons := parseFloats(*epsilonsStr)

	var c *client.Client
	if *envURL != "" {
		c = client.New(*envURL)
	} else {
		c = client.NewUnix(cfg.Socket)
	}

	info, err := c.Info()
	if err != nil {
		log.Fatalf("cannot connect to env: %v", err)
	}
	log.Printf("connected to env: %d arms", info.Arms)

	var results []runner.Result
	for _, eps := range epsilons {
		log.Printf("running epsilon=%.2f (%d runs x %d steps)...", eps, *runs, *steps)
		r, err := runner.Run(c, info.Arms, eps, *runs, *steps, *seed)
		if err != nil {
			log.Fatalf("epsilon=%.2f failed: %v", eps, err)
		}
		results = append(results, r)
	}

	runner.PrintSummary(os.Stdout, results, *steps)

	if *csvPath != "" {
		if err := writeCSV(*csvPath, results, *steps); err != nil {
			log.Fatalf("write csv: %v", err)
		}
		log.Printf("CSV written to %s", *csvPath)
	}
}

func parseFloats(s string) []float64 {
	parts := strings.Split(s, ",")
	fs := make([]float64, len(parts))
	for i, p := range parts {
		f, err := strconv.ParseFloat(strings.TrimSpace(p), 64)
		if err != nil {
			log.Fatalf("invalid epsilon %q: %v", p, err)
		}
		fs[i] = f
	}
	return fs
}

func writeCSV(path string, results []runner.Result, steps int) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	w := csv.NewWriter(f)
	defer w.Flush()

	w.Write([]string{"step", "epsilon", "avg_reward", "optimal_action_pct"})
	for _, r := range results {
		for s := 0; s < steps; s++ {
			w.Write([]string{
				strconv.Itoa(s + 1),
				fmt.Sprintf("%.2f", r.Epsilon),
				fmt.Sprintf("%.4f", r.AvgRewards[s]),
				fmt.Sprintf("%.2f", r.OptimalActionPct[s]),
			})
		}
	}
	return nil
}
