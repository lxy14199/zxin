package main

import (
	"zinx/znet"
)

func main() {
	s := znet.NewServer("demo")
	s.Serve()
}
