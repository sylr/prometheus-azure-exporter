package azure

import (
	"context"
	"fmt"
	"net/url"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/preview/subscription/mgmt/2018-03-01-preview/subscription"
	"github.com/Azure/azure-sdk-for-go/services/resources/mgmt/2018-05-01/resources"
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
	WalkBlob(*subscription.Model, *resources.Group, *storage.Account, *storage.ListContainerItem, *azblob.BlobItem)
}

// WalkStorageAccountContainer applies a function on all storage account containter blobs.
func WalkStorageAccountContainer(ctx context.Context, clients *AzureClients, subscription *subscription.Model, account *storage.Account, container *storage.ListContainerItem, walker StorageAccountContainerWalker) error {
	token, err := GetStorageToken(ctx)

	if err != nil {
		return err
	}

	details, _ := ParseResourceID(*account.ID)
	group, err := GetResourceGroup(ctx, clients, subscription, details.ResourceGroup)

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
	listOptions := azblob.ListBlobsSegmentOptions{
		Details: azblob.BlobListingDetails{
			Snapshots: true,
		},
	}

	for i := 0; ; i++ {
		t0 := time.Now()
		list, err := containerURL.ListBlobsFlatSegment(ctx, marker, listOptions)
		t1 := time.Since(t0).Seconds()

		if err != nil {
			ObserveAzureAPICallFailed(t1)
			ObserveAzureStorageAPICallFailed(t1, *subscription.DisplayName, *group.Name, *account.Name)
			return err
		}

		ObserveAzureAPICall(t1)
		ObserveAzureStorageAPICall(t1, *subscription.DisplayName, *group.Name, *account.Name)

		// Update request marker.
		marker = list.NextMarker

		walker.Lock()
		for _, blob := range list.Segment.BlobItems {
			walker.WalkBlob(subscription, group, account, container, &blob)
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
