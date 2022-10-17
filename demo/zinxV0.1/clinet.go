package main

import (
	"fmt"
	"net"
	"time"
)

func main() {
	fmt.Println("clinet start...")

	time.Sleep(1 * time.Second)

	conn, err := net.Dial("tcp", "127.0.0.1:9851")
	if err != nil {
		fmt.Println("clinet start err, exit!")
		return
	}

	for {
		_, err := conn.Write([]byte("Hello zinx V0.1..."))
		if err != nil {
			fmt.Println("wirte conn err", err)
			return
		}

		buf := make([]byte, 512)
		cnt, err := conn.Read(buf)
		if err != nil {
			fmt.Println("read buf error")
			return
		}

		fmt.Printf("server call back: %s, cnt = %d\n", buf, cnt)

		time.Sleep(1 * time.Second)
	}
	return
}
