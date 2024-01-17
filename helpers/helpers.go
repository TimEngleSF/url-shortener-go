package helpers

import (
	"math/rand"
	"net/url"

	"github.com/gin-gonic/gin"
)


func GenerateParam(length int) string {
	chars := "abcdefghijklmnopqrstuv0123456789"
	b := make([]byte, length)
	for i := range b {
		b[i] = chars[rand.Intn(len(chars))]
	}
	return string(b)
}

func FixURLString(u string) string {
	url := url.URL{
		Scheme: "https",
		Host: u,
	}
	return url.String()
}

func GetBaseURL(c *gin.Context) string {
	req := c.Request
	baseURL := "http://" + req.Host
	return baseURL
}