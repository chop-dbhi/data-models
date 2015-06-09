package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/julienschmidt/httprouter"
	"github.com/sirupsen/logrus"
)

// Global variables.
var (
	registeredRepos Repos
	reposDir        string
	secret          string
	googleAnalytics string
)

func main() {
	var (
		host     string
		port     int
		loglevel string
		interval time.Duration
	)

	// Bind and parse flags
	flag.StringVar(&loglevel, "log", "info", "Specify the log level.")
	flag.StringVar(&host, "host", "127.0.0.1", "Host or IP to bind to.")
	flag.IntVar(&port, "port", 8123, "Port to bind to.")
	flag.StringVar(&reposDir, "path", "data-models", "Local directory of the cloned repos")
	flag.DurationVar(&interval, "interval", time.Hour, "The interval for checking for updates.")
	flag.StringVar(&secret, "secret", "", "Secret for webhook integration.")
	flag.StringVar(&googleAnalytics, "ga", "", "Google Analytics tracking code.")
	flag.Var(&registeredRepos, "repo", "Git repository to include. Multiple values can be supplied.")

	flag.Parse()

	if lvl, err := logrus.ParseLevel(loglevel); err != nil {
		logrus.Fatalf("invalid log level")
	} else {
		logrus.SetLevel(lvl)
	}

	var err error

	reposDir, err = filepath.Abs(reposDir)

	if err = os.MkdirAll(reposDir, os.ModeDir|0775); err != nil {
		logrus.Fatalf("could not create repos directory: %s", err)
	}

	// Setup routes.
	router := httprouter.New()

	router.RedirectTrailingSlash = true
	router.RedirectFixedPath = true
	router.HandleMethodNotAllowed = true

	router.GET("/", httpIndex)
	router.GET("/repos", httpReposList)
	router.GET("/models", httpModels)
	router.GET("/models/:name", httpModel)
	router.GET("/models/:name/:version", httpModelVersion)
	router.GET("/models/:name/:version/:table", httpTable)
	router.GET("/models/:name/:version/:table/:field", httpField)
	router.GET("/compare/:name1/:version1/:name2/:version2", httpCompareModels)
	router.GET("/schemata/:name/:version", httpModelSchema)

	// Endpoint for webhook integration.
	router.POST("/_hook", httpUpdateRepos)

	// Update the repo on startup.
	go updateRepos()

	// Poll the repo.
	if interval > 0 {
		go pollRepos(interval)
	}

	// Listen.
	addr := fmt.Sprintf("%s:%d", host, port)

	logrus.Printf("Listening on %s...", addr)
	logrus.Fatal(http.ListenAndServe(addr, router))
}
