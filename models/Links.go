package links

import (
	"context"
	"fmt"
	"link-short/db"
	"link-short/helpers"

	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type LinkData struct {
	SiteURL  string `json:"siteURL" bson:"siteURL" validate:"required"`
	ParamURL string `json:"paramURL" bson:"paramURL"`
}

var validate = validator.New()

// AddLinkData inserts a linkData document into the database.
func (ld *LinkData) AddLinkData() (*LinkData, *mongo.InsertOneResult, error) {
	// Validate the struct
	if err := validate.Struct(ld); err != nil {
		return nil, nil, err
	}

	coll := db.GetDatabase().Collection("links")
	ld.ParamURL = helpers.GenerateParam(6)

	for {
		res, err := coll.InsertOne(context.TODO(), ld)
		if err == nil {
			return ld, res, nil
		}

		if writeErr, ok := err.(mongo.WriteException); ok {
			for _, e := range writeErr.WriteErrors {
				if e.Code == 11000 {
					ld.ParamURL = helpers.GenerateParam(6)
					continue
				}
			}
		}

		return nil, nil, err
	}
}

// FindSiteURL retrieves a linkData document from the database based on SiteURL.
func (ld *LinkData) FindSiteByURL() *LinkData {
	coll := db.GetDatabase().Collection("links")
	
	filter := bson.M{"siteURL": ld.SiteURL}

	var result LinkData

	if err := coll.FindOne(context.TODO(), filter).Decode(&result); err != nil {
		fmt.Println("err:", err)
		return nil
	}
	ld.ParamURL = result.ParamURL

	return ld
}

func (ld *LinkData) FindSiteByParam() (*LinkData, error) {
	coll := db.GetDatabase().Collection("links")
	fmt.Println("FINDBY", ld.ParamURL)
	filter := bson.M{"paramURL": ld.ParamURL}
	
	var result LinkData

	if err := coll.FindOne(context.TODO(), filter).Decode(&result); err != nil {
		fmt.Println("err:", err)
		return nil, err
	}

	ld.SiteURL = result.SiteURL

	return ld, nil
}
