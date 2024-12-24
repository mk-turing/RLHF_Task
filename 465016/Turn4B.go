package main

import (
	"errors"
	"fmt"
	"strings"
	"sync"
)

type Data map[string]interface{}

// Define custom error types
type ValidationError struct {
	Key string
	Err error
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("Validation error for key '%s': %v", e.Key, e.Err)
}

type TransformationError struct {
	Key string
	Err error
}

func (e *TransformationError) Error() string {
	return fmt.Sprintf("Transformation error for key '%s': %v", e.Key, e.Err)
}

// RollbackManager struct to handle concurrent rollback
type RollbackManager struct {
	mu       sync.Mutex
	stages   []map[string]interface{}
	rollback map[string]interface{}
}

func newRollbackManager() *RollbackManager {
	return &RollbackManager{
		stages:   make([]map[string]interface{}, 0),
		rollback: make(map[string]interface{}),
	}
}

// AddStage adds a new stage to the rollback manager
func (r *RollbackManager) AddStage() map[string]interface{} {
	r.mu.Lock()
	defer r.mu.Unlock()

	stage := make(map[string]interface{})
	r.stages = append(r.stages, stage)
	return stage
}

// PerformRollback rolls back all changes recorded in the rollback manager
func (r *RollbackManager) PerformRollback(data Data) {
	r.mu.Lock()
	defer r.mu.Unlock()

	// Rollback all stages in reverse order
	for i := len(r.stages) - 1; i >= 0; i-- {
		rollbackStage := r.stages[i]
		for key, oldValue := range rollbackStage {
			switch v := oldValue.(type) {
			case int:
				data[key] = v
			case string:
				data[key] = v
			default:
				delete(data, key)
			}
		}
	}
}

// Example transformation steps with concurrency handling
func transformAdd(data Data, r *RollbackManager, key string, value interface{}) error {
	rollbackStage := r.AddStage()
	if data[key] != nil {
		return &ValidationError{Key: key, Err: fmt.Errorf("Key %s already exists", key)}
	}
	data[key] = value
	rollbackStage[key] = nil
	return nil
}

func transformMultiply(data Data, r *RollbackManager, key string) error {
	rollbackStage := r.AddStage()
	value, ok := data[key].(int)
	if !ok {
		return &TransformationError{Key: key, Err: errors.New("Invalid type for multiplication")}
	}
	data[key] = value * 2
	rollbackStage[key] = value
	return nil
}

func transformStringToUpper(data Data, r *RollbackManager, key string) error {
	rollbackStage := r.AddStage()
	value, ok := data[key].(string)
	if !ok {
		return &TransformationError{Key: key, Err: errors.New("Invalid type for string conversion")}
	}
	data[key] = strings.ToUpper(value)
	rollbackStage[key] = value
	return nil
}

// Transformation worker that runs in parallel
func transformationWorker(wg *sync.WaitGroup, data Data, r *RollbackManager, stage int) {
	defer wg.Done()
	var err error