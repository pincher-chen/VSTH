package iris

import (
	"path"
	"reflect"
	"strconv"
	"strings"

	"github.com/kataras/iris/utils"
	"github.com/valyala/fasthttp"
)

type (
	// IParty is the interface which implements the whole Party of routes
	IParty interface {
		Handle(string, string, ...Handler)
		HandleFunc(string, string, ...HandlerFunc)
		HandleAnnotated(Handler) error
		Get(string, ...HandlerFunc)
		Post(string, ...HandlerFunc)
		Put(string, ...HandlerFunc)
		Delete(string, ...HandlerFunc)
		Connect(string, ...HandlerFunc)
		Head(string, ...HandlerFunc)
		Options(string, ...HandlerFunc)
		Patch(string, ...HandlerFunc)
		Trace(string, ...HandlerFunc)
		Any(string, ...HandlerFunc)
		Use(...Handler)
		UseFunc(...HandlerFunc)
		// Static serves a directory
		// accepts three parameters
		// first parameter is the request url path (string)
		// second parameter is the system directory (string)
		// third parameter is the level (int) of stripSlashes
		// * stripSlashes = 0, original path: "/foo/bar", result: "/foo/bar"
		// * stripSlashes = 1, original path: "/foo/bar", result: "/bar"
		// * stripSlashes = 2, original path: "/foo/bar", result: ""
		Static(string, string, int)
		StaticFS(string, string, int)
		Party(string, ...HandlerFunc) IParty // Each party can have a party too
		IsRoot() bool
	}

	// GardenParty  is the struct which makes all the job for registering routes and middlewares
	GardenParty struct {
		relativePath string
		station      *Iris // this station is where the party is happening, this station's Garden is the same for all Parties per Station & Router instance
		middleware   Middleware
		root         bool
	}
)

var _ IParty = &GardenParty{}

// IsRoot returns true if this is the root party ("/")
func (p *GardenParty) IsRoot() bool {
	return p.root
}

// Handle registers a route to the server's router
// if empty method is passed then registers handler(s) for all methods, same as .Any
func (p *GardenParty) Handle(method string, registedPath string, handlers ...Handler) {
	if method == "" { // then use like it was .Any
		for _, k := range AllMethods {
			p.Handle(k, registedPath, handlers...)
		}
		return
	}
	path := fixPath(p.relativePath + registedPath) // keep the last "/" as default ex: "/xyz/"
	if p.station.config.PathCorrection {
		// if we have path correction remove it with absPath
		path = fixPath(absPath(p.relativePath, registedPath)) // "/xyz"
	}
	middleware := JoinMiddleware(p.middleware, handlers)
	route := NewRoute(method, path, middleware)
	p.station.plugins.DoPreHandle(route)
	p.station.addRoute(route)
	p.station.plugins.DoPostHandle(route)
}

// HandleFunc registers and returns a route with a method string, path string and a handler
// registedPath is the relative url path
// handler is the iris.Handler which you can pass anything you want via iris.ToHandlerFunc(func(res,req){})... or just use func(c *iris.Context)
func (p *GardenParty) HandleFunc(method string, registedPath string, handlersFn ...HandlerFunc) {
	p.Handle(method, registedPath, ConvertToHandlers(handlersFn)...)
}

// HandleAnnotated registers a route handler using a Struct implements iris.Handler (as anonymous property)
// which it's metadata has the form of
// `method:"path"` and returns the route and an error if any occurs
// handler is passed by func(urstruct MyStruct) Serve(ctx *Context) {}
func (p *GardenParty) HandleAnnotated(irisHandler Handler) error {
	var method string
	var path string
	var errMessage = ""
	val := reflect.ValueOf(irisHandler).Elem()

	for i := 0; i < val.NumField(); i++ {
		typeField := val.Type().Field(i)

		if typeField.Anonymous && typeField.Name == "Handler" {
			tags := strings.Split(strings.TrimSpace(string(typeField.Tag)), " ")
			firstTag := tags[0]

			idx := strings.Index(string(firstTag), ":")

			tagName := strings.ToUpper(string(firstTag[:idx]))
			tagValue, unqerr := strconv.Unquote(string(firstTag[idx+1:]))

			if unqerr != nil {
				errMessage = errMessage + "\non getting path: " + unqerr.Error()
				continue
			}

			path = tagValue
			avalaibleMethodsStr := strings.Join(AllMethods[0:], ",")

			if !strings.Contains(avalaibleMethodsStr, tagName) {
				//wrong method passed
				errMessage = errMessage + "\nWrong method passed to the anonymous property iris.Handler -> " + tagName
				continue
			}

			method = tagName

		} else {
			errMessage = "\nStruct passed but it doesn't have an anonymous property of type iris.Hanndler, please refer to docs\n"
		}

	}

	if errMessage == "" {
		p.Handle(method, path, irisHandler)
	}

	var err error
	if errMessage != "" {
		err = ErrHandleAnnotated.Format(errMessage)
	}

	return err
}

