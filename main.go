package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"time"

	"form3/api"
	"form3/business/mongo"

	"github.com/go-chi/docgen"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

var routes = flag.Bool("doc", false, "Generate router documentation")

func main() {
	flag.Parse()

	if *routes {
		fmt.Println(docgen.MarkdownRoutesDoc((&api.Service{}).GetMux(), docgen.MarkdownOpts{
			ProjectPath: "form3-payments-api",
			Intro:       "Welcome to the form3-payments-api generated docs.",
		}))

		return
	}

	dsn, ok := os.LookupEnv("MONGO_DSN")
	if !ok {
		fatal("you need to set MONGO_DSN (https://godoc.org/gopkg.in/mgo.v2#Dial)", nil)
	}

	logger := getLogger()

	service, err := getService(dsn, logger)
	if err != nil {
		fatal("can't instantiate service", err)
	}

	fmt.Printf("url: %s\n", "http://localhost:8080")

	if err = http.ListenAndServe(":8080", service.GetMux()); err != nil {
		fatal("server crashed", err)
	}
}

func fatal(wrap string, err error) {
	if err != nil {
		err = errors.Wrap(err, wrap)
	} else {
		err = errors.New(wrap)
	}
	_, _ = fmt.Fprintln(os.Stderr, err)
	os.Exit(1)
}

func getService(dsn string, logger *logrus.Logger) (*api.Service, error) {
	persistenceLayer, err := mongo.New(dsn, 2*time.Second)
	if err != nil {
		return nil, err
	}

	return api.NewService(persistenceLayer, *logger), err
}

func getLogger() *logrus.Logger {
	logger := logrus.New()
	logger.Formatter = &logrus.TextFormatter{DisableTimestamp: true}

	return logger
}
