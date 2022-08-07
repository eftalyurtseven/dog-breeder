package main

import (
	"bytes"
	"context"
	"flag"
	"fmt"
	"image"
	"image/jpeg"
	"io"
	"io/ioutil"
	"net"
	"os"
	"time"

	"github.com/eftalyurtseven/dogbreed/insecure"
	"github.com/eftalyurtseven/dogbreed/models"
	pbBreeder "github.com/eftalyurtseven/dogbreed/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/grpclog"
)

func main() {
	var (
		addr  = flag.String("addr", "localhost", "grpc service address")
		port  = flag.String("port", "8080", "grpc service port")
		breed = flag.String("breed", "crow", "breed to search")
	)
	log := grpclog.NewLoggerV2(os.Stdout, ioutil.Discard, io.Discard)
	grpclog.SetLoggerV2(log)
	flag.Parse()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
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

	c := pbBreeder.NewBreederClient(conn)
	r, err := c.Search(ctx, &pbBreeder.DogReq{Breed: *breed})
	if err != nil {
		log.Fatalln("Failed to search: %v", err)
	}

	img, _, err := image.Decode(bytes.NewBuffer(r.Image))
	if err != nil {
		log.Fatalln("Failed to decode image: %v", err)
	}

	// Create path if not exists
	if _, err := os.Stat(models.IMAGE_PATH); os.IsNotExist(err) {
		if os.Mkdir(models.IMAGE_PATH, models.IMAGE_PATH_CHMOD) != nil {
			log.Fatalln("Failed to create directory: %v", err)
		}
	}

	out, err := os.Create(fmt.Sprintf("%s/%s.jpg", models.IMAGE_PATH, *breed))
	if err != nil {
		log.Fatalln("Failed to create file: %v", err)
	}
	defer out.Close()

	var opts jpeg.Options
	opts.Quality = 100

	err = jpeg.Encode(out, img, &opts)
	if err != nil {
		log.Fatalln("Failed to encode image: %v", err)
	}
}
