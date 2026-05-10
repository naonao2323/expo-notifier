
## Reference
https://github.com/oliveroneill/exponent-server-sdk-golang/tree/master

## Example

### Basic usage

```go
import (
    "context"
    "log"

    exponotifier "github.com/naonao2323/expo-notifier"
)

func main() {
    ctx := context.Background()

    notifier := exponotifier.NewNotifier()

    token, err := exponotifier.NewExponentPushToken("ExponentPushToken[xxxx]")
    if err != nil {
        log.Fatal(err)
    }

    err = notifier.Add(ctx, exponotifier.PushMessage{
        To:    token,
        Title: "Hello",
        Body:  "This is a notification",
    })
    if err != nil {
        log.Fatal(err)
    }

    notifier.Flush()
}
```

### Buffer settings

```go
import (
    "time"

    exponotifier "github.com/naonao2323/expo-notifier"
)

notifier := exponotifier.NewNotifier(
    exponotifier.WithBuffer(
        exponotifier.WithCountThreshold(50),
        exponotifier.WithDelayThreshold(2*time.Second),
        exponotifier.WithByteThreshold(512*1024), // 512KB
    ),
)
```
