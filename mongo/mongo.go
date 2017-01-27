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
	"fmt"
	"github.com/CloudyKit/framework/container"
	"reflect"
)

func MapManagerToCollection(c *container.IoC, tv interface{}, name string) error {

	typeOf, ok := tv.(reflect.Type)

	if !ok {
		typeOf = reflect.TypeOf(tv)
	}

	if typeOf.Kind() != reflect.Struct {
		panic(fmt.Errorf("type %q is not of kind Struct", typeOf))
	}

	c.MapInitializerFunc(typeOf, func(c *container.IoC, v reflect.Value) {
		db := LoadDatabase(c)
		collection := c.LoadType(CollectionType) // backups previous collection
		c.MapValue(CollectionType, db.C(name))   // map new collection
		c.InjectValue(v)                         // injects new collection into the value
		c.MapValue(CollectionType, collection)   // restore backed collection
	})

	return nil
}
