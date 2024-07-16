package model

type Revenue struct {
	Month   Month   `json:"month" gorm:"type:varchar(100);not null;unique"`
	Revenue float32 `json:"revenue" gorm:"type:float;not null"`
}

type Month string
type Revenues []Revenue

// Map to define the order of the months
var monthOrder = map[Month]int{
	"Jan": 1,
	"Feb": 2,
	"Mar": 3,
	"Apr": 4,
	"May": 5,
	"Jun": 6,
	"Jul": 7,
	"Aug": 8,
	"Sep": 9,
	"Oct": 10,
	"Nov": 11,
	"Dec": 12,
}

func (r Revenues) Len() int           { return len(r) }
func (r Revenues) Swap(i, j int)      { r[i], r[j] = r[j], r[i] }
func (r Revenues) Less(i, j int) bool { return monthOrder[r[i].Month] < monthOrder[r[j].Month] }
