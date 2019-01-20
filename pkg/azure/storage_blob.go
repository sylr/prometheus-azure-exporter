package azure

import (
	"context"
	"fmt"
	"net/url"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/storage/mgmt/2018-07-01/storage"
	"github.com/Azure/azure-storage-blob-go/azblob"
)

var (
	blobFormatString = `https://%s.blob.core.windows.net`
)

// StorageAccountContainerWalker in an interface that is to be implemented by struct you want
// to pass to Walking functions like WalkStorageAccount().
type StorageAccountContainerWalker interface {
	// Lock prevents concurrent process to WalkBlob().
	Lock()
	// Unlock releases the lock.
	Unlock()
	// WalkBlob is called for all blobs listed by the Walking function.
	WalkBlob(*storage.Account, *storage.ListContainerItem, *azblob.BlobItem)
}

// WalkStorageAccountContainer applies a function on all storage account containter blobs.
func WalkStorageAccountContainer(ctx context.Context, clients *AzureClients, account *storage.Account, container *storage.ListContainerItem, walker StorageAccountContainerWalker) error {
	token, err := GetStorageToken(ctx)

	if err != nil {
		return err
	}

	// ADAL credentials
	accessToken := token.Token().AccessToken
	credential := azblob.NewTokenCredential(accessToken, nil)

	// Preparing browsing container.
	pipeline := azblob.NewPipeline(credential, azblob.PipelineOptions{})
	url, _ := url.Parse(fmt.Sprintf(blobFormatString, *account.Name))
	serviceURL := azblob.NewServiceURL(*url, pipeline)
	containerURL := serviceURL.NewContainerURL(*container.Name)

	marker := azblob.Marker{}

	for i := 0; ; i++ {
		t0 := time.Now()
		list, err := containerURL.ListBlobsFlatSegment(
			ctx,
			marker,
			azblob.ListBlobsSegmentOptions{
				Details: azblob.BlobListingDetails{
					Snapshots: true,
				},
			},
		)
		t1 := time.Since(t0).Seconds()

		ObserveAzureAPICall(t1)
		ObserveAzureStorageAPICall(t1)

		if err != nil {
			ObserveAzureAPICallFailed(t1)
			ObserveAzureStorageAPICallFailed(t1)
			return err
		}

		// Update request marker.
		marker = list.NextMarker

		walker.Lock()
		for _, blob := range list.Segment.BlobItems {
			walker.WalkBlob(account, container, &blob)
		}
		walker.Unlock()

		// Continue iterating if we are not done.
		if !marker.NotDone() {
			break
		}
	}

	if err != nil {
		return err
	}

	return nil
}
