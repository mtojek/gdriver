# gdriver

[![Build Status](https://travis-ci.com/mtojek/gdriver.svg?branch=main)](https://travis-ci.com/mtojek/gdriver)

`gdriver` is a command-line tool, written in Go, used for downloading large personal files from Google Drive (API v3).
The tool provides file selection, integrity checks, transfer retries and requires a user-defined Cloud Platform project. 

[![asciicast](https://asciinema.org/a/372715.png)](https://asciinema.org/a/372715)

## Features

* Uses [Google Drive v3 API](https://developers.google.com/drive/api/v3/about-sdk)
* TUI-based file selection
* File integrity check (MD5 checksum)
* Transfer retries (app-internal and on application restart)
* [OAuth 2.0](https://developers.google.com/drive/api/v3/about-auth#OAuth2Authorizing) authorization protocol

## Requirements

* User-defined Cloud Platform project with enabled Drive API ([Quickstart](https://developers.google.com/drive/api/v3/quickstart/go#step_1_turn_on_the))

## Getting started

If you have installed [Go SDK](https://golang.org/doc/install#download), you can download and build the latest master of the `gdriver` binary:

```bash
go get github.com/mtojek/gdriver
```

Without Go SDK you can download prebuilt distribution from the [Releases](https://github.com/mtojek/gdriver/releases) page.
If you're working on a remote workstation, copy the URL link for particular release, use `curl` or `wget` to download the archive
and unpack it later (e.g `tar xvzf gdriver_X.Y.Z_linux_amd64.tar.gz`).

Run the `help` command to see available commands:

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

Hints:
* Use a meaningful project name as it will be presented as title in the Google authentication form.
* Select `Desktop application` type for the OAuth client.
* There were issues reported in the past with downloading and saving the credentials file in Firefox. In case of facing a similar issue,
please try to use Google Chrome.

Once you create the project, remember to download related login credentials.

The command will import the above credentials into the `~/.gdriver` directory and authenticate the Google user account:

```bash
gdriver auth --import-credentials credentials.json
```

You should be good to go now. Try to download first files using the `gdriver download` command, e.g.:

```bash
gdriver download <folderID> --select --output tmp
```

The `folderID` is the ID of a Drive folder (e.g. `Ax9h4tAyI53ZhqMSoa2opZ6o6m21OUyww`). The value can be easily copied from the
URL bar in the web browser. Open the directory in the [Google Drive](https://drive.google.com/) console and pick
the `folderID` part from the URL (e.g. `https://drive.google.com/drive/u/0/folders/Ax9h4tAyI53ZhqMSoa2opZ6o6m21OUyww`).

## Releases

Find latest revisions on the [Releases](https://github.com/mtojek/gdriver/releases) page.

## License

Apache License