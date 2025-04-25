// Package cli provides a framework for building command line interfaces
// with support for different output formats and nested commands.
// It allows easy creation and management of CLI commands, along with formatting
// outputs as JSON or plain text. This package supports command hierarchies and
// contextual execution.
//
// Example:
//
//	cmds := []*cli.Command[*Context]{
//	    {
//	        Use:     "version",
//	        Short:   "Print the version",
//	        Long:    "Print the version of the CLI",
//	        Example: "version",
//	        Run: func(cmd *cli.Command[*Context], args []string, ctx *Context) (cli.Data, error) {
//	            return &cli.DataMessage{
//	                Message: "Version: 1.0.0",
//	            }, nil
//	        },
//	    },
//	    {
//	        Use:   "users",
//	        Short: "Manage users",
//	        Long:  "Manage users in the system",
//	        Commands: []*cli.Command[*Context]{
//	            {
//	                Use:     "list",
//	                Short:   "List users",
//	                Long:    "List all users in the system",
//	                Example: "users list",
//	                Run: func(cmd *cli.Command[*Context], args []string, ctx *Context) (cli.Data, error) {
//	                    res, err := fetchUsers(ctx)
//	                    if err != nil {
//	                        return nil, err
//	                    }
//	                    data := &cli.DataList{
//	                        Title: "Users",
//	                        Items: []map[string]string{},
//	                    }
//	                    for _, user := range res.Items {
//	                        data.Items = append(data.Items, map[string]string{
//	                            "id":    user.Id,
//	                            "email": user.Email,
//	                        })
//	                    }
//	                    return data, nil
//	                },
//	            },
//	            {
//	                Use:     "create",
//	                Short:   "Create a user",
//	                Long:    "Create a user in the system",
//	                Example: "users create -username max.mustermann -email max.mustermann@talk-point.de",
//	                Run: func(cmd *cli.Command[*Context], args []string, ctx *Context) (cli.Data, error) {
//	                    user, err := createUser(ctx, args)
//	                    if err != nil {
//	                        return nil, err
//	                    }
//	                    return &cli.DataDetails{
//	                        Title: "User",
//	                        Item: map[string]string{
//	                            "id":      user.Id,
//	                            "email":   user.Email,
//	                        },
//	                    }, nil
//	                },
//	            },
//	        },
//	    },
//	}
//	ctx := &Context{}
//	shell := cli.Cli[*Context](ctx, cmds)
//	shell.Run()
package cli

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

// Format formats the given data and returns a string representation.
type Formatter interface {
	// Format formats the given data and returns a string representation.
	Format(data interface{}) (string, error)
	// Type returns the type of the formatter.
	Type() string
}

// JSONFormatter implements Formatter to output data in JSON format.
type JSONFormatter struct{}

func (j *JSONFormatter) Format(data interface{}) (string, error) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return "", err
	}
	return string(jsonData), nil
}

func (j *JSONFormatter) Type() string {
	return "json"
}

type TextFormatter struct{}

func (t *TextFormatter) Format(data interface{}) (string, error) {
	return fmt.Sprintf("%v", data), nil
}

func (t *TextFormatter) Type() string {
	return "text"
}

// Data is an interface for types that can be displayed using a Formatter.
// It requires a Display method that uses the provided formatter to create
// a string representation of the data.
type Data interface {
	Display(formatter Formatter) (string, error)
}

// DataMessage holds a simple text message. It is used to encapsulate a message
// that can be formatted and displayed.
type DataMessage struct {
	Message string `json:"message"`
}

// Display simply returns the message as a string without further formatting,
// regardless of the formatter used. This ensures that the message is displayed
// as intended without modification.
func (d *DataMessage) Display(formatter Formatter) (string, error) {
	return d.Message, nil
}

// DataList represents a structured list of items, each being a map of strings.
// It is typically used to present a collection of similar data objects.
type DataList struct {
	Title string              `json:"title"`
	Items []map[string]string `json:"items"`
}

