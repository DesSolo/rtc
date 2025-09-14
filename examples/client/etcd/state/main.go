package main

import (
	"context"
	"log"
	"net/http"
	"sync/atomic"

	"rtc/pkg/rtc"
	"rtc/pkg/rtc/loader"
	"rtc/pkg/rtc/providers/etcd"
)

const (
	newWearModelRTCKey = rtc.Key("new_wear_mode")
)

type server struct {
	newWearMode atomic.Bool
}

func (s *server) SetNewWearMode(newWearMode bool) {
	s.newWearMode.Store(newWearMode)
}

func (s *server) IsNewWearMode() bool {
	return s.newWearMode.Load()
}

func (s *server) Handler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		if s.IsNewWearMode() {
			w.Write([]byte("holidays!")) // nolint:errcheck
			return
		}

		w.Write([]byte("work")) // nolint:errcheck
	})
}

func main() {
	c, err := etcd.NewProvider([]string{"127.0.0.1:2379"}, "example", "dev", "latest")
	if err != nil {
		log.Fatalf("new provider err: %s", err.Error())
	}

	loader.SetDefault(c)

	ctx := context.Background()

	myServer := &server{}
	myServer.SetNewWearMode(loader.Get(ctx, newWearModelRTCKey).Bool())

	watchErr := loader.WatchValue(ctx, newWearModelRTCKey, func(_, newValue rtc.Value) {
		myServer.SetNewWearMode(newValue.Bool())
	})
	if watchErr != nil {
		log.Fatalf("watch err: %s", watchErr.Error())
	}

	mux := http.NewServeMux()
	mux.Handle("/", myServer.Handler())

	if err := http.ListenAndServe(":8081", mux); err != nil { // nolint:gosec
		log.Fatalf("listen err: %s", err.Error())
	}
}
