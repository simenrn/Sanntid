package eventmanager

import (
	. ".././definitions"
	. ".././elev"
	. ".././json"
	. ".././udp"
	"fmt"
	"time"
)

func EventManager(send_ch, receive_ch chan Udp_message, orderEventChannel, timeOut chan int) {
	go GetOrders(orderEventChannel)
	go ElevLights()
	//sendAliveTicker := time.NewTicker(time.Second*2)
	

	for {
		select {
		case UDP_Recv := <-receive_ch:
			if UDP_Recv.Raddr != Laddr.String() {
				fmt.Println("Laddr: ", Laddr.String())
				fmt.Println("Received from: ", UDP_Recv.Raddr)
				RecvRaddr = UDP_Recv.Raddr
				Dec_Msg := DecodeMsg(UDP_Recv.Data, UDP_Recv.Length)
				UpdateOrder(Dec_Msg, send_ch)
				fmt.Println(Dec_Msg)
			}
		case <-orderEventChannel:
			Udp_Msg.Data = EncodeMsg(ElevatorList[0])
			send_ch <- Udp_Msg
			ElevatorList[0].FirstMsg = false
			ElevatorList[0].MessageType = NOTHING

		case i := <- timeOut:
			ElevatorList[i].IsActive = false
		/*case <-sendAliveTicker.C:
			SendAliveSignal(orderEventChannel)
		*/}
	}
}
/*
func SendAliveSignal(orderEventChannel chan int){
		ElevatorList[0].MessageType = NOTHING
		orderEventChannel<- ON
}*/

func Timer(timeout chan int){
	for{
		var TimeNow = time.Now()
		for i := range(ElevatorList){
			if i != 0{
				if TimeNow.Sub(ElevatorList[i].LastCheckin) > 3 * time.Second {
					timeout<- i
				}
			}
		}
	}
}