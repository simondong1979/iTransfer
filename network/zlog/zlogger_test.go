package zlog_test

import (
	"github.com/simondong1979/iTransfer/network/zlog"
	"testing"
)

func TestStdZLog(t *testing.T) {

	//测试 默认debug输出
	zlog.Debug("network debug content1")
	zlog.Debug("network debug content2")

	zlog.Debugf(" network debug a = %d\n", 10)

	//设置log标记位，加上长文件名称 和 微秒 标记
	zlog.ResetFlags(zlog.BitDate | zlog.BitLongFile | zlog.BitLevel)
	zlog.Info("network info content")

	//设置日志前缀，主要标记当前日志模块
	zlog.SetPrefix("MODULE")
	zlog.Error("network error content")

	//添加标记位
	zlog.AddFlag(zlog.BitShortFile | zlog.BitTime)
	zlog.Stack(" network Stack! ")

	//设置日志写入文件
	zlog.SetLogFile("./log", "testfile.log")
	zlog.Debug("===> network debug content ~~666")
	zlog.Debug("===> network debug content ~~888")
	zlog.Error("===> network Error!!!! ~~~555~~~")

	//关闭debug调试
	zlog.CloseDebug()
	zlog.Debug("===> 我不应该出现~！")
	zlog.Debug("===> 我不应该出现~！")
	zlog.Error("===> network Error  after debug close !!!!")
}
