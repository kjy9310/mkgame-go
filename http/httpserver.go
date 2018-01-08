package main

import (
	"io/ioutil"
	"path/filepath"
	"net/http"
	"strings"
	//"fmt"
)

func main() {
	
	http.Handle("/", new(httpHandler))
	http.ListenAndServe(":5000", nil)
}

type httpHandler struct {
	http.Handler
}

func (h *httpHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	htmlPath := "/var/www/html"
	localPath := htmlPath + req.URL.Path
	if req.URL.Path == "/mksql" {
		w.Header().Add("Content-Type", "text/plain")
		DBhandler := initDB()
		var args []interface{}
		args = append(args, "2")
		success, resRow := DBhandler.selectRow("SELECT name, email FROM users where id = 1", nil,2)
		if success == false {
			w.Write([]byte("something is wrong"))
		}
		DBhandler.closeDB()		
		w.Write([]byte(resRow[0]))
		return
	}
	if req.URL.Path=="/api/login"{
		responseJson := login(w, req)
		w.Write(responseJson)
		return
	}else if req.URL.Path=="/api/logout" {
		responseJson := logout(w, req)
		w.Write(responseJson)
		return
	}else if strings.Contains(req.URL.Path,"/api/"){
		success, account := loginCheck(w, req)
		if success {
			w.Write([]byte(account+" login check ok"))
		}
		return
	}
	content, err := ioutil.ReadFile(localPath)
	if err != nil {
		w.WriteHeader(404)
		w.Write([]byte(http.StatusText(404)))
		return
	}
	contentType := getContentType(localPath)
	w.Header().Add("Content-Type", contentType)
	w.Write(content)
}

func getContentType(localPath string) string {
	var contentType string
	ext := filepath.Ext(localPath)

	switch ext {
	case ".html":
		contentType = "text/html"
	case ".css":
		contentType = "text/css"
	case ".js":
		contentType = "application/javascript"
	case ".png":
		contentType = "image/png"
	case ".jpg":
		contentType = "image/jpeg"
	default:
		contentType = "text/plain"
	}
	return contentType
}
