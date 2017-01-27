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

// gridFs uses the mgo mongo driver to store sessions on mongo gridfs
package gridFs

import (
	"github.com/CloudyKit/addons/mongo"
	"github.com/CloudyKit/framework/container"
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

func (sessionStore *Store) Writer(c *container.IoC, name string) (writer io.WriteCloser, err error) {
	writer, err = sessionStore.gridFs(mongo.LoadSession(c), name, true)
	return
}

func (sessionStore *Store) Reader(c *container.IoC, name string, t time.Time) (reader io.ReadCloser, err error) {
	reader, err = sessionStore.gridFs(mongo.LoadSession(c), name, false)
	return
}

func (sessionStore *Store) Remove(c *container.IoC, name string) (err error) {
	sess := mongo.LoadSession(c)
	defer sess.Close()
	return sess.DB(sessionStore.db).GridFS(sessionStore.prefix).Remove(name)
}

func (sessionStore *Store) GC(c *container.IoC, before time.Time) {
	sess := mongo.LoadSession(c)
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
