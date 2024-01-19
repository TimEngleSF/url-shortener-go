package controller

import (
	"fmt"
	"net/http"

	"link-short/helpers"
	links "link-short/models"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)



var validate = validator.New()

func Create(c *gin.Context) {
	var newLinkData links.LinkData
	if err := c.ShouldBindJSON(&newLinkData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Invalid siteURL string
	if !helpers.IsValidURL(newLinkData.SiteURL) {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "invalid link format"})	
		return
	}

	// Validate the struct
	if err := validate.Struct(newLinkData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	baseUrl := helpers.GetBaseURL(c)
	// Check if site already exists in db, if so respond with link
	existingSiteData := newLinkData.FindSiteByURL() 
	if existingSiteData != nil {
		c.JSON(http.StatusOK, gin.H{"link" : baseUrl + "/" + existingSiteData.ParamURL})
		return
	}

	// Create and store new link
	resultLinkData, res, err := newLinkData.AddLinkData()
	if err != nil {
		fmt.Println(err)
		return
	}
	c.JSON(http.StatusCreated, gin.H{"link" : baseUrl + "/" + resultLinkData.ParamURL})
	fmt.Println(resultLinkData, res)
	fmt.Printf("SiteURL: %s, ParamURL: %s\n", newLinkData.SiteURL, newLinkData.ParamURL)
	
}


