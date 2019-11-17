package main

import (
	"flag"
	"fmt"
	"net/http"

	"form3/api"
	"form3/business/sqlite"

	"github.com/go-chi/docgen"
	log "github.com/sirupsen/logrus"
)

var doImport = flag.Bool("import", false, "run the import before starting the server")
var addr = flag.String("addr", ":14051", "Server listen address")
var storage = flag.String("storage", ":memory:", "Either ':memory:' or a file path for sqlite")

func main() {
	flag.Parse()

	persistence, err := sqlite.New(*storage)
	if err != nil {
		log.Fatalf("can't initialize storage: %w", err)
	}
	logger := log.New()
	logger.Formatter = &log.TextFormatter{DisableTimestamp: true}

	service := api.NewService(persistence, *logger)

	if *doImport {
		if err := persistence.ImportData(); err != nil {
			log.Fatalf("Import failed: %w", err)
		}
	}

	fmt.Printf("starting server on %s with the following middlewares and routes:\n\n", *addr)
	fmt.Println(docgen.MarkdownRoutesDoc((&api.Service{}).GetMux(), docgen.MarkdownOpts{
		ProjectPath: "christmas-api",
	}))

	if err = http.ListenAndServe(*addr, service.GetMux()); err != nil {
		log.Fatalf("server crashed:\n%w", err)
	}
}
