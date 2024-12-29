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

type QuotaFile struct {
	file      *os.File
	quota     int64
	usedSpace int64
	mu        sync.Mutex
}

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

func (qf *QuotaFile) Close() error {
	return qf.file.Close()
}

func SetQuota(file *os.File, quota int64) *QuotaFile {
	return &QuotaFile{
		file:  file,
		quota: quota,
	}
}

func AdjustQuota(qf *QuotaFile, newQuota int64) {
	qf.mu.Lock()
	defer qf.mu.Unlock()
	qf.quota = newQuota
}

func TestDynamicFileSystemQuota(t *testing.T) {
	testCases := []struct {
		name        string
		initialQuota int64 // Initial quota in bytes
		writeSize   int   // Size of each write in bytes
		numWrites   int   // Number of writes to perform
		adjustQuota int64 // Quota adjustment (positive for increase, negative for decrease)
		expectedErr error // Expected error after the quota adjustment
	}{
		{
			name:        "Exceed Quota, Adjust Up",
			initialQuota: 1024,
			writeSize:   512,
			numWrites:   3,
			adjustQuota:  1024, // Increase quota by 1KB
			expectedErr: nil,
		},
		{
			name:        "Within Quota, Adjust Down",
			initialQuota: 1024,
			writeSize:   256,
			numWrites:   3,
			adjustQuota: -512, // Decrease quota by 512B
			expectedErr: ErrQuotaExceeded,
		},
		{
			name:        "Quota Zero, Adjust Up",
			initialQuota: 0,
			writeSize:   1,
			numWrites:   1,
			adjustQuota: 1,
			expectedErr: nil,
		},
		{
			name:        "Single Write Exceed Quota, Adjust Up",
			initialQuota: 1024,
			writeSize:   2048,
			numWrites:   1,
			adjustQuota: 2048,
			expectedErr: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tmpFile, err := os.CreateTemp("", "quota-test-*")
			if err != nil {
				t.Fatalf("Failed to create temporary file: %v", err)
			}
			defer os.Remove(tmpFile.Name())

			qf := SetQuota(tmpFile, tc.initialQuota)

			go func() {
				time.Sleep(2 * time.Second) // Simulate delay in quota adjustment
				AdjustQuota(qf, qf.quota+tc.adjustQuota)
			}()