package main

import (
	"fmt"
	"sync"
	"time"
)

type TrafficLightState struct {
	TrafficLight1 string
	TrafficLight2 string
	sync.Mutex
}

func (t *TrafficLightState) Update(light string, state string) {
	t.Lock()
	defer t.Unlock()
	switch light {
	case "TrafficLight1":
		t.TrafficLight1 = state
	case "TrafficLight2":
		t.TrafficLight2 = state
	default:
		fmt.Println("Invalid traffic light identifier")
	}
}

func (t *TrafficLightState) Status() string {
	t.Lock()
	defer t.Unlock()
	return fmt.Sprintf("TrafficLight1: %s, TrafficLight2: %s", t.TrafficLight1, t.TrafficLight2)
}

func main() {
	var tls TrafficLightState
	tls.TrafficLight1 = "Red"
	tls.TrafficLight2 = "Green"

	updateChan := make(chan struct{}) // Channel for update notification
	statusChan := make(chan string)   // Channel for status updates

	go func() {
		for range updateChan {
			fmt.Println(tls.Status())
		}
	}()

	go func() {
		states := []string{"Green", "Yellow", "Red"}
		for {
			time.Sleep(2 * time.Second)
			tls.Update("TrafficLight1", states[(tls.TrafficLight1 == "Green")+1])
			updateChan <- struct{}{} // Notify update
		}
	}()

	go func() {
		states := []string{"Red", "Yellow", "Green"}
		for {
			time.Sleep(2 * time.Second)
			tls.Update("TrafficLight2", states[(tls.TrafficLight2 == "Red")+1])
			updateChan <- struct{}{} // Notify update
		}
	}()

	time.Sleep(60 * time.Second) // Run for 60 seconds
	close(updateChan)
}
