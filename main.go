package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"os/user"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	"github.com/genuinetools/pkg/cli"
	"github.com/genuinetools/sshb0t/version"
	"github.com/sirupsen/logrus"
)

const (
	defaultKeyURI                string = "%s/%s.keys"
	defaultSSHAuthorizedKeysFile string = ".ssh/authorized_keys"
)

var (
	home               string
	authorizedKeysFile string
	enturl             string
	users              stringSlice

	interval time.Duration
	once     bool

	debug bool
)

// stringSlice is a slice of strings
type stringSlice []string

// implement the flag interface for stringSlice
func (s *stringSlice) String() string {
	return fmt.Sprintf("%s", *s)
}
func (s *stringSlice) Set(value string) error {
	*s = append(*s, value)
	return nil
}

func main() {
	var err error
	// get the home directory
	home, err = getHomeDir()
	if err != nil {
		logrus.Fatalf("getHomeDir failed: %v", err)
	}

	// Create a new cli program.
	p := cli.NewProgram()
	p.Name = "sshb0t"
	p.Description = "A bot for keeping your ssh authorized_keys up to date with user's GitHub keys"

	// Set the GitCommit and Version.
	p.GitCommit = version.GITCOMMIT
	p.Version = version.VERSION

	// Setup the global flags.
	p.FlagSet = flag.NewFlagSet("global", flag.ExitOnError)
	p.FlagSet.StringVar(&authorizedKeysFile, "keyfile", filepath.Join(home, defaultSSHAuthorizedKeysFile), "file to update the authorized_keys")
	p.FlagSet.StringVar(&enturl, "url", "https://github.com", "GitHub Enterprise URL")
	p.FlagSet.Var(&users, "user", "GitHub usernames for which to fetch keys")

	p.FlagSet.DurationVar(&interval, "interval", 30*time.Second, "update interval (ex. 5ms, 10s, 1m, 3h)")
	p.FlagSet.BoolVar(&once, "once", false, "run once and exit, do not run as a daemon")

	p.FlagSet.BoolVar(&debug, "d", false, "enable debug logging")

	// Set the before function.
	p.Before = func(ctx context.Context) error {
		// Set the log level.
		if debug {
			logrus.SetLevel(logrus.DebugLevel)
		}

		if len(users) < 1 {
			return errors.New("you must pass at least one username")
		}

		if len(authorizedKeysFile) < 1 {
			return errors.New("you must pass a file to save the authorized keys into or use the default")
		}
		return nil
	}

	// Set the main program action.
	p.Action = func(ctx context.Context, args []string) error {
		ticker := time.NewTicker(interval)

		// On ^C, or SIGTERM handle exit.
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt)
		signal.Notify(c, syscall.SIGTERM)
		go func() {
			for sig := range c {
				ticker.Stop()
				logrus.Infof("Received %s, exiting.", sig.String())
				os.Exit(0)
			}
		}()

		// If the user passed the once flag, just do the run once and exit.
		if once {
			run()
			os.Exit(0)
		}

		logrus.Infof("Starting bot to update %s every %s for users %s", authorizedKeysFile, interval, strings.Join(users, ", "))
		for range ticker.C {
			run()
		}

		return nil
	}

	// Run our program.
	p.Run()
}

func run() {
	// fetch the keys for each user
	var keys string
	for _, user := range users {
		// fetch the url
		uri := fmt.Sprintf(defaultKeyURI, enturl, user)
		baseURL, err := url.Parse(uri)
		if err != nil {
			logrus.Fatalf("parsing url %s failed: %v", uri, err)
		}
		uri = baseURL.String()
		logrus.Debugf("Fetching keys for user %s from %s", user, uri)
		resp, err := http.Get(uri)
		if err != nil {
			logrus.Warnf("Fetching keys for user %s from %s failed: %v", user, uri, err)
			continue
		}
		// make sure we got status 200
		if http.StatusOK != resp.StatusCode {
			logrus.Warnf("Expected status code 200 from %s but got %d for user %s", uri, resp.StatusCode, user)
			continue
		}
		// read the body
		b, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			logrus.Fatalf("Reading response body from %s for user %s failed: %v", uri, user, err)
			continue
		}
		// append to keys variable with a new line
		keys += string(b)
	}

	// update the authorized key file
	logrus.Infof("Updating authorized key file %s with keys from %s", authorizedKeysFile, strings.Join(users, ", "))
	if err := ioutil.WriteFile(authorizedKeysFile, []byte(keys), 0600); err != nil {
		logrus.Fatalf("Writing to file %s failed: %v", authorizedKeysFile, err)
	}
	logrus.Info("Successfully updated keys")
}

func getHomeDir() (string, error) {
	home := os.Getenv(homeKey)
	if home != "" {
		return home, nil
	}

	u, err := user.Current()
	if err != nil {
		return "", err
	}
	return u.HomeDir, nil
}
