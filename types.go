package tushare

import (
	"encoding/json"
	"fmt"
)

// Response Tushare API 响应结构
type Response struct {
	Code int           `json:"code"`
	Msg  string        `json:"msg"`
	Data *ResponseData `json:"data"`
}

// ResponseData 响应数据结构
type ResponseData struct {
	Fields  []string        `json:"fields"`
	Items   [][]interface{} `json:"items"`
	HasMore bool            `json:"has_more"` // 是否还有更多数据
}

// APIError API 错误
type APIError struct {
	Code int
	Msg  string
}

func (e *APIError) Error() string {
	return fmt.Sprintf("tushare api error: code=%d, msg=%s", e.Code, e.Msg)
}

// IsSuccess 判断请求是否成功
func (r *Response) IsSuccess() bool {
	return r.Code == 0
}

// ToRecords 将响应数据转换为记录列表（map 格式）
func (r *Response) ToRecords() []map[string]interface{} {
	if r.Data == nil {
		return nil
	}

	records := make([]map[string]interface{}, 0, len(r.Data.Items))
	for _, item := range r.Data.Items {
		record := make(map[string]interface{})
		for i, field := range r.Data.Fields {
			if i < len(item) {
				record[field] = item[i]
			}
		}
		records = append(records, record)
	}
	return records
}

// ToStruct 将响应数据转换为指定类型的切片
func (r *Response) ToStruct(v interface{}) error {
	records := r.ToRecords()
	data, err := json.Marshal(records)
	if err != nil {
		return fmt.Errorf("marshal records failed: %w", err)
	}
	return json.Unmarshal(data, v)
}

// DataFrame 简单的数据帧结构（类似 pandas DataFrame）
type DataFrame struct {
	Columns []string
	Data    []map[string]interface{}
}

// NewDataFrame 从响应创建 DataFrame
func NewDataFrame(resp *Response) *DataFrame {
	if resp.Data == nil {
		return &DataFrame{
			Columns: []string{},
			Data:    []map[string]interface{}{},
		}
	}
	return &DataFrame{
		Columns: resp.Data.Fields,
		Data:    resp.ToRecords(),
	}
}

// Len 返回数据行数
func (df *DataFrame) Len() int {
	return len(df.Data)
}

// Get 获取指定行和列的值
func (df *DataFrame) Get(row int, col string) (interface{}, bool) {
	if row < 0 || row >= len(df.Data) {
		return nil, false
	}
	val, ok := df.Data[row][col]
	return val, ok
}

// GetString 获取字符串值
func (df *DataFrame) GetString(row int, col string) string {
	val, ok := df.Get(row, col)
	if !ok || val == nil {
		return ""
	}
	switch v := val.(type) {
	case string:
		return v
	default:
		return fmt.Sprintf("%v", v)
	}
}

// GetFloat64 获取 float64 值
func (df *DataFrame) GetFloat64(row int, col string) float64 {
	val, ok := df.Get(row, col)
	if !ok || val == nil {
		return 0
	}
	switch v := val.(type) {
	case float64:
		return v
	case float32:
		return float64(v)
	case int:
		return float64(v)
	case int64:
		return float64(v)
	case string:
		var f float64
		fmt.Sscanf(v, "%f", &f)
		return f
	default:
		return 0
	}
}

// GetInt 获取 int 值
func (df *DataFrame) GetInt(row int, col string) int {
	return int(df.GetFloat64(row, col))
}
