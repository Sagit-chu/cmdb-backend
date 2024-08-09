package models

import (
	"database/sql"
	"encoding/json"
)

// NullString 扩展了 sql.NullString 以支持 JSON 和 SQL 解析
type NullString struct {
	sql.NullString
}

// UnmarshalJSON 自定义 JSON 解析
func (ns *NullString) UnmarshalJSON(data []byte) error {
	var s *string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	if s != nil {
		// 如果字符串非空且有效，将其设置为 Valid
		ns.Valid = true
		ns.String = *s
	} else {
		// 如果字符串是空或 null，设置为无效
		ns.Valid = false
		ns.String = ""
	}
	return nil
}

// MarshalJSON 自定义 JSON 序列化
func (ns NullString) MarshalJSON() ([]byte, error) {
	if ns.Valid {
		return json.Marshal(ns.String)
	}
	return json.Marshal(nil)
}
