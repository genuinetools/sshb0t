package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/signal"
	"os/user"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	"github.com/jessfraz/sshb0t/version"
	"github.com/sirupsen/logrus"
)

const (
	defaultGitHubKeyURI          string = "https://%s/%s.keys"
	defaultSSHAuthorizedKeysFile string = ".ssh/authorized_keys"

	// BANNER is what is printed for help/info output
	BANNER = `         _     _      ___  _
 ___ ___| |__ | |__  / _ \| |_
/ __/ __| '_ \| '_ \| | | | __|
\__ \__ \ | | | |_) | |_| | |_
|___/___/_| |_|_.__/ \___/ \__|
 A bot for keeping your ssh authorized_keys up to date with user's GitHub keys
 Version: %s
`
)

var (
	home               string
	authorizedKeysFile string
	gitURI             string
	users              stringSlice
	interval           string

	debug bool
	vrsn  bool
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

func init() {
	var err error
	// get the home directory
	home, err = getHomeDir()
	if err != nil {
		logrus.Fatalf("getHomeDir failed: %v", err)
	}

	// parse flags
	flag.StringVar(&authorizedKeysFile, "keyfile", filepath.Join(home, defaultSSHAuthorizedKeysFile), "file to update the authorized_keys")
	flag.StringVar(&gitURI, "gituri", "github.com", "Add custom git URI (ex. gitlab.com, github.com)")
	flag.Var(&users, "user", "GitHub usernames for which to fetch keys")
	flag.StringVar(&interval, "interval", "30s", "update interval (ex. 5ms, 10s, 1m, 3h)")

	flag.BoolVar(&vrsn, "version", false, "print version and exit")
	flag.BoolVar(&vrsn, "v", false, "print version and exit (shorthand)")
	flag.BoolVar(&debug, "d", false, "run in debug mode")

	flag.Usage = func() {
		fmt.Fprint(os.Stderr, fmt.Sprintf(BANNER, version.VERSION))
		flag.PrintDefaults()
	}

	flag.Parse()

	if vrsn {
		fmt.Printf("sshb0t version %s, build %s", version.VERSION, version.GITCOMMIT)
		os.Exit(0)
	}

	if flag.NArg() >= 1 {
		// parse the arg
		arg := flag.Args()[0]

		if arg == "help" {
			usageAndExit("", 0)
		}

		if arg == "version" {
			fmt.Printf("sshb0t version %s, build %s", version.VERSION, version.GITCOMMIT)
			os.Exit(0)
		}
	}

	if len(users) < 1 {
		usageAndExit("you must pass at least one username", 1)
	}

	if len(authorizedKeysFile) < 1 {
		usageAndExit("you must pass a file to save the authorized keys into or use the default", 1)
	}

	// set log level
	if debug {
		logrus.SetLevel(logrus.DebugLevel)
	}
}

func main() {
	var ticker *time.Ticker

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

	// parse the duration
	dur, err := time.ParseDuration(interval)
	if err != nil {
		logrus.Fatalf("parsing %s as duration failed: %v", interval, err)
	}
	ticker = time.NewTicker(dur)

	logrus.Infof("Starting bot to update %s every %s for users %s", authorizedKeysFile, interval, strings.Join(users, ", "))
	for range ticker.C {
		// fetch the keys for each user
		var keys string
		for _, user := range users {
			// fetch the url
			url := fmt.Sprintf(defaultGitHubKeyURI, gitURI, user)
			logrus.Debugf("Fetching keys for user %s from %s", user, url)
			resp, err := http.Get(url)
			if err != nil {
				logrus.Warnf("Fetching keys for user %s from %s failed: %v", user, url, err)
			}
			// make sure we got status 200
			if http.StatusOK != resp.StatusCode {
				logrus.Warnf("Expected status code 200 from %s but got %d for user %s", url, resp.StatusCode, user)
			}
			// read the body
			b, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				logrus.Fatalf("Reading response body from %s for user %s failed: %v", url, user, err)
			}
			// append to keys variable with a new line
			keys += string(b) + "\n"
		}

		// update the authorized key file
		logrus.Infof("Updating authorized key file %s with keys from %s", authorizedKeysFile, strings.Join(users, ", "))
		if err := ioutil.WriteFile(authorizedKeysFile, []byte(keys), 0600); err != nil {
			logrus.Fatalf("Writing to file %s failed: %v", authorizedKeysFile, err)
		}
		logrus.Info("Successfully updated keys")
	}
}

func usageAndExit(message string, exitCode int) {
	if message != "" {
		fmt.Fprintf(os.Stderr, message)
		fmt.Fprintf(os.Stderr, "\n\n")
	}
	flag.Usage()
	fmt.Fprintf(os.Stderr, "\n")
	os.Exit(exitCode)
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
