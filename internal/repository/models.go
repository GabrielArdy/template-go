package repository

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserDocs struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	UID       string             `bson:"uid"`
	FirstName string             `bson:"firstName"`
	LastName  string             `bson:"lastName"`
	TeacherID string             `bson:"teacherID"`
	Email     string             `bson:"email"`
	Grade     string             `bson:"grade"`
	Position  string             `bson:"position"`
	Status    string             `bson:"status"`
	CreatedAt time.Time          `bson:"createdAt"`
	UpdatedAt time.Time          `bson:"updatedAt"`
}

type AuthDocs struct {
	ID         primitive.ObjectID `bson:"_id,omitempty"`
	UID        string             `bson:"uid"`
	Email      string             `bson:"email"`
	Password   string             `bson:"password"`
	TeacherID  string             `bson:"teacherID"`
	Role       string             `bson:"role"`
	CreatedAt  time.Time          `bson:"createdAt"`
	UpdatedAt  time.Time          `bson:"updatedAt"`
	LastActive time.Time          `bson:"lastActive"`
}

type AttendanceDocs struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	UID       string             `bson:"uid"`
	Date      time.Time          `bson:"date"`
	CheckIn   time.Time          `bson:"checkIn"`
	CheckOut  time.Time          `bson:"checkOut"`
	Status    string             `bson:"status"`
	TeacherId string             `bson:"TeacherId"`
}

type ActivityDocs struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	UID       string             `bson:"uid"`
	Type      string             `bson:"type"`
	UserId    string             `bson:"userId"`
	CreatedAt time.Time          `bson:"createdAt"`
}
