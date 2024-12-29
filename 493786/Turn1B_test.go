package _93786

import (
	"errors"
	"os"
	"sync"
	"testing"
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

// Close closes the underlying file.
func (qf *QuotaFile) Close() error {
	return qf.file.Close()
}

// SetQuota creates a QuotaFile with the specified quota.
func SetQuota(file *os.File, quota int64) *QuotaFile {
	return &QuotaFile{
		file:  file,
		quota: quota,
	}
}

func TestFileSystemQuota(t *testing.T) {
	testCases := []struct {
		name        string
		quota       int64 // Quota in bytes
		writeSize   int   // Size of each write in bytes
		numWrites   int   // Number of writes to perform
		expectedErr error // Expected error
	}{
		{
			name:        "Exceed Quota",
			quota:       1024, // 1KB quota
			writeSize:   512,  // Write 512 bytes at a time
			numWrites:   3,    // Should exceed quota with 3 writes
			expectedErr: ErrQuotaExceeded,
		},
		{
			name:        "Within Quota",
			quota:       1024,
			writeSize:   256,
			numWrites:   4,
			expectedErr: nil,
		},
		{
			name:        "Quota Zero",
			quota:       0,
			writeSize:   1,
			numWrites:   1,
			expectedErr: ErrQuotaExceeded,
		},
		{
			name:        "Single Write Exceed Quota",
			quota:       1024,
			writeSize:   2048,
			numWrites:   1,
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

			// Perform writes
			for i := 0; i < tc.numWrites; i++ {
				buf := make([]byte, tc.writeSize)
				_, err := qf.Write(buf)
				if err != nil {
					if err != tc.expectedErr {
						t.Fatalf("Unexpected error on write %d: %v", i, err)
					}
					return // Expected error occurred, test passes
				}
			}

			// If no error occurred during writes and quota was not exceeded, fail the test
			if tc.expectedErr != nil {
				t.Fatalf("Expected error %v but did not occur", tc.expectedErr)
			}
		})
	}
}
