package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"

	"../RateLimiter/client"
	"../RateLimiter/models"
)

func main() {
	/*r, err := client.NewThrottleRateLimiter(
	&models.Config{
		Throttle: 1 * time.Second,
	})*/
	/*r, err := client.NewMaxConcurrencyLimiter(&models.Config{
		Limit:           2,
		TokenResetAfter: 10 * time.Second,
	})*/
	r, err := client.NewFixedWindowRateLimiter(&models.Config{
		Limit:         5,
		FixedInterval: 15 * time.Second,
	})
	if err != nil {
		panic(err)
	}
	var wg sync.WaitGroup
	rand.Seed(time.Now().UnixNano())
	doWork := func(id int) {
		defer wg.Done()
		token, err := r.Acquire()
		if err != nil {
			panic(err)
		}
		fmt.Printf("rate limit token is %s and acquired at %s", token.ID, token.CreatedAt)
		n := rand.Intn(5)
		fmt.Printf("Worke %d sleeping for %d\n", id, n)
		time.Sleep(time.Duration(n) * time.Second)
		fmt.Printf("Woker %d done\n", id)
		//r.Release(token)
	}
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go doWork(i)
	}
	wg.Wait()
}
