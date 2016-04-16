package elevator

import (
	. ".././definitions"
	. ".././driver"
	. ".././json"
	"errors"
	"fmt"
	"time"
)

func ElevInit() bool {

	// Initiate hardware
	if IoInit() == 0 {
		return false
	}
	// Zero all floor button lamps
	for i := 0; i < N_FLOORS; i++ {
		// Clearing all Call Down buttons
		if i != 0 {
			ElevSetButtonLamp(BUTTON_CALL_DOWN, i, 0)
		}

		// Clearing all Call UP buttons
		if i != N_FLOORS-1 {
			ElevSetButtonLamp(BUTTON_CALL_UP, i, 0)
		}

		ElevSetButtonLamp(BUTTON_COMMAND, i, 0)

	}
	ElevSetDoorOpenLamp(0)
	ElevSetFloorIndicator(0)
	if ElevGetFloorSensorSignal() != 0 {
		ElevSetMotorDirection(DIRN_DOWN)
		for {
			if ElevGetFloorSensorSignal() != -1 {
				ElevSetMotorDirection(DIRN_UP)
				time.Sleep(10 * time.Millisecond)
				ElevSetMotorDirection(DIRN_STOP)

				break
			}
		}
	}
	Msg.PrevDirn = DIRN_UP
	Msg.PrevFloor = Msg.CurrentFloor
	return true
}

func ElevSetMotorDirection(dirn int) {
	if dirn == 0 {
		IoWriteAnalog(MOTOR, 0)
	} else if dirn > 0 {
		IoClearBit(MOTORDIR)
		IoWriteAnalog(MOTOR, MOTOR_SPEED)
	} else {
		IoSetBit(MOTORDIR)
		IoWriteAnalog(MOTOR, MOTOR_SPEED)
	}
}

func ElevSetButtonLamp(button int, floor, value int) {
	if floor < 0 || floor >= N_FLOORS {
		errors.New("Floor is out of range.")
	} else if int(button) < 0 || int(button) >= N_BUTTONS {
		errors.New("Button is out of range")
	} else if value == 1 {
		IoSetBit(Lamp_channel_matrix[floor][button])
	} else {
		IoClearBit(Lamp_channel_matrix[floor][button])
	}
}

func ElevSetFloorIndicator(floor int) {
	if (floor < 0) || (floor >= N_FLOORS) {
		errors.New("Floor is out of range.")
	}
	switch floor {
	case 0:
		IoClearBit(LIGHT_FLOOR_IND1)
		IoClearBit(LIGHT_FLOOR_IND2)
	case 1:
		IoClearBit(LIGHT_FLOOR_IND1)
		IoSetBit(LIGHT_FLOOR_IND2)
	case 2:
		IoSetBit(LIGHT_FLOOR_IND1)
		IoClearBit(LIGHT_FLOOR_IND2)
	case 3:
		IoSetBit(LIGHT_FLOOR_IND1)
		IoSetBit(LIGHT_FLOOR_IND2)
	}
}

func ElevSetDoorOpenLamp(value int) {
	if value == 1 {
		IoSetBit(LIGHT_DOOR_OPEN)
	} else {
		IoClearBit(LIGHT_DOOR_OPEN)
	}
}

func ElevSetStopLamp(value int) {
	if value == 1 {
		IoSetBit(LIGHT_STOP)
	} else {
		IoClearBit(LIGHT_STOP)
	}
}

func ElevGetButtonSignal(button int, floor int) int {
	if IoReadBit(Button_channel_matrix[floor][button]) == 1 {
		return 1
	} else {
		return 0
	}
}

func ElevGetFloorSensorSignal() int {
	if IoReadBit(SENSOR_FLOOR1) == 1 {
		return 0
	} else if IoReadBit(SENSOR_FLOOR2) == 1 {
		return 1
	} else if IoReadBit(SENSOR_FLOOR3) == 1 {
		return 2
	} else if IoReadBit(SENSOR_FLOOR4) == 1 {
		return 3
	} else {
		return -1
	}
}

func ElevGetStopSignal() int {
	if IoReadBit(STOP) == 1 {
		return 1
	} else {
		return 0
	}
}

