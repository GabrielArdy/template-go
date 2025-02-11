package repository

import "go.mongodb.org/mongo-driver/mongo"

type CourseRepository struct {
	collection *mongo.Collection
	MaterialRepository
}

func NewCourseRepository(db *mongo.Database, collectionName string, mrp MaterialRepository) *CourseRepository {
	return &CourseRepository{
		collection:         db.Collection(collectionName),
		MaterialRepository: mrp,
	}
}

func (r *CourseRepository) InsertCourse(course Course) error {
	_, err := r.collection.InsertOne(nil, course)
	if err != nil {
		return err
	}
	return nil
}

func (r *CourseRepository) FindAllCourses() ([]Course, error) {
	cursor, err := r.collection.Find(nil, nil)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(nil)

	var results []Course
	for cursor.Next(nil) {
		var result Course
		err := cursor.Decode(&result)
		if err != nil {
			return nil, err
		}
		results = append(results, result)
	}

	return results, nil
}

func (r *CourseRepository) FindCourseById(id string) (Course, error) {
	var course Course
	err := r.collection.FindOne(nil, id).Decode(&course)
	if err != nil {
		return Course{}, err
	}
	return course, nil
}

func (r *CourseRepository) FindCourseByInstructor(instructor string) ([]Course, error) {
	cursor, err := r.collection.Find(nil, instructor)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(nil)

	var results []Course
	for cursor.Next(nil) {
		var result Course
		err := cursor.Decode(&result)
		if err != nil {
			return nil, err
		}
		results = append(results, result)
	}

	return results, nil
}

func (r *CourseRepository) DeleteCourse(id string) error {
	_, err := r.collection.DeleteOne(nil, id)
	if err != nil {
		return err
	}
	return nil
}

func (r *CourseRepository) UpdateCourse(course Course) error {
	_, err := r.collection.UpdateOne(nil, course.CourseID, course)
	if err != nil {
		return err
	}
	return nil
}

func (r *CourseRepository) InsertMaterialByCourseID(courseId string, materialId string) error {
	course, err := r.FindCourseById(courseId)
	if err != nil {
		return err
	}

	course.Materials = append(course.Materials, materialId)
	err = r.UpdateCourse(course)
	if err != nil {
		return err
	}
	return nil
}
