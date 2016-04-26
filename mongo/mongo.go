package mongo

import (
	"fmt"
	"github.com/CloudyKit/framework/context"
	"gopkg.in/mgo.v2"
	"reflect"
	"gopkg.in/mgo.v2/bson"
)

type ID struct {
	bson.ObjectId
}

func (b *ID) SetBSON(raw bson.Raw) error {
	return raw.Unmarshal(&b.ObjectId)
}

func (b *ID) GetBSON() (interface{}, error) {
	return b.ObjectId, nil
}

func (b ID) String() string {
	return b.Hex()
}

type IEntity interface {
	Id()
}

func AddEntityManager(c *context.Context, structtyp interface{}, name string) error {

	typeOf := reflect.TypeOf(structtyp)

	if typeOf.Kind() != reflect.Struct {
		panic(fmt.Errorf("Type %q is not struct kind", typeOf))
	}

	c.MapType(structtyp, func(c *context.Context, v reflect.Value) {
		db := c.Get((*mgo.Database)(nil)).(*mgo.Database)
		c = c.Child()
		defer c.Done()
		c.Map(db.C(name))
		c.InjectStructValue(v)
	})

	return nil
}