func ElevStopAtFloor(floor int, orderEventChannel chan int) {
	ElevSetMotorDirection(OFF)
	ElevSetButtonLamp(BUTTON_CALL_UP, floor, OFF)
	ElevSetButtonLamp(BUTTON_CALL_DOWN, floor, OFF)
	ElevSetButtonLamp(BUTTON_COMMAND, floor, OFF)
	ElevSetDoorOpenLamp(ON)
	ElevatorList[0].Que_Local[floor] = OFF
	ElevatorList[0].Que_Global_Down[floor] = OFF
	ElevatorList[0].Que_Global_Up[floor] = OFF
	ElevatorList[0].PrevFloor = ElevatorList[0].CurrentFloor
	if ElevatorList[0].NextFloor != nil {
	ElevatorList[0].NextFloor = ElevatorList[0].NextFloor[1:len(ElevatorList[0].NextFloor)]
	}
	if len(ElevatorList) > 1 { // To avoid sending remove message after crash
		ElevatorList[0].MessageType = REMOVE_ORDER
		orderEventChannel <- ON
	}
	time.Sleep(time.Second * 2)
	ElevSetDoorOpenLamp(OFF)
}

func GetInternalOrders() {
	var count = 0
	for floor := 0; floor < N_FLOORS; floor++ {
		if ElevGetButtonSignal(BUTTON_COMMAND, floor) == 1 {
			ElevatorList[0].Que_Local[floor] = 1
			if len(ElevatorList[0].NextFloor) == 0{
				ElevatorList[0].NextFloor = append(ElevatorList[0].NextFloor, floor)
			}
			for i := range ElevatorList[0].NextFloor{
				if floor == ElevatorList[0].NextFloor[i] {
					count++
				}
			}
			if count == 0 {
				ElevatorList[0].NextFloor = append(ElevatorList[0].NextFloor, floor)
			}
		}
	}
}

func GetExternalOrders(orderEventChannel chan int) {
	var count = 0
	for floor := 0; floor < N_FLOORS; floor++ {
		if ElevatorList[0].Que_Global_Up[floor] == OFF {
			if ElevGetButtonSignal(BUTTON_CALL_UP, floor) == 1 {
				ElevatorList[0].Que_Global_Up[floor] = 1
				if len(ElevatorList[0].NextFloor) == 0{
					ElevatorList[0].NextFloor = append(ElevatorList[0].NextFloor, floor)
				}
				for i := range ElevatorList[0].NextFloor{
					if floor == ElevatorList[0].NextFloor[i] {
						count++
					}
				}
				if count == 0 {
					ElevatorList[0].NextFloor = append(ElevatorList[0].NextFloor, floor)
				}
				ElevatorList[0].MessageType = ADD_ORDER
				orderEventChannel <- ON
			}
		}
		if ElevatorList[0].Que_Global_Down[floor] == OFF {
			if ElevGetButtonSignal(BUTTON_CALL_DOWN, floor) == 1 {
				ElevatorList[0].Que_Global_Down[floor] = 1
				if len(ElevatorList[0].NextFloor) == 0{
					ElevatorList[0].NextFloor = append(ElevatorList[0].NextFloor, floor)
				}
				for i := range ElevatorList[0].NextFloor{
					if floor == ElevatorList[0].NextFloor[i] {
						count++
					}
				}
				if count == 0 {
					ElevatorList[0].NextFloor = append(ElevatorList[0].NextFloor, floor)
				}
				ElevatorList[0].MessageType = ADD_ORDER
				orderEventChannel <- ON

			}
		}

	}
}

func GetOrders(orderEventChannel chan int) {
	for {
		GetInternalOrders()
		GetExternalOrders(orderEventChannel)
	}
}

func ElevLights() {
	for {
		if ElevGetFloorSensorSignal() != -1 {
			Current_Floor = ElevGetFloorSensorSignal()
			ElevSetFloorIndicator(ElevGetFloorSensorSignal())
			if Current_Floor != -1 {
				ElevatorList[0].CurrentFloor = Current_Floor
			}
		}
		for floor := 0; floor < N_FLOORS; floor++ {
			if ElevatorList[0].Que_Local[floor] == 1 {
				ElevSetButtonLamp(BUTTON_COMMAND, floor, 1)
			} else if ElevatorList[0].Que_Local[floor] == 0 {
				ElevSetButtonLamp(BUTTON_COMMAND, floor, 0)
			}
			if ElevatorList[0].Que_Global_Up[floor] == 1 {
				ElevSetButtonLamp(BUTTON_CALL_UP, floor, 1)
			} else if ElevatorList[0].Que_Global_Up[floor] == 0 {
				ElevSetButtonLamp(BUTTON_CALL_UP, floor, 0)
			}
			if ElevatorList[0].Que_Global_Down[floor] == 1 {
				ElevSetButtonLamp(BUTTON_CALL_DOWN, floor, 1)
			} else if ElevatorList[0].Que_Global_Down[floor] == 0 {
				ElevSetButtonLamp(BUTTON_CALL_DOWN, floor, 0)
			}
		}
	}
}

