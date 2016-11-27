package main

import (
	"github.com/kataras/iris"
	"encoding/base64"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"time"
	"strconv"
)
type Impressions struct {
	ID        		bson.ObjectId `bson:"_id,omitempty"`
	IP        		string
	CampaignPublisher	CampaignPublisher
	CampaignLink		CampaignLink
	Timestamp 		time.Time
}
type CampaignPublisher struct {
	ID        		bson.ObjectId `bson:"_id,omitempty"`
	Value			string
}
type CampaignLink struct {
	ID        		bson.ObjectId `bson:"_id,omitempty"`
	Value			string

}
const base64GifPixel = "R0lGODlhAQABAIAAAP///wAAACwAAAAAAQABAAACAkQBADs="

func main() {
	session, err := mgo.Dial("127.0.0.1")

	defer session.Close()

	session.SetMode(mgo.Monotonic, true)

	iris.Get("/pid/:pid/lid/:lid", func(ctx *iris.Context) {
		c := session.DB("tracking").C("requests")
		err = c.Insert(
			&Impressions{IP: ctx.RemoteAddr(),
				CampaignPublisher: CampaignPublisher{
					Value: ctx.Param("pid"),
				},
				CampaignLink: CampaignLink{
					Value: ctx.Param("lid"),
				},
				Timestamp: time.Now(),
			})
		if err != nil {
			panic(err)
		}

		collection := session.DB("tracking").C("requests")
		var impressionsCount = 0
		impressionsCount, err = collection.Find(bson.M{"campaignpublisher.value": ctx.Param("pid"), "campaignlink.value":ctx.Param("lid")}).Count()

		if err != nil {
			panic(err)
		}
		pixel,_ := base64.StdEncoding.DecodeString(base64GifPixel)
		ctx.Response.Header.Set("Content-Control", "no-cache")
		ctx.Response.Header.Set("Content-Type", "image/gif")
		ctx.Response.Header.Set("Content-Length", "43")
		ctx.Response.SetBody(pixel)
		ctx.HTML(200,"<h1>This link has: " + strconv.Itoa(impressionsCount) + " impressions.</h1>")


	})
	iris.Listen(":8080")
}