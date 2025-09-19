---
date: '2025-09-14T22:11:36+03:00'
draft: false
title: 'Clients library'
weight: 4
---

## golang
This is a basic example of using the library to get values from etcd

```go
package main

import (
	"context"
	"log"

	"github.com/DesSolo/rtc/pkg/rtc"
	"github.com/DesSolo/rtc/pkg/rtc/loader"
	"github.com/DesSolo/rtc/pkg/rtc/providers/etcd"
)

const (
	amountKey = rtc.Key("amount")
)

func main() {
	c, err := etcd.NewProvider([]string{"127.0.0.1:2379"}, "example", "dev", "latest")
	if err != nil {
		log.Fatal(err)
	}

	loader.SetDefault(c)

	amount := loader.Get(context.Background(), amountKey).Int()
	log.Printf("amount: %d", amount)
}
```

### Value interface

```go
type Value interface {
	String() string
	MaybeString() (string, error)

	Float64() float64
	MaybeFloat64() (float64, error)

	Bool() bool
	MaybeBool() (bool, error)

	Int() int
	MaybeInt() (int, error)
}
```

### Get
Get value by key
```go
amount := loader.Get(context.Background(), amountKey).Int()
log.Printf("amount: %d", amount)

```

### WatchValue
Subscribe to changes
```go
loader.WatchValue(ctx, amountKey, func(oldValue, newValue rtc.Value) {
    log.Printf("oldValue: %s, newValue: %s", oldValue.String(), newValue.String())
})
```

In conjunction with atomics, this could be a truly useful tool.
{{< details title="Example with atomics" closed="true" >}}
```go
package main

import (
	"context"
	"log"
	"net/http"
	"sync/atomic"

	"github.com/DesSolo/rtc/pkg/rtc"
	"github.com/DesSolo/rtc/pkg/rtc/loader"
	"github.com/DesSolo/rtc/pkg/rtc/providers/etcd"
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
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if s.IsNewWearMode() {
			w.Write([]byte("holidays!"))
			return
		}

		w.Write([]byte("work"))
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

	watchErr := loader.WatchValue(ctx, newWearModelRTCKey, func(oldValue, newValue rtc.Value) {
		myServer.SetNewWearMode(newValue.Bool())
	})
	if watchErr != nil {
		log.Fatalf("watch err: %s", watchErr.Error())
	}

	mux := http.NewServeMux()
	mux.Handle("/", myServer.Handler())

	http.ListenAndServe(":8081", mux)
}
```
{{</ details>}}