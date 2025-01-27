package fsm

import (
	"Heis/config"
	"Heis/singleElev/elevio"
	"Heis/singleElev/requests"
	"fmt"
	"time"
)

/**
 * @brief Implements the finite state machine (FSM) logic for elevator control.
 * @param buttons Channel for receiving button events.
 * @param floors Channel for receiving floor events.
 * @param obstr Channel for receiving obstruction events.
 * @param stop Channel for receiving stop events.
 * @param doorTimer Pointer to the door timer.
 * @param numFloors Total number of floors in the building.
 */
func Fsm(buttons chan elevio.ButtonEvent, floors chan int, obstr chan bool, stop chan bool, doorTimer *time.Timer, numFloors int) {
	requests.Clear_lights()
	for {
		select {
		case order := <-buttons:
			elevio.SetButtonLamp(order.Button, order.Floor, true)
			switch {
			case config.Our_elevator.Behaviour == config.EB_DoorOpen:
				if order.Floor == config.Our_elevator.Floor {
					elevio.SetDoorOpenLamp(true)
					requests.Clear_request_at_floor(&config.Our_elevator)
					doorTimer.Reset(time.Duration(3) * time.Second)
				} else {
					config.Our_elevator.Requests[order.Floor][order.Button] = 1

				}
			case config.Our_elevator.Behaviour == config.EB_Moving:
				config.Our_elevator.Requests[order.Floor][order.Button] = 1

			case config.Our_elevator.Behaviour == config.EB_Idle:
				if order.Floor == config.Our_elevator.Floor {

					elevio.SetDoorOpenLamp(true)
					requests.Clear_request_at_floor(&config.Our_elevator)
					config.Our_elevator.Behaviour = config.EB_DoorOpen
					doorTimer.Reset(time.Duration(3) * time.Second)

				} else {
					config.Our_elevator.Requests[order.Floor][order.Button] = 1
					if requests.Requests_above(config.Our_elevator) {
						config.Our_elevator.Dirn = elevio.MD_Up
						elevio.SetMotorDirection(config.Our_elevator.Dirn)
						config.Our_elevator.Behaviour = config.EB_Moving
					} else if requests.Requests_below(config.Our_elevator) {
						config.Our_elevator.Dirn = elevio.MD_Down
						elevio.SetMotorDirection(config.Our_elevator.Dirn)
						config.Our_elevator.Behaviour = config.EB_Moving
					}
				}
			}

		case floor := <-floors:
			config.Our_elevator.Floor = floor
			fmt.Printf("%+v\n", floor)
			fmt.Printf("retning før stop:")
			fmt.Print(config.Our_elevator.Dirn)

			if requests.Should_stop(config.Our_elevator) {

				fmt.Printf("retning:")
				fmt.Print(config.Our_elevator.Dirn)
				elevio.SetMotorDirection(elevio.MD_Stop)

				elevio.SetDoorOpenLamp(true)
				requests.Clear_request_at_floor(&config.Our_elevator)
				config.Our_elevator.Behaviour = config.EB_DoorOpen
				elevio.SetDoorOpenLamp(true)

				doorTimer.Reset(time.Duration(3) * time.Second)

			}
		case <-doorTimer.C:
			elevio.SetDoorOpenLamp(false)
			switch {
			case config.Our_elevator.Behaviour == config.EB_DoorOpen:
				requests.Requests_chooseDirection(&config.Our_elevator)
				elevio.SetMotorDirection(config.Our_elevator.Dirn)

				if config.Our_elevator.Dirn == elevio.MD_Stop {
					config.Our_elevator.Behaviour = config.EB_Idle
				} else {
					config.Our_elevator.Behaviour = config.EB_Moving

				}
			}

		case obstruction := <-obstr:
			if config.Our_elevator.Behaviour == config.EB_DoorOpen {
				if obstruction {
					if !doorTimer.Stop() {
						<-doorTimer.C
					}

				} else {
					doorTimer.Reset(time.Duration(3) * time.Second)
				}

			}

		case a := <-stop:
			if a { // Stop button pressed
				// Stop the elevator immediately
				elevio.SetMotorDirection(elevio.MD_Stop)

				// Clear all requests/orders
				requests.Clear_lights()
				requests.Clear_request_at_floor(&config.Our_elevator)

				// Turn on the stop button lamp
				elevio.SetStopLamp(true)

				// If the elevator is currently at a floor, open the door after three seconds
				if config.Our_elevator.Behaviour != config.EB_Moving {
					elevio.SetDoorOpenLamp(true)
					config.Our_elevator.Behaviour = config.EB_DoorOpen
					doorTimer.Reset(time.Duration(3) * time.Second)
				}
			} else { // Stop button released
				// Turn off the stop button lamp
				elevio.SetStopLamp(false)

				// If the elevator is currently at a floor, start the door timer to close the door after three seconds
				if config.Our_elevator.Behaviour != config.EB_Moving {
					doorTimer.Reset(time.Duration(3) * time.Second)
				}
			}
		}
	}
}