// Get registers a route for the Get http method
func (p *GardenParty) Get(path string, handlersFn ...HandlerFunc) {
	p.HandleFunc(MethodGet, path, handlersFn...)
}

// Post registers a route for the Post http method
func (p *GardenParty) Post(path string, handlersFn ...HandlerFunc) {
	p.HandleFunc(MethodPost, path, handlersFn...)
}

// Put registers a route for the Put http method
func (p *GardenParty) Put(path string, handlersFn ...HandlerFunc) {
	p.HandleFunc(MethodPut, path, handlersFn...)
}

// Delete registers a route for the Delete http method
func (p *GardenParty) Delete(path string, handlersFn ...HandlerFunc) {
	p.HandleFunc(MethodDelete, path, handlersFn...)
}

// Connect registers a route for the Connect http method
func (p *GardenParty) Connect(path string, handlersFn ...HandlerFunc) {
	p.HandleFunc(MethodConnect, path, handlersFn...)
}

// Head registers a route for the Head http method
func (p *GardenParty) Head(path string, handlersFn ...HandlerFunc) {
	p.HandleFunc(MethodHead, path, handlersFn...)
}

// Options registers a route for the Options http method
func (p *GardenParty) Options(path string, handlersFn ...HandlerFunc) {
	p.HandleFunc(MethodOptions, path, handlersFn...)
}

// Patch registers a route for the Patch http method
func (p *GardenParty) Patch(path string, handlersFn ...HandlerFunc) {
	p.HandleFunc(MethodPatch, path, handlersFn...)
}

// Trace registers a route for the Trace http method
func (p *GardenParty) Trace(path string, handlersFn ...HandlerFunc) {
	p.HandleFunc(MethodTrace, path, handlersFn...)
}

// Any registers a route for ALL of the http methods (Get,Post,Put,Head,Patch,Options,Connect,Delete)
func (p *GardenParty) Any(registedPath string, handlersFn ...HandlerFunc) {
	for _, k := range AllMethods {
		p.HandleFunc(k, registedPath, handlersFn...)
	}

}

// Use registers a Handler middleware
func (p *GardenParty) Use(handlers ...Handler) {
	p.middleware = append(p.middleware, handlers...)
}

// UseFunc registers a HandlerFunc middleware
func (p *GardenParty) UseFunc(handlersFn ...HandlerFunc) {
	p.Use(ConvertToHandlers(handlersFn)...)
}

// StaticHandlerFunc returns a HandlerFunc to serve static system directory
// Accepts 5 parameters
//
// first is the systemPath (string)
// Path to the root directory to serve files from.
//
// second is the  stripSlashes (int) level
// * stripSlashes = 0, original path: "/foo/bar", result: "/foo/bar"
// * stripSlashes = 1, original path: "/foo/bar", result: "/bar"
// * stripSlashes = 2, original path: "/foo/bar", result: ""
//
// third is the compress (bool)
// Transparently compresses responses if set to true.
//
// The server tries minimizing CPU usage by caching compressed files.
// It adds FSCompressedFileSuffix suffix to the original file name and
// tries saving the resulting compressed file under the new file name.
// So it is advisable to give the server write access to Root
// and to all inner folders in order to minimze CPU usage when serving
// compressed responses.
//
// fourth is the generateIndexPages (bool)
// Index pages for directories without files matching IndexNames
// are automatically generated if set.
//
// Directory index generation may be quite slow for directories
// with many files (more than 1K), so it is discouraged enabling
// index pages' generation for such directories.
//
// fifth is the indexNames ([]string)
// List of index file names to try opening during directory access.
//
// For example:
//
//     * index.html
//     * index.htm
//     * my-super-index.xml
//
func (p *GardenParty) StaticHandlerFunc(systemPath string, stripSlashes int, compress bool, generateIndexPages bool, indexNames []string) HandlerFunc {
	if indexNames == nil {
		indexNames = []string{}
	}
	fs := &fasthttp.FS{
		// Path to directory to serve.
		Root:       systemPath,
		IndexNames: indexNames,
		// Generate index pages if client requests directory contents.
		GenerateIndexPages: generateIndexPages,

		// Enable transparent compression to save network traffic.
		Compress: compress,
	}

	if stripSlashes > 0 {
		fs.PathRewrite = fasthttp.NewPathSlashesStripper(stripSlashes)
	}

	// Create request handler for serving static files.
	h := fs.NewRequestHandler()
	return func(ctx *Context) {
		h(ctx.RequestCtx)
		errCode := ctx.RequestCtx.Response.StatusCode()

		if errHandler := ctx.station.router.GetByCode(errCode); errHandler != nil {
			ctx.RequestCtx.Response.ResetBody()
			ctx.EmitError(errCode)
		}
		if ctx.pos < uint8(len(ctx.middleware))-1 {
			ctx.Next() // for any case
		}

	}
}

