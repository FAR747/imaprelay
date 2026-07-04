package target

import "context"

type Sender interface {
	Send(ctx context.Context, text string) error
}
