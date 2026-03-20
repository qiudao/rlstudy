package client

import (
	"net/http/httptest"
	"testing"

	"github.com/qiudao/rlstudy/pkg/env"
)

func TestClient_InfoAndStep(t *testing.T) {
	s := env.NewServer(10, 42)
	ts := httptest.NewServer(s)
	defer ts.Close()

	c := New(ts.URL)

	info, err := c.Info()
	if err != nil {
		t.Fatal(err)
	}
	if info.Arms != 10 {
		t.Errorf("expected 10 arms, got %d", info.Arms)
	}

	err = c.Reset()
	if err != nil {
		t.Fatal(err)
	}

	resp, err := c.Step(0)
	if err != nil {
		t.Fatal(err)
	}
	if resp.Reward != resp.Reward {
		t.Error("reward is NaN")
	}
}
