package main

import (
	. "./definitions"
	. "./elev"
	. "./eventmanager"
	. "./statemachine"
	. "./udp"
	"fmt"
	"time"
)

func main() {

	send_ch := make(chan Udp_message)
	receive_ch := make(chan Udp_message)
	OrderEventChannel := make(chan int)
	
	
	TimeOutChan := make(chan int)

	go Udp_init(LOCAL_LISTEN_PORT, BROADCAST_LISTEN_PORT, MESSAGE_SIZE, send_ch, receive_ch)

	if ElevInit() == false {
		fmt.Println("Elevator failed to initialize")
	}

	fmt.Println("Press STOP button to stop elevator and exit program\n")
	//ElevSetMotorDirection(1)
	Msg.ReadyToGo = true
	Msg.FirstMsg = true
	Msg.IsActive = true
	

	ElevatorList = append(ElevatorList, Msg)

	
	go StateMachine(OrderEventChannel)
	go EventManager(send_ch, receive_ch, OrderEventChannel, TimeOutChan)
	go Timer(TimeOutChan)
	go print()

	//tempIPSUGMEG := ElevatorList[0].MyIP[len(ElevatorList[0].MyIP)-4 : len(ElevatorList[0].MyIP)]
	//tempIPSTRINGTRUSE := string(tempIPSUGMEG)
	//if tempIPSTRINGTRUSE == "129.241.187.153" {
	//	fmt.Println("Du er HOMO")
	//}
	//fmt.Println("Byte til string kanskje?: ", tempIPSTRINGTRUSE)

	for {
		if ElevGetStopSignal() == 1 {
			ElevSetMotorDirection(0)
		}
	}
}

func print(){
	for{
		for i := range(ElevatorList){
			if i == 0{
				fmt.Println("This elevator: IP =\t",ElevatorList[i].MyIP[len(ElevatorList[i].MyIP)-1],", State: ",ElevatorList[i].State," , CurrentFloor: ",ElevatorList[i].CurrentFloor," , IsActive: ",ElevatorList[i].IsActive)
			} else {
				fmt.Println("Elevator ",i,": IP =\t",ElevatorList[i].MyIP[len(ElevatorList[i].MyIP)-1],", State: ",ElevatorList[i].State," , CurrentFloor: ",ElevatorList[i].CurrentFloor," , IsActive: ",ElevatorList[i].IsActive)
			}
		}
		fmt.Println("Local Que:\t\t",ElevatorList[0].Que_Local)
		fmt.Println("Global Up Que:\t\t",ElevatorList[0].Que_Global_Up)
		fmt.Println("Global Down Que:\t",ElevatorList[0].Que_Global_Down)
		fmt.Println("NextFloor:\t\t",ElevatorList[0].NextFloor,"\n")
		if NextDirection() == DIRN_UP{
			fmt.Println("Dirn: DIRN_UP")
		}
		if NextDirection() == DIRN_STOP{
			fmt.Println("Dirn: DIRN_STOP")
		}
		if NextDirection() == DIRN_DOWN{
			fmt.Println("Dirn: DIRN_DOWN")
		}

	time.Sleep(time.Second*3)
	}
}

