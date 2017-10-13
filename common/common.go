package common

import (
	"math/rand"
	"time"
)

const (
	alphabet string = "aAbBcCdDeEfFgGhHiIjJkKlLmMnNoOpPqQrRsStTuUvVwWxXyYzZ"
	number   string = "0123456789"
)

// RandomString creat a random string with length
func RandomString(strlen int, input string) string {
	rand.Seed(time.Now().UTC().UnixNano())

	rd := ""
	if input == "mixed" {
		mixed := alphabet + number
		for i := 0; i < strlen; i++ {
			index := rand.Intn(len(mixed))
			rd += mixed[index : index+1]
		}
	} else if input == "alphabet" {
		for i := 0; i < strlen; i++ {
			index := rand.Intn(len(alphabet))
			rd += alphabet[index : index+1]
		}
	} else if input == "number" {
		for i := 0; i < strlen; i++ {
			index := rand.Intn(len(number))
			rd += number[index : index+1]
		}
	}

	return rd
}
