package main

import (
	"flag"
	"io"
	"log"
	"net/http"
	_ "net/http/pprof"
	"os"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/sirupsen/logrus"
	"jingxi.cn/tools/shadow/internal/shadow"
)

type CmdLine struct {
	isDebug  bool
	httpAddr string
	dbDir    string
	logDir   string
	pprof    string
}

var cmdline CmdLine

func init() {
	flag.BoolVar(&cmdline.isDebug, "debug", false, "true")
	flag.StringVar(&cmdline.httpAddr, "http", "0.0.0.0:8080", "0.0.0.0:8080")
	flag.StringVar(&cmdline.pprof, "pprof", "", "0.0.0.0:6060")
	flag.StringVar(&cmdline.dbDir, "dbdir", "/app/db", "/app/db")
	flag.StringVar(&cmdline.logDir, "log", "/app/log", "/app/log")
}

func runProfServer() {
	err := http.ListenAndServe(cmdline.pprof, nil)
	if err != nil {
		logrus.Fatalf("Failed to start pprof server %s: %+v", cmdline.pprof, err)
	}
}

func initLog() {
	path := filepath.Join(cmdline.logDir, "shadow.log")
	writer, _ := rotatelogs.New(
		path+".%Y-%m-%d %H:%M",
		rotatelogs.WithLinkName(path),
		rotatelogs.WithMaxAge(48*time.Hour),
		rotatelogs.WithRotationTime(4*time.Hour))

	if cmdline.isDebug {
		writers := []io.Writer{
			os.Stdout,
			writer,
		}
		fileAndStdoutWriter := io.MultiWriter(writers...)
		logrus.SetOutput(fileAndStdoutWriter)
		logrus.SetLevel(logrus.TraceLevel)
		gin.DefaultWriter = fileAndStdoutWriter
	} else {
		logrus.SetOutput(writer)
		logrus.SetLevel(logrus.ErrorLevel)
		gin.DefaultWriter = writer
	}
	logrus.SetReportCaller(true)
}

func main() {
	flag.Parse()
	initLog()
	if len(cmdline.httpAddr) < 1 {
		logrus.Fatalf("http listen address empty")
		return
	}
	if len(cmdline.dbDir) < 1 {
		logrus.Fatalf("db directory empty")
		return
	}
	if len(cmdline.pprof) > 0 {
		log.Printf("pprof listen on: %s", cmdline.pprof)
		go runProfServer()
	}
	if !cmdline.isDebug {
		gin.SetMode(gin.ReleaseMode)
	}
	shadow.NewApp().Run(cmdline.httpAddr, cmdline.dbDir)
}
