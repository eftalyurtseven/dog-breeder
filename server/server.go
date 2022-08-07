package server

import (
	"context"
	"encoding/json"
	"fmt"
	"io"

	"github.com/eftalyurtseven/dogbreed/models"
	pbBreeder "github.com/eftalyurtseven/dogbreed/proto"
	"github.com/eftalyurtseven/dogbreed/protocol"
	"github.com/hashicorp/go-retryablehttp"
)

// Backend is the implementation of the backend service
type Backend struct {
	hc  *retryablehttp.Client
	url string
}

// Svc is the interface for the backend service
type Svc interface {
	Search(ctx context.Context, req *pbBreeder.DogReq) (*pbBreeder.DogRes, error)
}

// for type safety on build time.
var _ pbBreeder.BreederServer = (*Backend)(nil)

// NewBackend returns a new backend service
func New(hc *retryablehttp.Client) Svc {
	return &Backend{
		hc:  hc,
		url: models.DOGAPIURL,
	}
}

// Search returns the image for the breed passed in as a parameter
// It sends HTTP request to the dog api
// It returns the image and error if any
// Search can work with GRPC as well as HTTP
func (b *Backend) Search(ctx context.Context, req *pbBreeder.DogReq) (*pbBreeder.DogRes, error) {

	// make a request to the dog api and get the image for the breed
	// if the server is down, it will retry with RetryMax times
	res, err := b.hc.Get(fmt.Sprintf(b.url, req.Breed))
	if err != nil {
		return nil, fmt.Errorf("error fetching image: %v", err)
	}
	defer res.Body.Close()

	resp := new(protocol.DogApiResp)
	if err := json.NewDecoder(res.Body).Decode(resp); err != nil {
		return nil, fmt.Errorf("json.Decode failed with error: %v", err)
	}

	// if user requested a breed that is not found, dog api will return an error as a status
	// if the status is not success, return the error
	if resp.Status == models.DOGAPIERRSTATUS {
		return nil, fmt.Errorf("%v", resp.Message)
	}

	// fetch the image from the image url
	// if the image is not found, return the error
	image, err := b.fetchImage(ctx, resp.Message)
	if err != nil {
		return nil, fmt.Errorf("b.fetchImage failed w error: %v", err)
	}

	return &pbBreeder.DogRes{
		Image:  image,
		Status: true,
	}, nil
}

// fetchImage fetches the image from the url
// returns the image and error if any
func (b *Backend) fetchImage(ctx context.Context, image string) ([]byte, error) {
	res, err := b.hc.Get(image)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	return io.ReadAll(res.Body)
}
