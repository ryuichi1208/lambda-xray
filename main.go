package main

import (
	"context"
	"fmt"
	"io/ioutil"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-xray-sdk-go/xray"
	"golang.org/x/net/context/ctxhttp"
)

var (
	appName  = ""
	version  = ""
	revision = ""
)

type MyEvent struct {
	Name string `json:"name"`
}

func init() {
	xray.Configure(xray.Config{
		LogLevel:       "info",
		ServiceVersion: "1.2.3",
	})
}

func HandleRequest(ctx context.Context, evt MyEvent) (string, error) {
	// Start a segment
	ctx, seg := xray.BeginSegment(context.Background(), "api-request")
	_, err := getExample(ctx, "https://example.com")
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	// Close the segment
	_, err = getExample(ctx, "https://google.com")
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	seg.Close(nil)
	return fmt.Sprintf("%s-%s-%s, Hello %s!", appName, version, revision, evt.Name), nil
}

func getExample(ctx context.Context, url string) ([]byte, error) {
	resp, err := ctxhttp.Get(ctx, xray.Client(nil), url)
	if err != nil {
		return nil, err
	}
	return ioutil.ReadAll(resp.Body)
}

func main() {
	lambda.Start(HandleRequest)
}
