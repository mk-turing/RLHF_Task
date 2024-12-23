package main

import (
	"fmt"
	"log"
	"sync"

	"github.com/fsnotify/fsnotify"
)

// directoryWatcher holds the watcher for a specific directory
// along with its associated action function
type directoryWatcher struct {
	watcher *fsnotify.Watcher
	path    string
	action  func(event fsnotify.Event)
}

func newDirectoryWatcher(path string, action func(event fsnotify.Event)) (*directoryWatcher, error) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}
	dw := &directoryWatcher{
		watcher: watcher,
		path:    path,
		action:  action,
	}
	return dw, nil
}

// start starts the directory watcher and handles events
func (dw *directoryWatcher) start() {
	go func() {
		defer dw.watcher.Close()
		for {
			select {
			case event := <-dw.watcher.Events:
				// Handle events here
				dw.action(event)
			case err := <-dw.watcher.Errors:
				log.Println("error:", err)
			}
		}
	}()
}

func main() {
	// Define the actions for each directory
	actionForDir1 := func(event fsnotify.Event) {
		if event.Op&fsnotify.Write == fsnotify.Write {
			fmt.Println("Modified file in dir1:", event.Name)
		}
	}

	actionForDir2 := func(event fsnotify.Event) {
		if event.Op&fsnotify.Create == fsnotify.Create {
			fmt.Println("Created file in dir2:", event.Name)
		}
	}

	// Create watchers for multiple directories
	dir1Watcher, err := newDirectoryWatcher("dir", actionForDir1)
	if err != nil {
		log.Fatal(err)
	}

	dir2Watcher, err := newDirectoryWatcher("dir1", actionForDir2)
	if err != nil {
		log.Fatal(err)
	}

	// Add watchers to the watcher list
	watchers := []*directoryWatcher{dir1Watcher, dir2Watcher}

	// Start all watchers
	var wg sync.WaitGroup
	wg.Add(len(watchers))
	for _, dw := range watchers {
		if err := dw.watcher.Add(dw.path); err != nil {
			log.Println("Error adding watcher for", dw.path, ":", err)
			continue
		}
		dw.start()
		wg.Done()
	}

	// Wait for all watchers to start
	wg.Wait()

	fmt.Println("File watchers started. Press any key to exit...")
	fmt.Scanln()
}
