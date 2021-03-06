[![Build Status](https://travis-ci.org/clawio/authentication.svg?branch=master)](https://travis-ci.org/clawio/authentication)
[![GoDoc](https://godoc.org/github.com/clawio/authentication?status.svg)](https://godoc.org/github.com/clawio/authentication)
[![Go Report Card](https://goreportcard.com/badge/github.com/clawio/authentication)](https://goreportcard.com/report/github.com/clawio/authentication)
[![codecov](https://codecov.io/gh/clawio/authentication/branch/master/graph/badge.svg)](https://codecov.io/gh/clawio/authentication)



This repository contains the ClawIO Authentication Service.

This service authenticate users for using other services that require the user to be authenticated.
When the user has provided a valid set of credentials (username and password), the service returns an authentication token to be used in future requests.

Current implementations are as follows:

* Simple: uses a SQL database for persisting users and JWT for tokens.
* Memory: stores users in memory. For testing purposes.
