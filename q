[1mdiff --git pkg/azure/storage_blob.go pkg/azure/storage_blob.go[m
[1mindex 585961b..77abc49 100644[m
[1m--- pkg/azure/storage_blob.go[m
[1m+++ pkg/azure/storage_blob.go[m
[36m@@ -14,9 +14,10 @@[m [mvar ([m
 	blobFormatString = `https://%s.blob.core.windows.net`[m
 )[m
 [m
[31m-// ContainerWalker.[m
[32m+[m[32m// ContainerWalker in an interface that is to be implemented by struct you want[m
[32m+[m[32m// to pass to Walking functions like WalkStorageAccount().[m
 type ContainerWalker interface {[m
[31m-	// Lock prevents concurrent process to ObserveBlob().[m
[32m+[m	[32m// Lock prevents concurrent process to WalkBlob().[m
 	Lock()[m
 	// Unlock releases the lock.[m
 	Unlock()[m
[1mdiff --git pkg/metrics/storage.go pkg/metrics/storage.go[m
[1mindex 66fd741..0b9aa13 100644[m
[1m--- pkg/metrics/storage.go[m
[1m+++ pkg/metrics/storage.go[m
[36m@@ -66,7 +66,7 @@[m [mfunc UpdateStorageMetrics(ctx context.Context) {[m
 	blobsMetrics := StorageAccountContainerMetrics{BlobsSizeHistogram: hist}[m
 [m
 	// Loop over storage accounts.[m
[31m-	for _, account := range *storageAccounts {[m
[32m+[m	[32mfor i, account := range *storageAccounts {[m
 		accountLogger := contextLogger.WithFields(log.Fields{[m
 			"account": *account.Name,[m
 		})[m
[36m@@ -109,11 +109,18 @@[m [mfunc UpdateStorageMetrics(ctx context.Context) {[m
 			// https://play.golang.org/p/YRGEg4LS5jd[m
 			// https://github.com/golang/go/wiki/CommonMistakes#using-goroutines-on-loop-iterator-variables[m
 			// -------------------------------------------------[m
[32m+[m			[32mif key > 0 {[m
[32m+[m				[32mbreak[m
[32m+[m			[32m}[m
 		}[m
 [m
 		wg.Wait()[m
 [m
 		accountLogger.Debugf("Done updating storage account")[m
[32m+[m
[32m+[m		[32mif i > 0 {[m
[32m+[m			[32mbreak[m
[32m+[m		[32m}[m
 	}[m
 [m
 	// swapping current registered histogram with an updated copy[m
[1mdiff --git pkg/tools/wg.go pkg/tools/wg.go[m
[1mindex c146076..69d3394 100644[m
[1m--- pkg/tools/wg.go[m
[1m+++ pkg/tools/wg.go[m
[36m@@ -6,15 +6,19 @@[m [mimport ([m
 	"sync"[m
 )[m
 [m
[32m+[m[32m// BoundedWaitGroup is a wait group which has a limit boundary meaning it will[m
[32m+[m[32m// wait for Done() to be called before releasing Add(n) if the limit has been reached[m
 type BoundedWaitGroup struct {[m
 	wg sync.WaitGroup[m
 	ch chan struct{}[m
 }[m
 [m
[32m+[m[32m// NewBoundedWaitGroup returns a new BoundedWaitGroup[m
 func NewBoundedWaitGroup(cap int) BoundedWaitGroup {[m
 	return BoundedWaitGroup{ch: make(chan struct{}, cap)}[m
 }[m
 [m
[32m+[m[32m// Add ...[m
 func (bwg *BoundedWaitGroup) Add(delta int) {[m
 	for i := 0; i > delta; i-- {[m
 		<-bwg.ch[m
[36m@@ -25,6 +29,7 @@[m [mfunc (bwg *BoundedWaitGroup) Add(delta int) {[m
 	bwg.wg.Add(delta)[m
 }[m
 [m
[32m+[m[32m// Done ...[m
 func (bwg *BoundedWaitGroup) Done() {[m
 	bwg.Add(-1)[m
 }[m
