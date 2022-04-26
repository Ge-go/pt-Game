package util

import (
	"fmt"
	"math/rand"
	"time"
)

func RandomDigit() string {
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
	return fmt.Sprintf("%06v", rnd.Int31n(1000000))
}
