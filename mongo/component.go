// MIT License
//
// Copyright (c) 2017 José Santos <henrique_1609@me.com>
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

package mongo

import (
	"github.com/CloudyKit/framework/app"
	"github.com/CloudyKit/framework/container"
	"github.com/CloudyKit/framework/ensure"

	"gopkg.in/mgo.v2"
	"reflect"
)

var SessionType = reflect.TypeOf((*mgo.Session)(nil))
var DatabaseType = reflect.TypeOf((*mgo.Database)(nil))
var CollectionType = reflect.TypeOf((*mgo.Collection)(nil))

func LoadSession(cdi *container.IoC) *mgo.Session {
	return cdi.LoadType(SessionType).(*mgo.Session)
}

func LoadDatabase(cdi *container.IoC) *mgo.Database {
	return cdi.LoadType(DatabaseType).(*mgo.Database)
}

func LoadCollection(cdi *container.IoC) *mgo.Collection {
	return cdi.LoadType(CollectionType).(*mgo.Collection)
}

type Component struct {
	URL, DB      string
	DialInfo     *mgo.DialInfo
	Session      *mgo.Session
	Copy, Master bool
}

type magicSession mgo.Session

func (m *magicSession) Provide(di *container.IoC) interface{} {
	return (*mgo.Session)(m)
}

func (m *magicSession) Dispose() {
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
			a.IoC.MapProviderFunc(SessionType, func(cdi *container.IoC) interface{} {
				s := pp.Session.Copy()
				cdi.MapProvider(SessionType, (*magicSession)(s))
				return s
			})
		} else {
			a.IoC.MapProviderFunc(SessionType, func(cdi *container.IoC) interface{} {
				s := pp.Session.Clone()
				cdi.MapProvider(SessionType, (*magicSession)(s))
				return s
			})
		}

	} else {
		a.IoC.MapProviderFunc(SessionType, func(cdi *container.IoC) interface{} {
			s, err := pp.session()
			ensure.NilErr(err)
			cdi.MapProvider(SessionType, (*magicSession)(s))
			return s
		})
	}

	a.IoC.MapProviderFunc(DatabaseType, func(cdi *container.IoC) interface{} {
		return LoadSession(cdi).DB(pp.DB)
	})

}
