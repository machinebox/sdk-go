# Machine Box Go SDK [![GoDoc](https://godoc.org/github.com/machinebox/sdk-go?status.svg)](http://godoc.org/github.com/machinebox/sdk-go) [![Build Status](https://travis-ci.org/machinebox/sdk-go.svg?branch=master)](https://travis-ci.org/machinebox/sdk-go)

The official Machine Box Go SDK provides Go clients for each box.

## Usage

Go get the repo:

```
go get github.com/machinebox/sdk-go
```

Then import the package of the box you wish to use:

```go
import "github.com/machinebox/sdk-go/facebox"
```

Then create a client, providing the address of the running box.

(To get a box running locally, see the instructions at https://machinebox.io/account)

```go
faceboxClient := facebox.New("http://localhost:8080")
```