func FloorReached() int {
	for {
		if ElevGetFloorSensorSignal() != -1 {
			if ElevatorList[0].Dirn == DIRN_UP {
				if ElevatorList[0].PrevFloor == ElevatorList[0].CurrentFloor && ElevatorList[0].Que_Global_Up[ElevatorList[0].PrevFloor] == 1 {
					ElevatorList[0].PrevDirn = ElevatorList[0].Dirn
					fmt.Println("satan 5.0")
					return MOVING
				}
				if ElevatorList[0].Que_Local[ElevatorList[0].CurrentFloor] == 0 && ElevatorList[0].Que_Global_Up[ElevatorList[0].CurrentFloor] == 0 && ElevatorList[0].Que_Global_Down[ElevatorList[0].CurrentFloor] == 1 {
					if ElevatorList[0].CurrentFloor != N_FLOORS-1 {
						for floor := ElevatorList[0].CurrentFloor + 1; floor < N_FLOORS; floor++ {
							if ElevatorList[0].Que_Local[floor] == 1 || ElevatorList[0].Que_Global_Up[floor] == 1 {
								ElevatorList[0].PrevDirn = ElevatorList[0].Dirn
								fmt.Println("satan")
								return MOVING
							} 
							if ElevatorList[0].NextFloor[0] == floor {
								ElevatorList[0].PrevDirn = ElevatorList[0].Dirn
								fmt.Println("satan 2.0")
								return MOVING
							} 							
						}
						fmt.Println("satan i japan")
						ElevatorList[0].PrevDirn = ElevatorList[0].Dirn
						return DOOR_OPEN
					} else {
						return DOOR_OPEN
					}
				}
				if ElevatorList[0].Que_Global_Up[ElevatorList[0].CurrentFloor] == 1 || ElevatorList[0].Que_Local[ElevatorList[0].CurrentFloor] == 1 {
					ElevatorList[0].PrevDirn = ElevatorList[0].Dirn
					return DOOR_OPEN
				}
				if ElevatorList[0].CurrentFloor == N_FLOORS-1 && ElevatorList[0].Que_Local[ElevatorList[0].CurrentFloor] == 0 && ElevatorList[0].Que_Global_Up[ElevatorList[0].CurrentFloor] == 0 && ElevatorList[0].Que_Global_Down[ElevatorList[0].CurrentFloor] == 0 {
					ElevSetMotorDirection(DIRN_STOP)
					return IDLE
				}
			} else { // DIRN_DOWN
				if ElevatorList[0].Que_Global_Down[ElevatorList[0].PrevFloor] == 1 && ElevatorList[0].PrevFloor == ElevatorList[0].CurrentFloor{
					ElevatorList[0].PrevDirn = ElevatorList[0].Dirn
					fmt.Println("satan 6.0")
					return MOVING
				}
				if ElevatorList[0].Que_Local[ElevatorList[0].CurrentFloor] == 0 && ElevatorList[0].Que_Global_Down[ElevatorList[0].CurrentFloor] == 0 && ElevatorList[0].Que_Global_Up[ElevatorList[0].CurrentFloor] == 1 {
					if ElevatorList[0].CurrentFloor != 0 {
						for floor := ElevatorList[0].CurrentFloor - 1; floor >= 0; floor-- {
							if ElevatorList[0].Que_Local[floor] == 1 || ElevatorList[0].Que_Global_Down[floor] == 1 {
								fmt.Println("satan 3.0")
								return MOVING
							}
							if ElevatorList[0].NextFloor[0] == floor {
								ElevatorList[0].PrevDirn = ElevatorList[0].Dirn
								fmt.Println("satan 4.0")
								return MOVING
							} 							
						}
						fmt.Println("satan i fappan")
						ElevatorList[0].PrevDirn = ElevatorList[0].Dirn
						return DOOR_OPEN
					} else {
						return DOOR_OPEN
					}
				}
				if ElevatorList[0].Que_Global_Down[ElevatorList[0].CurrentFloor] == 1 || ElevatorList[0].Que_Local[ElevatorList[0].CurrentFloor] == 1{
					ElevatorList[0].PrevDirn = ElevatorList[0].Dirn
					return DOOR_OPEN
				}
				if ElevatorList[0].CurrentFloor == 0 && ElevatorList[0].Que_Local[ElevatorList[0].CurrentFloor] == 0 && ElevatorList[0].Que_Global_Up[ElevatorList[0].CurrentFloor] == 0 && ElevatorList[0].Que_Global_Down[ElevatorList[0].CurrentFloor] == 0 {
					ElevSetMotorDirection(DIRN_STOP)
					return IDLE
				}
			}
		}
	}
}

