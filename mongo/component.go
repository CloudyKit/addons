package mongo

import (
	"github.com/CloudyKit/framework/app"
	"github.com/CloudyKit/framework/assert"
	"github.com/CloudyKit/framework/cdi"

	"gopkg.in/mgo.v2"
	"reflect"
)

var SessionType = reflect.TypeOf((*mgo.Session)(nil))
var DatabaseType = reflect.TypeOf((*mgo.Database)(nil))
var CollectionType = reflect.TypeOf((*mgo.Collection)(nil))

func GetSession(cdi *cdi.Global) *mgo.Session {
	return cdi.GetByType(SessionType).(*mgo.Session)
}

func GetDatabase(cdi *cdi.Global) *mgo.Database {
	return cdi.GetByType(DatabaseType).(*mgo.Database)
}

func GetCollection(cdi *cdi.Global) *mgo.Collection {
	return cdi.GetByType(CollectionType).(*mgo.Collection)
}

type Component struct {
	URL, DB      string
	DialInfo     *mgo.DialInfo
	Session      *mgo.Session
	Copy, Master bool
}

type magicSession mgo.Session

func (m *magicSession) Provide(di *cdi.Global) interface{} {
	return (*mgo.Session)(m)
}

func (m *magicSession) Finalize() {
	(*mgo.Session)(m).Close()
}

func (pp *Component) session() (*mgo.Session, error) {
	if pp.DialInfo != nil {
		return mgo.DialWithInfo(pp.DialInfo)
	}
	return mgo.Dial(pp.URL)
}

func (pp *Component) Bootstrap(a *app.App) {

	if pp.Master {

		if pp.Session == nil {
			var err error
			pp.Session, err = pp.session()
			if err != nil {
				panic(err)
			}
		}

		if pp.Copy {
			a.Global.MapType(SessionType, func(cdi *cdi.Global) interface{} {
				s := pp.Session.Copy()
				cdi.MapType(SessionType, (*magicSession)(s))
				return s
			})
		} else {
			a.Global.MapType(SessionType, func(cdi *cdi.Global) interface{} {
				s := pp.Session.Clone()
				cdi.MapType(SessionType, (*magicSession)(s))
				return s
			})
		}

	} else {
		a.Global.MapType(SessionType, func(cdi *cdi.Global) interface{} {
			s, err := pp.session()
			assert.NilErr(err)
			cdi.MapType(SessionType, (*magicSession)(s))
			return s
		})
	}

	a.Global.MapType(DatabaseType, func(cdi *cdi.Global) interface{} {
		return GetSession(cdi).DB(pp.DB)
	})
}
