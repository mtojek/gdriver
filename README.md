# gdriver

`gdriver` is a command-line tool, written in Go, used for downloading large personal files from Google Drive (API v3).
The tool provides file selection, integrity checks, transfer retries and requires a user-defined Cloud Platform project. 

[![Build Status](https://travis-ci.com/mtojek/gdriver.svg?branch=main)](https://travis-ci.com/mtojek/gdriver)

## Features

* Uses [Google Drive v3 API](https://developers.google.com/drive/api/v3/about-sdk)
* TUI-based file selection
* File integrity check (MD5 checksum)
* Transfer retries (app-internal and on application restart)
* [OAuth 2.0](https://developers.google.com/drive/api/v3/about-auth#OAuth2Authorizing) authorization protocol

## Requirements

* User-defined Cloud Platform project with enabled Drive API ([Quickstart](https://developers.google.com/drive/api/v3/quickstart/go#step_1_turn_on_the))

## Getting started

Download and build the latest master of `gdriver` binary:

```bash
go get github.com/mtojek/gdriver
```

Alternatively, you can download built distribution from the [Releases](https://github.com/mtojek/gdriver/releases) page.

Run the `help` command and see available commands:

```bash
gdriver help

Use gdriver to download large files from Google Drive.

Usage:
  gdriver [command]

Available Commands:
  auth        Authenticate Google account
  check       Check files
  download    Download files
  help        Help about any command

Flags:
  -h, --help   help for gdriver

Use "gdriver [command] --help" for more information about a command.
```

### Run the application for the first time

Import client configuration (`credentials.json`) for the Cloud Platform project. If you haven't created the project
yet or enabled the Drive API, follow the [Quickstart](https://developers.google.com/drive/api/v3/quickstart/go#step_1_turn_on_the) steps.

The command will save the credentials in the `~/.gdriver` directory and authenticate the Google user account:

```bash
gdriver auth --import-credentials credentials.json
```

You should be good to go now. Try to download first files using the `gdriver download` command, e.g.:

```bash
gdriver download Ax9h4tAyI53ZhqMSoa2opZ6o6m21OUyww --select --output tmp
```

## Releases

Find latest revisions on the [Releases](https://github.com/mtojek/gdriver/releases) page.