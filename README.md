# sshb0t

[![Travis CI](https://img.shields.io/travis/genuinetools/sshb0t.svg?style=for-the-badge)](https://travis-ci.org/genuinetools/sshb0t)
[![GoDoc](https://img.shields.io/badge/godoc-reference-5272B4.svg?style=for-the-badge)](https://godoc.org/github.com/genuinetools/sshb0t)
[![Github All Releases](https://img.shields.io/github/downloads/genuinetools/sshb0t/total.svg?style=for-the-badge)](https://github.com/genuinetools/sshb0t/releases)

A bot for keeping your ssh `authorized_keys` up to date with user's GitHub keys
from `https://github.com/{username}.keys`.

**WARNING:** Only use this if you have two factor auth enabled for your GitHub
account and you make sure to delete old keys from your account.

**Table of Contents**

<!-- toc -->

<!-- tocstop -->

## Installation

#### Binaries

For installation instructions from binaries please visit the [Releases Page](https://github.com/genuinetools/sshb0t/releases).

#### Via Go

```console
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
sshb0t -  A bot for keeping your ssh authorized_keys up to date with user's GitHub keys.

Usage: sshb0t <command>

Flags:

  --url       GitHub Enterprise URL (default: https://github.com)
  --user      GitHub usernames for which to fetch keys (default: [])
  -d          enable debug logging (default: false)
  --interval  update interval (ex. 5ms, 10s, 1m, 3h) (default: 30s)
  --keyfile   file to update the authorized_keys (default: /home/jessie/.ssh/authorized_keys)
  --once      run once and exit, do not run as a daemon (default: false)

Commands:

  version  Show the version information.
```

[![Analytics](https://ga-beacon.appspot.com/UA-29404280-16/sshb0t/README.md)](https://github.com/genuinetools/sshb0t)
