package utils

import "time"

// 根据年份计算生肖
func GetChineseZodiac(birthDate time.Time) string {
	year := birthDate.Year()
	animals := [12]string{"鼠", "牛", "虎", "兔", "龙", "蛇", "马", "羊", "猴", "鸡", "狗", "猪"}
	return animals[(year-1900)%12]
}

// 根据月份和日期计算星座
func GetZodiac(birthDate time.Time) string {
	month := birthDate.Month()
	day := birthDate.Day()
	zodiacs := [12]string{"摩羯", "水瓶", "双鱼", "白羊", "金牛", "双子", "巨蟹", "狮子", "处女", "天秤", "天蝎", "射手"}

	switch month {
	case time.January:
		if day < 20 {
			return zodiacs[0]
		}
		return zodiacs[1]
	case time.February:
		if day < 19 {
			return zodiacs[1]
		}
		return zodiacs[2]
	case time.March:
		if day < 21 {
			return zodiacs[2]
		}
		return zodiacs[3]
	case time.April:
		if day < 20 {
			return zodiacs[3]
		}
		return zodiacs[4]
	case time.May:
		if day < 21 {
			return zodiacs[4]
		}
		return zodiacs[5]
	case time.June:
		if day < 22 {
			return zodiacs[5]
		}
		return zodiacs[6]
	case time.July:
		if day < 23 {
			return zodiacs[6]
		}
		return zodiacs[7]
	case time.August:
		if day < 23 {
			return zodiacs[7]
		}
		return zodiacs[8]
	case time.September:
		if day < 23 {
			return zodiacs[8]
		}
		return zodiacs[9]
	case time.October:
		if day < 24 {
			return zodiacs[9]
		}
		return zodiacs[10]
	case time.November:
		if day < 23 {
			return zodiacs[10]
		}
		return zodiacs[11]
	default:
		return zodiacs[11]
	}
}

// 根据当前年份和出生年份计算年龄
func GetAge(birthDate time.Time) int {
	return time.Now().Year() - birthDate.Year()
}

// 根据当前年份和出生年份计算年龄
func GetYear(birthDate time.Time) int {
	return birthDate.Year()
}
