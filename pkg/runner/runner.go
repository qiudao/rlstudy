package runner

import (
	"fmt"
	"io"

	"github.com/qiudao/rlstudy/pkg/agent"
	"github.com/qiudao/rlstudy/pkg/client"
)

type Result struct {
	Epsilon          float64
	AvgRewards       []float64
	OptimalActionPct []float64
}

func Run(c *client.Client, arms int, epsilon float64, runs int, steps int, seed int64) (Result, error) {
	res := Result{
		Epsilon:          epsilon,
		AvgRewards:       make([]float64, steps),
		OptimalActionPct: make([]float64, steps),
	}

	for r := 0; r < runs; r++ {
		if err := c.Reset(); err != nil {
			return res, fmt.Errorf("run %d reset: %w", r, err)
		}
		ag := agent.NewEpsilonGreedy(arms, epsilon, seed+int64(r))
		ag.Reset()

		for s := 0; s < steps; s++ {
			action := ag.SelectAction()
			resp, err := c.Step(action)
			if err != nil {
				return res, fmt.Errorf("run %d step %d: %w", r, s, err)
			}
			ag.Update(action, resp.Reward)
			res.AvgRewards[s] += resp.Reward
			if resp.Optimal {
				res.OptimalActionPct[s]++
			}
		}
	}

	for s := 0; s < steps; s++ {
		res.AvgRewards[s] /= float64(runs)
		res.OptimalActionPct[s] = res.OptimalActionPct[s] / float64(runs) * 100
	}

	return res, nil
}

func PrintSummary(w io.Writer, results []Result, steps int) {
	fmt.Fprintln(w, "=== Experiment Summary ===")
	for _, r := range results {
		start := steps - steps/10
		avgR := 0.0
		avgO := 0.0
		for s := start; s < steps; s++ {
			avgR += r.AvgRewards[s]
			avgO += r.OptimalActionPct[s]
		}
		n := float64(steps - start)
		fmt.Fprintf(w, "epsilon=%.2f  avg_reward=%.3f  optimal_action=%.1f%%\n",
			r.Epsilon, avgR/n, avgO/n)
	}
}
