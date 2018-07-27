package template

import (
	"fmt"
	"io"

	"github.com/klauspost/compress/gzip"

	"github.com/kataras/iris/config"
	"github.com/kataras/iris/context"
	"github.com/kataras/iris/render/template/engine/amber"
	"github.com/kataras/iris/render/template/engine/html"
	"github.com/kataras/iris/render/template/engine/jade"
	"github.com/kataras/iris/render/template/engine/markdown"
	"github.com/kataras/iris/render/template/engine/pongo"
	"github.com/kataras/iris/utils"
)

type (
	Engine interface {
		BuildTemplates() error
		ExecuteWriter(out io.Writer, name string, binding interface{}, layout string) error
	}

	Template struct {
		Engine        Engine
		IsDevelopment bool
		Gzip          bool
		ContentType   string
		Layout        string
		buffer        *utils.BufferPool // this is used only for RenderString
	}
)

// New creates and returns a Template instance which keeps the Template Engine and helps with render
func New(c config.Template) *Template {

	var e Engine
	// [ENGINE-2]
	switch c.Engine {
	case config.HTMLEngine:
		e = html.New(c) //  HTMLTemplate
	case config.PongoEngine:
		e = pongo.New(c) // Pongo2
	case config.MarkdownEngine:
		e = markdown.New(c) // Markdown
	case config.JadeEngine:
		e = jade.New(c) // Jade
	case config.AmberEngine:
		e = amber.New(c)
	default: // it's the config.NoEngine
		return nil
	}

	if err := e.BuildTemplates(); err != nil { // first build the templates, if error then panic because this is called before server's run
		panic(err)
	}

	compiledContentType := c.ContentType + "; charset=" + c.Charset

	t := &Template{
		Engine:        e,
		IsDevelopment: c.IsDevelopment,
		Gzip:          c.Gzip,
		ContentType:   compiledContentType,
		Layout:        c.Layout,
		buffer:        utils.NewBufferPool(64),
	}

	return t

}

func (t *Template) Render(ctx context.IContext, name string, binding interface{}, layout ...string) (err error) {

	if t == nil { // No engine was given but .Render was called
		ctx.WriteHTML(403, "<b> Iris </b> <br/> Templates are disabled via config.NoEngine, check your iris' configuration please.")
		return fmt.Errorf("[IRIS TEMPLATES] Templates are disabled via config.NoEngine, check your iris' configuration please.\n")
	}

	// build templates again on each render if IsDevelopment.
	if t.IsDevelopment {
		if err = t.Engine.BuildTemplates(); err != nil {
			return
		}
	}

	// I don't like this, something feels wrong
	_layout := ""
	if len(layout) > 0 {
		_layout = layout[0]
	}
	if _layout == "" {
		_layout = t.Layout
	}

	//
	ctx.GetRequestCtx().Response.Header.Set("Content-Type", t.ContentType)

	var out io.Writer

	if t.Gzip {
		out = gzip.NewWriter(ctx.GetRequestCtx().Response.BodyWriter()) // maybe this should have its own pool but how to write after? no .WriteTo method here, I will see it someday
		defer out.(*gzip.Writer).Close()                                // here we need to close the gzip writer, after we don't do that because it is the context's body writer, carefuly kataras!
		ctx.GetRequestCtx().Response.Header.Add("Content-Encoding", "gzip")
	} else {
		out = ctx.GetRequestCtx().Response.BodyWriter()
	}

	//means Minify was true when the template engine created

	/*Try2:
	if t.minifier != nil {
		// even this doesn't works, it removes some </tags>, sadly tdewolff/minify/html is buggy,I remove this from Iris. mw := t.minifier.Writer("text/html", out)
		if err = mw.Close(); err == nil {
			err = t.Engine.ExecuteWriter(mw, name, binding, _layout)
		}
	} else {
		err = t.Engine.ExecuteWriter(out, name, binding, _layout)
	}*/
	err = t.Engine.ExecuteWriter(out, name, binding, _layout)

	return
}

func (t *Template) RenderString(name string, binding interface{}, layout ...string) (result string, err error) {

	if t == nil { // No engine was given but .Render was called
		err = fmt.Errorf("[IRIS TEMPLATES] Templates are disabled via config.NoEngine, check your iris' configuration please.\n")
		return
	}

	// build templates again on each render if IsDevelopment.
	if t.IsDevelopment {
		if err = t.Engine.BuildTemplates(); err != nil {
			return
		}
	}

	// I don't like this, something feels wrong
	_layout := ""
	if len(layout) > 0 {
		_layout = layout[0]
	}
	if _layout == "" {
		_layout = t.Layout
	}

	out := t.buffer.Get()
	// if we have problems later consider that -> out.Reset()
	defer t.buffer.Put(out)
	err = t.Engine.ExecuteWriter(out, name, binding, _layout)
	if err == nil {
		result = out.String()
	}
	return
}
