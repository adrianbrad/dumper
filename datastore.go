package dumper

type DataStore struct {
	io     InsertOpener
	d      Dumper
	isOpen bool

	driver         string
	dataSourceName string
}

func NewDataStore(i InsertOpener, d Dumper, driver, dataSourceName string) *DataStore {
	ds := &DataStore{
		io:             i,
		d:              d,
		isOpen:         false,
		driver:         driver,
		dataSourceName: dataSourceName,
	}
	preper, err := ds.io.Open(driver, dataSourceName)
}

func (d *DataStore) StorePayload(p Payload) (err error) {
	err = d.inserter.Insert(p)
	if err != nil {
		// log error
		err = d.dumper.Dump(p)
		if err != nil {
			// this is the first thing to improve
			panic(err)
		}
	}
	return nil
}
