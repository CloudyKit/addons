package mongo

import (
	"github.com/CloudyKit/framework/app"
	"github.com/CloudyKit/framework/context"
	"gopkg.in/mgo.v2"
)

type plugin struct {
	url, db string
	master  *mgo.Session
	useCopy bool
}

func NewPlugin(url, db string) app.Plugin {
	return &plugin{url: url, db: db}
}

func NewPluginMasterSession(url, db string, useCopy bool) app.Plugin {
	master, _ := mgo.Dial(url)
	return &plugin{db: db, master: master, useCopy: useCopy}
}

type magicSession mgo.Session

func (m *magicSession) Provide(di *context.Context) interface{} {
	return (*mgo.Session)(m)
}

func (m *magicSession) Finalize() {
	(*mgo.Session)(m).Close()
}

func (pp *plugin) Init(c *context.Context) {
	if pp.master == nil {
		c.MapType(pp.master, func(cc *context.Context) interface{} {
			sess, err := mgo.Dial(pp.url)
			if err != nil {
				panic(err)
			}
			cc.MapType(sess, (*magicSession)(sess))
			return sess
		})
	} else {
		c.Map(pp.master)
		if pp.useCopy {
			c.MapType(pp.master, func(c *context.Context) interface{} {
				sess := pp.master.Copy()
				c.MapType(sess, (*magicSession)(sess))
				return sess
			})
		} else {
			c.MapType(pp.master, func(c *context.Context) interface{} {
				sess := pp.master.Clone()
				c.MapType(sess, (*magicSession)(sess))
				return sess
			})
		}
	}

	c.MapType((*mgo.Database)(nil), func(c *context.Context) interface{} {
		return c.Get(pp.master).(*mgo.Session).DB(pp.db)
	})
}
