/**
 * @file requests.go
 * @brief Contains functions related to elevator requests handling.
 */

package requests

import (
	"Heis/config"
	"Heis/singleElev/elevio"
	"fmt"
)

var Floors int = 4
var Buttons int = 4

/**
 * @brief Checks if there are any requests above the current floor.
 * @param e The current state of the elevator.
 * @return Returns true if there are requests above the current floor, otherwise false.
 */
func Requests_above(e config.Elevator) bool {
	for f := e.Floor + 1; f < Floors; f++ {
		for b := 0; b < Buttons; b++ {
			if e.Requests[f][b] == 1 {
				return true
			}
		}
	}
	return false
}

/**
 * @brief Checks if there are any requests below the current floor.
 * @param e The current state of the elevator.
 * @return Returns true if there are requests below the current floor, otherwise false.
 */
func Requests_below(e config.Elevator) bool {
	for f := e.Floor - 1; f >= 0; f-- {
		for b := 0; b < Buttons; b++ {
			if e.Requests[f][b] == 1 {
				return true
			}
		}
	}
	return false
}

/**
 * @brief Checks if there are any requests on the current floor.
 * @param e The current state of the elevator.
 * @return Returns true if there are requests on the current floor, otherwise false.
 */
func Requests_current_floor(e config.Elevator) bool {

	for b := 0; b < Buttons; b++ {
		if e.Requests[e.Floor][b] == 1 {
			return true
		}

	}
	return false
}

/**
 * @brief Determines if the elevator should stop at the current floor based on requests.
 * @param e The current state of the elevator.
 * @return Returns true if the elevator should stop at the current floor, otherwise false.
 */
func Should_stop(e config.Elevator) bool {
	if Requests_current_floor(e) {
		switch {
		case e.Dirn == elevio.MD_Down:
			if config.Our_elevator.Requests[e.Floor][elevio.BT_HallUp] == 1 && Requests_below(e) {
				return false
			}
		case e.Dirn == elevio.MD_Up:
			if config.Our_elevator.Requests[e.Floor][elevio.BT_HallDown] == 1 && Requests_above(e) {
				return false
			}

		}
		return true
	}
	return false
}

/**
 * @brief Clears all button lights in the elevator.
 */
func Clear_lights() {
	for f := 0; f < Floors; f++ {
		elevio.SetButtonLamp(0, f, false)
		elevio.SetButtonLamp(1, f, false)
		elevio.SetButtonLamp(2, f, false)

	}
}

/**
 * @brief Clears requests and button lights at the current floor.
 * @param e A pointer to the current state of the elevator.
 */
func Clear_request_at_floor(e *config.Elevator) {
	config.Our_elevator.Requests[e.Floor][int(elevio.BT_Cab)] = 0
	elevio.SetButtonLamp(elevio.BT_Cab, e.Floor, false)

	switch {
	case config.Our_elevator.Dirn == elevio.MD_Up:
		config.Our_elevator.Requests[e.Floor][int(elevio.BT_HallUp)] = 0
		elevio.SetButtonLamp(elevio.BT_HallUp, e.Floor, false)
		if !Requests_above(config.Our_elevator) {
			config.Our_elevator.Requests[e.Floor][int(elevio.BT_HallDown)] = 0
			elevio.SetButtonLamp(elevio.BT_HallDown, e.Floor, false)
		}
	case config.Our_elevator.Dirn == elevio.MD_Down:
		config.Our_elevator.Requests[e.Floor][int(elevio.BT_HallDown)] = 0
		elevio.SetButtonLamp(elevio.BT_HallDown, e.Floor, false)
		if !Requests_below(config.Our_elevator) {
			config.Our_elevator.Requests[e.Floor][int(elevio.BT_HallUp)] = 0
			elevio.SetButtonLamp(elevio.BT_HallUp, e.Floor, false)
		}

	}
}

/**
 * @brief Chooses the elevator direction based on current requests.
 * @param e A pointer to the current state of the elevator.
 */
func Requests_chooseDirection(e *config.Elevator) {
	fmt.Printf("retning, inni choose:")
	fmt.Print(e.Dirn)
	switch e.Dirn {
	case elevio.MD_Up:
		if Requests_above(*e) {
			e.Dirn = elevio.MD_Up
		} else if Requests_below(*e) {
			e.Dirn = elevio.MD_Down
		} else {
			e.Dirn = elevio.MD_Stop
		}
	case elevio.MD_Down:
		if Requests_below(*e) {
			e.Dirn = elevio.MD_Down
		} else if Requests_above(*e) {
			e.Dirn = elevio.MD_Up
		} else {
			e.Dirn = elevio.MD_Stop
		}
	case elevio.MD_Stop:
		if Requests_above(*e) {
			e.Dirn = elevio.MD_Up
		} else if Requests_below(*e) {
			e.Dirn = elevio.MD_Down
		} else {
			e.Dirn = elevio.MD_Stop
		}
	}

}
