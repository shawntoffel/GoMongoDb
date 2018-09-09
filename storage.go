package storage

import (
	"gopkg.in/mgo.v2"
)

type Storage interface {
	Close()
}

type Store struct {
	Session    *mgo.Session
	Collection *mgo.Collection
}

type DbConfig struct {
	DatabaseName   string
	CollectionName string
	Url            string
}

func NewStorage(dbConfig DbConfig) (*Store, error) {
	url := dbConfig.Url
	if url == "" {
		url = "localhost"
	}

	session, err := mgo.Dial(url)
	if err != nil {
		return nil, err
	}

	collection := session.DB(dbConfig.DatabaseName).C(dbConfig.CollectionName)

	return &Store{session, collection}, nil
}

func (store *Store) Close() {
	store.Session.Close()
}
