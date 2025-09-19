package main

import (
	"context"
	"fmt"

	"github.com/DesSolo/rtc/pkg/rtc"
	"github.com/DesSolo/rtc/pkg/rtc/providers/file"
	"github.com/DesSolo/rtc/pkg/rtc/providers/file/readers"
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

	fmt.Println(val.String()) // nolint:forbidigo
}
