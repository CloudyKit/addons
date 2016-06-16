package mongo

import (
	"fmt"
	"github.com/CloudyKit/framework/cdi"
	"reflect"
)

func AddEntityManager(c *cdi.Global, zeroval interface{}, name string) error {

	typeOf, ok := zeroval.(reflect.Type)

	if !ok {
		typeOf = reflect.TypeOf(zeroval)
	}

	if typeOf.Kind() != reflect.Struct {
		panic(fmt.Errorf("Type %q is not struct kind", typeOf))
	}

	c.MapType(typeOf, func(c *cdi.Global, v reflect.Value) {
		db := GetDatabase(c)
		collection := c.Val4Type(CollectionType) // backups previous collection
		c.MapType(CollectionType, db.C(name))    // map new collection
		c.InjectInStructValue(v)                 // injects new collection into the value
		c.MapType(CollectionType, collection)    // restore backed collection
	})

	return nil
}
