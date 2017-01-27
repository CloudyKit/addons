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

package main

import (
	"github.com/CloudyKit/addons/startmin"
	"github.com/CloudyKit/framework/app"
	"github.com/CloudyKit/framework/request"
	"github.com/CloudyKit/framework/view"
)

func main() {
	app.Default.Bootstrap(startmin.Component(""))
	view.DefaultSet.SetDevelopmentMode(true)

	menu := startmin.MenuItem{}

	menu.AppendMenu(&startmin.MenuItem{
		Key:  "dashboard",
		Text: "Dashboard",
		Icon: "dashboard",
	})

	menu.Get("catalog", "Catalog", "opencart").AppendMenus(
		[]*startmin.MenuItem{
			{
				Key:  "categories",
				Text: "Categories",
				Children: []*startmin.MenuItem{
					{
						Key:  "new-category",
						Text: "New Category",
					},
				},
			},
			{
				Key:  "new-product",
				Text: "New Product",
				Icon: "product-hunt",
			},
		},
	)

	view.DefaultSet.AddGlobal("menu", menu)

	app.Default.AddHandlerFunc("GET", "/", func(c *request.Context) {
		view.Render(c.IoC, "document", nil)
	})
	app.Default.RunServer(app.Env(":8082", "PORT"))
}
