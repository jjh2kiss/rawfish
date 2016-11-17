package rawfishnet

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"path"
	"strings"

	"github.com/jjh2kiss/rawfish/config"

	"github.com/jjh2kiss/rawfish/service"
	"github.com/jjh2kiss/rawfish/strings/limitless"
)

var htmlReplacer = strings.NewReplacer(
	"&", "&amp;",
	"<", "&lt;",
	">", "&gt;",
	// "&#34;" is shorter than "&quot;".
	`"`, "&#34;",
	// "&#39;" is shorter than "&apos;" and apos was not in HTML until HTML5.
	"'", "&#39;",
)

type RawfishHandler struct {
	filesystem           http.FileSystem
	prefix               string
	config               *config.Config
	service_type_checker *service.ServiceType
}

func NewRawfishHandler(config *config.Config, prefix string) *RawfishHandler {
	service_type_checker := service.NewServiceType(
		config.Root,
		service.Type(service.SERVICETYPE_NORMAL),
	)

	if service_type_checker == nil {
		return nil
	}

	return &RawfishHandler{
		filesystem:           http.Dir(config.Root),
		prefix:               prefix,
		config:               config,
		service_type_checker: service_type_checker,
	}
}

func (self *RawfishHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Printf("REQ %s %s", r.Method, r.URL.Path)
	self.serve(w, r, true)
}

func (self *RawfishHandler) GetLocalFilename(name string) string {
	return strings.TrimPrefix(name, self.prefix)
}

func (self *RawfishHandler) GetCleanLocalFilename(name string) string {
	return path.Clean(self.GetLocalFilename(name))
}

// localRedirect gives a Moved Permanently resoponse
// It does not convert relative paths to absolute paths like Redirect does.
func localRedirect(w http.ResponseWriter, r *http.Request, newPath string) {
	if q := r.URL.RawQuery; q != "" {
		newPath += "?" + q
	}

	w.Header().Set("Location", newPath)
	w.WriteHeader(http.StatusMovedPermanently)

	log.Printf("RES %d REDIRECT(%s->%s)",
		http.StatusMovedPermanently,
		r.URL.Path,
		newPath,
	)
}

func toHTTPError(err error) (msg string, httpStatus int) {
	if os.IsNotExist(err) {
		return "404 page not found", http.StatusNotFound
	}

	if os.IsPermission(err) {
		return "403 Forbidden", http.StatusForbidden
	}

	//Default:
	return "500 Internal Server Error", http.StatusInternalServerError
}

func dirList(w http.ResponseWriter, f http.File, base string) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprintf(w, "<pre>\n")
	for {
		dirs, err := f.Readdir(100)
		if err != nil || len(dirs) == 0 {
			break
		}
		for _, d := range dirs {
			name := d.Name()
			if d.IsDir() {
				name += "/"
			}
			// name may contain '?' or '#', which must be escaped to remain
			// part of the URL path, and not indicate the start of a query
			// string or fragment.
			url := url.URL{Path: base + name}
			fmt.Fprintf(w, "<a href=\"%s\">%s</a>\n", url.String(), htmlReplacer.Replace(name))
		}
	}
	fmt.Fprintf(w, "</pre>\n")
}

func force200Ok(w http.ResponseWriter, size int) {
	reader := bufio.NewReader(limitless.NewLimitlessReader(""))
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	io.CopyN(w, reader, int64(size))
}

func (self *RawfishHandler) Error(w http.ResponseWriter, err error) {
	if self.config.Force200Ok {
		log.Printf("RES 200 FORCE(len=%d)", self.config.Force200OkSize)
		reader := bufio.NewReader(limitless.NewLimitlessReader(""))
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		io.CopyN(w, reader, int64(self.config.Force200OkSize))
	} else {
		msg, code := toHTTPError(err)
		http.Error(w, msg, code)
		log.Printf("RES %s", msg)
	}
}

func (self *RawfishHandler) serveContent(w http.ResponseWriter, r *http.Request, size int64, content io.ReadSeeker) {
	hj, ok := w.(http.Hijacker)
	if !ok {
		http.Error(w, "webserver doesn't support hijacking", http.StatusInternalServerError)
		return
	}

	conn, rw, err := hj.Hijack()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	//Don't forget to close the connectin:
	defer conn.Close()
	if self.config.Rate > 0 {
		//conn is socket fd, it use unbuffered io
		CopyWithShapeIO(conn, content, int(size), self.config.Rate)
	} else {
		io.CopyN(rw, content, size)
	}
}

func (self *RawfishHandler) serveFile(w http.ResponseWriter, r *http.Request, f http.File, stat os.FileInfo) {
	filepath := self.GetCleanLocalFilename(r.URL.Path)
	t := self.service_type_checker.Get(path.Dir(filepath))

	if t.IsRawType() {
		//use original http.ResponseWriter
		log.Printf("RES 200 RAW %s", filepath)
		self.serveContent(w, r, stat.Size(), f)
	} else if t.IsNormalType() {
		//call golang/net/http/fs.go/ServeContent
		log.Printf("RES 200 NORMAL %s", filepath)
		http.ServeContent(w, r, stat.Name(), stat.ModTime(), f)
	} else {
		self.Error(w, nil)
	}
}

func (self *RawfishHandler) serveDir(w http.ResponseWriter, r *http.Request, f http.File) {
	base := r.URL.Path
	if base[len(base)-1] != '/' {
		base = base + "/"
	}
	log.Printf("RES 200 DIR Listing %s", base)
	dirList(w, f, base)
}

func (self *RawfishHandler) serve(w http.ResponseWriter, r *http.Request, redirect bool) {
	if self.config.Rate > 0 {
		w = NewShapeIOResponseWriter(w, self.config.Rate)
	}

	filepath := self.GetCleanLocalFilename(r.URL.Path)
	f, err := self.filesystem.Open(filepath)
	if err != nil {
		self.Error(w, err)
		return
	}

	defer f.Close()

	stat, err := f.Stat()
	if err != nil {
		self.Error(w, err)
		return
	}

	if redirect {
		url := r.URL.Path
		if stat.IsDir() {
			if url[len(url)-1] != '/' {
				localRedirect(w, r, path.Base(url)+"/")
				return
			}
		} else {
			if url[len(url)-1] == '/' {
				localRedirect(w, r, "../"+path.Base(url))
				return
			}
		}
	}

	if stat.IsDir() {
		self.serveDir(w, r, f)
	} else {
		self.serveFile(w, r, f, stat)
	}
}
