package subroutines

import (
	"context"
	"fmt"
	"time"
)

func routine(ctx context.Context) error {
	fmt.Println("Subroutines 1")

	go func() {
		fmt.Println("Subroutines 2")
		go func() {
			fmt.Println("Subroutines 3: Wait...")
			<-ctx.Done()
			fmt.Println("Subroutines 3: Done")
		}()

		fmt.Println("Subroutines 2: Loop...")

		for i := 1; i < 10; i++ {
			fmt.Printf("Subroutines 2.%d\n", i)
			time.Sleep(1 * time.Second)
		}

		fmt.Println("Subroutines 2: End")
	}()

	fmt.Println("Subroutines 1: End")

	return nil
}
