package notifier

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/grafana/alerting/images"
	alertingModels "github.com/grafana/alerting/models"
	alertingNotify "github.com/grafana/alerting/notify"

	"github.com/grafana/grafana/pkg/services/ngalert/models"
	"github.com/grafana/grafana/pkg/services/ngalert/store"
)

type imageProvider struct {
	store store.ImageStore
}

func newImageStore(store store.ImageStore) images.Provider {
	return &imageProvider{
		store: store,
	}
}

func (i imageProvider) GetImage(ctx context.Context, uri string) (*images.Image, error) {
	var (
		image *models.Image
		err   error
	)

	// Check whether the uri is a URL or a token to know how to query the DB.
	if strings.HasPrefix(uri, "http") {
		image, err = i.store.GetImageByURL(ctx, uri)
	} else {
		token := strings.TrimPrefix(uri, "token://")
		image, err = i.store.GetImage(ctx, token)
	}
	if err != nil {
		return nil, err
	}

	return &images.Image{
		Token:     image.Token,
		Path:      image.Path,
		URL:       image.URL,
		CreatedAt: image.CreatedAt,
	}, nil
}

func (i imageProvider) GetImageURL(ctx context.Context, alert *alertingNotify.Alert) (string, error) {
	fmt.Println("GetImageURL")
	uri, err := getImageURI(alert)
	if err != nil {
		return "", err
	}

	var image *models.Image

	// Check whether the uri is a URL or a token to know how to query the DB.
	if strings.HasPrefix(uri, "http") {
		image, err = i.store.GetImageByURL(ctx, uri)
	} else {
		token := strings.TrimPrefix(uri, "token://")
		image, err = i.store.GetImage(ctx, token)
	}
	if err != nil {
		return "", err
	}

	if !image.HasURL() {
		return "", images.ErrImagesNoURL
	}
	return image.URL, nil
}

func (i imageProvider) GetRawImage(ctx context.Context, alert *alertingNotify.Alert) (io.ReadCloser, string, error) {
	fmt.Println("GetRawImage")
	uri, err := getImageURI(alert)
	if err != nil {
		return nil, "", err
	}

	var image *models.Image
	// Check whether the uri is a URL or a token to know how to query the DB.
	if strings.HasPrefix(uri, "http") {
		image, err = i.store.GetImageByURL(ctx, uri)
	} else {
		token := strings.TrimPrefix(uri, "token://")
		image, err = i.store.GetImage(ctx, token)
	}
	if err != nil {
		return nil, "", err
	}

	if !image.HasPath() {
		return nil, "", images.ErrImagesNoPath
	}

	// Return image bytes and filename.
	readCloser, err := openImage(image.Path)
	if err != nil {
		return nil, "", err
	}
	filename := filepath.Base(image.Path)
	return readCloser, filename, nil
}

// getImageURI is a helper function to retrieve the image URI from the alert annotations as a string.
func getImageURI(alert *alertingNotify.Alert) (string, error) {
	uri, ok := alert.Annotations[alertingModels.ImageTokenAnnotation]
	if !ok {
		return "", fmt.Errorf("no image uri in annotations")
	}
	return string(uri), nil
}

// openImage returns an the io representation of an image from the given path.
func openImage(path string) (io.ReadCloser, error) {
	fp := filepath.Clean(path)
	_, err := os.Stat(fp)
	if os.IsNotExist(err) || os.IsPermission(err) {
		return nil, images.ErrImageNotFound
	}

	f, err := os.Open(fp)
	if err != nil {
		return nil, err
	}

	return f, nil
}
