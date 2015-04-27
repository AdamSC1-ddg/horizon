# Note this code is pre-alpha. 
It is definitely not ready yet for production.
 
# Horizon
[![Build Status](https://travis-ci.org/stellar/go-horizon.svg?branch=master)](https://travis-ci.org/stellar/go-horizon)

Horizon is the [client facing API](http://docs.stellarhorizon.apiary.io) server for the Stellar ecosystem.  See [an overview of the Stellar ecosystem](https://www.stellar.org/galaxy/getting-started/) for more details.

## Installing Dependencies

After cloning go-horizon into `$GOPATH/src/github.com/stellar/go-horizon`, cd into that directory and run `go get -t ./...`

## Running Tests

go-horizon uses [GoConvey](https://github.com/smartystreets/goconvey) for testing.  If you are going to use the `goconvey` tool for running your tests, you must ensure that package-based parallelism is turned off.  By default, GoConvey will run several packages test suites in parallel, but since go-horizon's constituent packages all write to a common database problems can arise.  You should launch `goconvey` like so:

```bash
goconvey -packages=1
```