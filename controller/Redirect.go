package controller

import (
	"fmt"
	"link-short/helpers"
	links "link-short/models"

	"github.com/gin-gonic/gin"
)

func Redirect(c *gin.Context) {
	paramURL := c.Param("paramURL")
	var newLinkData links.LinkData

	fmt.Println("PARAM", paramURL)

	if paramURL == ""{
		c.Redirect(302, helpers.GetBaseURL(c))
		return
	}
	
	newLinkData.ParamURL = paramURL

	
	redirectLinkData := newLinkData.FindSiteByParam()
	fmt.Println(redirectLinkData)


	c.Redirect(302, redirectLinkData.SiteURL)
}