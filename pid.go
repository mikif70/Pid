package pidlib

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"
)

// PID struttura
type PID struct {
	pid     string
	PIDFile string
	exe     string
	path    string
}

// New crea un nuovo PID
func New() *PID {
	p := PID{}
	p.init()
	return &p
}

func (p *PID) init() error {
	var err error
	p.path, err = filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		return err
	}
	p.exe = path.Base(os.Args[0])
	p.PIDFile = path.Join(p.path, p.exe) + ".run"

	return nil
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
	if !strings.Contains(string(cmd), p.exe) {
		fmt.Printf("PID %s used by %s\n", p, cmd)
		return true
	}

	return true
}

// Write scrive il PID nel PIDFile
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

// Remove Cancella il PIDFile
func (p *PID) Remove() {
	err := os.Remove(p.PIDFile)
	if err != nil {

	}
}
