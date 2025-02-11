package repository

import "go.mongodb.org/mongo-driver/mongo"

type MaterialRepository struct {
	collection *mongo.Collection
}

func NewMaterialRepository(db *mongo.Database, collectionName string) *MaterialRepository {
	return &MaterialRepository{
		collection: db.Collection(collectionName),
	}
}

func (r *MaterialRepository) InsertMaterialByCourseID(courseId string, material Material) error {
	_, err := r.collection.InsertOne(nil, material)
	if err != nil {
		return err
	}
	return nil
}

func (r *MaterialRepository) FindMaterialByCourseID(courseId string) ([]Material, error) {
	cursor, err := r.collection.Find(nil, courseId)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(nil)

	var results []Material
	for cursor.Next(nil) {
		var result Material
		err := cursor.Decode(&result)
		if err != nil {
			return nil, err
		}
		results = append(results, result)
	}

	return results, nil
}
