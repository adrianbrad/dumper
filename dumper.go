package dumper

import (
	log "github.com/sirupsen/logrus"
	"sync"
	"time"
)

type Dumper struct {
	pServ   Service
	pDumper Service

	dump            atomicFlag
	attemptingRecon atomicFlag

	lastWriteErr error

	reconTicker int

	mutex sync.Mutex
}

func New(reconTicker int, pServ, pDumper Service) *Dumper {
	return &Dumper{
		pServ:       pServ,
		pDumper:     pDumper,
		reconTicker: reconTicker,
	}
}

func (d *Dumper) Open() (err error) {
	return d.pDumper.Open()
}

func (d *Dumper) Close() (err error) {
	return d.pDumper.Close()
}

func (d *Dumper) Read(p []byte) (n int, err error) {
	if d.dump.Get() {
		return d.pDumper.Read(p)
	}
	return d.pServ.Read(p)
}

func (d *Dumper) attemptOpenService() {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	ticker := time.NewTicker(time.Duration(d.reconTicker) * time.Millisecond)
	for {
		<-ticker.C
		log.Info("Attempting reconnection to service")
		if err := d.pServ.Open(); err != nil {
			continue
		}
		d.dump.Set(false)
		log.Info("Reconnected to service")
		break
	}

	p := make([]byte, 100)
	var err error
	var n int
	for {
		n, err = d.pDumper.Read(p)
		if err != nil {
			break
		}

		p = p[:n]
		log.Info(string(p))
		_, _ = d.Write(p)
	}
	log.Errorf("Error while reading from dump, err: %s", err.Error())
}

func (d *Dumper) Write(p []byte) (n int, err error) {
	defer func() {
		dump := d.dump.Get()
		if dump {
			n, err = d.pDumper.Write(p)
			if err != nil {
				log.Error("Error while writing to dumper: %s", err.Error())
			}
			log.Info("Successfully wrote to dump")
		}
	}()

	if n, err = d.pServ.Write(p); err != nil {
		d.lastWriteErr = err
		d.dump.Set(true)
		d.pServ.Close()
		log.Errorf("Error while writing to service: %s", err.Error())

		if !d.attemptingRecon.Get() {
			go d.attemptOpenService()
			d.attemptingRecon.Set(true)
		}
		return
	}
	log.Info("Successfully wrote to service")
	return
}
