package repository

import "go.mongodb.org/mongo-driver/mongo"

type CrawlerRepository struct {
	collection *mongo.Collection
}

func NewCrawlerRepository(db *mongo.Database, collectionName string) *CrawlerRepository {
	return &CrawlerRepository{
		collection: db.Collection(collectionName),
	}
}

func (r *CrawlerRepository) InsertOne(crawlResult CrawlResult) error {
	_, err := r.collection.InsertOne(nil, crawlResult)
	if err != nil {
		return err
	}
	return nil
}

func (r *CrawlerRepository) FindAll() ([]CrawlResult, error) {
	cursor, err := r.collection.Find(nil, nil)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(nil)

	var results []CrawlResult
	for cursor.Next(nil) {
		var result CrawlResult
		err := cursor.Decode(&result)
		if err != nil {
			return nil, err
		}
		results = append(results, result)
	}

	return results, nil
}
