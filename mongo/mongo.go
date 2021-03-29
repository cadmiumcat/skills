package mongo

import (
	"context"
	"errors"
	dpMongodb "github.com/ONSdigital/dp-mongodb"
	dpMongoLock "github.com/ONSdigital/dp-mongodb/dplock"
	"github.com/ONSdigital/log.go/log"
	"github.com/cadmiumcat/books-api/apierrors"
	"github.com/cadmiumcat/books-api/config"
	"github.com/cadmiumcat/books-api/models"
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
)

// Mongo contains the information needed to create and interact with a mongo session
type Mongo struct {
	Collection string
	Database   string
	Session    *mgo.Session
	URI        string
	lockClient *dpMongoLock.Lock
}

// Init initialises a mongo session with the given configuration.
// It returns an error if the session already exists or if it cannot connect.
func (m *Mongo) Init(mongoConfig config.MongoConfig) (err error) {
	if m.Session != nil {
		return errors.New("session already exists")
	}

	if m.Session, err = mgo.Dial(mongoConfig.BindAddr); err != nil {
		return err
	}

	m.Collection = mongoConfig.Collection
	m.Database = mongoConfig.Database

	return nil
}

// Close closes the mongo session and returns any error
func (m *Mongo) Close(ctx context.Context) (err error) {
	m.lockClient.Close(ctx)
	return dpMongodb.Close(ctx, m.Session)
}

// AddBook adds a Book
func (m *Mongo) AddBook(book *models.Book) {
	session := m.Session.Copy()
	defer session.Close()

	collection := session.DB(m.Database).C(m.Collection)
	collection.Insert(book)

	return
}

// GetBook returns a models.Book for a given ID.
// It returns an error if the Book is not found
func (m *Mongo) GetBook(ctx context.Context, ID string) (*models.Book, error) {
	session := m.Session.Copy()
	defer session.Close()

	logData := log.Data{
		"book_id": ID,
		"database": m.Database,
		"collection": m.Collection}

	var book models.Book
	err := session.DB(m.Database).C(m.Collection).Find(bson.M{"_id": ID}).One(&book)

	if err != nil {
		if err == mgo.ErrNotFound {
			log.Event(ctx, apierrors.ErrBookNotFound.Error(), log.ERROR, log.Error(err), logData)
			return nil, apierrors.ErrBookNotFound
		}
		return nil, err
	}

	return &book, nil
}

// GetBooks returns all the existing models.Books.
// It returns an error if the models.Books cannot be listed.
func (m *Mongo) GetBooks(ctx context.Context) (models.Books, error) {

	session := m.Session.Copy()
	defer session.Close()

	logData := log.Data{
		"database": m.Database,
		"collection": m.Collection}

	list := session.DB(m.Database).C(m.Collection).Find(nil)

	books := &models.Books{}
	if err := list.All(&books.Items); err != nil {
		log.Event(ctx, "unable to retrieve books", log.ERROR, log.Error(err), logData)
		return models.Books{}, err
	}

	return *books, nil
}
