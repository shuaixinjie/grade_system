package csv

import (
	"encoding/csv"
	"grade_system/memory"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
)

// AddGrade 处理csv文件解析数据成实体
func AddGrade(file *os.File) (gs []*memory.Grade) {
	reader := csv.NewReader(file)
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			log.Fatal(err)
		}
		fields := strings.Fields(record[0])
		score, err := strconv.Atoi(fields[2])
		if err != nil {
			log.Fatal(err)
		}
		g := &memory.Grade{
			Class:   fields[0],
			Name:    fields[1],
			Score:   score,
			Subject: fields[3],
		}
		gs = append(gs, g)
	}
	return
}
