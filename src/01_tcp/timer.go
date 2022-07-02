// go run timer.go

package tcp

import (
	"fmt"
	"time"
)

func main() {
	interval := 0 * time.Second
	timer := time.NewTimer(interval)
	// blocking interval sec
	fmt.Println(<-timer.C)    // 2022-04-09 18:12:49.324362 +0900 KST m=+1.001357562
	fmt.Println(timer.Stop()) // false

	timer2 := time.NewTimer(100 * time.Second)
	fmt.Println(timer2.Stop()) // true

	timer3 := time.NewTimer(100 * time.Second)
	defer func() {
		stop := timer3.Stop()
		fmt.Println(stop) // true
	}()

}

/*

2022-06-11 17:00:54.369284 +0900 KST m=+0.000085569
false // timer1 is already stop

true // the call stops the timer2

true // the call stops the timer3

*/
