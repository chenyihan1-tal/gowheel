package text

import (
	"bytes"
	"io/ioutil"
	"regexp"
	"strings"

	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
)

func GbkToUtf8(src []byte) (str []byte, err error) {
	reader := transform.NewReader(bytes.NewReader(src), simplifiedchinese.GBK.NewDecoder())
	str, err = ioutil.ReadAll(reader)
	return
}

func Utf8ToGbk(src []byte) (str []byte, err error) {
	reader := transform.NewReader(bytes.NewReader(src), simplifiedchinese.GBK.NewEncoder())
	str, err = ioutil.ReadAll(reader)
	return
}

func RemoveHtml(src string) string {
	// 将HTML标签全转换成小写
	re, _ := regexp.Compile("<[\\S\\s]+?>")
	src = re.ReplaceAllStringFunc(src, strings.ToLower)

	// 去除STYLE
	re, _ = regexp.Compile("<style[\\S\\s]+?</style>")
	src = re.ReplaceAllString(src, "")

	// 去除SCRIPT
	re, _ = regexp.Compile("<script[\\S\\s]+?</script>")
	src = re.ReplaceAllString(src, "")

	// 去除所有尖括号内的HTML代码，并换成换行符
	re, _ = regexp.Compile("<[\\S\\s]+?>")
	src = re.ReplaceAllString(src, "\n")

	// 去除连续的换行符
	re, _ = regexp.Compile("\\s{2,}")
	src = re.ReplaceAllString(src, "\n")

	// 去除空格
	src = strings.ReplaceAll(src, "&nbsp;", "")

	return strings.TrimSpace(src)
}