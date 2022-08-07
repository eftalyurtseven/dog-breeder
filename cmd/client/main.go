package main

import (
	"context"
	"flag"
	"io"
	"io/ioutil"
	"net"
	"os"
	"time"

	"github.com/eftalyurtseven/dogbreed/client"
	"github.com/eftalyurtseven/dogbreed/insecure"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/grpclog"
)

func main() {
	var (
		addr  = flag.String("addr", "0.0.0.0", "grpc service address")
		port  = flag.String("port", "8080", "grpc service port")
		breed = flag.String("breed", "chow", "breed to search")
	)
	log := grpclog.NewLoggerV2(os.Stdout, ioutil.Discard, io.Discard)
	grpclog.SetLoggerV2(log)
	flag.Parse()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	conn, err := grpc.DialContext(
		ctx,
		net.JoinHostPort(*addr, *port),
		grpc.WithTransportCredentials(
			credentials.NewClientTLSFromCert(insecure.CertPool, ""),
		),
	)
	if err != nil {
		log.Fatalln("Failed to start server: %v", err)
	}

	defer conn.Close()

	// breeder client
	c := client.New(conn)
	r, err := c.Search(ctx, *breed)
	if err != nil {
		log.Fatalln("Failed to search: %v", err)
	}

	name, err := c.SaveImage(ctx, *breed, r.Image)
	if err != nil {
		log.Fatalln("Failed to save image: %v", err)
	}
	log.Info("Image saved at %s", name)
}
