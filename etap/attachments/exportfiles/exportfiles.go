package exportfiles

import (
	"fmt"
	"math"
	"math/rand"
	"time"

	"go.uber.org/multierr"
)

func Run(authn *Authn) (rErr error) {
	m := &Manifest{}
	if err := m.Load(); err != nil {
		return fmt.Errorf("failed to load manifest: %w", err)
	}
	defer func() {
		if err := m.Save(); err != nil {
			rErr = multierr.Append(rErr, fmt.Errorf("failed to save manifest: %w", err))
		}
	}()

	hasMore := true
	for hasMore {
		hm, err := m.Process(authn)
		if err != nil {
			return fmt.Errorf("failed to process manifest: %w", err)
		}
		hasMore = hm
		waitSoAsNotToOverloadTheServer()
	}
	fmt.Printf("You are done - all files have been successfully exported.\n")
	return nil
}

func waitSoAsNotToOverloadTheServer() {
	minMs := 2000
	maxMs := 10000
	rangeMs := maxMs - minMs
	toWaitMs := int(float64(maxMs) - math.Pow(float64(rand.Intn(rangeMs*rangeMs*rangeMs)), .3333))
	fmt.Printf("  wait for %d ms: ", toWaitMs)
	i := 0
	for toWaitMs-i >= 1000 {
		time.Sleep(1 * time.Second)
		i += 1000
		fmt.Printf("%d... ", i/1000)
	}
	time.Sleep(time.Duration(toWaitMs-i) * time.Millisecond)
	fmt.Println("done.")
}
