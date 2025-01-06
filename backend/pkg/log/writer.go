package log

import (
	"fmt"
	"os"
	"path"
	"sync"
	"time"
)

type (
	RotatePolicy string
	LoggerWriter interface {
	}

	FileWriter struct {
		sync.RWMutex

		// 配置
		policy         RotatePolicy
		fPath          string
		rotateByteCnt  uint64
		rotateDuration time.Duration
		backupFileCnt  int
		timeFormat     string

		// 状态
		writeByteCnt     uint64      // 已经写入byte
		writeFileNo      int         // 当前日志文件编号
		writerCreateTime time.Time   // 上次创建writer的时间
		fileNameChan     chan string // 文件名管道。 时间轮转时，使用带有buffer的chan。完成删除多余备份文件的功能
		writer           *os.File
	}

	StdoutWriter struct{}
)

const (
	RotatePolicyNone RotatePolicy = "none"
	RotatePolicySize RotatePolicy = "size"
	RotatePolicyTime RotatePolicy = "time"
)

func NewFileWriter(fpath string, policy RotatePolicy, rotateByte uint64, rotateDuration time.Duration, backupCnt int, timeFormat string) (*FileWriter, error) {
	_ = os.MkdirAll(path.Dir(fpath), 755)
	f, err := os.OpenFile(fpath, os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		return nil, err
	}

	w := &FileWriter{
		policy:         policy,
		fPath:          fpath,
		rotateByteCnt:  rotateByte,
		rotateDuration: rotateDuration,
		backupFileCnt:  backupCnt,
		timeFormat:     timeFormat,

		writeByteCnt:     0,
		writeFileNo:      0,
		writerCreateTime: time.Now(),
		fileNameChan:     make(chan string, backupCnt), // 保留多少个文件就缓存多大
		writer:           f,
	}

	return w, nil
}

func (w *FileWriter) Write(p []byte) (n int, err error) {
	var (
		pLen = uint64(len(p))
	)

	w.Lock()
	defer w.Unlock()
	for {
		// 检查是否需要轮转
		if !w.needRotate(pLen) {
			return w.writer.Write(p)
		}
		// 轮转
		err = w.renewWriter()
		if err != nil {
			return 0, err
		}
	}
}

func (w *FileWriter) needRotate(wantWriteLen uint64) bool {
	switch w.policy {
	case RotatePolicySize:
		if w.writeByteCnt+wantWriteLen > w.rotateByteCnt {
			return true
		}
		return false
	case RotatePolicyTime:
		now := time.Now()
		if w.writerCreateTime.Add(w.rotateDuration).Before(now) {
			return true
		}
		return false
	case RotatePolicyNone:
		fallthrough
	default:
		return false
	}
}

func (w *FileWriter) renewWriter() error {
	ext := path.Ext(w.fPath)
	other := w.fPath[:len(w.fPath)-len(ext)]
	switch w.policy {
	case RotatePolicySize:
		_ = w.writer.Close()
		// 移动文件
		newFname := fmt.Sprintf("%s%d.%s", other, w.writeFileNo, ext)
		_ = os.Rename(w.fPath, newFname)

		w.writeFileNo = (w.writeFileNo + 1) % w.backupFileCnt
		w.writeByteCnt = 0
		// 重置writer
		f, err := os.OpenFile(w.fPath, os.O_TRUNC|os.O_CREATE|os.O_RDWR, 0644)
		if err != nil {
			return err
		}
		w.writer = f
		return nil
	case RotatePolicyTime:
		_ = w.writer.Close()
		// 移动文件
		newFname := fmt.Sprintf("%s%s.%s", other, w.writerCreateTime.Format(w.timeFormat), ext)
		_ = os.Rename(w.fPath, newFname)

		w.writerCreateTime = time.Now()
		// 重置writer
		f, err := os.OpenFile(w.fPath, os.O_TRUNC|os.O_CREATE|os.O_RDWR, 0644)
		if err != nil {
			return err
		}
		w.writer = f
		go func() {
			// 通知删除一个文件
			w.notifyDeleteFile()
		}()
		// 将文件名写入管道
		w.fileNameChan <- newFname
		return nil
	case RotatePolicyNone:
		fallthrough
	default:
		return nil
	}

}

func (w *FileWriter) notifyDeleteFile() {
	filename := <-w.fileNameChan
	_ = os.RemoveAll(filename)
}

func NewStdoutWriter() *StdoutWriter {
	return &StdoutWriter{}
}

func (w *StdoutWriter) Write(p []byte) (n int, err error) {
	return os.Stdout.Write(p)
}
