package tools

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"math/rand"
	"strconv"
	"time"
)

func Rand(n int) string {
	const letter = "1234567890abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	letter_len := len(letter)

	b := make([]byte, n)
	for i := range b {
		b[i] = letter[rand.Intn(letter_len)]
	}
	return MD5(string(b)+fmt.Sprintf("%d,%d", time.Now().UnixNano(), rand.Int()))
}

func MD5(str string) string {
	m := md5.New()
	m.Write([]byte (str))
	return hex.EncodeToString(m.Sum(nil))
}

/**
 * @Desc  前置补0  fmt.Printf("%03d", a)
 * @Param d=数字，n=前面补多少个0
 * @return
 **/
func RepairSuff(d int, n int)string{
	return fmt.Sprintf("%0"+strconv.Itoa(n)+"d", d)
}
