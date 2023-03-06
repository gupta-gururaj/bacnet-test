package main

import (
	"fmt"

	bacnet "github.com/baetyl/baetyl-bacnet"
)

func main() {
	client, err := bacnet.NewClient("192.168.4.64", 47808)
	if err != nil {
		fmt.Println("err-->", err)
	}
}
