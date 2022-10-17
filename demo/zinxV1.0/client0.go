package main

import (
	"fmt"
	"io"
	"net"
	"time"
	"zinx/znet"
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
		dp := znet.NewDataPack()
		binaryMsg, err := dp.Pack(znet.NewMessage(0, []byte("zinxV0.6 client0 Test")))
		if err != nil {
			fmt.Println("Pack err:", err)
			return
		}
		if _, err := conn.Write(binaryMsg); err != nil {
			fmt.Println("write err:", err)
			return
		}

		binaryHead := make([]byte, dp.GetHeadLen())
		if _, err := io.ReadFull(conn, binaryHead); err != nil {
			fmt.Println("read head err:", err)
			break
		}
		msgHead, err := dp.UnPack(binaryHead)
		if err != nil {
			fmt.Println("client unpack msghead err:", err)
			break
		}
		if msgHead.GetMsgLen() > 0 {
			msg := msgHead.(*znet.Message)
			msg.Data = make([]byte, msg.GetMsgLen())
			if _, err := io.ReadFull(conn, msg.Data); err != nil {
				fmt.Println("read msg data error", err)
				return
			}
			fmt.Println("-----> Recv Server Message Id = ", msg.GetMsgId(), "msgLen = ", msg.GetMsgLen(), "data = ", string(msg.GetData()))
		}
		time.Sleep(1 * time.Second)
	}
	return
}
