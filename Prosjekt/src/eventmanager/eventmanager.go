package eventmanager

import (
	. ".././definitions"
	. ".././elev"
	. ".././json"
	. ".././udp"
	"fmt"
)

func EventManager(send_ch, receive_ch chan Udp_message, resendLostOrders, orderEventChannel chan int) {
	go GetOrders(orderEventChannel)
	go ElevLights()

	for {
		select {
		case UDP_Recv := <-receive_ch:
			if UDP_Recv.Raddr != Laddr.String() {
				fmt.Println("Laddr: ", Laddr.String())
				fmt.Println("Received from: ", UDP_Recv.Raddr)
				Dec_Msg := DecodeMsg(UDP_Recv.Data, UDP_Recv.Length)
				UpdateOrder(Dec_Msg, resendLostOrders)
				fmt.Println(Dec_Msg)
			}
		case <-orderEventChannel:
			Udp_Msg.Data = EncodeMsg(ElevatorList[0])
			send_ch <- Udp_Msg
			ElevatorList[0].FirstMsg = false
			ElevatorList[0].MessageType = NOTHING
		case <-resendLostOrders:
			send_ch <- Udp_Msg
		}
	}
}
