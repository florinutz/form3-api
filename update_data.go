package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	gzBoltHttp "github.com/florinutz/gz-boltdb/http"
)

const (
	dataFile   = "test-data"
	bucketName = "form3"
)

// no need for this anymore
func _main() {
	outputPath := dataFile
	if len(os.Args) > 1 {
		outputPath = strings.Join(os.Args[1:], " ")
	}

	reqs, err := generateRequests([]string{"http://mockbin.org/bin/41ca3269-d8c4-4063-9fd5-f306814ff03f"}, nil)
	if err != nil {
		log.Fatal(err)
	}

	errs := updateTestData(reqs, outputPath)
	for _, err := range errs {
		_, _ = fmt.Fprintln(os.Stderr, err)
	}
}

func generateRequests(urls []string, tweakFn func(req *http.Request)) (reqs []*http.Request, err error) {
	for _, u := range urls {
		var req *http.Request
		req, err = http.NewRequest("GET", u, nil)
		if err != nil {
			err = fmt.Errorf("could not create a request for url '%s'", u)
			return
		}
		if tweakFn != nil {
			tweakFn(req)
		}
		reqs = append(reqs, req)
	}
	return
}

func updateTestData(reqs []*http.Request, outputPath string) []error {
	return gzBoltHttp.DumpResponses(reqs, outputPath, bucketName, nil, func(response *http.Response) error {
		fmt.Printf("* received %s\n", response.Request.URL.String())
		return nil
	})
}
