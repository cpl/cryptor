# Cryptor Contributing Guide

## Welcome

I am very happy to see you've taken the interest to look at the guidelines for contributing to [Cryptor](git.cpl.li/cryptor). The topics covered below are simple guidelines and nothing is set in stone, use your judgment to save everyone time, always ask yourself "is X something I would do?", before proposing it.

Please note this is not a *Development Guide* (an introduction to bits of code and how to handle the project), but a set of instructions and recommendations for getting started with adding changes to the codebase and community. A Developer Guide will be provided on the [wiki](https://github.com/cpl/cryptor/wiki).

## Things you'll need

* [Git](https://git-scm.com), version control
* Copy of the project
  * Using git, `git clone https://git.cpl.li/cryptor`
  * Using go get, `go get cpl.li/go/cryptor`
  * Fork the project on [GitHub](https://git.cpl.li/cryptor)
  * Download source code as [ZIP](https://github.com/cpl/cryptor/archive/master.zip)
* [Go](https://golang.org) programming language
  * Installation instructions, source code and releases for the major operating system can be found [here](https://golang.org/doc/install)
  * Check your version by running `go version` in your terminal
    * Cryptor is being developed to work using the latest version (`v1.12`)
    * Soon, I'll post the oldest Go version Cryptor can run on
  * To get started with Go I recommend the [Go Tour](https://tour.golang.org/), an interactive tutorial
  * You can work from any IDE you prefer, there are no IDE dependencies
    * Certain [IDEs](https://golang.org/doc/editors.html) are recommended (or certain plugins for your favourite IDE)
  * The code style used in all Go source files must respect `vet` and `fmt` standards, see [this](https://golang.org/doc/effective_go.html) for more information
  * [Golang FAQ](https://golang.org/doc/faq)
  * Websites for researching [packages](https://golang.org/pkg/) or [documentation](https://godoc.org)

## Before contributing

* Check the [F.A.Q.](https://github.com/cpl/cryptor/wiki/FAQ) over on the wiki
* If you wish, reach out to the community over on [Discord](https://discord.gg/vGQ76Uz). 
* Check the [issues](https://github.com/cpl/cryptor/issues), [projects](https://github.com/cpl/cryptor/projects) and [PRs](https://github.com/cpl/cryptor/pulls) for any existing materials
* If in doubt, contact me `alexandru@cpl.li` (please use `Cryptor:` in your subject)
* Make sure you have your development setup working and you can run the Cryptor tests (`make test`)

## How to contribute?

### Report bugs and issues

In order to open new issues head over to the [GitHub issues page](https://github.com/cpl/cryptor/issues) and select "New Issue", use the "Bug report" template. Make sure the label `bug` is used and any other relevant labels (see below).

### Propose enhancement or optimisations

Before proposing any enhancement or optimisations have some arguments and *demo/prototype/benchmark/...* to improve the chances of being considered. These proposals are also created on the [GitHub issues page](https://github.com/cpl/cryptor/issues) using the "Enhancement/Optimisations Proposal" template. Make sure the right label is used (`optimisation` or `enhancement`, see explination below)

### Propose new features

Consider the security, privacy and costs of implementing the feature you dream of. After that open a new issue on the [GitHub issues page](https://github.com/cpl/cryptor/issues). Use the "Feature request" template and appropriate labels.

### Propose tests

Feel free to skip opening an issue for adding or chaning tests. You can skip straight to a PR.

### How to PR

Not many rules for Pull Requests. Make sure your version is compiling, respecting all language standards, all tests are passing and code is up to the Style Guidelines. After that open a PR and wait 3-5 days. Longer than that, feel free to `@` mention me or drop me an email.

## Non-code related contributions

### Documentation

New documentation, spell checking, overall layout improvements, etc, are all welcome too. Issues can be opened as `discussions` but I would rather you join the Discord server and discuss ideas openly there.

### Network volunteer

The Cryptor Network will requiere nodes to function and test the initial verions. Any node helps, no matter the hardware, it can be a Raspberry Pi, a supercomputer, a VPS or your home desktop/laptop.

## Continuous Integration

 <a href="https://travis-ci.org/cpl/cryptor">
   <img src="https://img.shields.io/travis/cpl/cryptor/master.svg" alt="Travis CI" />
 </a>
 <a href="https://coveralls.io/github/cpl/cryptor?branch=master">
   <img src="https://img.shields.io/coveralls/github/cpl/cryptor/master.svg" alt="Coverage Status" />
 </a>

The current continuous integration pipeline is simple:
* Push changes to GitHub
* Travis-CI runs tests and pushes code coverage to coveralls
* Coveralls detects any major drop in coverage

## Organisation

### Labels

All issues must be labeled accordingly. You can see the live list of issue labels over [here](https://github.com/cpl/cryptor/labels) with their respective descriptions. If in doubt don't worry, a colaborator will fix any misused labels.

### Projects

[Projects](https://github.com/cpl/cryptor/projects) and [milestones](https://github.com/cpl/cryptor/milestones_ give a sense of progress and direction to the overall Cryptor project. The boards are basic Kanbans with `Backlog` (todo), `In Progress` (doing) and `Done` for organising issues and other notes.

## Styleguide

### Commit Messages

* Pad your commit subject with the package in which the modifications are
  * `p2p/node: extend SetAddr argument validation to cover special case`
  * `foo/bar/zulu: refactor legacy code to v1.2.3 standards`
* Please **do not** include emojis (or any Unicode runes) in commit messages, I wish for the git history to be accessible in any environment
* Limit the commit subject to **72** characters or less
  * Keep the subject simple and clear
  * For any message longer than that use the commit description
    * `git commit -m "Subject string" -m "Description string"`
* Use the present tense and imperative mood
  * `add test for special case`, not `added test for ...`
  * `... to move files ...`, not `... that moves files ...`
* Use `[ci skip]` `[skip ci]` to skip certain CI builds which are mainly documentation or non-sourcecode related, [Travis-CI doc](https://docs.travis-ci.com/user/customizing-the-build/#skipping-a-build)

### Golang Code

[![Go Report Card](https://goreportcard.com/badge/cpl.li/go/cryptor)](https://goreportcard.com/report/cpl.li/go/cryptor)

* Respect `fmt` ([link](https://golang.org/doc/effective_go.html#formatting))
  >Gofmt formats Go programs. It uses tabs for indentation and blanks for
  >alignment. Alignment assumes that an editor is using a fixed-width font.
* Respect `vet`
  > Vet examines Go source code and reports suspicious constructs, such as
  > Printf calls whose arguments do not align with the format string. Vet uses
  > heuristics that do not guarantee all reports are genuine problems, but it
  > can find errors not caught by the compilers.
* Respect `godoc` ([link](https://golang.org/doc/effective_go.html#commentary))
* Respect `gofumpt` ([link](https://github.com/mvdan/gofumpt))
  * An unofficial backwards compatible (with the official `fmt`) formatting guideline

