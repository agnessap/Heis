/**
 * @file main.go
 * @brief Entry point for the elevator control program.
 */

package main

import (
	"Heis/singleElev/elevio"
	"Heis/singleElev/fsm"
	"time"
)

/**
 * @brief The entry point for the elevator control program.
 */
func main() {
	// Define the number of floors in the building
	numFloors := 4
	// Initialize elevator I/O
	elevio.Init("localhost:15657", numFloors)

	// Create channels for elevator I/O events
	drv_buttons := make(chan elevio.ButtonEvent)
	drv_floors := make(chan int)
	drv_obstr := make(chan bool)
	drv_stop := make(chan bool)
	doorTimer := time.NewTimer(time.Duration(3) * time.Second)

	// Start polling for elevator I/O events
	go elevio.PollButtons(drv_buttons)
	go elevio.PollFloorSensor(drv_floors)
	go elevio.PollObstructionSwitch(drv_obstr)
	go elevio.PollStopButton(drv_stop)

	// Start the finite state machine for single elevator control
	fsm.Fsm(drv_buttons, drv_floors, drv_obstr, drv_stop, doorTimer, numFloors)
}
