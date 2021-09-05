# npmsource [![Build Status](https://github.com/ffflorian/npmsource/workflows/Build/badge.svg)](https://github.com/ffflorian/npmsource/actions/)

Find (almost) every npm package's repository in an instant.

## Usage

Visit `npmsource.com/{packageName}` in your web browser, e.g. [`npmsource.com/nock`](https://npmsource.com/nock).

### Get the repository for a specific version

Visit `npmsource.com/{packageName}@{version}` in your web browser, e.g. [`npmsource.com/lodash@4.17.15`](https://npmsource.com/lodash@4.17.15). This also works with npm tags, e.g. [`npmsource.com/typescript@beta`](https://npmsource.com/typescript@beta)

If no version is specified, the latest version is assumed.

### Get the raw data

Visit `npmsource.com/{packageName}?raw` in your web browser, e.g. [`npmsource.com/commander?raw`](https://npmsource.com/commander?raw).

### Get source code for a specific version

Visit `npmsource.com/{packageName}?unpkg` in your web browser, e.g. [`npmsource.com/express@4.17.1?unpkg`](https://npmsource.com/express@4.17.1?unpkg). You can use the same features (`raw`, version, tags) as mentioned above.

## Server usage

### Local

Prerequisites:

- [Go](https://golang.org) >= 1.16

### Start the server

```
go run .
```
