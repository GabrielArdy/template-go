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

type Course struct {
	ID         primitive.ObjectID `bson:"_id,omitempty"`
	CourseID   string             `bson:"courseId"`
	CourseName string             `bson:"courseName"`
	Instructor string             `bson:"instructor"`
	Materials  []string           `bson:"materials"`
	Students   []string           `bson:"students"`
	MaxQuota   int                `bson:"maxQuota"`
	CreatedAt  time.Time          `bson:"createdAt"`
	UpdatedAt  time.Time          `bson:"updatedAt"`
}

type AssignmentSubmission struct {
	ID           primitive.ObjectID `bson:"_id,omitempty"`
	AssignmentID string             `bson:"assignmentId"`
	StudentID    string             `bson:"studentId"`
	Content      string             `bson:"content"`
	Images       []string           `bson:"images"`
	SubmittedAt  time.Time          `bson:"submittedAt"`
}

type Assignment struct {
	ID           primitive.ObjectID `bson:"_id,omitempty"`
	AssignmentID string             `bson:"assignmentId"`
	CourseID     string             `bson:"courseId"`
	Title        string             `bson:"title"`
	Content      string             `bson:"content"`
	Images       []string           `bson:"images"`
	DueDate      time.Time          `bson:"dueDate"`
}

type Content struct {
	ID           primitive.ObjectID `bson:"_id,omitempty"`
	CourseID     string             `bson:"courseId"`
	Title        string             `bson:"title"`
	Contents     []string           `bson:"content"`
	Images       []string           `bson:"images"`
	AssignmentID []string           `bson:"assignmentId"`
}

type Material struct {
	ID         primitive.ObjectID `bson:"_id,omitempty"`
	MaterialID string             `bson:"materialId"`
	Title      string             `bson:"title"`
	Content    []Content          `bson:"content"`
	CreatedAt  time.Time          `bson:"createdAt"`
	UpdatedAt  time.Time          `bson:"updatedAt"`
}

type GradeMarks struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	StudentID string             `bson:"studentId"`
	CourseID  string             `bson:"courseId"`
	Grade     float64            `bson:"grade"`
	Predicate string             `bson:"predicate"`
}

type CourseEnrollment struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	StudentID string             `bson:"studentId"`
	CourseID  string             `bson:"courseId"`
	Enrolled  time.Time          `bson:"enrolled"`
}
