package zpack

import (
	"fmt"
	"io"
	"net"
	"testing"
)

func TestDataPack_Unpack(t *testing.T) {
	listener, err := net.Listen("tcp", "127.0.0.1:7777")
	if err != nil {
		fmt.Println("client listen err:", err)
		return
	}
	go func() {
		for {
			conn, err := listener.Accept()
			if err != nil {
				fmt.Println("server accept err:", err)
			}

			// Handle client requests
			go func(conn net.Conn) {
				// Create a packet splitting and packaging object dp.
				dp := DataPackInstance
				for {
					// 1. Read the head part of the stream first.
					headData := make([]byte, dp.GetHeadLen())
					_, err := io.ReadFull(conn, headData) // ReadFull will fill msg until it's full
					if err != nil {
						fmt.Println("read head error:", err)
						return
					}
					// Unpack the headData byte stream into msg.
					msgHead, err := dp.Unpack(headData)
					if err != nil {
						fmt.Println("server unpack err:", err)
						return
					}

					if msgHead.GetDataLen() > 0 {
						// msg has data, read data again.
						msg := msgHead.(*Message)
						msg.Data = make([]byte, msg.GetDataLen())

						// Read the byte stream from io based on dataLen.
						_, err := io.ReadFull(conn, msg.Data)
						if err != nil {
							fmt.Println("server unpack data err:", err)
							return
						}

						fmt.Println("==> Recv Msg: ID=", msg.ID, ", len=", msg.DataLen, ", data=", string(msg.Data))
					}
				}
			}(conn)
		}
	}()

	// Block the client.
	select {}
}

func TestDataPack_Pack(t *testing.T) {
	conn, err := net.Dial("tcp", "127.0.0.1:7777")
	if err != nil {
		fmt.Println("client dial err:", err)
		return
	}

	// Create a packet splitting and packaging object dp.
	dp := DataPackInstance

	// Package msg1.
	msg1 := &Message{
		ID:      0,
		DataLen: 5,
		Data:    []byte{'h', 'e', 'l', 'l', 'o'},
	}

	sendData1, err := dp.Pack(msg1)
	if err != nil {
		fmt.Println("client pack msg1 err:", err)
		return
	}

	// Package msg2.
	msg2 := &Message{
		ID:      1,
		DataLen: 7,
		Data:    []byte{'w', 'o', 'r', 'l', 'd', '!', '!'},
	}
	sendData2, err := dp.Pack(msg2)
	if err != nil {
		fmt.Println("client temp msg2 err:", err)
		return
	}

	// Concatenate sendData1 and sendData2 to create a sticky packet.
	sendData1 = append(sendData1, sendData2...)

	// Write data to the server.
	conn.Write(sendData1)
}
