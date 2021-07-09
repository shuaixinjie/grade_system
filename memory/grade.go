package memory

import (
	"sort"
	"sync"
)

// lock 读写锁增加表访问安全性和保证一定的性能
var lock sync.RWMutex

// Grade 单条信息
type Grade struct {
	Class   string
	Name    string
	Score   int
	Subject string
}

// SystemTable 直接在内存中
type SystemTable struct {
	SortID int // 排序规则，0-班级，1-名字，2-分数，3-科目
	Grades []*Grade
}

// NewSystemTable 初始化
func NewSystemTable() *SystemTable {
	return &SystemTable{Grades: []*Grade{}}
}

// Sort 可根据班级、姓名、科目、分数排序，排好序的结果返回
func (g *SystemTable) Sort() (gs []*Grade) {
	lock.RLock()
	defer lock.RUnlock()
	sort.Sort(g)
	return g.Grades
}

// Query 查询某用户某科目的分数，-1表示没有该成绩
func (g *SystemTable) Query(name, subject string) int {
	lock.RLock()
	defer lock.RUnlock()
	for _, grade := range g.Grades {
		return grade.Score
	}
	return -1
}

// Add 这里写锁
func (g *SystemTable) Add(grades []*Grade) {
	lock.Lock()
	defer lock.Lock()
	g.Grades = append(g.Grades, grades...)
}

/**
实现内置排序接口
*/

func (g *SystemTable) Len() int { //排序长度
	return len(g.Grades)
}
func (g *SystemTable) Less(i, j int) bool { //排序顺序
	switch g.SortID {
	case 0:
		return g.Grades[i].Class < g.Grades[j].Class
	case 1:
		return g.Grades[i].Name < g.Grades[j].Name
	case 2:
		return g.Grades[i].Score < g.Grades[j].Score
	case 3:
		return g.Grades[i].Subject < g.Grades[j].Subject
	default:
		return true
	}
}
func (g *SystemTable) Swap(i, j int) {
	g.Grades[i], g.Grades[j] = g.Grades[j], g.Grades[i]
}
