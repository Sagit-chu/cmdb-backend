package models

import (
	"database/sql"
	"encoding/json"
	"fmt"
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
		// 确保无论字符串是否为空，只要有值（即不是 null），都应该设置为 Valid:true
		ns.Valid = true
		ns.String = *s
		fmt.Printf("Parsed valid string: %s\n", ns.String)
	} else {
		// 处理 JSON 中的 null 值
		ns.Valid = false
		ns.String = ""
		fmt.Println("Parsed null value")
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
