package main

import (
	"github.com/codecrafters-io/redis-starter-go/app/managers"
)

func main() {

	netManager := managers.NewNetManager()
	go netManager.Init()
	go netManager.Start()

	func() {
		for {
		}
	}()
}
