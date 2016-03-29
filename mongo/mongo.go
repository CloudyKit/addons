package mongo

import (
	"github.com/CloudyKit/framework/context"
	"gopkg.in/mgo.v2"
	"reflect"
)

func Mapto(c *context.Context, name string, typ interface{}) error {
	c.MapType(typ, func(c *context.Context, v reflect.Value) {
		db := c.Get((*mgo.Database)(nil)).(*mgo.Database)
		c = c.Child()
		c.Map(db.C(name))
		c.InjectStructValue(v)
		c.Done()
	})
	return nil
}
