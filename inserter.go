package dumper

import "context"

type InsertOpener interface {
	Opener
	Insert(ctx context.Context, p Payload) (err error)
}
