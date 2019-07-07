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
