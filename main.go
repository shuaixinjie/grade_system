package main

import (
	"fmt"
	"grade_system/csv"
	"grade_system/memory"
	"log"
	"net/http"
	"strconv"
)

var (
	table = memory.NewSystemTable()
)

func main() {

	http.HandleFunc("/uploadMore", uploadMore)
	http.HandleFunc("/query/one", queryOne)
	http.HandleFunc("/query/condition", queryCondition)
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal(err)
	}
}

// 处理多文件上传
func uploadMore(w http.ResponseWriter, r *http.Request) {
	files := UploadMoreInfo(w, r)
	for i := range files {
		go func() {
			grade := csv.AddGrade(files[i])
			table.Add(grade)
		}()
	}
	fmt.Fprintf(w, "upload success")
}

// queryOne 查询某个学生的成绩
func queryOne(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	name := query.Get("name")
	subject := query.Get("subject")
	score := table.Query(name, subject)
	fmt.Fprintf(w, `{score: %d}`, score)
}

// queryCondition 根据条件排序
func queryCondition(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	id, err := strconv.Atoi(query.Get("sortID"))
	if err != nil {
		fmt.Fprintf(w, "id is illegal")
	}
	table.SortID = id
	score := table.Sort()
	fmt.Fprintf(w, "%+v", score)
}



