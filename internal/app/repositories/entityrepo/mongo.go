package entityrepo

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"intro-rest/internal/app/errors"
	"intro-rest/internal/app/interfaces"
	"intro-rest/internal/app/models"
	"intro-rest/internal/app/services/configmanager"
	"intro-rest/internal/app/services/dbclient"
	"math/rand"
	"sync"
	"time"
)

type MongoEntityRepo struct {
	mu              sync.Mutex
	changedEntities map[string]int
	db              *dbclient.MongoDBClient
	config          *configmanager.Config
}

func NewMongoEntityRepo(db *dbclient.MongoDBClient, config *configmanager.Config) *MongoEntityRepo {
	return &MongoEntityRepo{
		changedEntities: make(map[string]int),
		db:              db,
		config:          config,
	}
}

func (r *MongoEntityRepo) Init() error {
	entities, err := r.GetAll()
	if err != nil {
		return err
	}
	rand.Seed(time.Now().UnixNano())
	for i := len(entities); i < r.config.MaxEntities; i++ {
		entity := models.Entity{
			Name: fmt.Sprintf("Entity %d", rand.Int63()),
		}
		_, err = r.Save(&entity)
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *MongoEntityRepo) Save(entity *models.Entity) (string, error) {
	ctx := context.Background()
	client, err := r.db.Connect(ctx)
	defer r.db.Disconnect(ctx)
	if err != nil {
		return "", err
	}

	entity.BeforeSave()
	var entityMongo Entity
	err = entityMongo.Import(entity)
	if err != nil {
		return "", err
	}
	c := client.Database("introvert").Collection("entities")
	filter := bson.D{{"_id", entityMongo.ID}}
	insData := bson.M{
		"name":    entityMongo.Name,
		"created": entityMongo.Created,
	}
	update := bson.D{{"$set", insData}}

	found := c.FindOne(ctx, filter)
	if found.Err() == mongo.ErrNoDocuments {
		res, err := c.InsertOne(ctx, insData)
		if err != nil {
			return "", err
		}
		entity.ID = res.InsertedID.(primitive.ObjectID).Hex()

		r.mu.Lock()
		r.changedEntities[entity.ID] = interfaces.EntityInserted
		r.mu.Unlock()
	} else {
		_, err := c.UpdateOne(ctx, filter, update)
		if err != nil {
			return "", err
		}

		r.mu.Lock()
		r.changedEntities[entity.ID] = interfaces.EntityUpdated
		r.mu.Unlock()
	}
	return entity.ID, nil
}

func (r *MongoEntityRepo) Delete(entityId string) error {
	ctx := context.Background()
	client, err := r.db.Connect(ctx)
	defer r.db.Disconnect(ctx)
	if err != nil {
		return err
	}

	c := client.Database("introvert").Collection("entities")
	id, err := primitive.ObjectIDFromHex(entityId)
	if err != nil {
		return err
	}

	one, err := c.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return err
	}
	if one.DeletedCount == 0 {
		return errors.ErrOnEntityDeleting
	}
	r.mu.Lock()
	r.changedEntities[entityId] = interfaces.EntityDeleted
	r.mu.Unlock()

	return nil
}

func (r *MongoEntityRepo) GetByID(entityID string) (*models.Entity, error) {
	ctx := context.Background()
	client, err := r.db.Connect(ctx)
	defer r.db.Disconnect(ctx)
	if err != nil {
		return nil, err
	}

	id, err := primitive.ObjectIDFromHex(entityID)
	if err != nil {
		return nil, err
	}

	c := client.Database("introvert").Collection("entities")
	result := c.FindOne(ctx, bson.M{"_id": id})

	var mongoEntity Entity
	err = result.Decode(&mongoEntity)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			err = errors.ErrEntityNotFound
		}
		return nil, err
	}

	entity := mongoEntity.Export()
	return entity, nil
}

func (r *MongoEntityRepo) GetAll() ([]models.Entity, error) {
	ctx := context.Background()
	client, err := r.db.Connect(ctx)
	defer r.db.Disconnect(ctx)
	if err != nil {
		return nil, err
	}

	c := client.Database("introvert").Collection("entities")
	cursor, err := c.Find(ctx, bson.M{}, nil)
	if err != nil {
		return nil, err
	}
	var mongoEntities []Entity
	if err = cursor.All(ctx, &mongoEntities); err != nil {
		return nil, err
	}
	results := make([]models.Entity, 0, len(mongoEntities))
	for _, me := range mongoEntities {
		results = append(results, *me.Export())
	}

	return results, nil
}

func (r *MongoEntityRepo) GetInsUpdated(changeStatus int) ([]models.Entity, error) {
	var ids []string
	for id, status := range r.changedEntities {
		if status == changeStatus {
			ids = append(ids, id)
		}
	}
	return r.getMany(ids)
}

func (r *MongoEntityRepo) GetDeleted() []string {
	var ids []string
	for id, status := range r.changedEntities {
		if status == interfaces.EntityDeleted {
			ids = append(ids, id)
		}
	}
	return ids
}

func (r *MongoEntityRepo) CleanChanged(entityId string) {
	r.mu.Lock()
	delete(r.changedEntities, entityId)
	r.mu.Unlock()
}

func (r *MongoEntityRepo) getMany(ids []string) ([]models.Entity, error) {
	if len(ids) == 0 {
		var results []models.Entity
		return results, nil
	}
	ctx := context.Background()
	client, err := r.db.Connect(ctx)
	defer r.db.Disconnect(ctx)
	if err != nil {
		return nil, err
	}
	mongoIds := make([]primitive.ObjectID, 0, len(ids))
	for _, entityId := range ids {
		id, err := primitive.ObjectIDFromHex(entityId)
		if err != nil {
			return nil, err
		}
		mongoIds = append(mongoIds, id)
	}
	c := client.Database("introvert").Collection("entities")
	cursor, err := c.Find(ctx, bson.M{"_id": bson.M{"$in": mongoIds}}, nil)
	if err != nil {
		return nil, err
	}
	var mongoEntities []Entity
	if err = cursor.All(ctx, &mongoEntities); err != nil {
		return nil, err
	}
	results := make([]models.Entity, 0, len(mongoEntities))
	for _, me := range mongoEntities {
		results = append(results, *me.Export())
	}
	return results, nil
}
