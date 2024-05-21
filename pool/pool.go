package pool

import "sync"

var (
	inst *PoolRunner
	once sync.Once
)

type PoolRunner struct {
	maxRunners int
}

// / Instance returns the singleton instance of the PoolRunner
func NewPool(options ...Option) {
	once.Do(func() {
		config := &PoolConfig{maxRunners: 10}

		for _, option := range options {
			option(config)
		}

		inst = &PoolRunner{
			maxRunners: config.maxRunners,
		}
	})

}

func Execute[IN any, OUT any](f func(IN) OUT, data []IN) []OUT {
	if inst == nil {
		panic("Didn't initialize the pool")
	}

	wg := sync.WaitGroup{}
	results := make(chan OUT, len(data))
	sem := make(chan struct{}, inst.maxRunners)

	for _, x := range data {
		wg.Add(1)
		sem <- struct{}{}
		go func(item IN) {
			defer func() {
				wg.Done()
				<-sem
			}()
			results <- f(item)
		}(x)
	}

	wg.Wait()
	close(results)

	var outputResults []OUT
	for result := range results {
		outputResults = append(outputResults, result)
	}

	return outputResults
}

type KeyedResult[K comparable, V any] struct {
	Key   K
	Value V
}

func RunInParallel[K comparable, OUT any](f ...func() KeyedResult[K, OUT]) []KeyedResult[K, OUT] {
	wg := sync.WaitGroup{}
	wg.Add(len(f))

	results := make(chan KeyedResult[K, OUT], len(f))
	for _, x := range f {
		go func(f func() KeyedResult[K, OUT]) {
			defer wg.Done()
			results <- f()
		}(x)
	}

	close(results)

	var out []KeyedResult[K, OUT]
	for x := range results {
		out = append(out, x)
	}
	return out
}

// FindByKey returns the first element in the data slice that has the key
func FindByKey[K comparable, OUT any](key K, data []KeyedResult[K, OUT]) *KeyedResult[K, OUT] {
	for _, x := range data {
		if x.Key == key {
			return &x
		}
	}
	return nil
}
