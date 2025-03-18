package books

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/theotruvelot/books/model"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type BookRepository struct {
	collection *mongo.Collection
	dbName     string
}

func NewBookRepository(db *mongo.Client, dbName string) *BookRepository {
	fmt.Printf("Creating BookRepository with database: %s\n", dbName)
	return &BookRepository{
		collection: db.Database(dbName).Collection("books"),
		dbName:     dbName,
	}
}

func (r *BookRepository) CreateBook(ctx context.Context, book *model.Book) error {
	if book.ID == uuid.Nil {
		book.ID = uuid.New()
		fmt.Printf("Generated new UUID: %s\n", book.ID.String())
	}

	fmt.Printf("Inserting book into %s.books with ID: %s: %+v\n", r.dbName, book.ID.String(), book)
	_, err := r.collection.InsertOne(ctx, book)
	return err
}

func (r *BookRepository) GetBook(ctx context.Context, id uuid.UUID) (*model.Book, error) {
	var book model.Book
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&book)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		}
		return nil, err
	}
	return &book, nil
}

func (r *BookRepository) UpdateBook(ctx context.Context, book *model.Book) error {
	_, err := r.collection.ReplaceOne(ctx, bson.M{"_id": book.ID}, book)
	return err
}

func (r *BookRepository) DeleteBook(ctx context.Context, id uuid.UUID) error {
	_, err := r.collection.DeleteOne(ctx, bson.M{"_id": id})
	return err
}

func (r *BookRepository) ListBooks(ctx context.Context, page, pageSize int) ([]*model.Book, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}

	skip := (page - 1) * pageSize

	opts := options.Find().
		SetLimit(int64(pageSize)).
		SetSkip(int64(skip))

	cursor, err := r.collection.Find(ctx, bson.M{}, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var books []*model.Book
	if err := cursor.All(ctx, &books); err != nil {
		return nil, err
	}

	return books, nil
}
