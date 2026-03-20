package agent

import "testing"

func TestEpsilonGreedy_Greedy(t *testing.T) {
	a := NewEpsilonGreedy(3, 0, 42)
	a.Update(0, 1.0)
	a.Update(1, 2.0)
	a.Update(2, 3.0)
	for i := 0; i < 100; i++ {
		if a.SelectAction() != 2 {
			t.Fatal("greedy agent should always pick action 2")
		}
	}
}

func TestEpsilonGreedy_Explores(t *testing.T) {
	a := NewEpsilonGreedy(10, 1.0, 42)
	counts := make([]int, 10)
	for i := 0; i < 10000; i++ {
		counts[a.SelectAction()]++
	}
	for i, c := range counts {
		if c < 500 || c > 1500 {
			t.Errorf("action %d picked %d times, expected ~1000", i, c)
		}
	}
}

func TestEpsilonGreedy_Reset(t *testing.T) {
	a := NewEpsilonGreedy(3, 0, 42)
	a.Update(0, 5.0)
	a.Reset()
	_ = a.SelectAction()
}
