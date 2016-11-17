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
	log.Printf(r.URL.Path)
	self.serveFile(w, r, true)
}

func (self *RawfishHandler) GetLocalFilename(name string) string {
	return strings.TrimPrefix(name, self.prefix)
}

// localRedirect gives a Moved Permanently resoponse
// It does not convert relative paths to absolute paths like Redirect does.
func localRedirect(w http.ResponseWriter, r *http.Request, newPath string) {
	if q := r.URL.RawQuery; q != "" {
		newPath += "?" + q
	}

	w.Header().Set("Location", newPath)
	w.WriteHeader(http.StatusMovedPermanently)
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

func (self *RawfishHandler) serveContent(w http.ResponseWriter, r *http.Request, name string, size int64, content io.ReadSeeker) {
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

func (self *RawfishHandler) serveFile(w http.ResponseWriter, r *http.Request, redirect bool) {
	name := r.URL.Path
	local_filename := self.GetLocalFilename(name)
	local_filename = path.Clean(local_filename)

	f, err := self.filesystem.Open(local_filename)
	if err != nil {
		if self.config.Force200Ok {
			if self.config.Rate > 0 {
				w = NewShapeIOResponseWriter(w, self.config.Rate)
			}

			reader := bufio.NewReader(limitless.NewLimitlessReader(""))
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			io.CopyN(w, reader, int64(self.config.Force200OkSize))
		} else {
			msg, code := toHTTPError(err)
			http.Error(w, msg, code)
		}
		return
	}
	defer f.Close()

	stat, err := f.Stat()
	if err != nil {
		msg, code := toHTTPError(err)
		http.Error(w, msg, code)
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
		base := r.URL.Path
		if base[len(base)-1] != '/' {
			base = base + "/"
		}

		dirList(w, f, base)
		return
	}

	t := self.service_type_checker.Get(path.Dir(local_filename))

	if t.IsRawType() {
		self.serveContent(w, r, stat.Name(), stat.Size(), f)
	} else if t.IsNormalType() {
		//call golang/net/http/fs.go/ServeContent
		if self.config.Rate > 0 {
			w = NewShapeIOResponseWriter(w, self.config.Rate)
		}
		http.ServeContent(w, r, stat.Name(), stat.ModTime(), f)
	} else {
		msg, code := toHTTPError(nil)
		http.Error(w, msg, code)
	}
}
