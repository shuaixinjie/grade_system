package main

import (
	"io"
	"log"
	"net/http"
	"os"
)

func UploadMoreInfo(w http.ResponseWriter, r *http.Request) (fs []*os.File) {
	if r.Method == "POST" {
		//设置内存大小
		r.ParseMultipartForm(32 << 20)
		//获取上传的文件组
		files := r.MultipartForm.File["file"]
		len := len(files)
		for i := 0; i < len; i++ {
			//打开上传文件
			file, err := files[i].Open()
			if err != nil {
				log.Fatal(err)
			}
			f, err := os.Create("temp")
			if err != nil {
				log.Fatal(err)
			}
			io.Copy(f, file)
			fs = append(fs, f)
			os.Remove("temp")
		}
		return
	}
	return
}
