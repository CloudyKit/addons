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

	menu := startmin.MenuBar{}

	menu.AddMenu(&startmin.MenuItem{
		Name: "dashboard",
		Text: "Dashboard",
		Icon: "dashboard",
	})

	menu.Get("catalog", "Catalog", "opencart").AddMenuList(
		[]*startmin.MenuItem{
			{
				Name: "categories",
				Text: "Categories",
				Children: []*startmin.MenuItem{
					{
						Name: "new-category",
						Text: "New Category",
					},
				},
			},
			{
				Name: "new-product",
				Text: "New Product",
				Icon: "product-hunt",
			},
		},
	)

	view.DefaultSet.AddGlobal("menu", menu)

	app.Default.AddHandlerFunc("GET", "/", func(c *request.Context) {
		view.Render(c.Global, "document", nil)
	})
	app.Default.RunServer(app.Env(":8082", "PORT"))
}
