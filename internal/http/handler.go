package http

import (
	"dumper"
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
)

type Handler struct {
	service dumper.Service
}

func New(port string, s dumper.Service) (err error) {
	h := &Handler{s}
	log.Info("Starting server on port: ", port)
	return http.ListenAndServe(":"+port, h)
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	var p dumper.Payload
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	err = json.Unmarshal(b, &p)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	if !p.Valid() {
		w.WriteHeader(http.StatusUnprocessableEntity)
		w.Write([]byte("payload not valid"))
		return
	}

	_, err = h.service.Write(b)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	return
}