func ZeroOrders() bool {
	for i := 0; i < N_FLOORS; i++ {
		if ElevatorList[0].Que_Local[i] == 1 || ElevatorList[0].Que_Global_Up[i] == 1 || ElevatorList[0].Que_Global_Down[i] == 1 {
			return false
		}
	}
	return true
}


func TotalOrdersInSameDirn() int{
	var ExternalOrders = 0
	var WithInternalOrders = 0
	var InternalCount = 0
	if ElevatorList[0].PrevDirn == DIRN_UP{
		for i := range ElevatorList{
			for j := ElevatorList[0].CurrentFloor+1; j < N_FLOORS; j++ {
				if ElevatorList[i].Que_Local[j] == ON{
					InternalCount = 1
				}
			}
			WithInternalOrders += InternalCount
		}
		for i := ElevatorList[0].CurrentFloor+1; i < N_FLOORS; i++{
			if ElevatorList[0].Que_Local[i] == ON{
				ExternalOrders ++
			}
		}
		
	}
	if ElevatorList[0].PrevDirn == DIRN_DOWN {
		for i := range ElevatorList{
			for j := ElevatorList[0].CurrentFloor-1; j >= 0; j-- {
				if ElevatorList[i].Que_Local[j] == ON{
					InternalCount = 1
				}
			}
			WithInternalOrders += InternalCount
		}
		for i := ElevatorList[0].CurrentFloor-1; i >= 0; i--{
			if ElevatorList[0].Que_Local[i] == ON{
				ExternalOrders ++
			}
		}
		
	}
	return WithInternalOrders + ExternalOrders
}

func ElevatorsInSameDirn() int {
	var LiftInSameDirection = 0
	for i := range ElevatorList {
		if ElevatorList[0].PrevDirn == DIRN_UP {
			if ElevatorList[i].IsActive == true && ElevatorList[i].ReadyToGo == true{
				LiftInSameDirection ++
			}
		}	
		if ElevatorList[0].PrevDirn == DIRN_DOWN {
			if ElevatorList[i].IsActive == true && ElevatorList[i].ReadyToGo == true {
				LiftInSameDirection++
			}
		}
	}
	return LiftInSameDirection
}

func AliveAndReadyElevators() int {
	var count = 0
	for i := range(ElevatorList){
		if ElevatorList[i].ReadyToGo == true {
			if ElevatorList[i].IsActive == true {
				count++
			}
		}
	}
	return count
}



