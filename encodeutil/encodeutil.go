package encodeutil

import "golang.org/x/text/encoding/simplifiedchinese"

// Charset 表示一个合法的字符集。
type Charset string

const (
	UTF8    = Charset("UTF-8")   // UTF-8 编码
	GB18030 = Charset("GB18030") // GB18030 编码
)

// ByteToString 将给定字符集的二进制数据转为 go 字符串。
func ByteToString(byte []byte, charset Charset) string {
	var str string
	switch charset {
	case GB18030:
		var decodeBytes, _ = simplifiedchinese.GB18030.NewDecoder().Bytes(byte)
		str = string(decodeBytes)
	case UTF8:
		fallthrough
	default:
		str = string(byte)
	}
	return str
}
