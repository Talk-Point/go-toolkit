# Go Toolkit

<div align="center">

![Go Version](https://img.shields.io/badge/Go-1.24+-00ADD8?style=for-the-badge&logo=go)
[![License](https://img.shields.io/github/license/Talk-Point/go-toolkit?style=for-the-badge)](LICENSE)
[![Go Report Card](https://goreportcard.com/badge/github.com/Talk-Point/go-toolkit?style=for-the-badge)](https://goreportcard.com/report/github.com/Talk-Point/go-toolkit)

A comprehensive utility library designed to enhance Go development workflows with production-ready tools for CLI applications, web security, event handling, and more.

</div>

## ✨ Features

- **🖥️ CLI Framework** - Build powerful command-line interfaces with nested commands, multiple output formats, and interactive input
- **🔐 CAPTCHA Integration** - Ready-to-use Cloudflare Turnstile CAPTCHA client with testing mode
- **📢 Event System** - Thread-safe signal dispatcher for decoupled component communication
- **🎨 Formatters** - Human-readable time formatting utilities
- **🛠️ Utilities** - Common helper functions for everyday Go development

## 📦 Installation

```bash
go get github.com/Talk-Point/go-toolkit@latest
```

## 🚀 Quick Start

### CLI Framework

Build sophisticated command-line applications with support for nested commands, different output formats, and contextual execution:

```go
package main

import (
    "fmt"
    "github.com/Talk-Point/go-toolkit/pkg/v2/cli"
)

type AppContext struct {
    Config string
    Debug  bool
}

func main() {
    ctx := &AppContext{Config: "app.conf"}
    
    commands := []*cli.Command[*AppContext]{
        {
            Use:   "version",
            Short: "Print the version",
            Run: func(cmd *cli.Command[*AppContext], args []string, ctx *AppContext) (cli.Data, error) {
                return &cli.DataMessage{Message: "v1.0.0"}, nil
            },
        },
        {
            Use:   "users",
            Short: "Manage users",
            Commands: []*cli.Command[*AppContext]{
                {
                    Use:   "list",
                    Short: "List all users",
                    Run: func(cmd *cli.Command[*AppContext], args []string, ctx *AppContext) (cli.Data, error) {
                        return &cli.DataList{
                            Title: "Active Users",
                            Items: []map[string]string{
                                {"id": "1", "name": "Alice", "role": "admin"},
                                {"id": "2", "name": "Bob", "role": "user"},
                            },
                        }, nil
                    },
                },
            },
        },
    }
    
    app := cli.Cli[*AppContext](ctx, commands)
    app.Run()
}
```

### CAPTCHA Integration

Integrate Cloudflare Turnstile CAPTCHA verification with ease:

```go
import "github.com/Talk-Point/go-toolkit/pkg/v2/captcha"

// Production mode with real Turnstile
cap := captcha.NewCaptchaTurnstile(siteKey, secretKey)

// Testing mode for development
cap := captcha.NewCaptchaTesting("test-site-key", "test-secret")

// Verify CAPTCHA token
err := cap.Verify(token, clientIP)
if err != nil {
    // Handle verification failure
}
```

### Event System

Implement loosely coupled communication between components:

```go
import "github.com/Talk-Point/go-toolkit/pkg/v2/signal"

// Create dispatcher
dispatcher := signal.NewSignalDispatcher()

// Register event handlers
dispatcher.Connect("user.created", func(sig signal.Signal, data interface{}) {
    user := data.(*User)
    fmt.Printf("New user created: %s\n", user.Name)
})

dispatcher.Connect("user.created", func(sig signal.Signal, data interface{}) {
    // Send welcome email
})

// Emit events
dispatcher.Emit("user.created", &User{Name: "Alice"})
```

### Time Formatting

Convert time values to human-readable formats:

```go
import "github.com/Talk-Point/go-toolkit/pkg/v2/formatter"

// Format duration
duration := formatter.TimePeriodHumanReadable(3665) // "1h 1m 5s"

// Format relative time
message := formatter.TimeAbsoluteFormatter(
    time.Now().Add(-3*24*time.Hour), 
    time.Now(),
) // "3 days ago"
```

## 📚 Module Reference

### CLI Package (`pkg/v2/cli`)

The CLI package provides a comprehensive framework for building command-line applications:

**Key Components:**
- `Command[T]` - Generic command structure supporting custom contexts
- `CliRoot[T]` - Root CLI handler with automatic help generation
- Multiple `Data` types for structured output (`DataMessage`, `DataList`, `DataDetails`, `DataError`)
- Built-in formatters (JSON, Text)
- Interactive input collection with validation

**Advanced Features:**
```go
// Interactive input with validation
type UserInput struct {
    Name     string `input:"name,required"`
    Email    string `input:"email,required,email"`
    Age      int    `input:"age,min=18,max=100"`
}

var input UserInput
err := cli.Input(&input, os.Args[1:])

// Custom output formatting
app.SetFormatter(&cli.JSONFormatter{})
```

### CAPTCHA Package (`pkg/v2/captcha`)

Provides a unified interface for CAPTCHA verification:

**Features:**
- Support for multiple CAPTCHA providers
- Built-in Cloudflare Turnstile integration
- Testing mode for development
- Thread-safe verification

**Example with HTTP handler:**
```go
func LoginHandler(w http.ResponseWriter, r *http.Request) {
    token := r.FormValue("cf-turnstile-response")
    ip := r.RemoteAddr
    
    if err := captcha.Verify(token, ip); err != nil {
        // Show error - CAPTCHA verification failed
        return
    }
    
    // Continue with login process
}
```

### Signal Package (`pkg/v2/signal`)

A thread-safe event dispatcher for decoupled communication:

**Features:**
- Concurrent callback execution
- Type-safe signal handling
- Simple publish-subscribe pattern
- No external dependencies

**Advanced Usage:**
```go
// Define custom signals
const (
    SignalUserCreated   signal.Signal = "user.created"
    SignalUserDeleted   signal.Signal = "user.deleted"
    SignalOrderComplete signal.Signal = "order.complete"
)

// Type-safe event data
type UserEvent struct {
    UserID   string
    Username string
    Action   string
}

// Register multiple handlers
dispatcher.Connect(SignalUserCreated, auditLogger)
dispatcher.Connect(SignalUserCreated, emailNotifier)
dispatcher.Connect(SignalUserCreated, analyticsTracker)
```

### Formatter Package (`pkg/v2/formatter`)

Utilities for formatting time-related data:

**Functions:**
- `TimePeriodHumanReadable(seconds int32) string` - Formats duration in human-readable format
- `TimeAbsoluteFormatter(date, reference time.Time) string` - Formats relative time

**Examples:**
```go
// Duration formatting
formatter.TimePeriodHumanReadable(90)     // "1m 30s"
formatter.TimePeriodHumanReadable(3600)   // "1h"
formatter.TimePeriodHumanReadable(86400)  // "1d"

// Relative time
now := time.Now()
formatter.TimeAbsoluteFormatter(now.Add(-time.Hour), now)        // "1 hour ago"
formatter.TimeAbsoluteFormatter(now.Add(time.Hour*24*7), now)    // "7 days from now"
```

### Shared Package (`pkg/v2/shared`)

Common utility functions:

**Available Functions:**
- `Contains(element string, slice []string) bool` - Check if string exists in slice

## 🏗️ Architecture

The toolkit follows these design principles:

- **Modularity**: Each package is independent and focused on a specific domain
- **Zero Dependencies**: Minimal external dependencies for maximum compatibility
- **Thread Safety**: Concurrent operations are safe where applicable
- **Generic Support**: Modern Go generics for type-safe, reusable code
- **Interface-Based**: Extensible design through well-defined interfaces

## 🧪 Testing

Run the test suite:

```bash
make test
```

Run tests with coverage:

```bash
go test -v -cover ./...
```

## 🤝 Contributing

We welcome contributions! Please feel free to submit a Pull Request. For major changes, please open an issue first to discuss what you would like to change.

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/AmazingFeature`)
3. Commit your changes (`git commit -m 'Add some AmazingFeature'`)
4. Push to the branch (`git push origin feature/AmazingFeature`)
5. Open a Pull Request

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🏢 About Talk-Point

This toolkit is developed and maintained by [Talk-Point](https://talk-point.de), enhancing our Go development workflow across projects.

---

<div align="center">
Made with ❤️ by Talk-Point
</div>