package web

import (
	"io/ioutil"
	"path/filepath"
	"net/http"
	"strings"
	"mkgame-go/controller"
	"mkgame-go/mysql"
	"log"
)

func ServerOn() {
	log.Println("server on start")
	go controller.ChatHub.Run()
	http.Handle("/", http.FileServer(http.Dir("./web/public")))
	http.HandleFunc("/ws", controller.ServeWs)
	// http.Handle("/", new(httpHandler))
	log.Fatal(http.ListenAndServe(":8080", nil))
	// http.ListenAndServe(":5000", nil)
	log.Println("server on!")
}

type httpHandler struct {
	http.Handler
}

func (h *httpHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	htmlPath := "/var/www/html"
	localPath := htmlPath + req.URL.Path
	if req.URL.Path == "/mksql" {
		w.Header().Add("Content-Type", "text/plain")
		DBhandler := mysql.Con
		var args []interface{}
		args = append(args, "2")
		success, resRow := DBhandler.SelectRow("SELECT name, email FROM users where id = 1", nil,2)
		if success == false {
			w.Write([]byte("something is wrong"))
		}
		// DBhandler.CloseDB()		
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
