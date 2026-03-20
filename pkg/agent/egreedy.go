package agent

import "math/rand"

// EpsilonGreedy implements sample-average epsilon-greedy action selection.
type EpsilonGreedy struct {
	k       int
	epsilon float64
	q       []float64 // estimated values
	n       []int     // action counts
	rng     *rand.Rand
}

func NewEpsilonGreedy(k int, epsilon float64, seed int64) *EpsilonGreedy {
	return &EpsilonGreedy{
		k:       k,
		epsilon: epsilon,
		q:       make([]float64, k),
		n:       make([]int, k),
		rng:     rand.New(rand.NewSource(seed)),
	}
}

func (e *EpsilonGreedy) SelectAction() int {
	if e.rng.Float64() < e.epsilon {
		return e.rng.Intn(e.k)
	}
	// Greedy: pick action with highest estimate (break ties randomly)
	best := 0
	for i := 1; i < e.k; i++ {
		if e.q[i] > e.q[best] {
			best = i
		}
	}
	return best
}

func (e *EpsilonGreedy) Update(action int, reward float64) {
	e.n[action]++
	e.q[action] += (reward - e.q[action]) / float64(e.n[action])
}

func (e *EpsilonGreedy) Reset() {
	for i := range e.q {
		e.q[i] = 0
		e.n[i] = 0
	}
}
