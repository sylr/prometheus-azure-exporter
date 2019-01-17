package metrics

import (
	"sync"

	"github.com/Azure/azure-sdk-for-go/services/storage/mgmt/2018-07-01/storage"
	"github.com/Azure/azure-storage-blob-go/azblob"
	"github.com/prometheus/client_golang/prometheus"
)

// StorageAccountContainerMetrics ...
type StorageAccountContainerMetrics struct {
	BlobsSizeHistogram *prometheus.HistogramVec
	mutex              sync.RWMutex
}

// Lock is here to make sure several Walkers do not update BlobsSizeHistogram
// at the same time.
func (s *StorageAccountContainerMetrics) Lock() {
	s.mutex.Lock()
}

// Unlock releases the lock.
func (s *StorageAccountContainerMetrics) Unlock() {
	s.mutex.Unlock()
}

// Reset resets the histogram data.
func (s *StorageAccountContainerMetrics) Reset() {
	s.Lock()
	s.BlobsSizeHistogram.Reset()
	s.Unlock()
}

// DeleteLabelValues deletes histogram's data associated with given labels.
func (s *StorageAccountContainerMetrics) DeleteLabelValues(labels ...string) {
	s.Lock()
	s.BlobsSizeHistogram.DeleteLabelValues(labels...)
	s.Unlock()
}

// WalkBlob is called over each blobs listed by the function walking the
// storage account container.
func (s *StorageAccountContainerMetrics) WalkBlob(account *storage.Account, container *storage.ListContainerItem, blob *azblob.BlobItem) {
	s.BlobsSizeHistogram.
		WithLabelValues(*account.Name, *container.Name).
		Observe(float64(*blob.Properties.ContentLength))
}
