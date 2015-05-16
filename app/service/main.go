package main

import (
	"flag"
	"fmt"
	"net/http"
	"path/filepath"
	"time"

	"github.com/julienschmidt/httprouter"
	"github.com/sirupsen/logrus"
)

const defaultRepoName = "https://github.com/chop-dbhi/data-models"

var (
	host       string
	port       int
	loglevel   string
	repoDir    string
	repoName   string
	repoBranch string
	interval   time.Duration
)

func main() {
	// Bind and parse flags
	flag.StringVar(&loglevel, "log", "info", "Specify the log level.")
	flag.StringVar(&host, "host", "127.0.0.1", "Host or IP to bind to")
	flag.IntVar(&port, "port", 8123, "Port to bind to")
	flag.StringVar(&repoDir, "path", "data-models", "Local name of the cloned repo")
	flag.StringVar(&repoName, "repo", defaultRepoName, "Remote path or URL of the repository.")
	flag.StringVar(&repoBranch, "branch", "master", "Branch to use.")
	flag.DurationVar(&interval, "interval", time.Hour, "The interval for checking for updates")

	flag.Parse()

	if lvl, err := logrus.ParseLevel(loglevel); err != nil {
		logrus.Fatalf("invalid log level")
	} else {
		logrus.SetLevel(lvl)
	}

	var err error

	repoDir, err = filepath.Abs(repoDir)

	if err != nil {
		logrus.Fatalf("bad repo path")
	}

	// Setup routes.
	router := httprouter.New()

	router.RedirectTrailingSlash = true
	router.RedirectFixedPath = true
	router.HandleMethodNotAllowed = true

	// API
	router.GET("/api/models", apiModels)
	router.GET("/api/models/:name/:version", apiModel)
	router.GET("/api/models/:name/:version/:table", apiTable)
	router.GET("/api/models/:name/:version/:table/:field", apiField)

	// Views.
	router.GET("/models.:ext", viewModels)
	router.GET("/models/:name/:version/full.:ext", viewModelFull)
	router.GET("/models/:name/:version/definition.:ext", viewModelDefinition)
	router.GET("/models/:name/:version/schema.:ext", viewModelSchema)
	router.GET("/models/:name/:version/mapping.:ext", viewModelMapping)

	// Update the repo on startup.
	go updateRepo()

	// Poll the repo.
	if interval > 0 {
		go pollRepo()
	}

	// Listen.
	addr := fmt.Sprintf("%s:%d", host, port)

	logrus.Printf("Listening on %s...", addr)
	logrus.Fatal(http.ListenAndServe(addr, router))
}
