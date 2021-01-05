package random

import (
	"math/rand"
	"time"

	"github.com/rs/xid"
)

func RandInt(max int) int {
	s1 := rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(s1)
	randomNum := r1.Intn(max)
	return randomNum
}

func ShortGUID() string {
	x := xid.New().String()[13:]
	return x
}
