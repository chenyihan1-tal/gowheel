package tool

import (
	"bytes"
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"reflect"
	"strings"
	"time"

	"github.com/satori/go.uuid"
)

// MD5 加密
func MD5(data interface{}) string {
	h := md5.New()

	if reflect.TypeOf(data).Name() == "string" {
		h.Write([]byte(data.(string)))
	} else {
		marshal, _ := json.Marshal(data)
		h.Write(marshal)
	}

	return hex.EncodeToString(h.Sum(nil))
}

// RandStr 获取指定长度的随机字符串
func RandStr(length int) string {
	str := uuid.NewV4().String()

	rStr := []rune(str)
	for i := 0; i < len(rStr); i++ {
		rand.Seed(time.Now().UnixNano())
		randI := rand.Intn(30)
		s := fmt.Sprintf("%c", rStr[i])
		if (randI-i)%2 == 0 {
			str = strings.Replace(str, s, strings.ToUpper(s), 1)
		}
	}
	if length > len([]rune(str)) {
		length = len([]rune(str))
	}
	return str[0:length]
}

// Request 发送HTTP请求
func Request(method, url string, data, header map[string]interface{}) (body []byte, err error) {
	marshal, _ := json.Marshal(data)

	client := &http.Client{}
	req, err := http.NewRequest(method, url, bytes.NewReader(marshal))
	if err != nil {
		return
	}

	for k, v := range header {
		req.Header.Add(k, fmt.Sprintf("%s", v))
	}

	res, err := client.Do(req)
	if err != nil {
		return
	}

	defer res.Body.Close()
	body, err = ioutil.ReadAll(res.Body)

	return
}

// DecBase64Img base64 图片解码
func DecBase64Img(base64Str string) (data []byte, extension string, err error) {
	if base64Str[11] == 'j' {
		extension = ".jpg"
		base64Str = base64Str[23:]
	} else if base64Str[11] == 'p' {
		base64Str = base64Str[22:]
		extension = ".png"
	} else if base64Str[11] == 'g' {
		base64Str = base64Str[22:]
		extension = ".gif"
	}

	data, err = base64.StdEncoding.DecodeString(base64Str)

	return
}

// InArray ...
func InArray(need interface{}, needArr []interface{}) bool {
	for _, v := range needArr {
		if need == v {
			return true
		}
	}
	return false
}
