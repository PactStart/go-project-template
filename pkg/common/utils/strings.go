package utils

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"orderin-server/pkg/common/customtypes"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"
)

func IntToString(i int) string {
	return strconv.FormatInt(int64(i), 10)
}

func StringToInt(i string) int {
	j, _ := strconv.Atoi(i)
	return j
}
func StringToInt64(i string) int64 {
	j, _ := strconv.ParseInt(i, 10, 64)
	return j
}
func StringToInt32(i string) int32 {
	j, _ := strconv.ParseInt(i, 10, 64)
	return int32(j)
}
func Int32ToString(i int32) string {
	return strconv.FormatInt(int64(i), 10)
}

func Uint32ToString(i uint32) string {
	return strconv.FormatInt(int64(i), 10)
}

// judge a string whether in the  string list
func IsContain(target string, List []string) bool {
	for _, element := range List {

		if target == element {
			return true
		}
	}
	return false
}

func IsContainAny(target string, List []string) bool {
	for _, element := range List {

		if strings.Contains(target, element) {
			return true
		}
	}
	return false
}

func IsContainInt32(target int32, List []int32) bool {
	for _, element := range List {
		if target == element {
			return true
		}
	}
	return false
}
func IsContainInt(target int, List []int) bool {
	for _, element := range List {
		if target == element {
			return true
		}
	}
	return false
}
func InterfaceArrayToStringArray(data []interface{}) (i []string) {
	for _, param := range data {
		i = append(i, param.(string))
	}
	return i
}
func StructToJsonString(param interface{}) string {
	dataType, _ := json.Marshal(param)
	dataString := string(dataType)
	return dataString
}

func StructToJsonBytes(param interface{}) []byte {
	dataType, _ := json.Marshal(param)
	return dataType
}

// The incoming parameter must be a pointer
func JsonStringToStruct(s string, args interface{}) error {
	err := json.Unmarshal([]byte(s), args)
	return err
}

func GetMsgID(sendID string) string {
	t := int64ToString(GetCurrentTimestampByNano())
	return Md5(t + sendID + int64ToString(rand.Int63n(GetCurrentTimestampByNano())))
}

func int64ToString(i int64) string {
	return strconv.FormatInt(i, 10)
}
func Int64ToString(i int64) string {
	return strconv.FormatInt(i, 10)
}

func RemoveDuplicateElement(idList []string) []string {
	result := make([]string, 0, len(idList))
	temp := map[string]struct{}{}
	for _, item := range idList {
		if _, ok := temp[item]; !ok {
			temp[item] = struct{}{}
			result = append(result, item)
		}
	}
	return result
}

func RemoveDuplicate[T comparable](arr []T) []T {
	result := make([]T, 0, len(arr))
	temp := map[T]struct{}{}
	for _, item := range arr {
		if _, ok := temp[item]; !ok {
			temp[item] = struct{}{}
			result = append(result, item)
		}
	}
	return result
}

func IsDuplicateStringSlice(arr []string) bool {
	t := make(map[string]struct{})
	for _, s := range arr {
		if _, ok := t[s]; ok {
			return true
		}
		t[s] = struct{}{}
	}
	return false
}

// 生成指定长度的随机字符串
func GenerateRandomString(length int) string {
	const charset = "$#@_-abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	charsetLength := len(charset)
	result := make([]byte, length)

	for i := 0; i < length; i++ {
		randomIndex := rand.Intn(charsetLength)
		result[i] = charset[randomIndex]
	}
	return string(result)
}

// 生成指定长度的随机字符串
func GenerateRandomStringExcludeSpecialCharacter(length int) string {
	const charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	charsetLength := len(charset)
	result := make([]byte, length)

	for i := 0; i < length; i++ {
		randomIndex := rand.Intn(charsetLength)
		result[i] = charset[randomIndex]
	}
	return string(result)
}

// 生成指定长度的随机数字字符串
func GenerateDigitalString(length int) string {
	const charset = "0123456789"
	charsetLength := len(charset)
	result := make([]byte, length)

	for i := 0; i < length; i++ {
		randomIndex := rand.Intn(charsetLength)
		result[i] = charset[randomIndex]
	}
	return string(result)
}

func FindNumbers(text string) []int {
	// 使用正则表达式匹配所有数字
	re := regexp.MustCompile(`\d+`)
	matches := re.FindAllString(text, -1)
	// 将匹配到的数字转换为整数
	numbers := make([]int, len(matches))
	for i, match := range matches {
		number, _ := strconv.Atoi(match)
		numbers[i] = number
	}
	// 对数字数组进行排序
	sort.Ints(numbers)
	return numbers
}

func FindFirstNumber(text string) *int {
	numbers := FindNumbers(text)
	if len(numbers) > 0 {
		return &numbers[0]
	}
	return nil
}

func FindPhone(text string) string {
	// 使用正则表达式匹配所有数字
	re := regexp.MustCompile(`\d{11}`)
	return re.FindString(text)
}

func FindIdCard(text string) string {
	// 使用正则表达式匹配所有数字
	re := regexp.MustCompile(`^[1-9]\d{5}(18|19|20)\d{2}((0[1-9]{1})|(10|11|12))(([0-2][1-9]{1})|(30|31))\d{3}(\d|X|x)$`)
	return re.FindString(text)
}

// 身份证号码中的出生日期格式为YYYYMMDD
func ParseBirthDateFromIdCard(id string) (*customtypes.Time, error) {
	// 正则表达式，用于匹配身份证号码中的出生年月日部分
	re := regexp.MustCompile(`(?m)^\d{6}(\d{4})(\d{2})(\d{2})`)
	// 使用FindStringSubmatch找到匹配的出生日期部分
	matches := re.FindStringSubmatch(id)
	if len(matches) < 4 {
		return nil, fmt.Errorf("invalid ID card number format")
	}
	// 提取出生日期字符串
	birthDateStr := fmt.Sprintf("%s%s%s", matches[1], matches[2], matches[3])
	// 将出生日期字符串转换为 time.Time 类型
	birthDate, err := time.Parse("20060102", birthDateStr)
	if err != nil {
		return nil, err
	}
	result := customtypes.Time(birthDate)
	return &result, nil
}
