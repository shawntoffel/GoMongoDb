package storage

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type (
	Storage struct {
		Session    *mgo.Session
		Collection *mgo.Collection
	}
)

type DbConfig struct {
	DatabaseName   string
	CollectionName string
	Url            string
}

func (storage *Storage) NewStorage(dbConfig DbConfig) (*Storage, error) {

	var url = dbConfig.Url
	var collectionName = dbConfig.CollectionName
	var databaseName = dbConfig.DatabaseName

	session, err := mgo.Dial(url)

	if err != nil {
		return nil, err
	}

	collection := session.DB(databaseName).C(collectionName)

	return &Storage{session, collection}, nil
}

func (storage *Storage) Insert(data interface{}) error {
	return storage.Collection.Insert(data)
}

func (storage *Storage) Find(key string, value string, outputType interface{}) interface{} {
	storage.Collection.Find(bson.M{key: value}).One(&outputType)

	return outputType
}

func (storage *Storage) GetRandom(result interface{}) (interface{}, error) {
	pipe := storage.Collection.Pipe([]bson.M{{"$sample": bson.M{"size": 1}}})

	var err = pipe.One(result)

	if err != nil {
		return nil, err
	}

	return result, nil
}

func (storage *Storage) Close() {
	storage.Session.Close()
}
