package uplink

import "context"

type Request struct {
	GatewayID string
	Payload   []byte
}

type Sender interface {
	Send(context.Context, []byte) error
}

func Handle(ctx context.Context, sender Sender, request Request) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	payload := append([]byte(nil), request.Payload...)
	return sender.Send(ctx, payload)
}
