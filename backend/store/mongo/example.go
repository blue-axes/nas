package mongo

import (
	"errors"
	"github.com/blue-axes/tmpl/pkg/context"
	"github.com/blue-axes/tmpl/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

type (
	example struct {
		ID   primitive.ObjectID `bson:"_id"`
		Name string             `bson:"name"`
	}
)

func (*example) TableName() string {
	return "example"
}

func (*example) DatabaseName() string {
	return "main"
}

func (m *example) ToExample() types.Example {
	return types.Example{
		ID:   m.ID.Hex(),
		Name: m.Name,
	}
}

func (m *example) FromExample(e types.Example) {
	m.ID = primitive.NewObjectIDFromTimestamp(time.Now())
	m.Name = e.Name
}

func (s *txStore) ListExamples(ctx *context.Context) (res []example, err error) {
	tmpRes, err := s.collection(&example{}).Find(ctx, bson.D{})
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return res, nil
		}
		return nil, err
	}
	err = tmpRes.All(ctx, &res)

	return res, err
}

func (s *txStore) CreateExamples(ctx *context.Context, e *types.Example) error {
	mdl := example{}
	mdl.FromExample(*e)
	tmp, err := s.collection(&example{}).InsertOne(ctx, &mdl)
	if err != nil {
		return err
	}
	id, ok := tmp.InsertedID.(primitive.ObjectID)
	if ok {
		e.ID = id.Hex()
	}
	return err
}
