# Machine Box Go SDK

The official Machine Box Go SDK provides Go clients for each box.

## Usage

Go get the repo:

```
go get github.com/machinebox/mb/exp/sdk-go
```

Then import the package of the box you wish to use:

```go
import "github.com/machinebox/mb/exp/sdk-go/facebox"
```

Then create a client, providing the address of the running box.

(To get a box running locally, see the instructions at https://machinebox.io/account)

```go
faceboxClient := facebox.New("http://localhost:8080")
```

It is recommended that you consider the startup time a box needs before it
is ready. The simplest approach is to use the boxutil.WaitForReady function:

```go
err := boxutil.WaitForReady(ctx, faceboxClient)
if err != nil {
    log.Fatalln("error waiting for box:", err)
}
```

A more advanced solution is to get notified whenever the status of a box changes
using the boxutil.StatusChan feature:

```go
go func(){
    statusChan := boxutil.StatusChan(ctx, faceboxClient)
    for {
        select {
        case status := <-statusChan:
            if !boxutil.IsReady(status) {
                log.Println("TODO: Pause work, the box isn't ready")
            } else {
                log.Println("TODO: resume work, the box is ready to go")
            }
        }
    }
}()
```
