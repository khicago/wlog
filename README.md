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
- **Multiple Logging Strategies**: Leaf, Branch, and NewTree strategies for different logging needs.
- **Development Mode**: Dedicated logging for development environments.
- **Global and Local Instances**: Support for both application-wide and localized logging.
- **Standardized Syntax**: Clear syntax for describing log usage and function position in call chains.

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

    // Using the builder pattern
    log, newCtx := wlog.B(ctx).
        WithStrategy(wlog.ForkBranch).
        WithFingerPrints("auth", "login").
        WithField("user_id", 12345).
        Build()
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

factory.Common("custom_module").Info("Logging with custom instance")
```

### Fingerprint Chain Management

WLog provides a standardized syntax for managing log chains and function positioning:

```go
// Assuming an existing chain: a/b/c

// Leaf Node: Appends to the chain without modifying context
// Prints: a/b/c/func_name
log := wlog.Leaf(ctx, "func_name")

// Branch Node: Appends to the chain and updates context
// Prints: a/b/c/func_name, and child nodes start from this new chain
log, ctx := wlog.Branch(ctx, "func_name")

// Detach New: Starts a new chain, context is passed but chain is reset
// Prints: func_name, and child nodes start from this new root
log, ctx := wlog.DetachNew(ctx, "func_name")
```

This syntax allows for clear and consistent logging across complex function call hierarchies.

### Context-Aware Logging

WLog provides powerful methods for context-aware logging:

```go
// Simple logging with context
wlog.Leaf(ctx, "service").Info("Processing request")

// Branching: creates a new log entry and updates the context
log, newCtx := wlog.Branch(ctx, "sub_service")
log.Info("Sub-service operation")

// Using the updated context
wlog.Leaf(newCtx, "final_step").Info("Operation completed")
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

    log.WithField("result", result).Info("Request processed successfully")

    wlog.Leaf(ctx, "finish_request").Info("Request handling completed")
}

func processRequest(ctx context.Context) (string, error) {
    wlog.Leaf(ctx, "process_request").Debug("Processing request")
    // ... processing logic ...
    return "result", nil
}
```

## Best Practices

1. **Use Single-Level Fingerprints**: For most cases, use a single fingerprint level for clarity.
2. **Leverage Builder Pattern for Complex Scenarios**: Use the Builder pattern when you need to add fields and customize the logging strategy.
3. **Prefer Convenience Methods for Simple Cases**: Use `Leaf`, `Branch`, or `DetachNew` for straightforward logging.
4. **Consistent Naming for Fingerprints**: Adopt a consistent naming convention, e.g., service or module names.
5. **Utilize Chain Management Appropriately**:
   - Use `Leaf` for appending to the log chain without modifying the context.
   - Use `Branch` when you want to extend the chain and update the context for child operations.
   - Use `DetachNew` when starting a completely new logical section or microservice boundary in your application.
6. **Leverage Dev() for Development-Only Logs**: Use `LDev.Log()` for development-only logs.
7. **Consistent Function Positioning**: Always use the standardized syntax at the beginning of your functions to clearly indicate their position in the call hierarchy.

## Advanced Initialization in Production

```go
func initWLog(fileLogger io.Writer) {
    logrus.SetFormatter(&logrus.JSONFormatter{
        PrettyPrint: true, // For better readability in development
    })
    logrus.SetLevel(logrus.TraceLevel)

    multiWriter := io.MultiWriter(os.Stdout, fileLogger)
    logrus.SetOutput(multiWriter)

    wlog.SetEntryGetter(func(ctx context.Context) *logrus.Entry {
        return logrus.WithContext(ctx)
    })
}
```

## Performance Comparison

While WLog may not always have the fastest raw performance, it offers unique features like chain management and context-aware logging that are not available in other libraries. The performance trade-off is often negligible in real-world scenarios, especially when considering the added functionality and improved log structure.

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
