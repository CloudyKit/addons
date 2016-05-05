package mongo

import (
	"fmt"
	"github.com/CloudyKit/framework/cdi"
	"gopkg.in/mgo.v2/bson"
	"reflect"
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

func AddEntityManager(c *cdi.DI, zeroval interface{}, name string) error {

	typeOf, ok := zeroval.(reflect.Type)

	if !ok {
		typeOf = reflect.TypeOf(zeroval)
	}

	if typeOf.Kind() != reflect.Struct {
		panic(fmt.Errorf("Type %q is not struct kind", typeOf))
	}

	c.MapType(typeOf, func(c *cdi.DI, v reflect.Value) {
		db := GetDatabase(c)
		collection := c.Val4Type(CollectionType) // backups previous collection
		c.MapType(CollectionType, db.C(name))    // map new collection
		c.InjectInStructValue(v)                 // injects new collection into the value
		c.MapType(CollectionType, collection)    // restore backuped collection
	})

	return nil
}