func NextDirection() int {
	ElevatorsInMyFloor := make([]int,0)
	var VisibleOrdersInSameDirection = 0
	var CloserThanCount = 0
	var FurtherThanCount= 0
	var LowestIpPlacing = 0
	var CloserThan = 0

	if ElevatorList[0].Que_Local[ElevatorList[0].CurrentFloor] == ON {
			return DIRN_STOP
	}
	if ElevatorList[0].Que_Global_Up[ElevatorList[0].CurrentFloor] == ON || ElevatorList[0].Que_Global_Down[ElevatorList[0].CurrentFloor] == ON {
		var temp byte
		for i := range ElevatorList {
			if ElevatorList[i].IsActive == true {
				if ElevatorList[i].CurrentFloor == ElevatorList[0].CurrentFloor {
					if ElevatorList[i].MyIP[len(ElevatorList[i].MyIP)-1] < ElevatorList[0].MyIP[len(ElevatorList[0].MyIP)-1] {
						temp = ElevatorList[i].MyIP[len(ElevatorList[i].MyIP)-1]
					} else {
						temp = ElevatorList[0].MyIP[len(ElevatorList[0].MyIP)-1]
					}
				}
			}
		}
		if temp == ElevatorList[0].MyIP[len(ElevatorList[0].MyIP)-1] {
			return DIRN_STOP
		}
	}


	if ElevatorList[0].PrevDirn == DIRN_UP{
		fmt.Println("a")
		for i := ElevatorList[0].CurrentFloor+1; i < N_FLOORS; i++{
			if ElevatorList[0].Que_Local[i] == ON || ElevatorList[0].Que_Global_Up[i] == ON || ElevatorList[0].Que_Global_Down[i] == ON {
				VisibleOrdersInSameDirection++
				fmt.Println("b")
			}
		}
		if VisibleOrdersInSameDirection > 0 && VisibleOrdersInSameDirection < TotalOrdersInSameDirn() {
				VisibleOrdersInSameDirection = TotalOrdersInSameDirn()
				fmt.Println("c")
		}
		if AliveAndReadyElevators() <= VisibleOrdersInSameDirection && VisibleOrdersInSameDirection !=0 {
			return DIRN_UP
			fmt.Println("d")
		} else {
			fmt.Println("e")
			for i := 0; i < AliveAndReadyElevators() - TotalOrdersInSameDirn(); i++{
				fmt.Println("f")
				for j := range(ElevatorList){
					if ElevatorList[0].CurrentFloor > ElevatorList[j].CurrentFloor{
						if ElevatorList[j].ReadyToGo == true || ElevatorList[j].IsActive == true {
							CloserThan = 1
							CloserThanCount++

						}
					}
					if ElevatorList[0].CurrentFloor < ElevatorList[j].CurrentFloor {
						if ElevatorList[j].ReadyToGo == true || ElevatorList[j].IsActive == true {
							FurtherThanCount++
						}
					}
				}
				
			}
			if CloserThan < AliveAndReadyElevators() - TotalOrdersInSameDirn() && VisibleOrdersInSameDirection !=0{
				return DIRN_UP
				fmt.Println("g")
			}
			var NeededElevators = FurtherThanCount + CloserThanCount - TotalOrdersInSameDirn()
			var TotalCount = FurtherThanCount + CloserThanCount
			
			if NeededElevators > FurtherThanCount + len(ElevatorsInMyFloor) {

				fmt.Println("h")
				for m := 0; m < NeededElevators - FurtherThanCount; m++{
					for k := range(ElevatorsInMyFloor){
						if ElevatorList[0].MyIP[len(ElevatorList[0].MyIP)-1] > ElevatorList[k].MyIP[len(ElevatorList[0].MyIP)-1]{
							LowestIpPlacing=1						
						}
					}
					LowestIpPlacing++
				}
				if LowestIpPlacing >NeededElevators - TotalCount{
					return DIRN_UP
				}
			}
		}
		ElevatorList[0].PrevDirn = DIRN_DOWN
	}
	if ElevatorList[0].PrevDirn == DIRN_DOWN{
		
		for i := ElevatorList[0].CurrentFloor-1; i >= 0; i--{
			if ElevatorList[0].Que_Local[i] == ON || ElevatorList[0].Que_Global_Up[i] == ON || ElevatorList[0].Que_Global_Down[i] == ON {
				VisibleOrdersInSameDirection++
			}
		}
		if VisibleOrdersInSameDirection > 0 {
			if VisibleOrdersInSameDirection < TotalOrdersInSameDirn(){
				VisibleOrdersInSameDirection = TotalOrdersInSameDirn()
			}
		}
		if AliveAndReadyElevators() <= VisibleOrdersInSameDirection && VisibleOrdersInSameDirection !=0 {
			return DIRN_DOWN
		} else {


			for i := 0; i <=  - ElevatorsInSameDirn() - AliveAndReadyElevators(); i++{
				for j := range(ElevatorList){
					if ElevatorList[0].CurrentFloor < ElevatorList[j].CurrentFloor{
						if ElevatorList[j].ReadyToGo == true || ElevatorList[j].IsActive == true {
							CloserThan = 1
							CloserThanCount++

						}
					}
					if ElevatorList[0].CurrentFloor > ElevatorList[j].CurrentFloor {
						if ElevatorList[j].ReadyToGo == true || ElevatorList[j].IsActive == true {
							FurtherThanCount++
						}
					}
				}
			}
			if CloserThan < AliveAndReadyElevators() - TotalOrdersInSameDirn() && VisibleOrdersInSameDirection !=0{
				return DIRN_UP
			}
			var NeededElevators = FurtherThanCount + CloserThanCount - TotalOrdersInSameDirn()
			var TotalCount = FurtherThanCount + CloserThanCount
			
			if NeededElevators > FurtherThanCount + len(ElevatorsInMyFloor) {
				for m := 0; m < NeededElevators - FurtherThanCount; m++{
					for k := range(ElevatorsInMyFloor){
						if ElevatorList[0].MyIP[len(ElevatorList[0].MyIP)-1] > ElevatorList[k].MyIP[len(ElevatorList[0].MyIP)-1]{
							LowestIpPlacing=1						
						}
					}
					LowestIpPlacing++
				}
				if LowestIpPlacing >NeededElevators - TotalCount{
					return DIRN_DOWN
				}
			}
		}
		ElevatorList[0].PrevDirn = DIRN_UP
	}
	return NOTHING
}



