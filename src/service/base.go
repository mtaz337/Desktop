package service

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"math/rand"
	"strings"
)

var (
	Sess *session.Session
)

func (q *Queue) randString() string {
	chars := []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz")
	length := 10
	var b strings.Builder
	for i := 0; i < length; i++ {
		b.WriteRune(chars[rand.Intn(len(chars))])
	}
	str := b.String()
	return str
}
