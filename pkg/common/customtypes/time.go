package customtypes

import (
	"database/sql/driver"
	"strconv"
	"strings"
	"time"
)

const (
	timeFormat  = "2006-01-02 15:04:05"
	timeFormat2 = "2006-01-02"
)

type Time time.Time

// 反序列化
func (t *Time) UnmarshalJSON(data []byte) (err error) {
	str := string(data)
	str = strings.Trim(str, "\"")
	if strings.Contains(str, "-") {
		if len(str) > 10 {
			if result, err2 := time.ParseInLocation(timeFormat, str, time.Local); err2 == nil {
				*t = Time(result)
			}
		} else {
			if result, err2 := time.ParseInLocation(timeFormat2, str, time.Local); err2 == nil {
				*t = Time(result)
			} else {
				err = err2
			}
		}
	} else {
		num, err := strconv.Atoi(str)
		if err != nil {
			return err
		}
		*t = Time(time.UnixMilli(int64(num)))
	}
	return
}

// 序列化
func (t Time) MarshalJSON() ([]byte, error) {
	return ([]byte)(strconv.FormatInt(time.Time(t).UnixMilli(), 10)), nil
}

// fmt.Println打印
func (t Time) String() string {
	b := make([]byte, 0, len(timeFormat))
	b = time.Time(t).AppendFormat(b, timeFormat)
	return string(b)
}

// gorm支持
func (t Time) Value() (driver.Value, error) {
	return time.Time(t), nil
}
