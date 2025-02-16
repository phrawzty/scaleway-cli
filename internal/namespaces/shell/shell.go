package shell

import (
	"context"
	"reflect"

	"github.com/scaleway/scaleway-cli/v2/internal/args"
	"github.com/scaleway/scaleway-cli/v2/internal/core"
)

func GetCommands() *core.Commands {
	return core.NewCommands(
		shellCommand(),
	)
}

func shellCommand() *core.Command {
	return &core.Command{
		Short:                "Start shell mode",
		Namespace:            "shell",
		AllowAnonymousClient: false,
		ArgsType:             reflect.TypeOf(args.RawArgs{}),
		Run: func(ctx context.Context, argsI interface{}) (interface{}, error) {
			return nil, nil
		},
	}
}
