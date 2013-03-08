package interpreter

import (
	"bytes"
	"errors"
	"io/ioutil"
	"net"
	"os"
	"os/exec"
	"path"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"syscall"
)

const RUBY_INTERPRETER_SOCKET_NAME = "/tmp/train.interpreter.socket"

type Interpreter struct {
	cmd     *exec.Cmd
	started bool
	mutex   sync.Mutex
}

var interpreter Interpreter

var Config struct {
	SASS struct {
		DebugInfo   bool
		LineNumbers bool
	}
}

func Compile(filePath string) (result string, err error) {
	fileExt := path.Ext(filePath)
	switch fileExt {
	case ".sass", ".scss", ".coffee":
		content, e := ioutil.ReadFile(filePath)
		if e != nil {
			panic(err)
		}
		result, err = interpreter.Render(strings.Replace(fileExt, ".", "", 1), content)
	default:
		err = errors.New("Unsupported format (" + filePath + "). Valid formats are: sass.")
	}

	return
}

func CloseInterpreter() {
	_, goFile, _, _ := runtime.Caller(0)
	pidFile := path.Dir(goFile) + "/interpreter.pid"

	if _, err := os.Stat(pidFile); err != nil && os.IsNotExist(err) {
		return
	}

	dat, err := ioutil.ReadFile(pidFile)
	if err != nil {
		panic(err)
	}

	err = exec.Command("rm", pidFile).Run()
	if err != nil {
		panic(err)
	}

	pid, err := strconv.Atoi(string(dat))
	if err != nil {
		panic(err)
	}

	err = syscall.Kill(pid, syscall.Signal(9))
	if err != nil {
		panic(err)
	}
}

func (this *Interpreter) Render(format string, content []byte) (result string, err error) {
	this.StartRubyInterpreter()

	conn, err := net.Dial("unix", RUBY_INTERPRETER_SOCKET_NAME)
	if err != nil {
		panic(err)
	}

	option := getOption()

	conn.Write([]byte(format + "<<" + option + "<<" + string(content)))
	var data bytes.Buffer
	data.ReadFrom(conn)
	conn.Close()

	compiled := strings.Split(data.String(), "<<")
	status := compiled[0]
	result = compiled[1]

	if status == "error" {
		err = errors.New("Could not compile " + format + ":\n" + result)
	}

	return
}

func (this *Interpreter) StartRubyInterpreter() {
	if this.started {
		return
	}

	this.mutex.Lock()

	_, goFile, _, _ := runtime.Caller(0)
	this.cmd = exec.Command("ruby", path.Dir(goFile)+"/interpreter.rb")
	waitForStarting := make(chan bool)
	this.cmd.Stdout = &StdoutCapturer{waitForStarting}
	go func() {
		err := this.cmd.Run()
		if err != nil {
			panic(err)
		}
	}()
	<-waitForStarting

	this.started = true
	this.mutex.Unlock()
}

type StdoutCapturer struct {
	waitForStarting chan bool
}

func (this *StdoutCapturer) Write(p []byte) (n int, err error) {
	if strings.Contains(string(p), "<<ready") {
		this.waitForStarting <- true
	}
	return
}

func getOption() string {
	if Config.SASS.LineNumbers {
		return "line_numbers"
	}
	if Config.SASS.DebugInfo {
		return "debug_info"
	}
	return ""
}
