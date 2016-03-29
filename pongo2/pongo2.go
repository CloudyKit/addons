package pongo2

import (
	"github.com/CloudyKit/framework/context"
	. "github.com/CloudyKit/framework/view"
	"gopkg.in/flosch/pongo2.v3"
	"io"
)

type pongo2plugin struct {
	Manager    *Manager
	extensions []string
	loader     *pongo2loader
}

func NewPlugin(baseDir string, extensions ...string) *pongo2plugin {
	return &pongo2plugin{loader: newPongo2Loader(baseDir), extensions: extensions}
}

func (plugin *pongo2plugin) Init(di *context.Context) {
	plugin.Manager.AddLoader(plugin.loader, plugin.extensions...)
}

func newPongo2Loader(baseDir string) *pongo2loader {
	newSet := pongo2.NewSet("viewloader")
	newSet.SetBaseDirectory(baseDir)
	return &pongo2loader{set: newSet}
}

type pongo2loader struct {
	set *pongo2.TemplateSet
}

type pongoRender pongo2.Template

func (tt *pongoRender) Execute(w io.Writer, c Data) error {
	return (*pongo2.Template)(tt).ExecuteWriter(pongo2.Context(c), w)
}

func (viewLoader *pongo2loader) View(name string) (view ViewRenderer, err error) {
	var viewRaw *pongo2.Template
	viewRaw, err = viewLoader.set.FromCache(name)
	view = (*pongoRender)(viewRaw)
	return
}
