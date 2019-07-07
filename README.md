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

[![asciicast](https://asciinema.org/a/INzHWa9uQe9ASeNhiGI4T5WBP.svg)](https://asciinema.org/a/INzHWa9uQe9ASeNhiGI4T5WBP)

## Installing

```
$ go get github.com/ricardomaraschini/crebain/cmd/crebain
```

## Running

```
$ crebain --path=/path/to/my/go/app
```

## Command line options

```
$ crebain --path=/path/to/my/go/app --exclude=.git --exclude=vendor
```
