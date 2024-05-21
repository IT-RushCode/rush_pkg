package utils

import "math/rand"

func GeneratePassword(length int, includeNumber bool) string {
	const charset = "LSKNDFKLN$Wlksdnfowinnf$DSFOI2NRFN3T04GNASIUDFH1"
	var password []byte
	var charSource string

	if includeNumber {
		charSource += "0123456789"
	}

	charSource += charset

	for i := 0; i < length; i++ {
		randNum := rand.Intn(len(charSource))
		password = append(password, charSource[randNum])
	}

	return string(password)
}
