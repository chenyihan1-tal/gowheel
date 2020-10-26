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
	var hBody []byte

	marshal, _ := json.Marshal(data)
	hBody = marshal

	if reflect.TypeOf(data).String() == "string" {
		hBody = []byte(data.(string))
	}
	if strings.Contains(reflect.TypeOf(data).String(), "int") {
		hBody = []byte(fmt.Sprintf("%v", data))
	}
	if strings.Contains(reflect.TypeOf(data).String(), "float") {
		hBody = []byte(fmt.Sprintf("%v", data))
	}
	if strings.Contains(reflect.TypeOf(data).String(), "[]uint8") {
		hBody = data.([]byte)
	}

	h.Write(hBody)

	return hex.EncodeToString(h.Sum(nil))
}

// RandStr 获取指定长度的随机字符串
func RandStr(length int) string {
	str := uuid.NewV4().String()

	str = MD5(str)
	rStr := []rune(str)

	for i := 0; i < len(rStr); i++ {
		rand.Seed(time.Now().UnixNano())
		randI := rand.Intn(len(rStr))
		s := fmt.Sprintf("%c", rStr[randI])
		str = strings.Replace(str, s, strings.ToUpper(s), 1)
	}

	if length > len([]rune(str)) {
		length = len([]rune(str))
	}

	return str[0:length]
}

// Request 发送HTTP请求
func Request(method, url string, data, header map[string]interface{}) (body []byte, err error) {
	url = strings.ReplaceAll(strings.ReplaceAll(strings.ReplaceAll(url, "\n", ""), " ", ""), "\r", "")

	marshal, _ := json.Marshal(data)

	client := &http.Client{}
	req, err := http.NewRequest(method, url, bytes.NewReader(marshal))
	if err != nil {
		err = fmt.Errorf("new request fail: %s", err.Error())
		return
	}

	for k, v := range header {
		req.Header.Add(k, fmt.Sprintf("%s", v))
	}

	res, err := client.Do(req)
	if err != nil {
		err = fmt.Errorf("do request fail: %s", err.Error())
		return
	}

	defer res.Body.Close()
	body, err = ioutil.ReadAll(res.Body)
	if err != nil {
		err = fmt.Errorf("read res body fail: %s", err.Error())
		return
	}

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
