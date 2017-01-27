package startmin

import (
	"github.com/CloudyKit/framework/app"
	"github.com/CloudyKit/framework/request"
	"github.com/CloudyKit/framework/view"
	"net/http"
	"os"
	"path"
	"runtime"
	"strings"
)

type MenuItem struct {
	Icon     string
	Key      string
	Text     string
	Href     string
	Children []*MenuItem
}

func (m *MenuItem) Get(name, text, icon string) *MenuItem {

	for i := 0; i < len(m.Children); i++ {
		m := m.Children[i]
		if m.Key == name {
			return m
		}
	}

	if text != "" {
		me := &MenuItem{Key: name, Text: text, Icon: icon}
		m.Children = append(m.Children, me)
		return me
	}

	return nil
}

func (m *MenuItem) GetByKey(key string) *MenuItem {
	return m.Get(key, "", "")
}

func (m *MenuItem) AppendMenu(menu ...*MenuItem) {
	m.AppendMenus(menu)
}

func (m *MenuItem) AppendMenus(menu []*MenuItem) {
	m.Children = append(m.Children, menu...)
}

func Component(publicPath string) app.ComponentFunc {
	return func(a *app.App) {

		if publicPath == "" {
			publicPath = "/startmin/public/"
		} else {
			if !strings.HasSuffix(publicPath, "/") {
				publicPath += "/"
			}
			if !strings.HasPrefix(publicPath, "/") {
				publicPath = "/" + publicPath
			}
		}

		jetSet := view.GetJetSet(a.IoC)
		jetSet.AddGlobal("startminPublicPath", a.Prefix+publicPath)
		jetSet.AddGopathPath("github.com/CloudyKit/addons/startmin/templates")

		var root = "./public/startmin"

		if f, err := os.Stat(root); err != nil || !f.IsDir() {
			_, root, _, _ = runtime.Caller(0)
			root = path.Join(path.Dir(root), "public")
		}

		a.AddHandlerFunc("GET", publicPath+"*publicFile", func(c *request.Context) {
			http.ServeFile(c.Response, c.Request, path.Join(root, c.ParamByName("publicFile")))
		})
	}
}
