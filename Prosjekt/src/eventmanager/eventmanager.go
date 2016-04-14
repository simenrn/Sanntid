package eventmanager

import(
	"fmt"
	. ".././json"
	. ".././udp"
	. ".././elev"
	. ".././definitions"
)

func EventManager(send_ch, receive_ch chan Udp_message, orderEventChannel chan int){
	go GetOrders(orderEventChannel)
	go ElevLights()
	


	for {
		select {
			case UDP_Recv := <- receive_ch:
				if UDP_Recv.Raddr != Laddr.String() {
					fmt.Println("Laddr: ", Laddr.String())
					fmt.Println("Received from: ", UDP_Recv.Raddr)
					Dec_Msg := DecodeMsg(UDP_Recv.Data, UDP_Recv.Length)
					UpdateOrder(Dec_Msg)
					fmt.Println(Dec_Msg)
				}
			case <- orderEventChannel:
				Udp_Msg.Data = EncodeMsg(Msg)
				send_ch <- Udp_Msg
				Msg.MessageType = NOTHING

		}
	}
}
