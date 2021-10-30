package utils

import "math/rand"

const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func RandomString(length int) string {
	str := ""
	for i := 1; i <= length; i++ {
		str += string(letters[rand.Intn(len(letters))])
	}
	return str
}
