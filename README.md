# closer

Processing interrupt functions of the main process of imported packages.

Support for various format completion functions.
Support for tracking keys, which implements the ability to stop the tracking process.
Support for error handling of termination functions.

### Usage example

```go
import (
  "fmt"
  "syscall"
  "time"
  "github.com/ds248a/closer"
)

// An example of functions for completing the main operation of an imported package
type Redis struct {}
type Postgre struct {}
type Cache struct {}

// 1. Value-returning functions
func (c *Redis) Close() error {
  time.Sleep(2 * time.Second)
  return fmt.Errorf("%s", fmt.Sprintln("Redis err"))
}

// 2. Void (nonvalue-returning) functions
func (c *Postgre) Close() {
  time.Sleep(2 * time.Second)
}

// 3. Function with unique name 
func (c *Cache) Clear() {
  time.Sleep(2 * time.Second)
}

func main() {
  cc := &Cache{}
  pg := &Postgre{}
  rd := &Redis{}

  // Execution time limit for all handlers 
  closer.SetTimeout(10 * time.Second)

  // Handler registration, with key saving 
  pKey := closer.Add(pg.Close)

  // Simple handler registration, no key saving 
  closer.Add(cc.Clear)

  // Cancel processing (interrupt tracking) by key 
  closer.Remove(pKey)

  // Handling errors on the application side, with saving the key 
  rKey := closer.Add(func() {
    err := rd.Close()
    if err != nil {
      fmt.Println(err.Error())
    }
  })

  // Out: c159f74d-8a6c-49fd-a181-83edc1d5d595
  fmt.Println(rKey)

  // Resetting all processing jobs 
  closer.Reset()

  // Called at the end of the main()
  closer.ListenSignal(syscall.SIGINT, syscall.SIGQUIT, syscall.SIGTERM)
}

```

### Object and cross-batch processing

```go
import (
  "github.com/ds248a/closer"
)

func main() {
  c := closer.NewCloser()
  c.Add(cc.Clear)
  c.Add(pg.Close)
  c.ListenSignal(syscall.SIGTERM)
}
```
