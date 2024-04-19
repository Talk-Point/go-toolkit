package cli

import (
	"testing"
)

type Context struct{}

func TestRunCommand(t *testing.T) {
	t.Run("Eaqsy", func(t *testing.T) {
		ctx := &Context{}
		cmds := []*Command[*Context]{
			{
				Use: "version",
				Run: func(cmd *Command[*Context], args []string, ctx *Context) (Data, error) {
					return &DataMessage{
						Message: "Version: 1.0.0",
					}, nil
				},
			},
		}

		c := Cli[*Context](ctx, cmds)
		c.RunWithCommand("version")
	})

	t.Run("Subcommands", func(t *testing.T) {
		ctx := &Context{}
		cmds := []*Command[*Context]{
			{
				Use: "version",
				Run: func(cmd *Command[*Context], args []string, ctx *Context) (Data, error) {
					return &DataMessage{
						Message: "Version: 1.0.0",
					}, nil
				},
			},
			{
				Use: "sub",
				Commands: []*Command[*Context]{
					{
						Use: "version",
						Run: func(cmd *Command[*Context], args []string, ctx *Context) (Data, error) {
							return &DataMessage{
								Message: "Sub Version: 1.0.0",
							}, nil
						},
					},
				},
			},
		}

		c := Cli[*Context](ctx, cmds)
		c.RunWithCommand("sub version")
	})

	t.Run("Help", func(t *testing.T) {
		ctx := &Context{}
		cmds := []*Command[*Context]{
			{
				Use: "version",
				Run: func(cmd *Command[*Context], args []string, ctx *Context) (Data, error) {
					return &DataMessage{
						Message: "Version: 1.0.0",
					}, nil
				},
			},
		}

		c := Cli[*Context](ctx, cmds)
		c.RunWithCommand("help")
	})

	t.Run("CommandNotFound", func(t *testing.T) {
		ctx := &Context{}
		cmds := []*Command[*Context]{
			{
				Use: "version",
				Run: func(cmd *Command[*Context], args []string, ctx *Context) (Data, error) {
					return &DataMessage{
						Message: "Version: 1.0.0",
					}, nil
				},
			},
		}

		c := Cli[*Context](ctx, cmds)
		c.RunWithCommand("notfound")
	})

	t.Run("JSON", func(t *testing.T) {
		ctx := &Context{}
		cmds := []*Command[*Context]{
			{
				Use: "version",
				Run: func(cmd *Command[*Context], args []string, ctx *Context) (Data, error) {
					return &DataMessage{
						Message: "Version: 1.0.0",
					}, nil
				},
			},
		}

		c := Cli[*Context](ctx, cmds)
		c.RunWithCommand("version --json")
	})
}

func TestFormatter(t *testing.T) {
	t.Run("Text", func(t *testing.T) {
		f := &TextFormatter{}
		data := &DataMessage{
			Message: "Version: 1.0.0",
		}
		f.Format(data)
	})

	t.Run("JSON", func(t *testing.T) {
		f := &JSONFormatter{}
		data := &DataMessage{
			Message: "Version: 1.0.0",
		}
		f.Format(data)
	})
}

func TestData(t *testing.T) {
	t.Run("Message", func(t *testing.T) {
		data := &DataMessage{
			Message: "Version: 1.0.0",
		}
		data.Display(&TextFormatter{})
		data.Display(&JSONFormatter{})
	})

	t.Run("Error", func(t *testing.T) {
		data := &DataError{
			Message: "Error: Something went wrong",
		}
		data.Display(&TextFormatter{})
		data.Display(&JSONFormatter{})
	})

	t.Run("DataDetails", func(t *testing.T) {
		data := &DataDetails{
			Title: "Details",
			Item:  map[string]string{"key": "value"},
		}
		data.Display(&TextFormatter{})
		data.Display(&JSONFormatter{})
	})

	t.Run("DataList", func(t *testing.T) {
		data := &DataList{
			Title: "List",
			Items: []map[string]string{
				{"key": "value"},
			},
		}
		data.Display(&TextFormatter{})
		data.Display(&JSONFormatter{})
	})
}
