package server

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	pbBreeder "github.com/eftalyurtseven/dogbreed/proto"
	"github.com/hashicorp/go-retryablehttp"
)

var dummyErr = errors.New("dummy error")

var tests = []struct {
	name        string
	status      string
	req         *pbBreeder.DogReq
	res         *pbBreeder.DogRes
	expectedRes *pbBreeder.DogRes
	expectedErr error
}{
	{
		name:   "success",
		status: "success",
		req: &pbBreeder.DogReq{
			Breed: "labrador",
		},
		res: &pbBreeder.DogRes{
			Image:  []byte("image"),
			Status: true,
		},
		expectedRes: &pbBreeder.DogRes{
			Image:  []byte("image"),
			Status: true,
		},
		expectedErr: nil,
	},
	{
		name:   "error check",
		status: "error",
		req: &pbBreeder.DogReq{
			Breed: "labrador",
		},
		res: &pbBreeder.DogRes{
			Image:  []byte("image"),
			Status: true,
		},
		expectedRes: &pbBreeder.DogRes{
			Image:  []byte("image"),
			Status: true,
		},
		expectedErr: dummyErr,
	},
}

// mockImageHTTPServer returns a mock server that returns the image passed in as a parameter
// to the handler.
// TODO: move to mockHTTPServer func with url.path
func mockImageHTTPServer(expected []byte) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write(expected)
	}))
}

// mockHTTPServer returns a mock server that returns the status and message passed in as parameters
// to the handler.
func mockHTTPServer(status, imageSrvURI string) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(fmt.Sprintf(`{"message":"%s","status":"%s","code":0}`, imageSrvURI, status)))
	}))
}

// TestSearch tests the Search method of the server.
// It tests the following scenarios:
// 1. Successful search
// 2. Error search
func TestSearch(t *testing.T) {
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a mock server that returns image data
			mockImageSvr := mockImageHTTPServer(tt.res.Image)
			defer mockImageSvr.Close()

			// Create a mock server that returns status and message
			var server *httptest.Server
			if tt.status == "success" {
				server = mockHTTPServer(tt.status, mockImageSvr.URL)
			} else {
				server = mockHTTPServer(tt.status, dummyErr.Error())
			}
			defer server.Close()

			// Create a new backend server
			b := &Backend{
				hc:  retryablehttp.NewClient(),
				url: server.URL + "/%s", // %s is a placeholder for the breed
			}

			res, err := b.Search(context.Background(), tt.req)
			if err != nil && err.Error() != tt.expectedErr.Error() {
				t.Errorf("expected error %v, got %v", tt.expectedErr, err)
			}

			if res == nil {
				return
			}
			if !bytes.Equal(res.Image, tt.expectedRes.Image) {
				t.Errorf("unexpected image: %v, %v", string(res.Image), string(tt.expectedRes.Image))
			}
		})
	}
}
