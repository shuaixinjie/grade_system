package memory

import (
	"math"
	"sort"
	"sync"
)

// lock 读写锁增加表访问安全性和保证一定的性能
var (
	lock sync.RWMutex
)

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

// AddOne 这里写锁
func (g *SystemTable) AddOne(grade *Grade) {
	lock.Lock()
	defer lock.Lock()
	g.Grades = append(g.Grades, grade)
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

// QueryAverage 查询某班级及科目的平均分数
func (g *SystemTable) QueryAverage(class, subject string) (float64, int) {
	lock.RLock()
	defer lock.RUnlock()
	// 总分数，及人次
	sumScore, num := 0, 0
	for _, g := range g.Grades {
		if g.Class == class && g.Subject == subject {
			sumScore += g.Score
			num++
		}
	}
	return float64(sumScore / num), num
}

// QueryMaxOrMin 查询某班级及科目的最大或者最小分数
func (g *SystemTable) QueryMaxOrMin(class, subject string) (int, int) {
	lock.RLock()
	defer lock.RUnlock()
	// 最大分数， 最小分数
	maxScore, minScore := 0, 0
	for _, g := range g.Grades {
		if g.Class == class && g.Subject == subject {
			if g.Score > maxScore {
				maxScore = g.Score
			}
			if minScore < g.Score {
				minScore = g.Score
			}

		}
	}
	return maxScore, minScore
}

// QueryMedian 查询某班级及科目的中位数
func (g *SystemTable) QueryMedian(class, subject string) float64 {
	lock.RLock()
	defer lock.RUnlock()
	s := make([]int, 0, 10)
	for _, g := range g.Grades {
		if g.Class == class && g.Subject == subject {
			s = append(s, g.Score)
		}
	}
	return getMedian(s)
}

// getMedian 写个冒泡算中位数
func getMedian(s []int) float64 {
	length := len(s)
	for i := 0; i < length/2+1; i++ {
		for j := length - 1; j > i; j-- {
			if s[j+1] < s[j] {
				s[j], s[i] = s[i], s[j]
			}
		}
	}
	if length%2 == 1 {
		return float64((s[length/2] + s[length/2+1]) / 2)
	} else {
		return float64(s[length/2+1])
	}
}

// QueryZ 查询某班级及科目的标准差 s=sqrt(((x1-x)^2 +(x2-x)^2 +......(xn-x)^2)/（n-1))
func (g *SystemTable) QueryZ(class, subject string) float64 {
	lock.RLock()
	defer lock.RUnlock()
	average, n := g.QueryAverage(class, subject)
	var sum float64
	for _, g := range g.Grades {
		if g.Class == class && g.Subject == subject {
			sum += math.Abs((average - float64(g.Score)) * (average - float64(g.Score)))
		}
	}
	return math.Sqrt(sum / float64(n-1))
}
