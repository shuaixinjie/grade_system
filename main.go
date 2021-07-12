package main

import (
	"encoding/json"
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
	// 录入单个学生成绩
	http.HandleFunc("/student", addScore)
	http.HandleFunc("/query/average", queryAvg)
	http.HandleFunc("/query/maxmin", queryMaxOrMin)
	http.HandleFunc("/query/median", queryMedian)
	http.HandleFunc("/query/z", queryZ)
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


// studentReq 单个分数上传
type studentReq struct {
	Name    string `json:"name"`
	Class   string `json:"class"`
	Subject string `json:"subject"`
	Score   int    `json:"score"`
}

// addScore 录入单个学生成绩
func addScore(w http.ResponseWriter, r *http.Request) {
	l := r.ContentLength
	body := make([]byte, l)
	_, err := r.Body.Read(body)
	if err != nil {
		log.Fatal(err)
	}
	stu := new(studentReq)
	err = json.Unmarshal(body, stu)
	if err != nil {
		log.Fatal(err)
	}
	g := new(memory.Grade)
	g.Name = stu.Name
	g.Class = stu.Class
	g.Subject = stu.Subject
	g.Score = stu.Score
	table.AddOne(g)
	fmt.Fprintf(w, "add success")
}


// queryAvg 查询平均分班级 科目
func queryAvg(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	class := query.Get("class")
	sub := query.Get("subject")

	average, _ := table.QueryAverage(class, sub)
	fmt.Fprintf(w, "%s %s 的平均分数为%f", class, sub, average)
}

// queryMaxOrMin 查询最高分数和最低分数
func queryMaxOrMin(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	class := query.Get("class")
	sub := query.Get("subject")

	max, min := table.QueryMaxOrMin(class, sub)
	fmt.Fprintf(w, "%s %s 的最高分数分数为%f，最低分数为%f", class, sub, max, min)
}

// queryMedian 查询中位数班级 科目
func queryMedian(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	class := query.Get("class")
	sub := query.Get("subject")

	m := table.QueryMedian(class, sub)
	fmt.Fprintf(w, "%s %s 的中位数分数为%f", class, sub, m)
}

// queryZ 查询标准差 班级 科目
func queryZ(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	class := query.Get("class")
	sub := query.Get("subject")

	m := table.QueryZ(class, sub)
	fmt.Fprintf(w, "%s %s 的标准差分数为%f", class, sub, m)
}