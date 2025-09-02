package main

import (
	"context"
	"log"
	"path"
	"strconv"
	"time"

	clientv3 "go.etcd.io/etcd/client/v3"

	"rtc/pkg/rtc"
	"rtc/pkg/rtc/loader"
	"rtc/pkg/rtc/providers/etcd"
)

const (
	amountKey = rtc.Key("amount")
)

var (
	etcdAddresses = []string{"127.0.0.1:2379"}
)

const (
	project = "example"
	env     = "prod"
	release = "latest"
)

func main() {
	go setValue()
	go getValue()

	time.Sleep(5 * time.Second)
}

func getValue() {
	c, err := etcd.NewProvider(etcdAddresses, project, env, release)
	if err != nil {
		log.Fatal(err)
	}

	loader.SetDefault(c)

	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		amount := loader.Get(context.Background(), amountKey).Int()
		log.Printf("amount: %d", amount)
	}
}

func setValue() {
	client, err := clientv3.New(clientv3.Config{
		Endpoints: etcdAddresses,
	})
	if err != nil {
		log.Fatal(err)
	}

	keyPath := path.Join("rtc", project, env, release, string(amountKey))

	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	var item int

	for range ticker.C {
		if _, err := client.Put(context.Background(), keyPath, strconv.Itoa(item)); err != nil {
			log.Fatal(err)
		}
		item++
	}
}
