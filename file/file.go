package file

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"os"
	"path"
	"runtime"
	"strings"
)

// GetSize get the file size
func GetSize(f multipart.File) (int, error) {
	content, err := ioutil.ReadAll(f)

	return len(content), err
}

// GetExt get the file ext
func GetExt(fileName string) string {
	return path.Ext(fileName)
}

// CheckNotExist check if the file exists
func CheckNotExist(src string) bool {
	_, err := os.Stat(src)

	return os.IsNotExist(err)
}

// CheckPermission check if the file has permission
func CheckPermission(src string) bool {
	_, err := os.Stat(src)

	return os.IsPermission(err)
}

// IsNotExistMkDir create a directory if it does not exist
func IsNotExistMkDir(src string) error {
	if notExist := CheckNotExist(src); notExist == true {
		if err := MkDir(src); err != nil {
			return err
		}
	}

	return nil
}

// MkDir create a directory
func MkDir(src string) error {
	err := os.MkdirAll(src, os.ModePerm)
	if err != nil {
		return err
	}

	return nil
}

// Open a file according to a specific mode
func Open(name string, flag int, perm os.FileMode) (*os.File, error) {
	f, err := os.OpenFile(name, flag, perm)
	if err != nil {
		return nil, err
	}

	return f, nil
}

// MustOpen maximize trying to open the file
func MustOpen(fileName, filePath string) (*os.File, error) {
	fileName = "/" + fileName

	perm := CheckPermission(filePath)
	if perm == true {
		return nil, fmt.Errorf("file.CheckPermission Permission denied src: %s", filePath)
	}

	err := IsNotExistMkDir(filePath)
	if err != nil {
		return nil, fmt.Errorf("file.IsNotExistMkDir src: %s, err: %v", filePath, err)
	}

	f, err := Open(filePath+fileName, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		return nil, fmt.Errorf("fail to open file :%v", err)
	}

	return f, nil
}

// ReadName ...
func ReadName(filePath string) (file []string, err error) {
	if CheckNotExist(filePath) {
		err = fmt.Errorf("file is not exists")
		return
	}

	var fileRoot string
	var dotCount int

	if runtime.GOOS == "windows" || runtime.GOOS == "plan9" {
		dotCount = strings.Count(filePath, "..\\")
	} else {
		dotCount = strings.Count(filePath, "../")
	}
	for i := 0; i <= dotCount; i++ {
		if runtime.GOOS == "windows" || runtime.GOOS == "plan9" {
			fileRoot += "..\\"
		} else {
			fileRoot += "../"
		}
	}

	stat, err := os.Stat(filePath)
	if err != nil {
		return
	}
	if !stat.IsDir() {
		file = append(file, fileRoot+stat.Name())
		return
	}
	// read dir
	dir, err := os.Open(filePath)
	if err != nil {
		return
	}
	fis, err := dir.Readdir(0)
	if err != nil {
		return
	}
	for _, fi := range fis {
		if fi.IsDir() {
			tempFile, err := ReadName(dir.Name() + "/" + fi.Name())
			if err != nil {
				return nil, err
			}
			file = append(file, tempFile...)
			continue
		}
		file = append(file, dir.Name()+"/"+fi.Name())
	}

	return
}

// Tar
func Tar(src, dst string) (err error) {

	if !strings.Contains(dst, "tar.gz") {
		dst += ".tar.gz"
	}

	var dstDir string
	var dstName string

	dstDir = path.Dir(dst)
	dstName = path.Base(dst)

	// if the file exists
	if CheckNotExist(src) {
		err = fmt.Errorf("src is not exists")
		return
	}

	// file write
	fw, err := MustOpen(dstName, dstDir)
	if err != nil {
		return
	}

	defer fw.Close()
	// gzip write
	gw := gzip.NewWriter(fw)
	defer gw.Close()
	// tar write
	tw := tar.NewWriter(gw)
	defer tw.Close()

	file, err := ReadName(src)
	if err != nil {
		return
	}

	var rootPath string
	stat, err := os.Stat(src)
	if err != nil {
		return
	}
	if stat.IsDir() {
		rootPath = path.Clean(src)
	} else {
		if runtime.GOOS == "windows" || runtime.GOOS == "plan9" {
			rootPath = "..\\"
		} else {
			rootPath = "../"
		}
	}

	for _, v := range file {
		// open file
		fr, err := os.Open(v)
		if err != nil {
			return err
		}
		stat, err := os.Stat(v)
		if err != nil {
			return err
		}
		// header
		h := new(tar.Header)
		h.Name = strings.Replace(v, rootPath, "", 1)
		h.Size = stat.Size()
		h.Mode = int64(stat.Mode())
		h.ModTime = stat.ModTime()
		// write header
		err = tw.WriteHeader(h)
		if err != nil {
			return err
		}
		// write file
		_, err = io.Copy(tw, fr)
		if err != nil {
			return err
		}
		fmt.Println(strings.Replace(v, rootPath, "", 1) + " ...")
	}
	return
}

// UnTar
func UnTar(src, dst string) (err error) {
	// if the file exists
	if CheckNotExist(src) {
		err = fmt.Errorf("src is not exists")
		return
	}
	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()
	gr, err := gzip.NewReader(srcFile)
	if err != nil {
		return err
	}
	defer gr.Close()
	tr := tar.NewReader(gr)
	for {
		hdr, err := tr.Next()
		if err != nil {
			if err == io.EOF {
				break
			} else {
				return err
			}
		}
		fmt.Println(dst+hdr.Name, "...")
		file, err := MustOpen(path.Base(hdr.Name), dst+path.Dir(hdr.Name))
		if err != nil {
			return err
		}
		io.Copy(file, tr)
	}
	return
}
