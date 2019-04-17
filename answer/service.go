package answer

import (
	"errors"
	"github.com/godebug/ks/history"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"time"
)

type Answer struct {
	h *history.History
}

func (e *Answer) Answer(path string) error {
	h, err := history.NewHistory(path)
	if err != nil {
		return errors.New("Load history: " + err.Error())
	}
	e.h = h
	return nil
}

func (e *Answer) Serve(port string) error {
	r := mux.NewRouter()
	r.HandleFunc("/ask/{balls}", ask)
	http.Handle("/", r)
	srv := &http.Server{
		Handler:      r,
		Addr:         port,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	return srv.ListenAndServe()
}

func ask(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	w.WriteHeader(http.StatusOK)
	log.Printf("Balls: %s\n", vars["balls"])
}
