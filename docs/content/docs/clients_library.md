---
date: '2025-09-14T22:11:36+03:00'
draft: true
title: 'Clients library'
---

## golang
```go
package main

import (
	"context"
	"log"

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
		log.Fatal(err)
	}

	loader.SetDefault(c)

	amount := loader.Get(context.Background(), amountKey).Int()
	log.Printf("amount: %d", amount)
}
```