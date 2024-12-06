package user

import (
	"context"
	"errors"
	"fmt"

	"github.com/claytonssmint/desafio-leilao-go/configuration/logger"
	"github.com/claytonssmint/desafio-leilao-go/internal/entity/user_entity"
	"github.com/claytonssmint/desafio-leilao-go/internal/internal_error"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserEntityMongo struct {
	Id   string `bson:"_id"`
	Name string `bson:"name"`
}

type UserRepository struct {
	Collection *mongo.Collection
}

func NewUserRepository(database *mongo.Database) *UserRepository {
	return &UserRepository{
		Collection: database.Collection("users"),
	}
}

func (ur *UserRepository) FindUserById(ctx context.Context, userId string) (*user_entity.User, *internal_error.InternalError) {
	filter := bson.M{"_id": userId}

	var UserEntityMongo UserEntityMongo
	err := ur.Collection.FindOne(ctx, filter).Decode(&UserEntityMongo)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			logger.Error(fmt.Sprintf("User not found with this id = %d", userId), err)
			return nil, internal_error.NewNotFoundError(fmt.Sprintf("User not found with this id = %d", userId))
		}

		logger.Error("Error when trying to find user by id", err)
		return nil, internal_error.NewInternalServerError("Error when trying to find user by userId")
	}

	userEntity := &user_entity.User{
		Id:   UserEntityMongo.Id,
		Name: UserEntityMongo.Name,
	}

	return userEntity, nil
}
