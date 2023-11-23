package customType

import (
	"database/sql/driver"
	"fmt"
	"strings"
	"time"

	corecode "github.com/tzRex/freely-handle/coreCode"
)

type CustomTime struct {
	time.Time
}

var (
	TimeSchemaToDatetime = "2006-01-02 15:04:05"
	TimeSchemaToDate     = "2006-01-02"
	TimeSchemaToTime     = "15:04:05"
)

// 序列化时的回调：json.Marshal
func (ct *CustomTime) MarshalJSON() ([]byte, error) {
	if ct == nil || ct.Time.IsZero() {
		return []byte("null"), nil
	}
	formatted := fmt.Sprintf("\"%s\"", ct.Time.Format(TimeSchemaToDatetime))
	return []byte(formatted), nil
}

// 反序列化时的回调：json.Unmarshal
func (ct *CustomTime) UnmarshalJSON(b []byte) error {
	str := strings.ReplaceAll(string(b), "\"", "")

	if str == "" || str == "null" || str == "None" {
		return nil
	}

	// 使用指定格式解析时间字符串
	t, err := time.Parse(TimeSchemaToDatetime, str)
	if err != nil {
		return fmt.Errorf("date.format.incorrect example: %s", TimeSchemaToDatetime)
	}

	ct.Time = t
	return nil
}

// gorm写入数据的回调
func (ct CustomTime) Value() (driver.Value, error) {
	if ct.Time.IsZero() {
		return nil, nil
	}
	formatted := ct.Time.Format(TimeSchemaToDatetime)
	return []byte(formatted), nil
}

// gorm读取数据的回调
func (ct *CustomTime) Scan(value interface{}) error {
	b, ok := value.(time.Time)
	if !ok {
		return corecode.ErrColumnTypeFail
	}

	t, err := time.Parse(TimeSchemaToDatetime, b.Format(TimeSchemaToDatetime))
	if err != nil {
		return err
	}

	ct.Time = t
	return nil
}
