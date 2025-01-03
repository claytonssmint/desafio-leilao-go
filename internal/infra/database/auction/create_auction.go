package auction

import (
	"context"
	"os"
	"time"

	"github.com/claytonssmint/desafio-leilao-go/configuration/logger"
	"github.com/claytonssmint/desafio-leilao-go/internal/entity/auction_entity"
	"github.com/claytonssmint/desafio-leilao-go/internal/internal_error"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type AuctionEntityMongo struct {
	Id          string                          `bson:"_id"`
	ProductName string                          `bson:"product_name"`
	Category    string                          `bson:"category"`
	Description string                          `bson:"description"`
	Condition   auction_entity.ProductCondition `bson:"condition"`
	Status      auction_entity.AuctionStatus    `bson:"status"`
	Timestamp   int64                           `bson:"timestamp"`
}

type AuctionRepository struct {
	Collection *mongo.Collection
}

func NewAuctionRepository(database *mongo.Database) *AuctionRepository {
	repo := &AuctionRepository{
		Collection: database.Collection("auctions"),
	}
	go repo.startAuctionClosingRoutine()
	return repo
}

func (ar *AuctionRepository) CreateAuction(ctx context.Context, auctionEntity *auction_entity.Auction) *internal_error.InternalError {
	auctionEntityMongo := &AuctionEntityMongo{
		Id:          auctionEntity.Id,
		ProductName: auctionEntity.ProductName,
		Category:    auctionEntity.Category,
		Description: auctionEntity.Description,
		Condition:   auctionEntity.Condition,
		Status:      auctionEntity.Status,
		Timestamp:   auctionEntity.Timestamp.Unix(),
	}

	_, err := ar.Collection.InsertOne(ctx, auctionEntityMongo)
	if err != nil {
		logger.Error("Error trying to insert auction", err)
		return internal_error.NewInternalServerError("Error trying to insert auction")
	}

	return nil
}

func getAuctionDuration() time.Duration {
	durationStr := os.Getenv("AUCTION_DURATION")
	duration, err := time.ParseDuration(durationStr)
	if err != nil {
		return 24 * time.Hour
	}
	return duration
}

func (ar *AuctionRepository) startAuctionClosingRoutine() {
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			ar.closeExpiredAuctions()
		}
	}
}

func (ar *AuctionRepository) closeExpiredAuctions() {
	ctx := context.Background()
	now := time.Now().Unix()
	filter := bson.M{
		"status":    auction_entity.Active,
		"timestamp": bson.M{"$lt": now - int64(getAuctionDuration().Seconds())},
	}
	update := bson.M{
		"$set": bson.M{"status": auction_entity.Completed},
	}

	_, err := ar.Collection.UpdateMany(ctx, filter, update)
	if err != nil {
		logger.Error("Error trying to close expired auctions", err)
	}
}
