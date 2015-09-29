// pid
package pid

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strconv"
	"strings"
)

type PID struct {
	pid     string
	PIDFile string
	exe     string
}

func New() *PID {
	p := PID{}
	p.init()
	return &p
}

func (p *PID) init() {
	p.exe = path.Base(os.Args[0])
	p.PIDFile = "./" + p.exe + ".run"
}

// verifica se esiste il PIDFile;
// se esiste legge il PID e controlla se e' running il processo associato
func (p *PID) check() bool {
	bpid, err := ioutil.ReadFile(p.PIDFile)
	p.pid = strings.TrimRight(string(bpid), "\n")
	if err == nil && p.readCmd() {
		return true
	}
	return false
}

// controlla se esiste il processo associato al PID,
// se il cmd e' lo stesso e se e' in esecuzione.
func (p *PID) readCmd() bool {
	bcmd, err := ioutil.ReadFile(path.Join("/proc", p.pid, "cmdline"))
	// non esiste la dir relativa al PID su /proc
	if err != nil {
		fmt.Println("cmdline error: ", err)
		return false
	}
	cmd := bytes.Trim(bcmd, "\x00")
	if strings.Contains(string(cmd), p.exe) {
		return true
	} else {
		fmt.Printf("PID %s used by %s\n", p, cmd)
		return true
	}
	return true
}

// scrive il PID nel PIDFile
func (p *PID) Write() {

	if p.check() {
		fmt.Println("Running: ", p.pid)
		os.Exit(-6)
	}

	fpid, err := os.OpenFile(p.PIDFile, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		fmt.Println("PID file error: ", err.Error())
		os.Exit(-5)
	}
	fpid.WriteString(strconv.Itoa(os.Getpid()))
	fpid.Close()
}

// Cancella il PIDFile
func (p *PID) Remove() {
	err := os.Remove(p.PIDFile)
	if err != nil {

	}
}
