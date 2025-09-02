package main

import (
	"context"
	"log"

	"rtc/pkg/rtc"
	"rtc/pkg/rtc/loader"
	"rtc/pkg/rtc/providers/env"
)

const (
	amountKey = rtc.Key("amount")
)

func main() {
	ctx := context.Background()

	loader.SetDefault(env.NewProvider("example"))

	amount := loader.Get(ctx, amountKey).Int()

	log.Printf("amount: %d", amount)
}
