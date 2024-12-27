package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"sync"
	"time"
)

type queryValidator struct {
	params      map[string][]validator
	rateLimiter *rateLimiter
}

type validator interface {
	Validate(string, *http.Request) error
}

type requiredValidator struct{}

func (v requiredValidator) Validate(s string, _ *http.Request) error {
	if s == "" {
		return fmt.Errorf("parameter is required")
	}
	return nil
}

type intValidator struct{}

func (v intValidator) Validate(s string, _ *http.Request) error {
	_, err := strconv.Atoi(s)
	return err
}

type allowedValuesValidator struct {
	allowed []string
}

func (v allowedValuesValidator) Validate(s string, _ *http.Request) error {
	for _, a := range v.allowed {
		if s == a {
			return nil
		}
	}
	return fmt.Errorf("parameter value is not allowed")
}

type conditionalValidator struct {
	dependency string
	value      string
	validator  validator
}

func (v conditionalValidator) Validate(s string, r *http.Request) error {
	dependencyValues, ok := r.URL.Query()[v.dependency]
	if !ok || len(dependencyValues) == 0 {
		return nil
	}

	for _, dependencyValue := range dependencyValues {
		if dependencyValue == v.value {
			return v.validator.Validate(s, r)
		}
	}

	return nil
}

type rateLimiter struct {
	limits map[string][]int // Change to store timestamps as a slice
	mu     sync.Mutex
}

func newRateLimiter(limits map[string]int) *rateLimiter {
	// Initialize limits with empty slices for each key
	rlLimits := make(map[string][]int)
	for key := range limits {
		rlLimits[key] = []int{}
	}
	return &rateLimiter{limits: rlLimits}
}

func (rl *rateLimiter) Allow(key string) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	limit, ok := rl.limits[key]
	if !ok {
		return false // No limit defined for this key
	}

	now := time.Now()
	count := 0
	// Iterate over timestamps (now using time Unix format)
	for _, reqTime := range limit {
		if now.Sub(time.Unix(int64(reqTime), 0)) < time.Minute {
			count++
		}
	}

	// Check if count is less than limit
	if count < limit[0] {
		// Append the current timestamp
		rl.limits[key] = append(rl.limits[key], int(now.Unix()))
		return true
	}
	return false
}

func newQueryValidator(limits map[string]int) *queryValidator {
	rl := newRateLimiter(limits)
	return &queryValidator{params: make(map[string][]validator), rateLimiter: rl}
}

func (v *queryValidator) AddRule(param string, validators ...validator) {
	v.params[param] = append(v.params[param], validators...)
}

func (v *queryValidator) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		for param, validators := range v.params {
			values, ok := r.URL.Query()[param]
			if !ok || len(values) == 0 {
				continue
			}

			for _, value := range values {
				key := fmt.Sprintf("%s:%s", param, value)
				if !v.rateLimiter.Allow(key) {
					http.Error(w, fmt.Sprintf("Rate limit exceeded for parameter %s with value %s", param, value), http.StatusTooManyRequests)
					return
				}

				for _, validator := range validators {
					if err := validator.Validate(value, r); err != nil {
						http.Error(w, fmt.Sprintf("Invalid query parameter %s: %v", param, err), http.StatusBadRequest)
						return
					}
				}
			}
		}

		next.ServeHTTP(w, r)
	})
}

func main() {
	limits := map[string]int{
		"limit:100":  10, // Allows 10 requests per minute for limit=100
		"order:desc": 5,  // Allows 5 requests per minute for order=desc
	}

	v := newQueryValidator(limits)
	v.AddRule("limit", requiredValidator{}, intValidator{})
	v.AddRule("order", allowedValuesValidator{allowed: []string{"asc", "desc"}})

	// Conditional validation
	v.AddRule("skip", conditionalValidator{
		dependency: "action",
		value:      "create",
		validator:  requiredValidator{},
	})

	v.AddRule("sort", conditionalValidator{
		dependency: "action",
		value:      "list",
		validator:  requiredValidator{},
	})

	http.Handle("/api/data", v.Middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Data endpoint reached")
	})))

	log.Fatal(http.ListenAndServe(":8080", nil))
}
