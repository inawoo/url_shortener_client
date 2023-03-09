package client

import (
	"net/url"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type (
	URLCollection struct {
		ID        primitive.ObjectID `json:"id" bson:"id" example:"5f9e9b9b9b9b9b9b9b9b9b9b" description:"id"`
		Host      string             `json:"host" bson:"host" example:"http://localhost:8080/" description:"host url"`
		Path      string             `json:"url" bson:"url" example:"/test" description:"url path"`
		Params    map[string]string  `json:"params" bson:"params" example:"{\"test\": \"test\"}" description:"query params"`
		Code      string             `json:"code" bson:"code" example:"3432" description:"short url code"` // indexed
		Count     int64              `json:"count" bson:"count" example:"3" description:"count of redirection"`
		EventID   *string            `json:"event_id" bson:"event_id" example:"5f9e9b9b9b9b9b9b9b9b9b98" description:"event id"`
		UserID    *string            `json:"user_id" bson:"user_id" example:"5f9e9b9b9b9b9b9b9b9b9b9b" description:"user id"`
		CreatedAt time.Time          `json:"created_at" bson:"created_at" example:"2020-10-31T00:00:00Z"`
		UpdatedAt *time.Time         `json:"updated_at" bson:"updated_at" example:"2020-10-31T00:00:00Z"`
		DeletedAt *time.Time         `json:"deleted_at" bson:"deleted_at" example:"2020-10-31T00:00:00Z"`
	}

	ShortenURLRequest struct {
		URL     string  `json:"url" example:"http://localhost:8080/test?test=test" description:"url to shorten"`
		EventID *string `json:"event_id" example:"5f9e9b9b9b9b9b9b9b9b9b98" description:"event id"`
		UserID  *string `json:"user_id" example:"5f9e9b9b9b9b9b9b9b9b9b9b" description:"user id"`
	}
)

func (fetchedUrl URLCollection) ComposeURLString() string {
	var redirectUrl = url.URL{
		Scheme: "https",
		Host:   fetchedUrl.Host,
		Path:   fetchedUrl.Path,
	}
	var params = make(url.Values, len(fetchedUrl.Params))
	for k, v := range fetchedUrl.Params {
		params.Add(k, v)
	}
	redirectUrl.RawQuery = params.Encode()
	return redirectUrl.String()
}