func (d *DataList) Display(formatter Formatter) (string, error) {
	return formatter.Format(d)
}

func (d *DataList) Error() string {
	a := []string{}

	a = append(a, d.Title)

	for _, item := range d.Items {
		for k, v := range item {
			a = append(a, fmt.Sprintf("%s: %s", k, v))
		}
	}

	return strings.Join(a, "\n")
}

// DataDetails holds detailed information about a single item, typically used
// for displaying detailed views of a specific entity.
type DataDetails struct {
	Title string            `json:"title"`
	Item  map[string]string `json:"item"`
}

func (d *DataDetails) Display(formatter Formatter) (string, error) {
	return formatter.Format(d)
}

func (d *DataDetails) Error() string {
	a := []string{}
	a = append(a, d.Title)
	for k, v := range d.Item {
		a = append(a, fmt.Sprintf("%s: %s", k, v))
	}
	return strings.Join(a, "\n")
}

// DataError is used to represent errors as data. This allows error messages to be formatted
// and displayed using the same mechanisms as other data types.
type DataError struct {
	Message string `json:"error"`
}

func (d *DataError) Error() string {
	return d.Message
}

func (d *DataError) Display(formatter Formatter) (string, error) {
	return formatter.Format(d)
}

type Command[T any] struct {
	Use      string
	Short    string
	Long     string
	Run      func(cmd *Command[T], args []string, ctx T) (Data, error)
	Commands []*Command[T]
	Example  string
}

type CliRoot[T any] struct {
	Ctx       T
	Commands  []*Command[T]
	Formatter Formatter
}

func (c *CliRoot[T]) Run() {
	data, err := c.runCommand(c.Commands, os.Args[1:])
	if err != nil {
		data := &DataError{
			Message: err.Error(),
		}
		v, err := data.Display(c.Formatter)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		fmt.Fprintln(os.Stderr, v)
		os.Exit(1)
	}
	if data != nil {
		v1, _ := data.Display(c.Formatter)
		fmt.Println(v1)
	}
}

func (c *CliRoot[T]) RunWithCommand(command string) (Data, error) {
	commandArgs := strings.Fields(command)
	return c.runCommand(c.Commands, commandArgs)
}

func (c *CliRoot[T]) runCommand(commands []*Command[T], args []string) (Data, error) {
	filteredArgs := []string{}
	for _, arg := range args {
		if !strings.HasPrefix(arg, "-json") && !strings.HasPrefix(arg, "--json") {
			filteredArgs = append(filteredArgs, arg)
		} else {
			if arg == "-json" || arg == "--json" {
				c.Formatter = &JSONFormatter{}
			}
		}
	}

	if len(filteredArgs) == 0 {
		return c.Help(commands)
	}
	// check if first argument is -help
	if filteredArgs[0] == "-help" || filteredArgs[0] == "--help" {
		return c.Help(commands)
	}

	for _, cmd := range commands {
		if cmd.Use == filteredArgs[0] {
			if cmd.Commands == nil {
				data, err := cmd.Run(cmd, filteredArgs[1:], c.Ctx)
				return data, err
			} else {
				return c.runCommand(cmd.Commands, filteredArgs[1:])
			}
		}
	}

	return nil, fmt.Errorf("command %s not found", filteredArgs[0])
}

func (c *CliRoot[T]) Help(commands []*Command[T]) (Data, error) {
	if c.Commands == nil {

		return &DataMessage{
			Message: "No commands found",
		}, nil
	}
	data := &DataList{
		Title: "Available commands",
		Items: []map[string]string{},
	}

	for _, cmd := range commands {
		data.Items = append(data.Items, map[string]string{
			"Use":   cmd.Use,
			"Short": cmd.Short,
		})
	}

	return data, nil

}

func Cli[T any](ctx T, cmds []*Command[T]) *CliRoot[T] {
	return &CliRoot[T]{
		Ctx:       ctx,
		Commands:  cmds,
		Formatter: &TextFormatter{},
	}
}
