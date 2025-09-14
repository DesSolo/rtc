package main

import (
	"context"
	"log"
	"time"

	"rtc/pkg/rtc"
	"rtc/pkg/rtc/loader"
	"rtc/pkg/rtc/providers/etcd"
)

const (
	amountKey = rtc.Key("amount")
)

func main() {
	c, err := etcd.NewProvider([]string{"127.0.0.1:2379"}, "example", "dev", "latest")
	if err != nil {
		log.Fatalf("new provider err: %s", err.Error())
	}

	loader.SetDefault(c)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	amount := loader.Get(ctx, amountKey).Int()
	log.Printf("amount: %d", amount)

	watchErr := loader.WatchValue(ctx, amountKey, func(oldValue, newValue rtc.Value) {
		log.Printf("old value: %s, new value: %s", oldValue.String(), newValue.String())
	})

	if watchErr != nil {
		log.Fatalf("watch err: %s", watchErr.Error())
	}

	<-ctx.Done()
}