/*
func TotalOrdersInSameDirn() int {
	var Orders = 0
	var InternalOrderCount = 0
	var ExternalOrderCount = 0
	for i := range ElevatorList {
		if ElevatorList[i].IsActive == true {

			if ElevatorList[0].PrevDirn == DIRN_UP {
				for floor := ElevatorList[i].CurrentFloor + 1; floor < N_FLOORS; floor++ {
					if ElevatorList[i].Que_Local[floor] == ON {
						InternalOrderCount++
					}
				}
				for floor := ElevatorList[i].CurrentFloor + 1; floor < N_FLOORS; floor++ {
					if ElevatorList[i].Que_Global_Up[floor] == ON || ElevatorList[i].Que_Global_Down[floor] == ON {
						ExternalOrderCount++
					}
				}
				Orders = InternalOrderCount + ExternalOrderCount
				
			}
			if ElevatorList[0].PrevDirn == DIRN_DOWN {
				for floor := ElevatorList[i].CurrentFloor - 1; floor >= 0; floor-- {
					if ElevatorList[i].Que_Local[floor] == ON {
						InternalOrderCount++
					}
				}
				for floor := ElevatorList[i].CurrentFloor - 1; floor >= 0; floor-- {
					if ElevatorList[i].Que_Global_Up[floor] == ON || ElevatorList[i].Que_Global_Down[floor] == ON{
						ExternalOrderCount++
					}
				}
				Orders = InternalOrderCount + ExternalOrderCount
				
			}
		}
	}
	return Orders
}

func ElevatorsInSameDirn() int {
	var LiftInSameDirection = 0
	for i := range ElevatorList {
		if ElevatorList[0].PrevDirn == DIRN_UP {
			LiftInSameDirection += ElevatorList[i].ReadyToGoUp
		}
		if ElevatorList[0].PrevDirn == DIRN_DOWN {
			LiftInSameDirection += ElevatorList[i].ReadyToGoDown
		}
	}
	return LiftInSameDirection
}


func NextDirection() int {
	if ElevatorList[0].Que_Local[ElevatorList[0].PrevFloor] == ON || ElevatorList[0].Que_Global_Up[ElevatorList[0].PrevFloor] == ON || ElevatorList[0].Que_Global_Down[ElevatorList[0].PrevFloor] == ON {
		return DIRN_STOP
	} else if ZeroOrders() == false {
		if ElevatorList[0].PrevDirn == DIRN_UP {
			for i := ElevatorList[0].PrevFloor; i < N_FLOORS; i++ {
				if ElevatorList[0].Que_Local[i] == ON || ElevatorList[0].Que_Global_Up[i] == ON || ElevatorList[0].Que_Global_Down[i] == ON {
					return DIRN_UP
				}
			}
			return DIRN_DOWN
		} else if ElevatorList[0].PrevDirn == DIRN_DOWN {
			for i := ElevatorList[0].PrevFloor; i >= 0; i-- {
				if ElevatorList[0].Que_Local[i] == ON || ElevatorList[0].Que_Global_Up[i] == ON || ElevatorList[0].Que_Global_Down[i] == ON {
					return DIRN_DOWN
				}
			}
		}
		return DIRN_UP
	} else {
		return NOTHING
	}
}

func NextDirection() int {
	if len(ElevatorList)>0{
		var ReadyElevators = 0
		var OrderInSameDirection = 0
		//var temp byte
		for f := range(ElevatorList){
			if ElevatorList[f].State == IDLE || ElevatorList[f].State == DOOR_OPEN{
				ReadyElevators ++
			}
		}
		if ElevatorList[0].Que_Local[ElevatorList[0].CurrentFloor] == ON {
			return DIRN_STOP
		}
/*
		var VisibleOrderInSameDirection = 0
		for i := range(ElevatorList) {
			if ElevatorList[i].IsActive == true{
				if ElevatorList[i].State == IDLE || ElevatorList[i].State == DOOR_OPEN{

					if ElevatorList[0].Que_Global_Up[ElevatorList[0].CurrentFloor] == ON || ElevatorList[0].Que_Global_Down[ElevatorList[0].CurrentFloor] == ON {
						if ElevatorList[i].CurrentFloor == ElevatorList[0].CurrentFloor {
							if ElevatorList[i].MyIP[len(ElevatorList[i].MyIP)-1] < ElevatorList[0].MyIP[len(ElevatorList[0].MyIP)-1] {
								temp = ElevatorList[i].MyIP[len(ElevatorList[i].MyIP)-1]
							} else {
								temp = ElevatorList[0].MyIP[len(ElevatorList[0].MyIP)-1]
							}
						}
					}
				}
			}
		}
		if temp == ElevatorList[0].MyIP[len(ElevatorList[0].MyIP)-1] {
			return DIRN_STOP
		}
		


		if ElevatorList[0].PrevDirn == DIRN_UP {
						for i := ElevatorList[0].CurrentFloor; i < N_FLOORS; i++ {
							if ElevatorList[0].Que_Local[i] == ON || ElevatorList[0].Que_Global_Up[i] == ON || ElevatorList[0].Que_Global_Down[i] == ON {
							OrderInSameDirection++
							}
						}
					}
					if ReadyElevators <= OrderInSameDirection{
						return DIRN_UP
					}










		if ElevatorList[0].Que_Local[ElevatorList[0].CurrentFloor] == ON {
			return DIRN_STOP
		}
		
		if ElevatorList[0].Que_Global_Up[ElevatorList[0].CurrentFloor] == ON || ElevatorList[0].Que_Global_Down[ElevatorList[0].CurrentFloor] == ON {
			var temp byte
			for i := range ElevatorList {
				if ElevatorList[i].IsActive == true {
					if ElevatorList[i].CurrentFloor == ElevatorList[0].CurrentFloor {
						if ElevatorList[i].MyIP[len(ElevatorList[i].MyIP)-1] < ElevatorList[0].MyIP[len(ElevatorList[0].MyIP)-1] {
							temp = ElevatorList[i].MyIP[len(ElevatorList[i].MyIP)-1]
						} else {
							temp = ElevatorList[0].MyIP[len(ElevatorList[0].MyIP)-1]
						}
					}
				}
			}
			if temp == ElevatorList[0].MyIP[len(ElevatorList[0].MyIP)-1] {
				return DIRN_STOP
			}
		}

		if ElevatorList[0].PrevDirn == DIRN_UP {
			for i := ElevatorList[0].CurrentFloor; i < N_FLOORS; i++ {
				if ElevatorList[0].Que_Local[i] == ON || ElevatorList[0].Que_Global_Up[i] == ON || ElevatorList[0].Que_Global_Down[i] == ON {
					OrderInSameDirection++
				}
				if ReadyElevators <= OrderInSameDirection{
					return DIRN_UP
				}
			
			//fmt.Println("ordersamedi: ", OrderInSameDirection)
			//fmt.Println("Totalordersame: ", TotalOrdersInSameDirn())
			if ElevatorsInSameDirn() <= OrderInSameDirection {
			//	fmt.Println("Heiser i samme retning: ", ElevatorsInSameDirn())
			//	fmt.Println("Hølet til morra di")
				
				return DIRN_UP
			} 


		}
			//fmt.Println("ststja")
			ElevatorList[0].PrevDirn = DIRN_DOWN
			ElevatorList[0].ReadyToGoUp = OFF
			ElevatorList[0].ReadyToGoDown = ON
		}
		if ElevatorList[0].PrevDirn == DIRN_DOWN {
			for i := ElevatorList[0].CurrentFloor; i >= 0; i-- {
				if ElevatorList[0].Que_Local[i] == ON || ElevatorList[0].Que_Global_Up[i] == ON || ElevatorList[0].Que_Global_Down[i] == ON {
					OrderInSameDirection++
				}
			}
			//fmt.Println("NED ordersamedi: ", OrderInSameDirection)
			//fmt.Println("NED Totalordersame: ", TotalOrdersInSameDirn())
			if OrderInSameDirection < TotalOrdersInSameDirn() {
				OrderInSameDirection++
			}
			if ElevatorsInSameDirn() <= OrderInSameDirection {
			//	fmt.Println("Heiser i samme retning: ", ElevatorsInSameDirn())
			//	fmt.Println("Føde hølet til horemora di")
				
				return DIRN_DOWN
			} 
			//fmt.Println("mongo")
			ElevatorList[0].PrevDirn = DIRN_UP
			ElevatorList[0].ReadyToGoUp = ON
			ElevatorList[0].ReadyToGoDown = OFF
		}
	}
	return NOTHING
}*/

