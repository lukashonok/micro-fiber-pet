package migrations

import (
	"context"

	migrate "github.com/xakep666/mongo-migrate"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func init() {
	migrate.MustRegister(
		// Up: apply migration
		func(ctx context.Context, db *mongo.Database) error {
			// Create "users" collection
			if err := db.CreateCollection(ctx, "users"); err != nil {
				// ignore if already exists
				if !mongo.IsDuplicateKeyError(err) {
					return err
				}
			}

			// Add index on "email" in users
			userIndex := mongo.IndexModel{
				Keys:    bson.D{{Key: "email", Value: 1}},
				Options: options.Index().SetUnique(true).SetName("idx_users_email"),
			}
			if _, err := db.Collection("users").Indexes().CreateOne(ctx, userIndex); err != nil {
				return err
			}

			// Create "books" collection
			if err := db.CreateCollection(ctx, "books"); err != nil {
				if !mongo.IsDuplicateKeyError(err) {
					return err
				}
			}

			// Add index on "title" in books
			bookIndex := mongo.IndexModel{
				Keys:    bson.D{{Key: "title", Value: 1}},
				Options: options.Index().SetName("idx_books_title"),
			}
			if _, err := db.Collection("books").Indexes().CreateOne(ctx, bookIndex); err != nil {
				return err
			}

			return nil
		},
		// Down: rollback migration
		func(ctx context.Context, db *mongo.Database) error {
			if err := db.Collection("users").Drop(ctx); err != nil {
				return err
			}
			if err := db.Collection("books").Drop(ctx); err != nil {
				return err
			}
			return nil
		},
	)
}
