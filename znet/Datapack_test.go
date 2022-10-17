package znet

import (
	"fmt"
	"io"
	"net"
	"testing"
)

func TestDataPack(t *testing.T) {
	listenner, err := net.Listen("tcp", "127.0.0.1:9851")
	if err != nil {
		fmt.Println("Server listen err: ", err)
		return
	}
	go func() {
		for {
			conn, err := listenner.Accept()
			if err != nil {
				fmt.Println("server accept err: ", err)
			}

			go func(conn net.Conn) {
				dp := NewDataPack()
				for {
					headData := make([]byte, dp.GetHeadLen())
					_, err := io.ReadFull(conn, headData)
					if err != nil {
						fmt.Println("read head err: ", err)
						break
					}
					msgHead, err := dp.UnPack(headData)
					if err != nil {
						fmt.Println("server unpack err: ", err)
						break
					}
					if msgHead.GetMsgLen() > 0 {
						msg := msgHead.(*Message)
						msg.Data = make([]byte, msg.GetMsgLen())
						_, err := io.ReadFull(conn, msg.Data)
						if err != nil {
							fmt.Println("server unpack data err: ", err)
							return
						}
						fmt.Println("----> recv MsgID: ", msg.Id, " DataLen:", msg.DataLen, " data:", msg.Data)
					}
				}
			}(conn)
		}
	}()

	conn, err := net.Dial("tcp", "127.0.0.1:9851")
	if err != nil {
		fmt.Println("client dial err:", err)
		return
	}
	dp := NewDataPack()
	msg1 := &Message{
		Id:      1,
		DataLen: 4,
		Data:    []byte{'z', 'i', 'n', 'x'},
	}
	sendData1, err := dp.Pack(msg1)
	if err != nil {
		fmt.Println("client pack err:", err)
		return
	}
	msg2 := &Message{
		Id:      2,
		DataLen: 8,
		Data:    []byte{'h', 'e', 'l', 'l', 'o', 'l', 'x', 'y'},
	}
	sendData2, err := dp.Pack(msg2)
	if err != nil {
		fmt.Println("client pack err:", err)
		return
	}
	sendData1 = append(sendData1, sendData2...)
	conn.Write(sendData1)
	select {}
}
