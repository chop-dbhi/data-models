package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
)

var ErrInvalidRepo = errors.New("repo: invalid repo URI")

const defaultRepoName = "https://github.com/chop-dbhi/data-models@master"

type Repos []*Repo

// Repos is a string or repo that implement the flag.Value interface.
func (r *Repos) String() string {
	return defaultRepoName
}

func (r *Repos) Set(s string) error {
	p, err := ParseRepo(s)

	if err != nil {
		return err
	}

	*r = append(*r, p)

	return nil
}

type Repo struct {
	URL    string
	Branch string

	FetchTime  time.Time
	CommitSHA1 string
	CommitTime time.Time

	// For updating.
	sync.Mutex

	// SHA1 of the previous state.
	prevSHA1 string
	updating bool
	path     string
}

func (r *Repo) String() string {
	return fmt.Sprintf("%s@%s", r.URL, r.Branch)
}

func (r *Repo) MarshalJSON() ([]byte, error) {
	aux := map[string]interface{}{
		"uri":       r.URL,
		"branch":    r.Branch,
		"fetchTime": r.FetchTime,
		"commit": map[string]interface{}{
			"sha1": r.CommitSHA1,
			"time": r.CommitTime,
		},
	}

	return json.Marshal(aux)
}

func (r *Repo) info() {
	cmd := exec.Command("git", "log", "-1", "--format=%H|%ct")

	buf := bytes.NewBuffer(nil)

	cmd.Dir = r.path
	cmd.Stdout = buf

	if err := cmd.Run(); err != nil {
		logrus.Errorf("repo: error getting commit info: %s", err)
	}

	v := strings.TrimSpace(buf.String())
	parts := strings.Split(v, "|")

	ts, err := strconv.Atoi(parts[1])

	if err != nil {
		logrus.Errorf("repo: error parsing timestamp: %s", err)
	}

	r.prevSHA1 = r.CommitSHA1
	r.CommitSHA1 = parts[0]
	r.CommitTime = time.Unix(int64(ts), 0)
	r.FetchTime = time.Now()
}

func (r *Repo) clone() {
	cmd := exec.Command("git", "clone", "--branch", r.Branch, r.URL, r.path)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		logrus.Fatalf("problem cloning repo: %s", err)
	}

	logrus.Debugf("repo: cloneed repo %s", r)
	r.info()
}

func (r *Repo) pull() {
	remote := fmt.Sprintf("origin/%s", r.Branch)
	cmd := exec.Command("git", "pull", ".", remote)

	cmd.Dir = r.path
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		logrus.Fatalf("problem pulling repo: %s", err)
		return
	}

	logrus.Debugf("repo: pulled repo %s", r)
	r.info()
}

// updateRepo clones or updates the repo and returns true
// if an update occurred.
func (r *Repo) update() bool {
	// Update already in progress
	if r.updating {
		return false
	}

	r.Lock()

	defer func() {
		r.updating = false
		r.Unlock()
	}()

	r.updating = true

	gitDir := filepath.Join(r.path, ".git")

	if _, err := os.Stat(gitDir); err != nil {
		r.clone()
	} else {
		r.pull()
	}

	return r.CommitSHA1 != r.prevSHA1
}

func ParseRepo(uri string) (*Repo, error) {
	toks := strings.SplitN(uri, "@", 2)

	r := &Repo{
		Branch: "master",
	}

	uri = toks[0]

	// A remote URL must be absolute.
	if purl, err := url.ParseRequestURI(uri); err == nil {
		r.URL = uri

		// Go-style namespacing e.g. github.com/chop-dbhi/data-models
		r.path = filepath.Join(reposDir, purl.Host, purl.Path)
	} else if uri, err = filepath.Abs(uri); err == nil {
		gitDir := filepath.Join(uri, ".git")

		if _, err = os.Stat(gitDir); err != nil {
			return nil, err
		}

		r.URL = uri
		r.path = uri
	} else {
		return nil, ErrInvalidRepo
	}

	if len(toks) > 1 {
		r.Branch = toks[1]
	}

	return r, nil
}

// Update all the repos.
func updateRepos() {
	wg := sync.WaitGroup{}
	wg.Add(len(registeredRepos))

	var changed bool

	for _, r := range registeredRepos {
		go func(r *Repo) {
			if r.update() {
				changed = true
			}

			wg.Done()
		}(r)
	}

	wg.Wait()

	// Rebuild the cache if any of the repos changed.
	if changed {
		rebuildCache()
	}
}

// pollRepo periodically checks the repo for updates.
func pollRepos() {
	// Check for updates every hour.
	t := time.NewTicker(interval)

	for {
		select {
		case <-t.C:
			updateRepos()
		}
	}
}
