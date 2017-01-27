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
