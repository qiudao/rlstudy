package env

import "math/rand"

// Bandit is a k-armed bandit environment.
// Each arm has a true value q*(a) ~ N(0,1).
// Pulling arm a returns reward ~ N(q*(a), 1).
type Bandit struct {
	k     int
	qStar []float64
	rng   *rand.Rand
}

func NewBandit(k int, seed int64) *Bandit {
	b := &Bandit{k: k, rng: rand.New(rand.NewSource(seed))}
	b.Reset()
	return b
}

func (b *Bandit) Reset() {
	b.qStar = make([]float64, b.k)
	for i := range b.qStar {
		b.qStar[i] = b.rng.NormFloat64()
	}
}

func (b *Bandit) Arms() int { return b.k }

func (b *Bandit) Step(action int) float64 {
	return b.qStar[action] + b.rng.NormFloat64()
}

func (b *Bandit) OptimalAction() int {
	best := 0
	for i := 1; i < b.k; i++ {
		if b.qStar[i] > b.qStar[best] {
			best = i
		}
	}
	return best
}

func (b *Bandit) QStar(action int) float64 {
	return b.qStar[action]
}
