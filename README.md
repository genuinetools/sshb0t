# sshb0t

[![Travis CI](https://travis-ci.org/genuinetools/sshb0t.svg?branch=master)](https://travis-ci.org/genuinetools/sshb0t)

A bot for keeping your ssh `authorized_keys` up to date with user's GitHub keys
from `https://github.com/{username}.keys`.

**WARNING:** Only use this if you have two factor auth enabled for your GitHub
account and you make sure to delete old keys from your account.

## Installation

#### Binaries

- **darwin** [386](https://github.com/genuinetools/sshb0t/releases/download/v0.3.3/sshb0t-darwin-386) / [amd64](https://github.com/genuinetools/sshb0t/releases/download/v0.3.3/sshb0t-darwin-amd64)
- **freebsd** [386](https://github.com/genuinetools/sshb0t/releases/download/v0.3.3/sshb0t-freebsd-386) / [amd64](https://github.com/genuinetools/sshb0t/releases/download/v0.3.3/sshb0t-freebsd-amd64)
- **linux** [386](https://github.com/genuinetools/sshb0t/releases/download/v0.3.3/sshb0t-linux-386) / [amd64](https://github.com/genuinetools/sshb0t/releases/download/v0.3.3/sshb0t-linux-amd64) / [arm](https://github.com/genuinetools/sshb0t/releases/download/v0.3.3/sshb0t-linux-arm) / [arm64](https://github.com/genuinetools/sshb0t/releases/download/v0.3.3/sshb0t-linux-arm64)
- **solaris** [amd64](https://github.com/genuinetools/sshb0t/releases/download/v0.3.3/sshb0t-solaris-amd64)
- **windows** [386](https://github.com/genuinetools/sshb0t/releases/download/v0.3.3/sshb0t-windows-386) / [amd64](https://github.com/genuinetools/sshb0t/releases/download/v0.3.3/sshb0t-windows-amd64)

#### Via Go

```bash
$ go get github.com/genuinetools/sshb0t
```

#### Running with Docker

```console
$ docker run -d --restart always \
    --name sshb0t \
    -v ${HOME}/.ssh/authorized_keys:/root/.ssh/authorized_keys \
    r.j3ss.co/sshb0t --user genuinetools --keyfile /root/.ssh/authorized_keys
```

## Usage

```console
$ sshb0t -h
         _     _      ___  _
 ___ ___| |__ | |__  / _ \| |_
/ __/ __| '_ \| '_ \| | | | __|
\__ \__ \ | | | |_) | |_| | |_
|___/___/_| |_|_.__/ \___/ \__|

 A bot for keeping your ssh authorized_keys up to date with user's GitHub keys
 Version: v0.3.3
 Build: 996b83e

  -d    run in debug mode
  -interval duration
        update interval (ex. 5ms, 10s, 1m, 3h) (default 30s)
  -keyfile string
        file to update the authorized_keys (default "~/.ssh/authorized_keys")
  -once
        run once and exit, do not run as a daemon
  -url string
        GitHub Enterprise URL (default "https://github.com")
  -user value
        GitHub usernames for which to fetch keys
  -v    print version and exit (shorthand)
  -version
        print version and exit
```



[![Analytics](https://ga-beacon.appspot.com/UA-29404280-16/sshb0t/README.md)](https://github.com/genuinetools/sshb0t)
