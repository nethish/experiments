package main

import (
	"fmt"
	"net/http"
	"sync"
	"time"
)

type CircuitBreaker struct {
	failures      int
	successes     int
	state         string // "closed", "open", "half-open"
	lock          sync.Mutex
	lastFailureAt time.Time
}

func NewCircuitBreaker() *CircuitBreaker {
	return &CircuitBreaker{
		state: "closed",
	}
}

func (cb *CircuitBreaker) beforeRequest() bool {
	cb.lock.Lock()
	defer cb.lock.Unlock()

	switch cb.state {
	case "open":
		// Allow retry after 3 seconds
		if time.Since(cb.lastFailureAt) > 3*time.Second {
			cb.state = "half-open"
			return true
		}
		return false
	default:
		return true
	}
}

func (cb *CircuitBreaker) afterRequest(success bool) {
	cb.lock.Lock()
	defer cb.lock.Unlock()

	if success {
		cb.successes++
		cb.failures = 0
		if cb.state == "half-open" && cb.successes >= 2 {
			cb.state = "closed"
			fmt.Println("[CircuitBreaker] Transition to CLOSED")
		}
	} else {
		cb.failures++
		cb.successes = 0
		cb.lastFailureAt = time.Now()
		if cb.failures >= 3 {
			cb.state = "open"
			fmt.Println("[CircuitBreaker] Transition to OPEN")
		}
	}
}

func main() {
	cb := NewCircuitBreaker()
	client := http.Client{
		Timeout: 2 * time.Second,
	}

	for {
		if !cb.beforeRequest() {
			fmt.Println("Circuit breaker OPEN, skipping request...")
			time.Sleep(1 * time.Second)
			continue
		}

		resp, err := client.Get("http://localhost:8080/")
		if err != nil {
			fmt.Println("Request failed:", err)
			cb.afterRequest(false)
		} else {
			if resp.StatusCode == http.StatusOK {
				fmt.Println("Success:", resp.StatusCode)
				cb.afterRequest(true)
			} else {
				fmt.Println("Server error:", resp.StatusCode)
				cb.afterRequest(false)
			}
			resp.Body.Close()
		}
		time.Sleep(1 * time.Second)
	}
}
