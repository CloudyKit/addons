// mgoStore uses the mgo mongo driver to store sessions on mongo gridfs
package gridFs

import (
	"github.com/CloudyKit/addons/mongo"
	"github.com/CloudyKit/framework/cdi"
	"github.com/CloudyKit/framework/session"
	"github.com/jhsx/qm"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"io"
	"time"
)

// New returns a new store
func New(db, prefix string) session.Store {
	return &Store{
		db:     db,
		prefix: prefix,
	}
}

type Store struct {
	db, prefix string
}

func (sessionStore *Store) gridFs(session *mgo.Session, name string, create bool) (*mgo.GridFile, error) {
	session.SetMode(mgo.Strong, false)
	gridFs := session.DB(sessionStore.db).GridFS(sessionStore.prefix)
	if create {
		gridFs.Remove(name)
		gridFile, err := gridFs.Create(name)
		return gridFile, err
	}
	gridFile, err := gridFs.Open(name)
	return gridFile, err
}

func (sessionStore *Store) Writer(c *cdi.Global, name string) (writer io.WriteCloser, err error) {
	writer, err = sessionStore.gridFs(mongo.GetSession(c), name, true)
	return
}

func (sessionStore *Store) Reader(c *cdi.Global, name string) (reader io.ReadCloser, err error) {
	reader, err = sessionStore.gridFs(mongo.GetSession(c), name, false)
	return
}

func (sessionStore *Store) Remove(c *cdi.Global, name string) (err error) {
	sess := mongo.GetSession(c)
	defer sess.Close()
	return sess.DB(sessionStore.db).GridFS(sessionStore.prefix).Remove(name)
}

func (sessionStore *Store) Gc(c *cdi.Global, before time.Time) {
	sess := mongo.GetSession(c)
	gridFs := sess.DB(sessionStore.db).GridFS(sessionStore.prefix)

	var fileId struct {
		Id bson.ObjectId `bson:"_id"`
	}

	inter := gridFs.Find(qm.Lt("uploadDate", before)).Iter()

	for inter.Next(&fileId) {
		gridFs.RemoveId(fileId.Id)
	}

	inter.Close()
	sess.Close()
}
