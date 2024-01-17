package controller

import (
	"context"
	"fmt"

	// "io"
	"math/rand"
	"net/http"
	"net/url"

	"link-short/db"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type linkData struct {
	SiteURL  string `json:"siteURL" bson:"siteURL" validate:"required"`
	ParamURL string `json:"paramURL bson:"paramURL"`
}

var validate = validator.New()

func Create(c *gin.Context) {
	var newLinkData linkData
	if err := c.ShouldBindJSON(&newLinkData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Invalid siteURL string
	if _, err := url.ParseRequestURI(newLinkData.SiteURL); err != nil {
		newLinkData.SiteURL = fixURL(newLinkData.SiteURL)		
	}

	// Validate the struct
	if err := validate.Struct(newLinkData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	baseUrl := getBaseURL(c)
	// Check if site already exists in db, if so respond with link
	existingSiteData := findSiteURL(&newLinkData) 
	if existingSiteData != nil {
		c.JSON(http.StatusAccepted, gin.H{"link" : baseUrl + "/" + existingSiteData.ParamURL})
		return
	}

	// Create and store new link
	resultLinkData, res := addLinkData(&newLinkData)
	c.JSON(http.StatusAccepted, gin.H{"link" : baseUrl + "/" + resultLinkData.ParamURL})
	fmt.Println(resultLinkData, res)
	fmt.Printf("SiteURL: %s, ParamURL: %s\n", newLinkData.SiteURL, newLinkData.ParamURL)
	
}

func fixURL(u string) string {
	url := url.URL{
		Scheme: "https",
		Host: u,
	}
	return url.String()
}

func findSiteURL(ld *linkData) *linkData {
	coll := db.GetDatabase().Collection("links")

	filter := bson.M{"siteURL": *&ld.SiteURL}

	var result linkData

	if err := coll.FindOne(context.TODO(), filter).Decode(&result); err != nil{
		fmt.Println("err:", err)
		return nil
	}
	ld.ParamURL = result.ParamURL
	return ld
}

func generateParam(length int) string {
	chars := "abcdefghijklmnopqrstuv0123456789"
	b := make([]byte, length)
	for i := range b {
		b[i] = chars[rand.Intn(len(chars))]
	}
	return string(b)
}

func addLinkData(ld *linkData) (*linkData, *mongo.InsertOneResult) {
	coll := db.GetDatabase().Collection("links")
	*&ld.ParamURL = generateParam(6)

	for {
		res, err := coll.InsertOne(context.TODO(), ld)
		if err == nil {
			return ld, res
		}

		if writeErr, ok := err.(mongo.WriteException); ok {
			for _, e := range writeErr.WriteErrors {
				if e.Code == 11000 {
					*&ld.ParamURL = generateParam(6)
					continue
				}
			}
		}

		return nil, nil
	}
}

func getBaseURL(c *gin.Context) string {
	req := c.Request
	baseURL := "http://" + req.Host
	return baseURL
}
