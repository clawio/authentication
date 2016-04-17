[![Build Status](https://drone.io/github.com/clawio/authentication/status.png)](https://drone.io/github.com/clawio/authentication/latest)
[![GoDoc](https://godoc.org/github.com/clawio/authentication?status.svg)](https://godoc.org/github.com/clawio/authentication)
[![Go Report Card](https://goreportcard.com/badge/github.com/clawio/authentication)](https://goreportcard.com/report/github.com/clawio/authentication)
[![codecov.io](https://codecov.io/github/clawio/authentication/coverage.svg?branch=master)](https://codecov.io/github/clawio/authentication?branch=master)

This repository contains the ClawIO Authentication Service.

This service authenticate users for using other services that require the user to be authenticated.
When the user has provided a valid set of credentials (username and password), the service returns an authentication token to be used in future requests.

Current implementations are as follows:

* Simple: uses a SQLite3 database for persisting users and JWT for tokens.
