package client

import (
	"context"
	"errors"
	"net"
	"testing"

	"github.com/eftalyurtseven/dogbreed/insecure"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

const (
	addr = "0.0.0.0" // TODO: change this to docker container ip
	port = "8080"
)

var tests = []struct {
	name        string
	breed       string
	expectedErr error
}{
	{
		name:        "valid",
		breed:       "affenpinscher",
		expectedErr: nil,
	},
	{
		name:        "invalid",
		breed:       "eftal",
		expectedErr: errors.New("rpc error: code = Unknown desc = Breed not found (master breed does not exist)"),
	},
}

func getConnection(ctx context.Context) (*grpc.ClientConn, error) {
	return grpc.DialContext(
		ctx,
		net.JoinHostPort(addr, port),
		grpc.WithTransportCredentials(credentials.NewClientTLSFromCert(insecure.CertPool, "")), // TODO: use cert pool for integration tests
	)
}

func TestSearchIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	ctx := context.Background()
	conn, err := getConnection(ctx)
	if err != nil {
		t.Fatalf("Failed to start server: %v", err)
	}

	client := New(conn)
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			_, err := client.Search(ctx, test.breed)
			if err != nil && err.Error() != test.expectedErr.Error() {
				t.Errorf("expected error %v, got %v", test.expectedErr, err)
			}
		})
	}
}
