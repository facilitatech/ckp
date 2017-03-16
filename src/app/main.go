// dependency-check-php - Analisa dependências de programas php
// https://github.com/totalbr/dependency-check-php for the canonical source repository
// Copyright (c) facilita.tech - 2016-2017 (http://facilita.tech)

package main

import (
	"io/ioutil"
	"os"
	"strings"
	"bufio"
	"log"
	"github.com/agtorre/gocolorize"
	"syscall"
	"unsafe"
)

var (
	INFO     *log.Logger
	WARNING  *log.Logger
	CRITICAL *log.Logger
	BLUE     *log.Logger
	path     string
	i 	 func(v ...interface{}) string
	w 	 func(v ...interface{}) string
	c 	 func(v ...interface{}) string
	b 	 func(v ...interface{}) string
	winsize  int
)

type Winsize struct {
	Row    uint16
	Col    uint16
	Xpixel uint16
	Ypixel uint16
}

func main() {

	winsize = getWidth()

	info := gocolorize.NewColor("green+h:black")
	blue := gocolorize.NewColor("black+i:black")
	warning := gocolorize.NewColor("yellow+i:black")
	critical := gocolorize.NewColor("black+i:red")

	i = info.Paint
	w = warning.Paint
	c = critical.Paint
	b = blue.Paint

	INFO     = log.New(os.Stdout, i("SCANNING       "), log.Ldate|log.Lmicroseconds|log.Lshortfile)
	WARNING  = log.New(os.Stdout, w("FOUND          "), log.Ldate|log.Lmicroseconds|log.Lshortfile)
	CRITICAL = log.New(os.Stdout, c("NOT FOUND      "), log.Ldate|log.Lmicroseconds|log.Lshortfile)
	BLUE     = log.New(os.Stdout, b("               "), log.Ldate|log.Lmicroseconds|log.Lshortfile)

	if len(os.Args) == 3 {
		if os.Args[1] == "--check"  && os.Args[2] != "" {
			path = os.Args[2]
			err := readDir(os.Args[2], false)
			if err != nil {
				panic(err)
			}
		}
	}
}

func readDir(diretorio string, signal bool) (error) {
	files, err := ioutil.ReadDir(diretorio)
	if err != nil {
		return err
	}
	for _, file := range files {
		// check for files with extension .php
		if strings.Contains(file.Name(), ".php") {
			filename := file.Name()
			if signal {
				filename = diretorio + "/" + file.Name()
			}
			readFile(filename, "null", signal)
		} else if file.IsDir() {
			readDir(diretorio +"/"+ file.Name(), true)
		}
	}
	return nil
}

func generateLog(dependencia, fileorigem string) {
	// Check if file exists
	_, err := os.Stat("dependency_logs.txt")
	if err != nil {
		INFO.Println("Create file for log generation")
		// Create a new file
		file, err := os.Create("dependency_logs.txt")
		if err != nil {
			log.Fatalln(err)
		}
		if err := file.Close(); err != nil {
			log.Fatalln(err)
		}
		return
	}
	file, err := os.OpenFile("dependency_logs.txt", os.O_CREATE|os.O_RDWR, 0755)
	if err != nil {
		log.Fatalln(err)
	}
	// Close and remove the file after main finishes execution
	defer func() {
		if err := file.Close(); err != nil {
			log.Fatalln(err)
		}
	}()
	newtext := generateSpaces(" " + dependencia + " origin -> " + fileorigem)
	CRITICAL.Println(c(newtext))
}

func readFile(file, anterior string, signal bool) {
	pathFile := path +"/"+ file

	if signal {
		pathFile = file
	}

	nFile, err := os.Open(pathFile)
	if err != nil {
		generateLog(pathFile, anterior)
	} else {
		newtext := generateSpaces(" file: " + pathFile)
		INFO.Println(i(newtext))

		scanner := bufio.NewScanner(nFile)
		scanner.Split(bufio.ScanLines)

		newtext = generateSpaces(" dependency")
		INFO.Println(i(newtext))
		for scanner.Scan() {
			text := scanner.Text()
			indexRequire := strings.Index(text, "require")
			if indexRequire == 0 {
				split := strings.Split(text, "\"")
				if len(split) == 3 {
					newtext = generateSpaces(" [ require ] found: " + split[1] + " in file -> " + pathFile)
					WARNING.Println(w(newtext))
					if strings.Contains(split[1], ".php") {
						readFile(split[1], pathFile, false)
					}
				}
			}
			indexInclude := strings.Index(text, "include")
			if indexInclude == 0 {
				split := strings.Split(text, "\"")
				if len(split) == 3 {
					newtext = generateSpaces(" [ include ] found: " + split[1] + " in file -> " + pathFile)
					WARNING.Println(w(newtext))
					if strings.Contains(split[1], ".php") {
						readFile(split[1], pathFile, false)
					}
				}
			}
		}
	}
}

func generateSpaces(str string) string {
	length := (winsize-55)-len(str)
	s3 := []byte(str)
	for i := 0; i < length; i++ {
		s3 = append(s3, '\u0020')
	}
	return string(s3)
}

func getWidth() int {
	ws := &Winsize{}
	retCode, _, errno := syscall.Syscall(syscall.SYS_IOCTL,
		uintptr(syscall.Stdin),
		uintptr(syscall.TIOCGWINSZ),
		uintptr(unsafe.Pointer(ws)))

	if int(retCode) == -1 {
		panic(errno)
	}
	return int(ws.Col)
}