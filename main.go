package main

import (
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"os"

	"github.com/eftalyurtseven/dogbreed/insecure"
	pbBreeder "github.com/eftalyurtseven/dogbreed/proto"
	"github.com/eftalyurtseven/dogbreed/server"
	grpc_validator "github.com/grpc-ecosystem/go-grpc-middleware/validator"
	"github.com/hashicorp/go-retryablehttp"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/grpclog"
)

func main() {
	// flag
	var gRPCPort = flag.Int("grpc-port", 8080, "gRPC port")

	// log definition and initialization
	log := grpclog.NewLoggerV2(os.Stdout, ioutil.Discard, io.Discard)
	grpclog.SetLoggerV2(log)

	// flag parsing
	flag.Parse()

	addr := fmt.Sprintf("localhost:%d", *gRPCPort)
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalln("Failed to start server: %v", err)
	}

	// create grpc server that uses TLS
	s := grpc.NewServer(
		grpc.Creds(credentials.NewServerTLSFromCert(&insecure.Cert)),
		grpc.UnaryInterceptor(grpc_validator.UnaryServerInterceptor()),
		grpc.StreamInterceptor(grpc_validator.StreamServerInterceptor()),
	)
	// create http client for dog breed service
	hc := getHttpClient()
	// register the service with the server
	pbBreeder.RegisterBreederServer(s, server.New(hc))
	log.Info("Serving gRPC on https://", addr)

	go func() {
		log.Fatal(s.Serve(lis))
	}()
	select {}
}

func getHttpClient() *retryablehttp.Client {
	httpClient := retryablehttp.NewClient()
	httpClient.RetryMax = 3
	return httpClient
}
