package downlink

import "context"

type Command struct {
	DeviceID string
	Payload  []byte
	Priority int
}

type CommandSender interface {
	Send(context.Context, Command) error
}

func Handle(ctx context.Context, sender CommandSender, command Command) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	command.Payload = append([]byte(nil), command.Payload...)
	return sender.Send(ctx, command)
}
