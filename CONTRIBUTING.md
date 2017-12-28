# Cryptor Contribution Guide

## Welcome
I am very happy to see you're intrested in developing Cryptor. We are looking for
any passionate developer (backend, frontend, networking, design, pr, etc...), every
little contribution helps.

[![GitHub forks](https://img.shields.io/github/forks/thee-engineer/cryptor.svg?style=social&label=Fork&maxAge=2592000)](https://github.com/thee-engineer/cryptor/)

## Workflow
0. Start by [installing Go](https://golang.org/doc/install) and getting familiar with the syntax
    * Run `go version` first, who knows you might already have it installed.
    * We recommend Go 1.5 or higher
1. Now in order to get used to the codebase you could get a copy of the project in one of the following ways:
    * Fork the project to your profile (best option for contributing)
        * Then clone it to your machine
    * Clone the project using `git`
        * `git clone git@github.com:thee-engineer/cryptor.git`
        * `git clone https://github.com/thee-engineer/cryptor.git`
    * Obtain cryptor using `go`
        * `go get github.com/thee-engineer/cryptor`
2. Read the documentation
    * Have a look at [Cryptor's documentation](https://godoc.org/github.com/thee-engineer/cryptor)
    * Explore the code and read the comments
3. Contribute
    * Add a new feature (please open an issue with the `feature` tag before starting work)
    * Look over the [issues](https://github.com/thee-engineer/cryptor/issues)
    * Open new issues (respect the [issue template](https://github.com/thee-engineer/cryptor/blob/master/.github/ISSUE_TEMPLATE.md))
4. Commit your changes & submit a [pull request](https://help.github.com/articles/about-pull-requests/)
    * Before doing so, please run the test suite using `make test` and `make cover && make view`
    * If your request contains failed tests it will be ignored
    * Respect the [pull request template](https://github.com/thee-engineer/cryptor/blob/master/.github/PULL_REQUEST_TEMPLATE.md)

## Communication
* [Telegram](https://t.me/joinchat/Gheirw6fIN_dDFwTAcwURA)
* Slack (Upcoming)
* IRC (Upcoming)
* Gitter (Upcoming)

## Resources
  * [Go packages](https://golang.org/pkg/) always check here for standard packages and documentation
  * [Go Doc](https://godoc.org/) package search and documentation
    * [Cryptor's documentation](https://godoc.org/github.com/thee-engineer/cryptor)
  * Go Tutorials, Introduction and Standards ([here](https://golang.org/doc/))
    * [Installing Go](https://golang.org/doc/install) for all platforms
    * [Tour of Go](https://tour.golang.org/) is an interactive introduction tutorial
    * [IDEs for Go](https://golang.org/doc/editors.html) or packages for your IDE
    * [Go Standards](https://golang.org/doc/effective_go.html) must be used in all Go code, no exceptions
    * [Go FAQ](https://golang.org/doc/faq) some interesting reads
  * Continuous integration used by Cryptor
    * [Travis CI](https://travis-ci.org/thee-engineer/cryptor) runs the test suite on multiple Go versions
    * [Codecov](https://codecov.io/gh/thee-engineer/cryptor) provides code coverage
