package model

type Test struct {
	Id           uint   `gorm:"primarykey;not null" json:"id" form:"id"`
	ExampleField bool   `gorm:"not null" json:"example_field" form:"example_field"`
	Testcol      string `gorm:"varchar(45);not null" json:"testcol" form:"testcol"`
}

func (Test) TableName() string {
	return "test"
}
