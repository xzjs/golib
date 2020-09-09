package golib

import (
	"strconv"
	"strings"
)

// GetID 从字符串中获取uint的id
func GetID(str string) uint {
	tmp, err := strconv.Atoi(str)
	if err != nil {
		return 0
	}
	return uint(tmp)
}

//ConvertToByte 转换编码格式
// func ConvertToByte(src string, srcCode string, targetCode string) []byte {
// 	srcCoder := mahonia.NewDecoder(srcCode)
// 	srcResult := srcCoder.ConvertString(src)
// 	tagCoder := mahonia.NewDecoder(targetCode)
// 	_, cdata, _ := tagCoder.Translate([]byte(srcResult), true)
// 	return cdata
// }

//UnescapeUnicode 转换编码
func UnescapeUnicode(raw []byte) ([]byte, error) {
	str, err := strconv.Unquote(strings.Replace(strconv.Quote(string(raw)), `\\u`, `\u`, -1))
	if err != nil {
		return nil, err
	}
	return []byte(str), nil
}
