package env

import (
	"math"
	"testing"
)

func TestBandit_Arms(t *testing.T) {
	b := NewBandit(10, 42)
	if b.Arms() != 10 {
		t.Errorf("expected 10 arms, got %d", b.Arms())
	}
}

func TestBandit_Reset(t *testing.T) {
	b := NewBandit(10, 42)
	opt1 := b.OptimalAction()
	b.Reset()
	opt2 := b.OptimalAction()
	if opt2 < 0 || opt2 >= 10 {
		t.Errorf("optimal action out of range: %d", opt2)
	}
	_ = opt1
}

func TestBandit_Step_ReturnsReward(t *testing.T) {
	b := NewBandit(10, 42)
	n := 10000
	sum := 0.0
	for i := 0; i < n; i++ {
		sum += b.Step(0)
	}
	mean := sum / float64(n)
	if math.Abs(mean-b.QStar(0)) > 0.1 {
		t.Errorf("mean reward %.3f too far from q*(0)=%.3f", mean, b.QStar(0))
	}
}

func TestBandit_OptimalAction(t *testing.T) {
	b := NewBandit(10, 42)
	opt := b.OptimalAction()
	for a := 0; a < 10; a++ {
		if b.QStar(a) > b.QStar(opt) {
			t.Errorf("action %d has higher q* than optimal %d", a, opt)
		}
	}
}
