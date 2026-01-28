You Only Log.... Once
=====================

Our logging library inspired by [loggingsucks.com](https://loggingsucks.com).

## Usage


```go
logger := yolog.NewLogger("example")
defer logger.Log(os.Stdout)

// Adding attributes
logger.SetAttr("event", "user_signup")

// Adding errors
logger.SetErrorCode(500)
logger.AddError(fmt.Errorf("An error occurred and it was really bad"))

// Demonstrates how a struct can be utilised.
logger.SetAttr("user", struct {
    Name  string `json:"name"`
    Email string `json:"email"`
}{
    Name:  "Jane Doe",
    Email: "jane@example.com",
})

// Sleep so we can see a good elapsed number.
// The logger also prints a elapsed_time_ms field for performance debugging.
time.Sleep(time.Second)
```

## Sampling

* All errors will be logged
* All non error log events can be sampled using the below configuration

```bash
# Log all events
export YOLOG_SAMPLE_RATE=1.0

# Only log half
export YOLOG_SAMPLE_RATE=0.5
```
