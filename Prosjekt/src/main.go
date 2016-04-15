package main

import (
	. "./definitions"
	. "./elev"
	. "./eventmanager"
	. "./statemachine"
	. "./udp"
	"fmt"
	//"time"
)

func main() {

	send_ch := make(chan Udp_message)
	receive_ch := make(chan Udp_message)
	OrderEventChannel := make(chan int)
	ResendLostOrders := make(chan int)

	go Udp_init(LOCAL_LISTEN_PORT, BROADCAST_LISTEN_PORT, MESSAGE_SIZE, send_ch, receive_ch)

	if ElevInit() == false {
		fmt.Println("Elevator failed to initialize")
	}

	fmt.Println("Press STOP button to stop elevator and exit program\n")
	//ElevSetMotorDirection(1)
	Msg.ReadyToGoUp = ON
	OtherLift_1.MyIP = nil
	OtherLift_2.MyIP = nil
	Msg.FirstMsg = true

	ElevatorList = append(ElevatorList, Msg)

	go StateMachine(OrderEventChannel)
	go EventManager(send_ch, receive_ch, ResendLostOrders, OrderEventChannel)
	tempIPSUGMEG := ElevatorList[0].MyIP[len(ElevatorList[0].MyIP)-4 : len(ElevatorList[0].MyIP)]
	tempIPSTRINGTRUSE := string(tempIPSUGMEG)
	if tempIPSTRINGTRUSE == "129.241.187.153" {
		fmt.Println("Du er HOMO")
	}
	fmt.Println("Byte til string kanskje?: ", tempIPSTRINGTRUSE)
	fmt.Println("My IP: ", ElevatorList[0].MyIP[len(ElevatorList[0].MyIP)-4:len(ElevatorList[0].MyIP)])
	fmt.Println("Lengden på IP: ", len(ElevatorList[0].MyIP))
	fmt.Println("Siste byte av IP: ", ElevatorList[0].MyIP[len(ElevatorList[0].MyIP)-1])

	for {

		//fmt.Println("Other lift 1 IP: ", OtherLift_1.MyIP)
		//fmt.Println("Other lift 2 IP: ", OtherLift_2.MyIP)
		//fmt.Println("Local Address: ", Laddr)
		//fmt.Println("Broadcast Address: ", Baddr)
		/*
			fmt.Println("Interne ordre: ", ElevatorList[0].Que_Local)
			fmt.Println("Opp ordre: ", ElevatorList[0].Que_Global_Up)
			fmt.Println("Ned ordre: ", ElevatorList[0].Que_Global_Down)
			fmt.Println("ElevatorList[0].PrevFloor: ", ElevatorList[0].PrevFloor)
			fmt.Println("ElevatorList[0].Dirn: ", ElevatorList[0].Dirn)
			fmt.Println("ElevatorList[0].State", ElevatorList[0].State)
			time.Sleep(time.Second * 3)
		*/
		if ElevGetStopSignal() == 1 {
			ElevSetMotorDirection(0)
		}
	}
}
