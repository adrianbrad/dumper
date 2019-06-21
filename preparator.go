package dumper

type Preparator interface {
	Prepare(query string) (Executor, error)
}
