# go-toolkit

GoToolkit is a versatile utility library tailored to enhance the Go development workflow at Talk-Point.

## Modules

### CLI Helper

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

Client for [Cloudflare Turnstile Captcha](https://www.cloudflare.com/products/turnstile/)

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

Dispatcher signal library 

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