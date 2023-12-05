package util

import (
	"math/rand"
	"strings"
	"time"
)

var randGen *rand.Rand

const alphabet = "abcdefghijklmnopqrstuvwxyz"

func init() {
	source := rand.NewSource(time.Now().UnixNano())
	randGen = rand.New(source)
}

// random integer from min to max
func RandomInt(min, max int64) int64 {
	return min + randGen.Int63n(max-min+1)
}

// random string of len n
func RandomStr(n int) string {
	var sb strings.Builder
	k := len(alphabet)

	for i := 0; i < n; i++ {
		sb.WriteByte(alphabet[randGen.Intn(k)])
	}

	return sb.String()
}

// random money
func RandomMoney() int64 {
	return RandomInt(0, 1000)
}

// random currency
func RandomCurrency() string {
	currencies := []string{"USD", "EUR"}
	n := len(currencies)
	return currencies[rand.Intn(n)]
}
