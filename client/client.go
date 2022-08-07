package client

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"image"
	"image/jpeg"
	"os"

	"github.com/eftalyurtseven/dogbreed/models"
	pbBreeder "github.com/eftalyurtseven/dogbreed/proto"
	"google.golang.org/grpc"
)

type Client struct {
	bc pbBreeder.BreederClient
}

func New(conn *grpc.ClientConn) *Client {
	bc := pbBreeder.NewBreederClient(conn)
	return &Client{bc: bc}
}

func (c *Client) Search(ctx context.Context, breed string) (*pbBreeder.DogRes, error) {
	r, err := c.bc.Search(ctx, &pbBreeder.DogReq{Breed: breed})
	if err != nil {
		return nil, err
	}

	if !r.Status {
		return nil, errors.New("no result found")
	}

	return r, nil
}

func (c *Client) SaveImage(ctx context.Context, breed string, imgByte []byte) (string, error) {

	// Create path if not exists
	if _, err := os.Stat(models.IMAGE_PATH); os.IsNotExist(err) {
		if os.Mkdir(models.IMAGE_PATH, models.IMAGE_PATH_CHMOD) != nil {
			return "", err
		}
	}

	img, _, err := image.Decode(bytes.NewBuffer(imgByte))
	if err != nil {
		return "", err
	}

	out, err := os.Create(fmt.Sprintf("%s/%s.jpg", models.IMAGE_PATH, breed))
	if err != nil {
		return "", err
	}
	defer out.Close()

	var opts jpeg.Options
	opts.Quality = 100

	if err := jpeg.Encode(out, img, &opts); err != nil {
		return "", err
	}

	return out.Name(), nil
}
