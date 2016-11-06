package storage

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type Storage interface {
	Insert(data interface{}) error
	Upsert(key string, value string, data interface{}) error
	Find(key string, value string, outputType interface{}) (interface{}, error)
	Remove(key string, value string) error
	ListAllWithSort(outputType interface{}, sort string) (interface{}, error)
	ListAll(outputType interface{}) (interface{}, error)
	GetRandom(result interface{}) (interface{}, error)
	Close()
}

type (
	storage struct {
		Session    *mgo.Session
		Collection *mgo.Collection
	}
)

type DbConfig struct {
	DatabaseName   string
	CollectionName string
	Url            string
}

func NewStorage(dbConfig DbConfig) (Storage, error) {

	var url = dbConfig.Url
	var collectionName = dbConfig.CollectionName
	var databaseName = dbConfig.DatabaseName

	session, err := mgo.Dial(url)

	if err != nil {
		return nil, err
	}

	collection := session.DB(databaseName).C(collectionName)

	return &storage{session, collection}, nil
}

func (store *storage) Insert(data interface{}) error {
	return store.Collection.Insert(data)
}

func (store *storage) Upsert(key string, value string, data interface{}) error {
	_, err := store.Collection.Upsert(bson.M{key: value}, data)

	return err
}

func (store *storage) Find(key string, value string, outputType interface{}) (interface{}, error) {
	err := store.Collection.Find(bson.M{key: value}).One(outputType)

	return outputType, err
}

func (store *storage) Remove(key string, value string) error {
	err := store.Collection.Remove(bson.M{key: value})

	return err
}

func (store *storage) ListAllWithSort(outputType interface{}, sort string) (interface{}, error) {
	err := store.Collection.Find(nil).Sort(sort).All(outputType)

	return outputType, err
}

func (store *storage) ListAll(outputType interface{}) (interface{}, error) {
	err := store.Collection.Find(nil).All(outputType)

	return outputType, err
}

func (store *storage) GetRandom(result interface{}) (interface{}, error) {
	pipe := store.Collection.Pipe([]bson.M{{"$sample": bson.M{"size": 1}}})

	var err = pipe.One(result)

	if err != nil {
		return nil, err
	}

	return result, nil
}

func (store *storage) Close() {
	store.Session.Close()
}
