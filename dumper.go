package dumper

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"time"
)

type Dumper struct {
	pServ   Service
	pDumper Service

	dump atomicFlag

	lastWriteErr error

	reconTicker int
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
	//no-op
	return 0, fmt.Errorf("this is a no-op")
}

func (d *Dumper) attemptOpenService() {
	ticker := time.NewTicker(time.Duration(d.reconTicker) * time.Millisecond)
	for {
		<-ticker.C
		if err := d.pServ.Open(); err != nil {
			continue
		}
		log.Info("Reconnected to service")
		break
	}

	var p []byte
	for _, err := d.pDumper.Read(p); err != nil; {
		_, _ = d.pServ.Write(p)
	}
}

func (d *Dumper) Write(p []byte) (n int, err error) {
	defer func() {
		dump := d.dump.Get()
		if dump {
			n, err = d.pDumper.Write(p)
			if err != nil {
				log.Error("Error while writing to dumper: %s", err.Error())
			}
		}
	}()

	if n, err = d.pServ.Write(p); err != nil {
		d.lastWriteErr = err
		d.dump.Set(true)
		d.pServ.Close()
		log.Errorf("Error while writing to service: %s", err.Error())
		log.Info("Service down, attempting reconnection")
		go d.attemptOpenService()
	}
	return
}