// Static registers a route which serves a system directory
// this doesn't generates an index page which list all files
// no compression is used also, for these features look at StaticFS func
// accepts three parameters
// first parameter is the request url path (string)
// second parameter is the system directory (string)
// third parameter is the level (int) of stripSlashes
// * stripSlashes = 0, original path: "/foo/bar", result: "/foo/bar"
// * stripSlashes = 1, original path: "/foo/bar", result: "/bar"
// * stripSlashes = 2, original path: "/foo/bar", result: ""
func (p *GardenParty) Static(relative string, systemPath string, stripSlashes int) {
	if relative[len(relative)-1] != SlashByte { // if / then /*filepath, if /something then /something/*filepath
		relative += "/"
	}

	h := p.StaticHandlerFunc(systemPath, stripSlashes, false, false, nil)

	p.Get(relative+"*filepath", h)
	p.Head(relative+"*filepath", h)
}

// StaticFS registers a route which serves a system directory
// generates an index page which list all files
// uses compression which file cache, if you use this method it will generate compressed files also
// think this function as small fileserver with http
// accepts three parameters
// first parameter is the request url path (string)
// second parameter is the system directory (string)
// third parameter is the level (int) of stripSlashes
// * stripSlashes = 0, original path: "/foo/bar", result: "/foo/bar"
// * stripSlashes = 1, original path: "/foo/bar", result: "/bar"
// * stripSlashes = 2, original path: "/foo/bar", result: ""
func (p *GardenParty) StaticFS(relative string, systemPath string, stripSlashes int) {
	if relative[len(relative)-1] != SlashByte {
		relative += "/"
	}

	h := p.StaticHandlerFunc(systemPath, stripSlashes, true, true, nil)
	p.Get(relative+"*filepath", h)
	p.Head(relative+"*filepath", h)
}

// StaticWeb same as Static but if index.html exists and request uri is '/' then display the index.html's contents
// accepts three parameters
// first parameter is the request url path (string)
// second parameter is the system directory (string)
// third parameter is the level (int) of stripSlashes
// * stripSlashes = 0, original path: "/foo/bar", result: "/foo/bar"
// * stripSlashes = 1, original path: "/foo/bar", result: "/bar"
// * stripSlashes = 2, original path: "/foo/bar", result: ""
func (p *GardenParty) StaticWeb(relative string, systemPath string, stripSlashes int) {
	if relative[len(relative)-1] != SlashByte { // if / then /*filepath, if /something then /something/*filepath
		relative += "/"
	}

	hasIndex := utils.Exists(systemPath + utils.PathSeparator + "index.html")
	serveHandler := p.StaticHandlerFunc(systemPath, 1, false, !hasIndex, nil) // if not index.html exists then generate index.html which shows the list of files
	indexHandler := func(ctx *Context) {
		if len(ctx.Param("filepath")) < 2 && hasIndex {
			ctx.Request.SetRequestURI("index.html")
		}
		ctx.Next()

	}
	p.Get(relative+"*filepath", indexHandler, serveHandler)
	p.Head(relative+"*filepath", indexHandler, serveHandler)
}

// Party is just a group joiner of routes which have the same prefix and share same middleware(s) also.
// Party can also be named as 'Join' or 'Node' or 'Group' , Party chosen because it has more fun
func (p *GardenParty) Party(path string, handlersFn ...HandlerFunc) IParty {
	middleware := ConvertToHandlers(handlersFn)
	if path[0] != SlashByte && strings.Contains(path, ".") {
		//it's domain so no handlers share (even the global ) or path, nothing.
	} else {
		// set path to parent+child
		path = absPath(p.relativePath, path)
		// append the parent's +child's handlers
		middleware = JoinMiddleware(p.middleware, middleware)
	}

	return &GardenParty{relativePath: path, station: p.station, middleware: middleware}
}

func absPath(rootPath string, relativePath string) (absPath string) {

	if relativePath == "" {
		absPath = rootPath
	} else {
		absPath = path.Join(rootPath, relativePath)
	}

	return
}

// fixPath fix the double slashes, (because of root,I just do that before the .Handle no need for anything else special)
func fixPath(str string) string {

	strafter := strings.Replace(str, "//", Slash, -1)

	if strafter[0] == SlashByte && strings.Count(strafter, ".") >= 2 {
		//it's domain, remove the first slash
		strafter = strafter[1:]
	}

	return strafter
}
