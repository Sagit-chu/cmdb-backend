package models

import (
	"database/sql"
	"encoding/json"
)

// NullString 扩展了 sql.NullString 以支持 JSON 和 SQL 解析
type NullString struct {
	sql.NullString
}

func (ns *NullString) UnmarshalJSON(data []byte) error {
	var s *string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}

	if s != nil {
		ns.Valid = true
		ns.String = *s
	} else {
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
