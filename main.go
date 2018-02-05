// ckp - Check PHP files
// https://github.com/facilitatech/ckp/ for the canonical source repository
// Copyright (c) facilita.tech - 2016-2018 (http://facilita.tech)

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
	"strconv"
	"fmt"
)

var (
	scanning      *log.Logger
	found         *log.Logger
	notFound      *log.Logger
	result        *log.Logger
	empty         *log.Logger
	path          string
	scanningPrint func(v ...interface{}) string
	foundPrint    func(v ...interface{}) string
	notFoundPrint func(v ...interface{}) string
	resultPrint   func(v ...interface{}) string
	winsize       int
	logger        []string
	dir           []string
	files         []string
)

type Winsize struct {
	Row    uint16
	Col    uint16
	Xpixel uint16
	Ypixel uint16
}

func main() {

	// get size of window
	winsize       =  getWidth()

	scanningColor := gocolorize.NewColor("green+h:black")
	resultColor   := gocolorize.NewColor("white+h:black")
	foundColor    := gocolorize.NewColor("black+i:yellow")
	notFoundColor := gocolorize.NewColor("black+i:red")
	scanningPrint =  scanningColor.Paint
	foundPrint    =  foundColor.Paint
	notFoundPrint =  notFoundColor.Paint
	resultPrint   =  resultColor.Paint
	scanning      =  log.New(os.Stdout, scanningPrint("Scanning  -->  "), 0)
	found         =  log.New(os.Stdout, foundPrint("Found          "), 0)
	notFound      =  log.New(os.Stdout, notFoundPrint("Not found      "), 0)
	result        =  log.New(os.Stdout, resultPrint("Result    -->  "), 0)
	empty         =  log.New(os.Stdout, resultPrint("               "), 0)

	if len(os.Args) == 3 {
		if os.Args[1] == "--check"  && os.Args[2] != "" {
			path = os.Args[2]
			// initiate read directories
			readDir(os.Args[2], false)
			resultDisplay()
		}
	}
}

func resultDisplay() {
	// scan result
	for j:=0;j<2;j++ {
		line := generateSpaces(" ")
		empty.Println(resultPrint(line))
	}
	line := generateSpaces(" Broken dependencies:")
	empty.Println(resultPrint(line))

	line = generateSpaces(" ")
	empty.Println(resultPrint(line))

	err := writeLog()
	if err != nil {
		newtext := generateSpaces("Error on write dependency_logs.txt log")
		result.Println(resultPrint(newtext))
	}

	for j:=0;j<2;j++ {
		line := generateSpaces(" ")
		empty.Println(resultPrint(line))
	}
	// scan Details
	total := generateSpaces(
		" Broken dependencies: " + strconv.Itoa(len(logger)) +
			"   |   " +
			"Directories scanned: " + strconv.Itoa(len(dir)) +
			"   |   " +
			"Files opened: " + strconv.Itoa(len(files)),
	)
	empty.Println(resultPrint(total))

	footer := generateSpaces(" ")
	empty.Println(resultPrint(footer))
}

func writeLog() error {
	file, err := os.OpenFile("dependency_logs.txt", os.O_CREATE|os.O_RDWR, 0755)
	if err != nil {
		newtext := generateSpaces("File dependency_logs.txt not found")
		result.Println(resultPrint(newtext))
	}
	defer file.Close()
	w := bufio.NewWriter(file)

	for i := 0; i < len(logger); i++ {
		newtext := generateSpaces(" " + logger[i])
		result.Println(resultPrint(newtext))
		fmt.Fprintf(w, "%v\n", newtext)
	}
	w.Flush()

	return err
}

func readDir(directory string, signal bool) {
	files, err := ioutil.ReadDir(directory)
	if err != nil {
		panic(err)
	}
	for _, file := range files {
		// check for files with extension .php
		if strings.Contains(file.Name(), ".php") {
			filename := file.Name()
			if signal {
				filename = directory + "/" + file.Name()
			}
			readFile(filename, "null", signal)
		} else if file.IsDir() {
			if len(dir) == 0 {
				dir = append(dir, directory +"/"+ file.Name())
			}
			registerDir(directory +"/"+ file.Name())
			readDir(directory +"/"+ file.Name(), true)
		}
	}
	return
}

func generateLog(dependencia, fileorigem string) {
	// Check if file exists
	_, err := os.Stat("dependency_logs.txt")
	if err != nil {
		scanning.Println("Create file for log generation")
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
	notFound.Println(notFoundPrint(newtext))

	text := dependencia + " origin -> " + fileorigem
	if len(logger) == 0 {
		logger = append(logger, text)
	}
	logger = registerLog(text, logger)
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
		if len(files) == 0 {
			files = append(files, pathFile)
		}
		checkScann := registerFile(pathFile)
		if !checkScann {
			//return
			newtext := generateSpaces(" " + pathFile)
			scanning.Println(scanningPrint(newtext))

			scanner := bufio.NewScanner(nFile)
			scanner.Split(bufio.ScanLines)

			// Only scan for "require*" or "include*" entries
			// @todo improvement for "use" namespaces
			for scanner.Scan() {
				text := scanner.Text()
				indexRequire := strings.Index(text, "require") // require or require_once
				if indexRequire != -1 {
					split := strings.Split(text, "\"")
					if len(split) == 3 {
						newtext = generateSpaces(" [ require ] found: " + split[1] + " in file -> " + pathFile)
						found.Println(foundPrint(newtext))
						if strings.Contains(split[1], ".php") {
							// only files *.php
							readFile(split[1], pathFile, false)
						}
					}
				}
				indexInclude := strings.Index(text, "include") // include or include_once
				if indexInclude != -1 {
					split := strings.Split(text, "\"")
					if len(split) == 3 {
						newtext = generateSpaces(" [ include ] found: " + split[1] + " in file -> " + pathFile)
						found.Println(foundPrint(newtext))
						if strings.Contains(split[1], ".php") {
							readFile(split[1], pathFile, false)
						}
					}
				}
			}
		}
	}
	return
}

func registerLog(text string, logger []string) []string {
	// check if logger exists
	exists, _ := inArray(text, logger)
	if !exists {
		// @todo improvement for scanning this routes ../
		index := strings.Index(text, "../")
		if index == -1 {
			logger = append(logger, text)
		}
	}
	return logger
}

func registerFile(name string) bool {
	exists, _ := inArray(name, files)
	if !exists {
		files = append(files, name)
	}
	return exists
}

func registerDir(name string) {
	exists, _ := inArray(name, dir)
	if !exists {
		dir = append(dir, name)
	}
}

func inArray(val string, array []string) (exists bool, index int) {
	exists = false
	index = -1;
	for i, v := range array {
		if val == v {
			index = i
			exists = true
			return
		}
	}
	return
}

func generateSpaces(str string) string {
	length := (winsize-15)-len(str)
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