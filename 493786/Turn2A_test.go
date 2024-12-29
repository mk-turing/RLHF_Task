package _93786

import (
	"errors"
	"os"
	"sync"
	"testing"
	"time"
)

var (
	ErrQuotaExceeded = errors.New("quota exceeded")
)

// QuotaFile simulates a file with a quota.
type QuotaFile struct {
	file      *os.File
	quota     int64
	usedSpace int64
	mu        sync.Mutex
}

// Write simulates writing to a file with quota enforcement.
func (qf *QuotaFile) Write(p []byte) (n int, err error) {
	qf.mu.Lock()
	defer qf.mu.Unlock()

	writeSize := int64(len(p))
	if qf.usedSpace+writeSize > qf.quota {
		return 0, ErrQuotaExceeded
	}

	qf.usedSpace += writeSize
	return len(p), nil
}

// AdjustQuota dynamically adjusts the quota of the QuotaFile.
func (qf *QuotaFile) AdjustQuota(newQuota int64) {
	qf.mu.Lock()
	defer qf.mu.Unlock()
	qf.quota = newQuota
}

// SetQuota creates a QuotaFile with the specified quota.
func SetQuota(file *os.File, quota int64) *QuotaFile {
	return &QuotaFile{
		file:  file,
		quota: quota,
	}
}

func TestDynamicFileSystemQuota(t *testing.T) {
	testCases := []struct {
		name        string
		quota       int64
		writeSize   int
		numWrites   int
		adjustments []struct {
			time     time.Duration
			newQuota int64
		}
		expectedErr error
	}{
		{
			name:      "Exceed Quota and Adjust",
			quota:     1024,
			writeSize: 512,
			numWrites: 4,
			adjustments: []struct {
				time     time.Duration
				newQuota int64
			}{
				{time.Second * 1, 2048}, // Adjust quota after 1 second
			},
			expectedErr: nil,
		},
		{
			name:      "Adjust Before Exceeding Quota",
			quota:     1024,
			writeSize: 512,
			numWrites: 3,
			adjustments: []struct {
				time     time.Duration
				newQuota int64
			}{
				{time.Second * 1, 2048}, // Adjust quota after 1 second
			},
			expectedErr: nil,
		},
		{
			name:      "Multiple Adjustments",
			quota:     1024,
			writeSize: 512,
			numWrites: 6,
			adjustments: []struct {
				time     time.Duration
				newQuota int64
			}{
				{time.Second * 1, 2048}, // Adjust quota after 1 second
				{time.Second * 3, 1024}, // Adjust quota back to original after 3 seconds
			},
			expectedErr: ErrQuotaExceeded,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create a temporary file
			tmpFile, err := os.CreateTemp("", "quota-test-*")
			if err != nil {
				t.Fatalf("Failed to create temporary file: %v", err)
			}
			defer os.Remove(tmpFile.Name())

			// Set quota on the file
			qf := SetQuota(tmpFile, tc.quota)

			// Perform quota adjustments in a goroutine
			var wg sync.WaitGroup
			for _, adjustment := range tc.adjustments {
				wg.Add(1)
				go func(adjustment struct {
					time     time.Duration
					newQuota int64
				}) {
					defer wg.Done()
					time.Sleep(adjustment.time)
					qf.AdjustQuota(adjustment.newQuota)
				}(adjustment)
			}

			// Perform writes
			for i := 0; i < tc.numWrites; i++ {
				buf := make([]byte, tc.writeSize)
				_, err := qf.Write(buf)
				if err != nil {
					if err == ErrQuotaExceeded {
						// Wait for adjustments to complete before continuing
						wg.Wait()
						_, err = qf.Write(buf)
						if err != nil && err != tc.expectedErr {
							t.Fatalf("Unexpected error on write %d: %v", i, err)
						}
						return // Expected error occurred, test passes
					}

					if err != tc.expectedErr {
						t.Fatalf("Unexpected error on write %d: %v", i, err)
					}
					return
				}
			}

			// Wait for adjustments to complete
			wg.Wait()

			// If no error occurred during writes and quota was not exceeded, fail the test
			if tc.expectedErr != nil {
				t.Fatalf("Expected error %v but did not occur", tc.expectedErr)
			}
		})
	}
}
