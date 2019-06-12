package main

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func lee(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Method", r.Method)
	fmt.Println("url", r.URL)
	fmt.Println("header", r.Header)
	fmt.Println("body", r.Body)
	w.Write([]byte("hello go"))
}

func main() {

	fileHandler := http.FileServer(http.Dir("E:/javaEE/vlog/src/main/video"))

	http.Handle("/video/", http.StripPrefix("/video/", fileHandler))

	//注册
	http.HandleFunc("/lee", lee)
	http.HandleFunc("/api/upload", uploadHandler)
	http.HandleFunc("/api/list", getFileListHandler)

	http.ListenAndServe(":8000", nil)

}

//上传文件
func uploadHandler(w http.ResponseWriter, r *http.Request) {

	//限制文件上传大小

	r.Body = http.MaxBytesReader(w, r.Body, 10*1024*1024)

	err := r.ParseMultipartForm(10 * 1024 * 1024)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	//获取上传文件
	file, fileHeader, err := r.FormFile("uploadFile")

	//检查是否是mp4格式
	ret := strings.HasSuffix(fileHeader.Filename, ".mp4")

	if ret == false {
		http.Error(w, "not mp4", http.StatusInternalServerError)
		return
	}

	//重命名文件
	md5Byte := md5.Sum([]byte(fileHeader.Filename + time.Now().String()))
	md5Str := fmt.Sprintf("%x", md5Byte)
	newFileName := md5Str + ".mp4"

	dst, err := os.Create("E:/javaEE/vlog/src/main/video/" + newFileName)

	defer dst.Close()

	if _, err := io.Copy(dst, file); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	return

}

//获取文件列表
func getFileListHandler(w http.ResponseWriter, r *http.Request) {

	files, _ := filepath.Glob("E:/javaEE/vlog/src/main/video/*")
	var ret [] string
	for _, file := range files {
		ret = append(ret, "http://"+r.Host+"/video/"+filepath.Base(file))
	}
	retJson, _ := json.Marshal(ret)
	w.Write(retJson)
	return

}
