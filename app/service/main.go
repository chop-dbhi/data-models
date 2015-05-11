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

var (
	host       string
	port       int
	loglevel   string
	repoDir    string
	repoBranch string
	interval   time.Duration
)

func main() {
	// Bind and parse flags
	flag.StringVar(&loglevel, "log", "info", "Specify the log level.")
	flag.StringVar(&host, "host", "127.0.0.1", "Host or IP to bind to")
	flag.IntVar(&port, "port", 8123, "Port to bind to")
	flag.StringVar(&repoDir, "repo", "repo", "Path to the cloned repo")
	flag.StringVar(&repoBranch, "branch", "master", "Repo branch")
	flag.DurationVar(&interval, "interval", time.Hour, "Fetch inteval")

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

	// Update the repo on startup.
	go updateRepo()

	// Start the interval.
	go pollRepo()

	// Setup routes.
	router := httprouter.New()

	router.RedirectTrailingSlash = true
	router.RedirectFixedPath = true
	router.HandleMethodNotAllowed = true

	router.GET("/models", apiDataModels)
	router.GET("/models/:name/:version", apiDataModel)

	// Integrations.
	router.GET("/github", githubWebhook)

	// Listen.
	addr := fmt.Sprintf("%s:%d", host, port)
	logrus.Printf("Listening on %s...", addr)
	logrus.Fatal(http.ListenAndServe(addr, router))
}
