package httpserver

import (
	"context"
	"encoding/json"
	"github.com/gazercloud/gazernode/logger"
	"github.com/gazercloud/gazernode/system/system"
	"github.com/gorilla/mux"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

type HttpServer struct {
	srv      *http.Server
	r        *mux.Router
	system   *system.System
	rootPath string
}

func CurrentExePath() string {
	dir, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	return dir
}

func NewHttpServer(sys *system.System) *HttpServer {
	var c HttpServer
	c.rootPath = CurrentExePath() + "/www"
	c.system = sys
	return &c
}

func (c *HttpServer) Start() {
	logger.Println("HttpServer start")
	c.r = mux.NewRouter()

	// API
	c.r.HandleFunc("/api/request", c.processApiRequest)

	// Static files
	c.r.NotFoundHandler = http.HandlerFunc(c.processFile)

	c.srv = &http.Server{Addr: ":8084"} // 127.0.0.1
	c.srv.Handler = c.r
	go c.thListen()
}

func (c *HttpServer) thListen() {
	logger.Println("HttpServer thListen begin")
	err := c.srv.ListenAndServe()
	if err != nil {
		logger.Println("HttpServer thListen error: ", err)
	}
	logger.Println("HttpServer thListen end")
}

func (c *HttpServer) Stop() error {
	var err error

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	if err = c.srv.Shutdown(ctx); err != nil {
		logger.Println(err)
	}
	return err
}

/*func (c *HttpServer) Request(requestText string) (string, error) {
	var err error
	var respBytes []byte

	type Request struct {
		Function string `json:"func"`
		Path     string `json:"path"`
		Layer    string `json:"layer"`
	}
	var req Request
	err = json.Unmarshal([]byte(requestText), &req)
	if err != nil {
		return "", err
	}

	type Response struct {
		Value    string `json:"v"`
		DateTime string `json:"t"`
		Error    string `json:"e"`
	}

	var resp Response
	resp.Value = "123"
	resp.DateTime = time.Now().Format("2006-01-02 15-04-05.999")
	resp.Error = "ok"

	respBytes, err = json.MarshalIndent(resp, "", " ")
	if err != nil {
		return "", err
	}

	return string(respBytes), nil
}*/

func (c *HttpServer) processApiRequest(w http.ResponseWriter, r *http.Request) {
	var err error
	var responseText []byte
	requestJson := r.FormValue("rj")
	function := r.FormValue("func")

	if r.Method == "POST" {
		//body, err := ioutil.ReadAll(r.Body)
		//logger.Println("ParseMultipartForm: ", err, body)
		//requestJson = string(body)
	}

	if len(function) > 0 {
		responseText, err = c.requestJson(function, []byte(requestJson))
		//logger.Println("call ", function, "req:", requestJson)
	}

	if err != nil {
		type ErrorObject struct {
			Error string `json:"error"`
		}

		var errObj ErrorObject
		errObj.Error = err.Error()

		w.WriteHeader(500)
		b, _ := json.Marshal(errObj)
		_, _ = w.Write(b)
		return
	}
	_, _ = w.Write([]byte(responseText))
}

func (c *HttpServer) processDiagnostics(w http.ResponseWriter, r *http.Request) {
	var responseText string

	//responseText = c.system.Diag()

	_, _ = w.Write([]byte(responseText))
}

func SplitRequest(path string) []string {
	return strings.FieldsFunc(path, func(r rune) bool {
		return r == '/'
	})
}

func (c *HttpServer) contentTypeByExt(ext string) string {
	var builtinTypesLower = map[string]string{
		".css":  "text/css; charset=utf-8",
		".gif":  "image/gif",
		".htm":  "text/html; charset=utf-8",
		".html": "text/html; charset=utf-8",
		".jpeg": "image/jpeg",
		".jpg":  "image/jpeg",
		".js":   "text/javascript; charset=utf-8",
		".mjs":  "text/javascript; charset=utf-8",
		".pdf":  "application/pdf",
		".png":  "image/png",
		".svg":  "image/svg+xml",
		".wasm": "application/wasm",
		".webp": "image/webp",
		".xml":  "text/xml; charset=utf-8",
	}

	logger.Println("Ext: ", ext)

	if ct, ok := builtinTypesLower[ext]; ok {
		return ct
	}
	return "text/plain"
}

func (c *HttpServer) processFile(w http.ResponseWriter, r *http.Request) {
	var err error
	var fileContent []byte
	var writtenBytes int

	realIP := getRealAddr(r)
	logger.Println("Real IP: ", realIP)
	logger.Println("HttpServer processFile: ", r.URL.String())
	URL := *r.URL

	if URL.Path == "/" || URL.Path == "" {
		URL.Path = "/index.html"
	}

	//fileContent, err = httpdata.Asset("www" + URL.Path)

	//fileContent, err = ioutil.ReadFile(url)
	fileContent = []byte("GAZER")

	if strings.Contains(URL.Path, "light") {
		fileContent, _ = ioutil.ReadFile("d:\\2\\index.html")
	}

	if err == nil {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		writtenBytes, err = w.Write(fileContent)
		if err != nil {
			logger.Println("HttpServer sendError w.Write error:", err)
		}
		if writtenBytes != len(fileContent) {
			logger.Println("HttpServer sendError w.Write data size mismatch. (", writtenBytes, " / ", len(fileContent))
		}
	} else {
		logger.Println("HttpServer processFile error: ", err)
		w.WriteHeader(404)
	}
}

func getRealAddr(r *http.Request) string {

	remoteIP := ""
	// the default is the originating ip. but we try to find better options because this is almost
	// never the right IP
	if parts := strings.Split(r.RemoteAddr, ":"); len(parts) == 2 {
		remoteIP = parts[0]
	}
	// If we have a forwarded-for header, take the address from there
	if xff := strings.Trim(r.Header.Get("X-Forwarded-For"), ","); len(xff) > 0 {
		addrs := strings.Split(xff, ",")
		lastFwd := addrs[len(addrs)-1]
		if ip := net.ParseIP(lastFwd); ip != nil {
			remoteIP = ip.String()
		}
		// parse X-Real-Ip header
	} else if xri := r.Header.Get("X-Real-Ip"); len(xri) > 0 {
		if ip := net.ParseIP(xri); ip != nil {
			remoteIP = ip.String()
		}
	}

	return remoteIP

}

func (c *HttpServer) sendError(w http.ResponseWriter, errorToSend error) {
	var err error
	var writtenBytes int
	var b []byte
	w.WriteHeader(500)
	b, err = json.Marshal(errorToSend.Error())
	if err != nil {
		logger.Println("HttpServer sendError json.Marshal error:", err)
	}
	writtenBytes, err = w.Write(b)
	if err != nil {
		logger.Println("HttpServer sendError w.Write error:", err)
	}
	if writtenBytes != len(b) {
		logger.Println("HttpServer sendError w.Write data size mismatch. (", writtenBytes, " / ", len(b))
	}
}

func (c *HttpServer) fullpath(url string, demo bool) (string, error) {
	result := ""

	if demo {
		result = c.rootPath + "/demo/" + url
	} else {
		result = c.rootPath + "/" + url
	}

	fi, err := os.Stat(result)
	if err == nil {
		if fi.IsDir() {
			result += "/index.html"
		}
	}

	return result, err
}

func (c *HttpServer) redirect(w http.ResponseWriter, r *http.Request, url string) {
	w.Header().Set("Cache-Control", "no-cache, private, max-age=0")
	w.Header().Set("Expires", time.Unix(0, 0).Format(http.TimeFormat))
	w.Header().Set("Pragma", "no-cache")
	w.Header().Set("X-Accel-Expires", "0")
	http.Redirect(w, r, url, 307)
}

func (c *HttpServer) processTemplate(tmp []byte, demo bool) []byte {
	tmpString := string(tmp)
	re := regexp.MustCompile(`\{#.*?#\}`)
	reResults := re.FindAllString(tmpString, 100)
	for _, reString := range reResults {
		filePath := strings.ReplaceAll(reString, "{#", "")
		filePath = strings.ReplaceAll(filePath, "#}", "")
		url, err := c.fullpath(filePath, demo)
		if err != nil {
			logger.Println("processTemplate - c.fullpath(filePath) - ", err)
			continue
		}
		fileContent, err := ioutil.ReadFile(url)
		if err != nil {
			fileContent = []byte("-")
		} else {
			fileContent = c.processTemplate(fileContent, demo)
		}
		tmpString = strings.ReplaceAll(tmpString, reString, string(fileContent))
	}
	return []byte(tmpString)
}
