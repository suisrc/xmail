package mgo_test

import (
	"context"
	"errors"
	"testing"
	"time"
	"vkc/mgo"

	"github.com/go-playground/assert/v2"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	//readpref "go.mongodb.org/mongo-driver/mongo/readpref"
)

func TestStore(t *testing.T) {
	mgo.C.Mongodb.Host = "web3-mongo"
	mgo.C.Mongodb.Port = "27017"
	mgo.C.Mongodb.Username = "xx"
	mgo.C.Mongodb.Password = "xx"
	mgo.C.Mongodb.Database = "test"
	mgo.C.Mongodb.RawOptions = "w=majority&authSource=admin"

	dbx, clx, err := mgo.NewDefaultDatabase()
	assert.Equal(t, nil, err)
	defer clx()
	storer := &Storer{cll: dbx.Collection("test")}

	key := "test"
	ctx := context.Background()

	// err = storer.Set(ctx, key, "test", 60*time.Second)
	// assert.Equal(t, nil, err)

	b, err := storer.Exists(ctx, key)
	assert.Equal(t, nil, err)
	assert.Equal(t, true, b)

	val, b, err := storer.Get(ctx, key)
	assert.Equal(t, nil, err)
	assert.Equal(t, true, b)
	assert.Equal(t, "test", val)

	err = storer.Delete(ctx, key)
	assert.Equal(t, nil, err)

}

// Storer 存储
type Storer struct {
	cll *mongo.Collection
}

// Get ... 从存储中获取数据
func (s *Storer) Get(ctx context.Context, key string) (string, bool, error) {
	filter := bson.M{"id": key}
	res := s.cll.FindOne(ctx, filter)
	if err := res.Err(); err != nil {
		return "", false, err
	}
	data := bson.M{}
	if err := res.Decode(&data); err != nil {
		return "", false, err
	}
	return data["value"].(string), true, nil
}

// Set ... 设置存储
func (s *Storer) Set(ctx context.Context, key, value string, expiration time.Duration) error {
	data := bson.M{"id": key, "value": value, "created_time": time.Now(), "expired_time": time.Now().Add(expiration)}
	_, err := s.cll.InsertOne(ctx, data)
	if err != nil {
		return err
	}
	return nil
}

// Exists ... 判断是否存在
func (s *Storer) Exists(ctx context.Context, key string) (bool, error) {
	_, exist, err := s.Get(ctx, key)
	if err != nil {
		return false, err
	}
	return exist, nil
}

// Delete ... 删除
func (s *Storer) Delete(ctx context.Context, key string) error {
	filter := bson.M{"id": key}
	res, err := s.cll.DeleteOne(ctx, filter)
	if err != nil {
		return err
	} else if res.DeletedCount == 0 {
		return errors.New("not exist")
	}
	return nil
}

// =====================================================
// 将本地文件中的用户同步到mongodb中

var (
	mgo_path = "shconf/__mongo.json"
	mgo_coll = "user"
)

func TestMongo(t *testing.T) {
	cli, clx, err := mgo.NewDatabaseByFile(mgo_path)
	if err != nil {
		panic(err) // 直接终止程序
	}
	defer clx()
	// 连接到指定的集合
	cll := cli.Collection(mgo_coll)
	ctx := context.TODO()
	// 增加索引
	CreateIndex(cll, ctx)

}

// CreateIndex 创建索引
func CreateIndex(cll *mongo.Collection, ctx context.Context) {
	idx := mongo.IndexModel{Keys: bson.M{"username": 1}, Options: options.Index().SetUnique(true).SetSparse(true).SetName("udx_username")}
	_, err := cll.Indexes().CreateOne(ctx, idx)
	if err != nil {
		panic(err)
	}
	logrus.Info("create index success")
}
