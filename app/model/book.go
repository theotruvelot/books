package model

import (
	"time"

	"github.com/google/uuid"
)

type Book struct {
	ID               uuid.UUID `json:"id" bson:"_id" validate:"omitempty"`
	Title            string    `json:"title" bson:"title" validate:"required"`
	ISBN             string    `json:"isbn" bson:"isbn" validate:"required"`
	PageCount        int       `json:"pageCount" bson:"pageCount" validate:"required,gt=0"`
	PublishedDate    time.Time `json:"publishedDate" bson:"publishedDate" validate:"required"`
	ThumbnailURL     string    `json:"thumbnailUrl" bson:"thumbnailUrl" validate:"omitempty,url"`
	ShortDescription string    `json:"shortDescription" bson:"shortDescription" validate:"omitempty"`
	LongDescription  string    `json:"longDescription" bson:"longDescription" validate:"omitempty"`
	Status           string    `json:"status" bson:"status" validate:"required,oneof=PUBLISH MEAP UPCOMING"`
	Authors          []string  `json:"authors" bson:"authors" validate:"required,min=1,dive,required"`
	Categories       []string  `json:"categories" bson:"categories" validate:"required,min=1,dive,required"`
}
