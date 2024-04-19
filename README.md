# go-toolkit

GoToolkit is a versatile utility library tailored to enhance the Go development workflow at Talk-Point.

## Modules

### CLI Helper

Package cli provides a framework for building command line interfaces with support for different output formats and nested commands. It allows easy creation and management of CLI commands, along with formatting outputs as JSON or plain text. This package supports command hierarchies and contextual execution.

```go
import (
	"context"

	"github.com/Talk-Point/go-toolkit/pkg/v2/cli"
)

cmds := []*cli.Command[*Context]{
    {
        Use:     "version",
        Short:   "Print the version",
        Long:    "Print the version of the CLI",
        Example: "version",
        Run: func(cmd *cli.Command[*Context], args []string, ctx *Context) (cli.Data, error) {
            return &cli.DataMessage{
                Message: "Version: 1.0.0",
            }, nil
        },
        {
			Use:   "users",
			Short: "Manage users",
			Long:  "Manage users in the system",
			Commands: []*cli.Command[*Context]{
				{
					Use:     "list",
					Short:   "List users",
					Long:    "List all users in the system",
					Example: "users list",
					Run: func(cmd *cli.Command[*Context], args []string, ctx *Context) (cli.Data, error) {
                        res, err := fetchUsers(ctx)
                        if err != nil {
							return nil, err
						}
						data := &cli.DataList{
							Title: "Users",
							Items: []map[string]string{},
						}
						for _, user := range res.Items {
							data.Items = append(data.Items, map[string]string{
								"id":    user.Id,
								"email": user.Email,
							})
						}
						return data, nil
					},
				},
				{
					Use:     "create",
					Short:   "Create a user",
					Long:    "Create a user in the system",
					Example: "users create -username max.mustermann -email max.mustermann@talk-point.de",
					Run: func(cmd *cli.Command[*Context], args []string, ctx *Context) (cli.Data, error) {
                        user, err := createUser(ctx, args)
                        if err != nil {
                            return nil, err
                        }
                        return &cli.DataDetails{
							Title: "User",
							Item: map[string]string{
								"id":       user.Id,
								"email":    user.Email,
							},
						}, nil
					},
				},
    },
}

ctx := &Context{
    // the context
}
shell := cli.Cli[*Context](ctx, cmds)
shell.Run()
```

### Captcha

Package captcha provides a client for [Cloudflare Turnstile Captcha](https://www.cloudflare.com/products/turnstile/).

The package allows you to create and verify captchas of different types. Currently, it supports two types of captchas: Turnstile and Testing. The Turnstile captcha is a real captcha that requires verification, while the Testing captcha is a dummy captcha used for testing purposes.

Each captcha is represented by a Captcha struct, which contains the necessary information for captcha verification, such as the site key, secret, and the type of captcha.

The package provides two functions for creating captchas: NewCaptchaTurnstile and NewCaptchaTesting. These functions return a pointer to a Captcha struct with the specified site key, secret, and type.

The Captcha struct has a Verify method that takes a token and an IP address, and verifies the captcha based on its type. If the captcha type is not supported, the method returns an error.

```go
cfToken := c.FormValue("cf-turnstile-response")
ip := c.RealIP()
cap := captcha.NewCaptchaTurnstile(h.CaptchaConfig.SiteKey, h.CaptchaConfig.SecretKey)
if cfToken == "" {
    return Render(c, http.StatusOK, login.LoginForm(cap, h.AuthenticatorId, login.LoginErrors{
        Email:              email,
        InvalidCredentials: "Verification failed",
    }))
}
```

### Signal

Package signal implements a simple event dispatch system that allows components within an application to communicate with each other in a loosely coupled manner. It provides a SignalDispatcher which maintains a mapping of signals (events) to callbacks that are to be executed when a signal is emitted.

A Signal in this context is a string identifier that represents a specific type of event. Callback functions registered to a signal are called with the signal and any accompanying data when that signal is emitted.

```go
dispatcher := signal.NewSignalDispatcher()
dispatcher.Connect("order-created", func(signal signal.Signal, data interface{}) {
    order, ok := data.(*models.Order)
    if !ok {
        fmt.Println("Invalid order data")
        return
    }
})
```