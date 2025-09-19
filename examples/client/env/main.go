package main

import (
	"context"
	"log"

	"github.com/DesSolo/rtc/pkg/rtc"
	"github.com/DesSolo/rtc/pkg/rtc/loader"
	"github.com/DesSolo/rtc/pkg/rtc/providers/env"
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
