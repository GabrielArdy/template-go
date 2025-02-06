package repository

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CrawlResult struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	URL       string             `bson:"url"`
	Title     string             `bson:"title"`
	Content   string             `bson:"content"`
	Links     []string           `bson:"links"`
	CrawledAt time.Time          `bson:"crawledAt"`
}

type UserAuth struct {
	ID           primitive.ObjectID `bson:"_id,omitempty"`
	UserID       string             `bson:"userId"`
	Username     string             `bson:"username"`
	Email        string             `bson:"email"`
	Password     string             `bson:"password"`
	LastLogin    time.Time          `bson:"lastLogin"`
	RefreshToken string             `bson:"refreshToken"`
	AccessToken  string             `bson:"accessToken"`
}

type User struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	UserID    string             `bson:"userId"`
	FirstName string             `bson:"firstName"`
	LastName  string             `bson:"lastName"`
	Phone     string             `bson:"phone"`
	Email     string             `bson:"email"`
	CreatedAt time.Time          `bson:"createdAt"`
	UpdatedAt time.Time          `bson:"updatedAt"`
}
