package elevator

import (
	. ".././driver"
	. ".././definitions"
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
	Msg.Que_Local[floor] = OFF
	Msg.Que_Global_Down[floor] = OFF
	Msg.Que_Global_Up[floor] = OFF
	Msg.MessageType = REMOVE_ORDER
	orderEventChannel <- ON
	time.Sleep(time.Second * 3)
	ElevSetDoorOpenLamp(OFF)
}

func GetInternalOrders() {
	for floor := 0; floor < N_FLOORS; floor++ {
		if ElevGetButtonSignal(BUTTON_COMMAND, floor) == 1 {
			Msg.Que_Local[floor] = 1
		}
	}
}

func GetExternalOrders(orderEventChannel chan int) {
	for floor := 0; floor < N_FLOORS; floor++ {
		if Msg.Que_Global_Up[floor] == OFF{
			if ElevGetButtonSignal(BUTTON_CALL_UP, floor) == 1 {
				Msg.Que_Global_Up[floor] = 1
				Msg.MessageType = ADD_ORDER
				orderEventChannel<-ON
			}
		}
		if Msg.Que_Global_Down[floor] == OFF{
			if ElevGetButtonSignal(BUTTON_CALL_DOWN, floor) == 1 {
				Msg.Que_Global_Down[floor] = 1
				Msg.MessageType = ADD_ORDER
				orderEventChannel<-ON
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
				Msg.PrevFloor = Current_Floor
			}
		}
		for floor := 0; floor < N_FLOORS; floor++ {
			if Msg.Que_Local[floor] == 1 {
				ElevSetButtonLamp(BUTTON_COMMAND, floor, 1)
			} else if Msg.Que_Local[floor] == 0 {
				ElevSetButtonLamp(BUTTON_COMMAND, floor, 0)
			}
			if Msg.Que_Global_Up[floor] == 1 {
				ElevSetButtonLamp(BUTTON_CALL_UP, floor, 1)
			} else if Msg.Que_Global_Up[floor] == 0 {
				ElevSetButtonLamp(BUTTON_CALL_UP, floor, 0)
			}
			if Msg.Que_Global_Down[floor] == 1 {
				ElevSetButtonLamp(BUTTON_CALL_DOWN, floor, 1)
			} else if Msg.Que_Global_Down[floor] == 0 {
				ElevSetButtonLamp(BUTTON_CALL_DOWN, floor, 0)
			}
		}
	}
}

func FloorReached() int {
	for {
		if ElevGetFloorSensorSignal() != -1 {
			if Msg.Dirn == DIRN_UP {
				if Msg.Que_Local[Msg.PrevFloor] == 1 || Msg.Que_Global_Up[Msg.PrevFloor] == 1 {
					Msg.PrevDirn = Msg.Dirn
					return DOOR_OPEN
				}
				if Msg.Que_Local[Msg.PrevFloor] == 0 && Msg.Que_Global_Up[Msg.PrevFloor] == 0 && Msg.Que_Global_Down[Msg.PrevFloor] == 1 {
					if Msg.PrevFloor != N_FLOORS-1 {
						for floor := Msg.PrevFloor + 1; floor < N_FLOORS; floor++ {
							if Msg.Que_Local[floor] == 1 || Msg.Que_Global_Up[floor] == 1 {
								fmt.Println("satan")
								return MOVING
							} else {
								fmt.Println("satan i japan")
								Msg.PrevDirn = Msg.Dirn
								return DOOR_OPEN
							}
						}
					} else {
						return DOOR_OPEN
					}
				}
			} else if Msg.Dirn == DIRN_DOWN {
				if Msg.Que_Local[Msg.PrevFloor] == 1 || Msg.Que_Global_Down[Msg.PrevFloor] == 1 {
					Msg.PrevDirn = Msg.Dirn
					return DOOR_OPEN
				}
				if Msg.Que_Local[Msg.PrevFloor] == 0 && Msg.Que_Global_Down[Msg.PrevFloor] == 0 && Msg.Que_Global_Up[Msg.PrevFloor] == 1 {
					if Msg.PrevFloor != 0 {
						for floor := Msg.PrevFloor - 1; floor >= 0; floor-- {
							if Msg.Que_Local[floor] == 1 || Msg.Que_Global_Down[floor] == 1 {
								return MOVING
							} else {
								Msg.PrevDirn = Msg.Dirn
								return DOOR_OPEN
							}
						}
					} else {
						return DOOR_OPEN
					}
				}
			} else {
				if Msg.Dirn == DIRN_STOP {
					return DOOR_OPEN
				}
			}
		}
	}
}

func ZeroOrders() bool {
	for i := 0; i < N_FLOORS; i++ {
		if Msg.Que_Local[i] == 1 || Msg.Que_Global_Up[i] == 1 || Msg.Que_Global_Down[i] == 1 {
			return false
		}
	}
	return true
}

func NextDirection() int {
	if Msg.Que_Local[Msg.PrevFloor] == ON || Msg.Que_Global_Up[Msg.PrevFloor] == ON || Msg.Que_Global_Down[Msg.PrevFloor] == ON {
		return DIRN_STOP
	} else if ZeroOrders() == false {
		if Msg.PrevDirn == DIRN_UP {
			for i := Msg.PrevFloor; i < N_FLOORS; i++ {
				if Msg.Que_Local[i] == ON || Msg.Que_Global_Up[i] == ON || Msg.Que_Global_Down[i] == ON {
					return DIRN_UP
				}
			}
			return DIRN_DOWN
		} else if Msg.PrevDirn == DIRN_DOWN {
			for i := Msg.PrevFloor; i >= 0; i-- {
				if Msg.Que_Local[i] == ON || Msg.Que_Global_Up[i] == ON || Msg.Que_Global_Down[i] == ON {
					return DIRN_DOWN
				}
			}
		}
		return DIRN_UP
	} else {
		return NOTHING
	}
}

func UpdateOrder(otherLift MSG) {
	switch otherLift.MessageType{
		case ADD_ORDER:
			for floor := 0; floor<N_FLOORS; floor++{
				if otherLift.Que_Global_Up[floor] == ON{
					Msg.Que_Global_Up[floor] = ON
				}
				if otherLift.Que_Global_Down[floor] == ON{
					Msg.Que_Global_Down[floor] = ON
				}
			}
		case REMOVE_ORDER:
			for floor := 0; floor<N_FLOORS; floor++{
				if otherLift.Que_Global_Up[floor] == OFF{
					Msg.Que_Global_Up[floor] = OFF
				}
				if otherLift.Que_Global_Down[floor] == OFF{
					Msg.Que_Global_Down[floor] = OFF
				}
			}
	}
}

