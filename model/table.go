package model

// DynamicModel 定义一个通用的结构体模板
type DynamicModel struct {
	TableName  string  `json:"table"`
	FieldsInfo []Field `json:"fields"`
}

type Field struct {
	Name        string      `json:"name"`
	Type        string      `json:"type"` // 对应数据库类型，例如："varchar(255)", "tinyint(1)"等
	IsNullable  bool        `json:"is_nullable"`
	Constraints Constraints `json:"constraints"` // 可选，包含长度、默认值等额外约束信息
}

type Constraints struct {
	Length  uint   `json:"length,omitempty"`
	Default string `json:"default,omitempty"`
}
