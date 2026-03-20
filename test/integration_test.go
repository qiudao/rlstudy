package test

import (
	"net/http/httptest"
	"testing"

	"github.com/qiudao/rlstudy/pkg/client"
	"github.com/qiudao/rlstudy/pkg/env"
	"github.com/qiudao/rlstudy/pkg/runner"
)

func TestIntegration_SmallExperiment(t *testing.T) {
	s := env.NewServer(10, 42)
	ts := httptest.NewServer(s)
	defer ts.Close()

	c := client.New(ts.URL)
	result, err := runner.Run(c, 10, 0.1, 10, 100, 42)
	if err != nil {
		t.Fatal(err)
	}

	// After 100 steps with epsilon=0.1, avg reward should be positive
	last := result.AvgRewards[99]
	if last < 0 {
		t.Errorf("expected positive avg reward at step 100, got %.3f", last)
	}

	// Optimal action % should be above chance (10%)
	lastOpt := result.OptimalActionPct[99]
	if lastOpt < 10 {
		t.Errorf("expected optimal%% > 10%%, got %.1f%%", lastOpt)
	}
}
