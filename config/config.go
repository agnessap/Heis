/**
 * @file config.go
 * @brief Contains neccessary datatypes
 */

package config

import (
	"Heis/singleElev/elevio"
)

const (
	NumFloors  = 4
	NumButtons = 4
)

var Our_elevator Elevator
var Pair DirnBehaviourPair

type ElevatorBehaviour int

const (
	EB_Idle ElevatorBehaviour = iota
	EB_DoorOpen
	EB_Moving
)

type DirnBehaviourPair struct {
	Dirn      elevio.MotorDirection
	Behaviour ElevatorBehaviour
}

type ClearRequestVariant int

const (
	CV_All ClearRequestVariant = iota
	CV_InDirn
)

type Elevator struct {
	Floor     int
	NextDest  int
	Dirn      elevio.MotorDirection
	Requests  [4][4]int
	Behaviour ElevatorBehaviour

	Config struct {
		ClearRequestVariant ClearRequestVariant
		DoorOpenDurationS   float64
	}
}
