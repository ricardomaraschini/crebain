[![CircleCI](https://circleci.com/gh/ricardomaraschini/crebain/tree/master.svg?style=svg)](https://circleci.com/gh/ricardomaraschini/crebain/tree/master)
[![Go Report Card](https://goreportcard.com/badge/github.com/axcdnt/snitch)](https://goreportcard.com/report/github.com/axcdnt/snitch)

![crebain](https://raw.githubusercontent.com/ricardomaraschini/crebain/master/logo/crebain.png)


## About

Crebain watches a directory and automatically executes `go test` everytime a
go file within the directory(or subdirectories) changes.

## Demo

Text user interface
[![asciicast](https://asciinema.org/a/258335.svg)](https://asciinema.org/a/258335)

Regular/Simple interface
[![asciicast](https://asciinema.org/a/INzHWa9uQe9ASeNhiGI4T5WBP.svg)](https://asciinema.org/a/INzHWa9uQe9ASeNhiGI4T5WBP)

## Installing

```
$ go get github.com/ricardomaraschini/crebain/cmd/crebain
```

## Running

```
$ crebain --path=/path/to/my/go/app --tui
```

## Command line options

```
$ crebain -h
Usage of crebain:
  -exclude value
        regex rules for excluding paths from watching (default ^\.)
  -path string
        the path to be watched (default "current directory")
  -tui
        enable text user interface
```

## Navigating on TUI

Use keys `j`, `k` to select a test result. To navigate on the test result, use `J`, `K`, `L` and `H` keys.
