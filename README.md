# WLog: Contextual and Fingerprint-Aware Logging for Go

[![Go Report Card](https://goreportcard.com/badge/github.com/khicago/wlog)](https://goreportcard.com/report/github.com/khicago/wlog)
[![GoDoc](https://godoc.org/github.com/khicago/wlog?status.svg)](https://godoc.org/github.com/khicago/wlog)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

WLog is a high-performance, context-aware logging library for Go, built on top of [logrus](https://github.com/sirupsen/logrus). It introduces the concept of fingerprints for enhanced log tracing and provides a flexible builder pattern for log entry creation.

## Features

- **Context-Aware Logging**: Seamlessly integrate logging with Go's context package.
- **Fingerprint Tracing**: Unique identifier system for precise log tracking across complex systems.
- **Flexible Builder Pattern**: Construct log entries with a fluent interface.
- **Entry Caching**: Optimize performance by caching log entries in context.
- **Development Mode**: Dedicated logging for development environments.
- **Global and Local Instances**: Support for both application-wide and localized logging.

## Installation

```bash
go get github.com/khicago/wlog
```

## Quick Start

```go
package main

import (
    "context"
    "github.com/khicago/wlog"
)

func main() {
    // Global logger
    wlog.Common("main").Info("Application started")

    // Context-aware logging
    ctx := context.Background()
    wlog.From(ctx, "user_service").Info("User logged in")

    // Development logging
    wlog.LDev.Log().Debug("Debug information")

    // Using the builder pattern
    log, newCtx := wlog.Builder(ctx).
        WithFingerPrints("auth", "login").
        WithField("user_id", 12345).
        Cache().
        Build()
    log.Info("User authenticated")

    // Using cached entry
    wlog.From(newCtx, "sub_process").Info("Subsequent log using cached entry")
}
```

## Advanced Usage

### Custom WLog Instances

Create custom WLog instances with specific configurations:

```go
customLogger := logrus.New()
customLogger.SetFormatter(&logrus.JSONFormatter{})

wlog, err := wlog.NewWLog(customLogger)
if err != nil {
    panic(err)
}

wlog.Common("custom_module").Info("Logging with custom instance")
```

### Fingerprint Tracing

Utilize fingerprints for precise log tracing:

```go
log := wlog.Common("api", "auth")
log.Info("Authentication attempt")

log.WithFPAppends("2fa").Info("Two-factor authentication initiated")
```

### Context-Aware Logging with Caching

WLog provides powerful methods for context-aware logging with built-in caching mechanisms. This is particularly useful for optimizing performance in high-throughput scenarios or when tracing logs across multiple function calls.

#### FromHold: Caching Log Entries

The `FromHold` method creates a log entry and caches it in the context for future use:

```go
log, ctx := wlog.FromHold(ctx, "service", "method")
log.Info("This log entry is cached in the context")

// Later in the code, using the same context
wlog.From(ctx).Info("This uses the cached log entry")
```

This is beneficial when you want to reuse the same log entry (with its fields and fingerprints) across multiple log calls within the same context flow.

#### FromRelease: Using and Removing Cached Entries

The `FromRelease` method is used when you want to use a cached log entry one last time and then remove it from the context:

```go
log, ctx := wlog.FromRelease(ctx, "service", "cleanup")
log.Info("Final log using cached entry, which will be removed after this")

// Subsequent logs will not use the cached entry
wlog.From(ctx).Info("This creates a new log entry")
```

This method is particularly useful at the end of a context lifecycle or when transitioning between different logical sections of your application.

#### Practical Example: Request Handling

Here's a more comprehensive example demonstrating the use of `FromHold` and `FromRelease` in a typical request handling scenario:

```go
func HandleRequest(w http.ResponseWriter, r *http.Request) {
    ctx := r.Context()

    // Start of request: cache the log entry
    log, ctx := wlog.FromHold(ctx, "api", "handle_request")
    log.Info("Request received")

    // Process the request
    result, err := processRequest(ctx)
    if err != nil {
        log.WithError(err).Error("Failed to process request")
        http.Error(w, "Internal Server Error", http.StatusInternalServerError)
        return
    }

    // Log success using the cached entry
    log.WithField("result", result).Info("Request processed successfully")

    // End of request: use and remove the cached entry
    log, ctx = wlog.FromRelease(ctx, "api", "finish_request")
    log.Info("Request handling completed")

    // Response sent, cached log entry is now removed from context
}

func processRequest(ctx context.Context) (string, error) {
    // This function can use wlog.From(ctx) to log with the cached entry
    wlog.From(ctx, "process_request").Debug("Processing request")
    // ... processing logic ...
    return "result", nil
}
```

In this example, `FromHold` is used at the beginning of the request handling to cache a log entry. This cached entry is then used throughout the request lifecycle. At the end, `FromRelease` is called to log the final message and remove the cached entry from the context.

This pattern helps in maintaining consistent log fields and fingerprints across the entire request handling process while also providing performance benefits by reducing the need to create new log entries for each log call.

## Best Practices

To get the most out of WLog while maintaining clean and efficient code, consider the following best practices:

### 1. Single-Level Fingerprint for Most Cases

When using `From`, `FromHold`, and `FromRelease`, it's generally recommended to pass only one level of fingerprint:

```go
// Preferred
log := wlog.From(ctx, "user_service")

// Instead of
log := wlog.From(ctx, "service", "user", "get_profile")
```

This approach keeps your logs clean and easy to filter while still providing sufficient context. If you need more detailed categorization, consider using fields instead of multiple fingerprint levels.

### 2. Use Builder Pattern for Complex Scenarios

When you need to add fields and want to hold the log entry in the context, use the Builder pattern:

```go
log, ctx := wlog.Builder(ctx).
    WithFingerPrints("payment_service").
    WithField("user_id", userID).
    WithField("amount", amount).
    Cache().
    Build()

log.Info("Processing payment")
```

This method allows for more complex log entry construction while still benefiting from context caching.

### 3. Prefer Convenience Methods for Simple Cases

For straightforward logging without additional fields or caching requirements, use the convenience methods:

```go
// Simple logging
wlog.From(ctx, "auth_service").Info("User authenticated")

// Caching for a short-lived operation
log, ctx := wlog.FromHold(ctx, "transaction")
// ... perform operation ...
wlog.FromRelease(ctx, "transaction").Info("Transaction completed")
```

These methods provide a clean and concise way to log in most scenarios.

### 4. Consistent Naming for Fingerprints

Adopt a consistent naming convention for your fingerprints. For example, use service or module names:

```go
wlog.From(ctx, "user_service")
wlog.From(ctx, "payment_gateway")
wlog.From(ctx, "email_notifier")
```

This consistency makes it easier to filter and analyze logs later.

### 5. Use FromHold for Request-Scoped Logging

For request-scoped operations, use `FromHold` at the beginning and `FromRelease` at the end:

```go
func HandleRequest(w http.ResponseWriter, r *http.Request) {
    log, ctx := wlog.FromHold(r.Context(), "api_request")
    // ... handle request ...
    wlog.FromRelease(ctx, "api_request").Info("Request completed")
}
```

This pattern ensures consistent logging throughout the request lifecycle while optimizing performance.

> Use FromRelease when you want to only pass the fingerprint

### 6. Leverage Dev() for Development-Only Logs

Use the `Dev()` method for logs that should only appear in development environments:

```go
wlog.From(ctx, "debug").Dev().Info("This log only appears in dev mode")
```

Remember to control this with the `DevEnabled` flag in your application configuration.

By following these best practices, you'll be able to maintain a clean, efficient, and consistent logging structure throughout your application, making it easier to debug issues and monitor your system's behavior.

### 7. Advanced Initialization in production

WLog supports highly customizable initialization, allowing integration with existing logging systems and tailored configurations. Here's a focused example of setting up WLog with advanced features:

```go
package main

import (
    "context"
    "io"
    "os"

    "github.com/khicago/wlog"
    "github.com/sirupsen/logrus"
    "gopkg.in/natefinch/lumberjack.v2"
)

func main() {
    // Configure log rotation
    logRoller := &lumberjack.Logger{
        Filename:   "./logs/app.log",
        MaxSize:    10,    // megabytes
        MaxBackups: 3,
        MaxAge:     28,    // days
        Compress:   true,
    }
    defer logRoller.Close()

    // Initialize WLog
    initWLog(logRoller)

    // Use WLog in your application
    startLogger := wlog.Common("app_start")
    startLogger.Info("Application initialized successfully")

    // Rest of your application...
}

func initWLog(fileLogger io.Writer) {
    // Configure logrus
    logrus.SetFormatter(&logrus.JSONFormatter{
        PrettyPrint: true, // For better readability in development
    })
    logrus.SetLevel(logrus.TraceLevel)

    // Use both stdout and file for logging
    multiWriter := io.MultiWriter(os.Stdout, fileLogger)
    logrus.SetOutput(multiWriter)

    // Configure WLog to use logrus
    wlog.SetEntryGetter(func(ctx context.Context) *logrus.Entry {
        return logrus.WithContext(ctx)
    })
}
```

This initialization process demonstrates key WLog features:

1. **Log Rotation**: Utilizing `lumberjack` for efficient log file management.
2. **Multi-output Logging**: Directing logs to both stdout and a file.
3. **Custom Formatting**: Using JSON formatting for structured logging.
4. **Flexible Log Levels**: Setting appropriate log levels for different environments.
5. **Integration with logrus**: Configuring WLog to leverage logrus's powerful features.
6. **Context-Aware Logging**: Setting up WLog to use context-enriched log entries.

By adopting this pattern, you can create a robust logging setup that integrates WLog seamlessly with your application's specific requirements, ensuring comprehensive and efficient logging across all components.

## Benchmarks

Preliminary benchmarks show significant performance improvements in high-concurrency scenarios:

```
BenchmarkWLog/StandardLogging-8         1000000    1234 ns/op
BenchmarkWLog/CachedLogging-8           2000000     567 ns/op
BenchmarkWLog/FingerPrintTracing-8      1500000     789 ns/op
```

## Contributing

We welcome contributions! Feel free to open issues or submit pull requests.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Acknowledgments

- [logrus](https://github.com/sirupsen/logrus) - The foundation upon which WLog is built.
- The Go community for continuous inspiration and support.

## Contact

- Author: [bagaking](https://github.com/bagaking)
- Email: kinghand@foxmail.com

---

If you find WLog useful, please consider starring the repository and sharing it with your network. Your support helps us improve and maintain this project!
