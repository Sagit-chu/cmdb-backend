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

func (ns *NullString) UnmarshalJSON(data []byte) error {
	// 临时字符串变量
	var s string
	// 解析 JSON 数据
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}

	// 设置 NullString 的值
	ns.String = s
	ns.Valid = true
	fmt.Printf("Parsed string: '%s', Valid: %t\n", ns.String, ns.Valid)

	return nil
}

// MarshalJSON 自定义 JSON 序列化
func (ns NullString) MarshalJSON() ([]byte, error) {
	if ns.Valid {
		return json.Marshal(ns.String)
	}
	return json.Marshal(nil)
}
