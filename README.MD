<h1 align="center">GOWHEEL</h1>

<p align="center">
	<a><img width="200px" src="https://github.com/xavierror/gowheel/blob/master/gopher.png?raw=true"/></a>
</p>

<p align="center">Go package | Go 轮子 | 自用 | FOR PRIVATE USE</p>

## GETTING STARTED

```bash
go get github.com/xavierror/gowheel
```

## USING

#### logs
```go
package main

import (
	"fmt"
	"os"

	"github.com/xavierror/gowheel/logs"
)

func main() {
	logs.Debug("xxx")
	logs.Debug("xxx", "ooo")
	logs.Debug(fmt.Sprintf("debug: %s", os.Getenv("debug")))
	logs.Info("xxx")
	logs.Warn("xxx")
	logs.Error("xxx")
	logs.Fatal("xxx")
}
// output:
// [2020/09/18 14:46:57.510][D][main.go:11][xxx]
// [2020/09/18 14:46:57.511][D][main.go:12][xxx ooo]
// [2020/09/18 14:46:57.512][D][main.go:13][debug: true]
// [2020/09/18 14:46:57.513][I][main.go:14][xxx]
// [2020/09/18 14:46:57.515][W][main.go:15][xxx]
// [2020/09/18 14:46:57.516][E][main.go:16][xxx]
// [2020/09/18 14:46:57.518][F][main.go:17][xxx]

```

#### tool
```go
package main

import (
	"fmt"

	"github.com/xavierror/gowheel/tool"
)

func main()  {
	fmt.Println(tool.MD5("123456"))
	fmt.Println(tool.MD5(123456))
        fmt.Println(tool.RandStr(16))
	fmt.Println(tool.RandStr(32))
}

// output:
// e10adc3949ba59abbe56e057f20f883e
// e10adc3949ba59abbe56e057f20f883e
// 38d55b7b04b0c149
// 7A2778216AdbFb47dd0A93eFA054b24F
```

#### MORE
...