# WLog: Advanced Contextual Logging for Go

[![Go Report Card](https://goreportcard.com/badge/github.com/khicago/wlog)](https://goreportcard.com/report/github.com/khicago/wlog)
[![GoDoc](https://godoc.org/github.com/khicago/wlog?status.svg)](https://godoc.org/github.com/khicago/wlog)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

WLog is a high-performance, context-aware logging library for Go, built on top of [logrus](https://github.com/sirupsen/logrus). It introduces advanced concepts like fingerprint chains and column-based fields for enhanced log tracing and structured logging.

## Features

- **Context-Aware Logging**: Seamless integration with Go's context package.
- **Fingerprint Chain**: Hierarchical identifier system for precise log tracking across complex systems.
- **Column-Based Fields**: Efficient storage and manipulation of log fields.
- **Flexible Builder Pattern**: Construct log entries with a fluent interface.
- **Multiple Logging Strategies**: Leaf, Branch, and Detach strategies for different logging needs.
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
    wlog.Leaf(ctx, "user_service").Info("User logged in")

    // Development logging
    wlog.LDev.Log().Debug("Debug information")

    // Using the builder pattern with fields
    log, newCtx := wlog.By(ctx, "auth").
        Field("user_id", 12345).
        Branch()
    log.Info("User authenticated")

    // Using updated context
    wlog.Leaf(newCtx, "sub_process").Info("Subsequent log using updated context")
}
```

## Advanced Usage

### Custom Factory Instances

Create custom Factory instances with specific configurations:

```go
customLogger := logrus.New()
customLogger.SetFormatter(&logrus.JSONFormatter{})

factory, err := wlog.NewFactory(customLogger)
if err != nil {
    panic(err)
}

factory.NewBuilder(context.Background()).Name("custom_module").Leaf().Info("Logging with custom instance")
```

### Fingerprint Chain Management

WLog provides three strategies for managing log chains:

1. **Leaf**: Appends to the chain without modifying context.
2. **Branch**: Appends to the chain and updates context.
3. **Detach**: Starts a new chain, resetting the context.

```go
// Leaf: Appends to the chain without modifying context
log := wlog.Leaf(ctx, "func_name")
log.Info("Leaf log entry")

// Branch: Appends to the chain and updates context
log, newCtx := wlog.Branch(ctx, "func_name")
log.Info("Branch log entry")

// Detach: Starts a new chain, resetting the context
log, newCtx := wlog.Detach(ctx, "func_name")
log.Info("Detached log entry")
```

### Adding Fields to Log Entries

When you need to add fields to your log entries, use the `By` method with the builder pattern:

```go
log := wlog.By(ctx, "payment_service").
    Field("user_id", userID).
    Field("amount", amount).
    Leaf()
log.Info("Processing payment")
```

### Practical Example: Request Handling

```go
func HandleRequest(w http.ResponseWriter, r *http.Request) {
    log, ctx := wlog.Branch(r.Context(), "api_request")
    log.Info("Request received")

    result, err := processRequest(ctx)
    if err != nil {
        log.WithError(err).Error("Failed to process request")
        http.Error(w, "Internal Server Error", http.StatusInternalServerError)
        return
    }

    wlog.By(ctx, "finish_request").Field("result", result).Leaf().Info("Request processed successfully")
}

func processRequest(ctx context.Context) (string, error) {
    wlog.Leaf(ctx, "process_request").Debug("Processing request")
    // ... processing logic ...
    return "result", nil
}
```

## Best Practices

1. **Use Single-Level Fingerprints**: For most cases, use a single fingerprint level for clarity.
2. **Leverage Builder Pattern for Fields**: Use `wlog.By(ctx, "name").Field(key, value)...` when adding fields.
3. **Choose Appropriate Strategy**:
   - Use `Leaf` for simple logging without context updates.
   - Use `Branch` when you want to update the context for child operations.
   - Use `Detach` when starting a new logical section in your application.
4. **Consistent Naming for Fingerprints**: Adopt a consistent naming convention, e.g., service or module names.
5. **Utilize Development Logging**: Use `LDev.Log()` for development-only logs.

## Contributing

We welcome contributions! Please feel free to submit issues or pull requests.

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