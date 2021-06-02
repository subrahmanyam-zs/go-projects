package dbmigration

import (
	"context"
	"strconv"
	"time"

	"developer.zopsmart.com/go/gofr/pkg/datastore"
	"developer.zopsmart.com/go/gofr/pkg/errors"
	"developer.zopsmart.com/go/gofr/pkg/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Mongo struct {
	datastore.MongoDB
	coll          *mongo.Collection
	newMigrations []gofrMigration // tracks all the migrations run with it's startTime and endTime
}

func NewMongo(m datastore.MongoDB) *Mongo {
	coll := m.Collection("gofr_migrations")
	return &Mongo{m, coll, make([]gofrMigration, 0)}
}

func (md *Mongo) preRun(app, method, name string) error {
	if md.isDirty(app) {
		return &errors.Response{Reason: "dirty migration check failed"}
	}

	ver, _ := strconv.Atoi(name)

	md.newMigrations = append(md.newMigrations, gofrMigration{
		App:       app,
		Version:   int64(ver),
		StartTime: time.Now(),
		Method:    method,
	})

	return nil
}

func (md *Mongo) Run(m Migrator, app, name, method string, logger log.Logger) (err error) {
	err = md.preRun(app, method, name)
	if err != nil {
		return err
	}

	ds := &datastore.DataStore{MongoDB: md}

	if method == UP {
		err = m.Up(ds, logger)
	} else {
		err = m.Down(ds, logger)
	}

	if err != nil {
		md.deleteRow(app, method, name)
		return err
	}

	err = md.postRun(app, method, name)
	if err != nil {
		md.deleteRow(app, method, name)
		return err
	}

	return
}

func (md *Mongo) postRun(app, method, name string) error {
	ver, _ := strconv.Atoi(name)

	for i, v := range md.newMigrations {
		if v.App == app && v.Method == method && v.Version == int64(ver) {
			md.newMigrations[i].EndTime = time.Now()
		}
	}

	return nil
}

func (md *Mongo) LastRunVersion(app, method string) int {
	mt := gofrMigration{}
	opts := options.FindOne().SetSort(bson.D{{Key: "version", Value: -1}})

	err := md.coll.FindOne(context.TODO(), bson.D{{Key: "app", Value: app}, {Key: "method", Value: method}}, opts).Decode(&mt)
	if err != nil {
		return 0
	}

	return int(mt.Version)
}

func (md *Mongo) isDirty(app string) bool {
	mt := gofrMigration{}
	_ = md.coll.FindOne(context.TODO(), bson.D{{Key: "app", Value: app}, {Key: "endtime", Value: time.Time{}}}).Decode(&mt)

	return mt.Version != 0
}

func (md *Mongo) deleteRow(app, method, name string) {
	ver, _ := strconv.Atoi(name)
	_, _ = md.coll.DeleteOne(context.TODO(), bson.D{{Key: "app", Value: app}, {Key: "method", Value: method}, {Key: "version", Value: ver}})
}

func (md *Mongo) GetAllMigrations(app string) (upMigrations, downMigrations []int) {
	var mt gofrMigration

	cur, err := md.coll.Find(context.TODO(), bson.D{{Key: "app", Value: app}}, nil)
	if err != nil {
		return
	}

	for cur.Next(context.TODO()) {
		err := cur.Decode(&mt)
		if err != nil {
			return
		}

		if mt.Method == UP {
			upMigrations = append(upMigrations, int(mt.Version))
		} else {
			downMigrations = append(downMigrations, int(mt.Version))
		}
	}

	return
}

func (md *Mongo) FinishMigration() error {
	for _, v := range md.newMigrations {
		_, _ = md.coll.InsertOne(context.TODO(), gofrMigration{v.App, v.Version, v.StartTime, v.EndTime, v.Method})
	}

	return nil
}
