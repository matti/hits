package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/paulbellamy/ratecounter"
)

type Handler struct{}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/readyz" || r.URL.Path == "/healthz" {
		w.WriteHeader(http.StatusOK)
		return
	}

	hostname, _ := os.Hostname()

	rate.Incr(1)

	body := fmt.Sprintf("%s hit at path %s from %s - rate is %d requests in %s\n", hostname, r.URL.Path, r.RemoteAddr, rate.Rate(), interval)
	log.Print(body)

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(body))
}

var rate *ratecounter.RateCounter
var interval time.Duration

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	interval = 10 * time.Second
	if s := os.Getenv("INTERVAL"); s != "" {
		if d, err := time.ParseDuration(s); err != nil {
			panic(err)
		} else {
			interval = d
		}
	}
	rate = ratecounter.NewRateCounter(interval)

	var port = "8080"
	if p := os.Getenv("PORT"); p != "" {
		port = p
	}

	addr := ":" + port

	handler := &Handler{}
	server := &http.Server{Addr: addr, Handler: handler}
	go func() {
		if err := server.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			panic(err)
		}
	}()

	log.Println("hits listens in", addr)
	s := <-sigs
	log.Println("received signal", s)

	if err := server.Shutdown(ctx); err != nil {
		panic(err)
	}

	log.Println("bye")
}
