package model

import (
	"math/big"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func generateCode(len int) string {
	min := int('a')
	max := int('z')
	bytes := make([]byte, len)
	for i := 0; i < len; i++ {
		bytes[i] = byte(min + rand.Intn(max-min))
	}
	return string(bytes)
}

func validateID(id string) bool {
	if len(id) != 10 {
		return false
	}
	sum := 0
	for i := 0; i < 9; i++ {
		sum += (int(id[i]) - '0') * (10 - i)
	}
	mod := sum % 11
	c := 11 - mod
	if mod < 2 {
		c = mod
	}
	return int(id[9])-'0' == c
}

func validateCardNumber(str string) bool {
	if len(str) != 16 {
		return false
	}
	sum := 0
	for i, c := range str {
		multiplier := 2 - i%2
		num := int(c-'0') * multiplier
		if num > 9 {
			num -= 9
		}
		sum += num
	}
	return sum%10 == 0
}

func validateSheba(sheba string) bool {
	sheba = strings.ToUpper(sheba)
	if len(sheba) != 26 {
		return false
	}

	sheba = sheba[4:] + sheba[0:4]
	for i := 'A'; i <= 'Z'; i++ {
		sheba = strings.ReplaceAll(sheba, string(i), strconv.Itoa(int(i-'A'+10)))
	}

	shebaInt, success := new(big.Int).SetString(sheba, 10)
	if !success {
		return false
	}
	num := big.NewInt(97)
	res := new(big.Int)
	res.Mod(shebaInt, num)

	num.SetInt64(1)
	return res.Cmp(num) == 0
}
