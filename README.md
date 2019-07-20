``` 
 ██████╗██████╗ ███████╗██████╗  █████╗ ██╗███╗   ██╗
██╔════╝██╔══██╗██╔════╝██╔══██╗██╔══██╗██║████╗  ██║
██║     ██████╔╝█████╗  ██████╔╝███████║██║██╔██╗ ██║
██║     ██╔══██╗██╔══╝  ██╔══██╗██╔══██║██║██║╚██╗██║
╚██████╗██║  ██║███████╗██████╔╝██║  ██║██║██║ ╚████║
╚═════╝╚═╝  ╚═╝╚══════╝╚═════╝ ╚═╝  ╚═╝╚═╝╚═╝  ╚═══╝
```
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
