# jackpot-server

[![Build Status](https://travis-ci.org/solefaucet/jackpot-server.svg?branch=master)](https://travis-ci.org/solefaucet/jackpot-server)
[![Go Report Card](http://goreportcard.com/badge/solefaucet/jackpot-server)](http://goreportcard.com/report/solefaucet/jackpot-server)
[![codecov.io](https://codecov.io/github/solefaucet/jackpot-server/coverage.svg?branch=master)](https://codecov.io/github/solefaucet/jackpot-server?branch=master)

======

## Requirement

* go1.6
* mysql5.7

## Installation

```bash
# easy enough by go get
$ go get -u github.com/solefaucet/jackpot-server
```

## DB Migration

#### Requirement

goose is needed for DB migration

```bash
$ go get bitbucket.org/liamstask/goose/cmd/goose
```

#### How to

```bash
# Migrate DB to the most recent version available
$ goose up

# Roll back version by 1
$ goose down

# Create a new migration
$ goose create SomeThingDescriptiveEnoughForYourChangeToDB sql
```

## Development

#### Dependency Management

```bash
# After third party library is introduced or removed
$ godep save ./...
```

#### Lint

```bash
$ make metalint
```

#### Test

```bash
$ make test
```

## CONTRIBUTE
* fork it
* create an issue that describes what you are going to work on
* create a new branch in your own repo and do the job
* use [commitizen](https://github.com/commitizen/cz-cli) to write commit message so we can generate change log in the future with [conventional-changelog](https://github.com/commitizen/cz-conventional-changelog)
* create a pull request that connects to the issue

### how to write commit message 
```
# install commitzen and cz-conventional-changelog first
$ npm install -g commitizen --verbose
$ npm install

# use git cz instead of git commit in the future
# follow the instructions to write commit message
$ git add files && git cz
```

## License

_jackpot-server_ is released under version 3.0 of the [GNU General Public License](COPYING).