func UpdateOrder(externUpdate MSG, send_ch chan Udp_message) {
	var isSame = 0
	for j := range ElevatorList {
		if externUpdate.FirstMsg == true {
			if externUpdate.MyIP[len(externUpdate.MyIP)-1] == ElevatorList[j].MyIP[len(ElevatorList[j].MyIP)-1] {
				isSame++
				tempIP := externUpdate.MyIP[len(externUpdate.MyIP)-4 : len(externUpdate.MyIP)]
				tempIPString := string(tempIP)
				Udp_Msg.Raddr = tempIPString
				ElevatorList[j].MessageType = LOST_ORDERS
				ElevatorList[j].Que_Global_Up = ElevatorList[0].Que_Global_Up
				ElevatorList[j].Que_Global_Down = ElevatorList[0].Que_Global_Down
				ElevatorList[j].LastCheckin = time.Now()
				Udp_Msg.Raddr = RecvRaddr
				Udp_Msg.Data = EncodeMsg(ElevatorList[j])
				send_ch <- Udp_Msg
			}
		} else {
			if externUpdate.MyIP[len(externUpdate.MyIP)-1] == ElevatorList[j].MyIP[len(ElevatorList[j].MyIP)-1] {
				if ElevatorList[j].IsActive == true {
					ElevatorList[j] = externUpdate
					ElevatorList[j].LastCheckin = time.Now()
					isSame++
				}
			}
		}
	}

	if isSame == 0 {
		lift := new(MSG)
		*lift = externUpdate
		lift.LastCheckin = time.Now()
		ElevatorList = append(ElevatorList, *lift)

	}

	switch externUpdate.MessageType {
	case ADD_ORDER:
		for floor := 0; floor < N_FLOORS; floor++ {
			if externUpdate.Que_Global_Up[floor] == ON {
				ElevatorList[0].Que_Global_Up[floor] = ON
			}
			if externUpdate.Que_Global_Down[floor] == ON {
				ElevatorList[0].Que_Global_Down[floor] = ON
			}
		}
	case REMOVE_ORDER:
		for floor := 0; floor < N_FLOORS; floor++ {
			if externUpdate.Que_Global_Up[floor] == OFF {
				ElevatorList[0].Que_Global_Up[floor] = OFF
			}
			if externUpdate.Que_Global_Down[floor] == OFF {
				ElevatorList[0].Que_Global_Down[floor] = OFF
			}
		}
	case LOST_ORDERS:
		for floor := 0; floor < N_FLOORS; floor++ {
			if externUpdate.Que_Global_Up[floor] == ON {
				ElevatorList[0].Que_Global_Up[floor] = ON
			}
			if externUpdate.Que_Global_Down[floor] == ON {
				ElevatorList[0].Que_Global_Down[floor] = ON
			}
			if externUpdate.Que_Local[floor] == ON {
				ElevatorList[0].Que_Local[floor] = ON
			}
		}
	case NOTHING:
		for i := range ElevatorList {
			if externUpdate.MyIP[len(externUpdate.MyIP)-1] == ElevatorList[i].MyIP[len(ElevatorList[i].MyIP)-1] {

			}
		}
	}
}
