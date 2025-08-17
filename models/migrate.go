package models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
)

// StringArray 自定义类型，用于处理 []string 与 JSON 的转换
type StringArray []string

// Scan 实现 sql.Scanner 接口，从数据库读取 JSON 到 Go 类型
func (sa *StringArray) Scan(value interface{}) error {
	if value == nil {
		*sa = StringArray{}
		return nil
	}

	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("类型断言为 []byte 失败")
	}

	return json.Unmarshal(bytes, sa)
}

// Value 实现 driver.Valuer 接口，将 Go 类型转换为数据库存储的 JSON
func (sa StringArray) Value() (driver.Value, error) {
	if len(sa) == 0 {
		return "[]", nil // 返回空数组的 JSON 表示
	}
	return json.Marshal(sa)
}

// ChatIDHistory 表示用户的聊天 ID 历史记录
type ChatIDHistory struct {
	ID       uint        `gorm:"primarykey"`
	Username string      `gorm:"type:varchar(100);not null"`
	ChatIDs  StringArray `gorm:"type:json"` // 使用 json 类型存储字符串数组
}
