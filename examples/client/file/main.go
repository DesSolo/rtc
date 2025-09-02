package main

import (
	"context"
	"fmt"

	"rtc/pkg/rtc"
	"rtc/pkg/rtc/providers/file"
	"rtc/pkg/rtc/providers/file/readers"
)

const (
	amountKey = rtc.Key("amount")
)

func main() {
	c, err := file.NewProvider("examples/client/file/config.yaml", readers.SimpleYAML())
	if err != nil {
		panic(err)
	}

	val, err := c.Value(context.Background(), amountKey)
	if err != nil {
		panic(err)
	}

	fmt.Println(val.String())
}
