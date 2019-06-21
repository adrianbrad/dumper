package dumper

type Opener interface {
	Open(driverName, dataSourceName string) (Preparator, error)
}
