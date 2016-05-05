package mongo

import (
	"github.com/CloudyKit/framework/app"
	"github.com/CloudyKit/framework/cdi"
	"gopkg.in/mgo.v2"
	"reflect"
)

var SessionType = reflect.TypeOf((*mgo.Session)(nil))
var DatabaseType = reflect.TypeOf((*mgo.Database)(nil))
var CollectionType = reflect.TypeOf((*mgo.Collection)(nil))

func GetSession(cdi *cdi.DI) *mgo.Session {
	return cdi.Val4Type(SessionType).(*mgo.Session)
}

func GetDatabase(cdi *cdi.DI) *mgo.Database {
	return cdi.Val4Type(DatabaseType).(*mgo.Database)
}

func GetCollection(cdi *cdi.DI) *mgo.Collection {
	return cdi.Val4Type(CollectionType).(*mgo.Collection)
}

type Boot struct {
	URL, DB      string
	Session      *mgo.Session
	Copy, Master bool
}

type magicSession mgo.Session

func (m *magicSession) Provide(di *cdi.DI) interface{} {
	return (*mgo.Session)(m)
}

func (m *magicSession) Finalize() {
	(*mgo.Session)(m).Close()
}

func (pp *Boot) Bootstrap(a *app.App) {

	if pp.Master {

		if pp.Session == nil {
			var err error
			pp.Session, err = mgo.Dial(pp.URL)
			if err != nil {
				panic(err)
			}
		}

		if pp.Copy {
			a.Global.MapType(SessionType, func(cdi *cdi.DI) interface{} {
				s := pp.Session.Copy()
				cdi.MapType(SessionType, (*magicSession)(s))
				return s
			})
		} else {
			a.Global.MapType(SessionType, func(cdi *cdi.DI) interface{} {
				s := pp.Session.Clone()
				cdi.MapType(SessionType, (*magicSession)(s))
				return s
			})
		}
	} else {
		a.Global.MapType(SessionType, func(cdi *cdi.DI) interface{} {
			s, err := mgo.Dial(pp.URL)
			if err != nil {
				panic(err)
			}
			cdi.MapType(SessionType, (*magicSession)(s))
			return s
		})
	}

	a.Global.MapType(DatabaseType, func(cdi *cdi.DI) interface{} {
		return GetSession(cdi).DB(pp.DB)
	})
}
