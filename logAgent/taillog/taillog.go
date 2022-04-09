package taillog

import (
	"fmt"
	"github.com/hpcloud/tail"
	"time"
)

//专门从日志文件收集日志的模块
func main() {
	fileName := "./my.log" //日志文件路径和日志名称
	config := tail.Config{
		ReOpen:    true,                                 //重新打开，日志切割
		Follow:    true,                                 //是否跟随上次关闭的文件
		Location:  &tail.SeekInfo{Offset: 0, Whence: 2}, //从文件哪个地方开始读
		MustExist: false,                                //文件不存在不报错
		Poll:      true,
	}
	tails, err := tail.TailFile(fileName, config)
	if err != nil {
		fmt.Println("tail file failed, err:", err)
		return
	}
	var (
		line *tail.Line
		ok   bool
	)
	for {
		line, ok = <-tails.Lines
		if !ok {
			fmt.Printf("tail file close reopen, filename: %s\n", tails.Filename)
			time.Sleep(time.Second)
			continue
		}
		fmt.Println("line:", line.Text)
	}
}
