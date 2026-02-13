// TODO

package main

import (
	"fmt"
	"sync"
	"time"
)

type ParkingLot struct {
	slots chan struct{}
}

func NewParkingLot(slots int) *ParkingLot {
	return &ParkingLot{
		slots: make(chan struct{}, slots),
	}
}

func (p *ParkingLot) Park(carID int64) {
	p.slots <- struct{}{}
	fmt.Printf("Car %d has been parked.\n", carID)
	time.Sleep(time.Second)
	fmt.Printf("Car %d has left.\n", carID)
	<-p.slots
}

func main() {
	fmt.Println("main start")

	parking := NewParkingLot(3)
	var wg sync.WaitGroup
	carIDs := []int64{1, 2, 3, 4, 5, 6}

	for _, carID := range carIDs {
		wg.Add(1)
		go func(id int64) {
			defer wg.Done()
			parking.Park(id)
		}(carID)
	}
	wg.Wait()

	fmt.Println("All cars are parked")
}
