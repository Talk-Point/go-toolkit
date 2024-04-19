package cli

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

type Formatter interface {
	Format(data interface{}) (string, error)
	Type() string
}

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

type Data interface {
	Display(formatter Formatter) (string, error)
}

type DataMessage struct {
	Message string `json:"message"`
}

func (d *DataMessage) Display(formatter Formatter) (string, error) {
	return d.Message, nil
}

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

	return nil, fmt.Errorf("command " + filteredArgs[0] + " not found")
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
