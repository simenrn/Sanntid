package main

import (
	. "./elev"
	. "./statemachine"
	. "./udp"
	. "./definitions"
	"fmt"
	. "./eventmanager"
	//"time"
)



func main() {

	send_ch := make(chan Udp_message)
	receive_ch := make(chan Udp_message)
	OrderEventChannel := make(chan int)

	go Udp_init(LOCAL_LISTEN_PORT, BROADCAST_LISTEN_PORT, MESSAGE_SIZE, send_ch, receive_ch)

	if ElevInit() == false {
		fmt.Println("Elevator failed to initialize")
	}

	fmt.Println("Press STOP button to stop elevator and exit program\n")
	//ElevSetMotorDirection(1)
	
	go StateMachine(OrderEventChannel)
	go EventManager(send_ch, receive_ch, OrderEventChannel)

	


	for {
		//fmt.Println("Local Address: ", Laddr)
		//fmt.Println("Broadcast Address: ", Baddr)
/*
		fmt.Println("Interne ordre: ", Msg.Que_Local)
		fmt.Println("Opp ordre: ", Msg.Que_Global_Up)
		fmt.Println("Ned ordre: ", Msg.Que_Global_Down)
		fmt.Println("Msg.PrevFloor: ", Msg.PrevFloor)
		fmt.Println("Msg.Dirn: ", Msg.Dirn)
		fmt.Println("Msg.State", Msg.State)
		time.Sleep(time.Second * 3)
*/
		if ElevGetStopSignal() == 1 {
			ElevSetMotorDirection(0)
		}
	}
}
