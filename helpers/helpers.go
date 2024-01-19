package helpers

import (
	"math/rand"
	"net/url"
	"regexp"
	"strings"

	"github.com/gin-gonic/gin"
)

func GenerateParam(length int) string {
	chars := "abcdefghijklmnopqrstuvwxyz0123456789" + 
		"ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	b := make([]byte, length)
	for i := range b {
		b[i] = chars[rand.Intn(len(chars))]
	}
	return string(b)
}

func AddURLPrefix(u string) string {
	if strings.HasPrefix(u, "http://") || strings.HasPrefix(u, "https://") {
		return u
	}
	
	url := url.URL{
		Scheme: "https",
		Host:   u,
	}
	return url.String()
}
// TODO:// Change http:// to https when we get SSL certs
func GetBaseURL(c *gin.Context) string {
	req := c.Request
	baseURL := "http://" + req.Host
	return baseURL
}


func IsValidURL(url string) bool {
	pattern := `^(https?|ftp):\/\/[^\s\/$.?#].[^\s]*\.[^\s]{2,}.*$`
	regex := regexp.MustCompile(pattern)

	return regex.MatchString(url)
}