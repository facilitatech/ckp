// ckp - Check PHP files
// https://github.com/facilitatech/ckp/ for the canonical source repository
// Copyright (c) facilita.tech - 2016-2018 (http://facilita.tech)

package main

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/agtorre/gocolorize"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
	"unsafe"
)

var (
	scanning        *log.Logger
	found           *log.Logger
	notFound        *log.Logger
	result          *log.Logger
	empty           *log.Logger
	path            string
	scanningPrint   func(v ...interface{}) string
	foundPrint      func(v ...interface{}) string
	notFoundPrint   func(v ...interface{}) string
	resultPrint     func(v ...interface{}) string
	winsize         int
	logger          []string
	dir             []string
	files           []string
	filesExists     []string
	filesDontExists []string
	filesDiffers    []string
	IgnoreFolders   []string
	params          = new(Params)
)

const (
	nameDirDiffs = "diffs"
)

type Winsize struct {
	Row    uint16
	Col    uint16
	Xpixel uint16
	Ypixel uint16
}

type Params struct {
	options []string
}

func (p *Params) Set(params []string) {
	for i := range params {
		p.options = append(p.options, params[i])
	}
}

func (p *Params) GetAll() []string {
	return p.options
}

func (p *Params) Get(name int) string {
	return p.options[name]
}

func (p *Params) Count() int {
	return len(p.options)
}

func (p *Params) hasExport() bool {
	for i := range p.options {
		if p.options[i] == "--export" {
			return true
		}
	}
	return false
}

func main() {

	// Set arguments passed to ckp
	params.Set(os.Args)

	// Get size of terminal window where the ckp is running.
	winsize = getWidth()

	scanningColor := gocolorize.NewColor("green+h:black")
	resultColor := gocolorize.NewColor("white+h:black")
	foundColor := gocolorize.NewColor("black+i:yellow")
	notFoundColor := gocolorize.NewColor("black+i:red")
	scanningPrint = scanningColor.Paint
	foundPrint = foundColor.Paint
	notFoundPrint = notFoundColor.Paint
	resultPrint = resultColor.Paint
	scanning = log.New(os.Stdout, scanningPrint("Scanning  -->  "), 0)
	found = log.New(os.Stdout, foundPrint("Found          "), 0)
	notFound = log.New(os.Stdout, notFoundPrint("Not found      "), 0)
	result = log.New(os.Stdout, resultPrint("Result    -->  "), 0)
	empty = log.New(os.Stdout, resultPrint("               "), 0)

	// using this only for analysis of the dependencies of
	// the programs PHP at the moment
	if len(os.Args) == 3 {
		if os.Args[1] == "--check-dependencies" && os.Args[2] != "" {
			path = os.Args[2]
			// initiate read directories
			readDir(os.Args[2], false, "php")
			params.ResultDisplay()
		}
	}

	// This session initialize diff analysis and your options
	// --ignore Ignore folders who which are not part of the process
	// --export Export the data obtained from the diffs
	if len(os.Args) >= 4 {
		if os.Args[1] == "--diff" && os.Args[2] != "" && os.Args[3] != "" {
			path = os.Args[2]

			if len(os.Args) >= 6 {
				if os.Args[4] == "--ignore" && os.Args[5] != "" {
					split := strings.Split(os.Args[5], ",")
					for i := range split {
						removeSpace := strings.Trim(split[i], " ")
						IgnoreFolders = append(IgnoreFolders, removeSpace)
					}
				}
				if len(os.Args) == 7 {
					if os.Args[6] == "--export" {
						pwd, err := os.Getwd()
						if err != nil {
							panic(err)
						}
						err = os.MkdirAll(pwd+"/"+nameDirDiffs, 0755)
						if err != nil {
							panic(err)
						}
					}
				}
			}
			// initiate read directories
			pwd, err := os.Getwd()
			if err != nil {
				panic(err)
			}
			path := pwd + "/" + os.Args[2]
			readRecursiveDir(path, os.Args[2], os.Args[3])
			params.ResultDisplay()
		}
	}
}

func (p *Params) ResultDisplay() {
	// scan result
	for j := 0; j < 2; j++ {
		line := generateSpaces(" ")
		empty.Println(resultPrint(line))
	}

	if len(logger) != 0 {
		writeLog("dependency_logs.txt", logger)
	}

	if len(filesDiffers) != 0 {
		writeLog("differ_logs.txt", filesDiffers)
		if params.hasExport() {
			folderExported := generateSpaces("Folder exported: " + nameDirDiffs)
			empty.Println(resultPrint(folderExported))
		}
	}

	brokenDependencies := generateSpaces("Broken dependencies: " + strconv.Itoa(len(logger)))
	empty.Println(resultPrint(brokenDependencies))
	directoriesScanned := generateSpaces("Directories scanned: " + strconv.Itoa(len(dir)))
	empty.Println(resultPrint(directoriesScanned))
	filesOpened := generateSpaces("Files opened: " + strconv.Itoa(len(files)))
	empty.Println(resultPrint(filesOpened))
	filesDiffers := generateSpaces("Files differs: " + strconv.Itoa(len(filesDiffers)))
	empty.Println(resultPrint(filesDiffers))

	footer := generateSpaces(" ")
	empty.Println(resultPrint(footer))
}

func writeLog(fileToWrite string, data []string) error {
	openFile, err := os.OpenFile(fileToWrite, os.O_CREATE|os.O_RDWR, 0755)
	if err != nil {
		panic(err)
	}
	defer openFile.Close()
	w := bufio.NewWriter(openFile)

	for i := 0; i < len(data); i++ {
		fmt.Fprintf(w, "%v\n", data[i])
	}
	w.Flush()

	space := generateSpaces(" ")
	result.Println(resultPrint(space))

	pwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	newtext := generateSpaces("Log: " + pwd + "/" + fileToWrite)
	empty.Println(resultPrint(newtext))

	return err
}

func readDir(directory string, signal bool, extension string) {
	files, err := ioutil.ReadDir(directory)
	if err != nil {
		panic(err)
	}
	for _, file := range files {
		// check for files with extension .php
		if strings.Contains(file.Name(), "."+extension) {
			filename := file.Name()
			if signal {
				filename = directory + "/" + file.Name()
			}
			readFileDependencie(filename, "null", signal)
		} else if file.IsDir() {
			if len(dir) == 0 {
				dir = append(dir, directory+"/"+file.Name())
			}
			registerDir(directory + "/" + file.Name())
			readDir(directory+"/"+file.Name(), true, extension)
		}
	}
	return
}

func readRecursiveDir(directory, dirComFirst, dirComSecond string) {
	files, err := ioutil.ReadDir(directory)
	if err != nil {
		panic(err)
	}
	for _, file := range files {
		fileOrDirName := directory + "/" + file.Name()
		if file.IsDir() {
			// Check if file aren't in the []IgnoreFolders
			ignore, _ := inArray(file.Name(), IgnoreFolders)
			if !ignore {
				registerDir(fileOrDirName)
				readRecursiveDir(fileOrDirName, dirComFirst, dirComSecond)
			}
			// if is continue to another record
			continue
		}
		// if is not a folder put on into  -> openFiles -> compareBetweenTwoFiles
		params.CompareBetweenTwoFiles(openTwoFiles(fileOrDirName, dirComFirst, dirComSecond))
	}
}

func openTwoFiles(file, dirComFirst, dirComSecond string) ([]byte, []byte, string, string) {
	// Register file for doesn't scan again
	checkScann := registerFile(file)
	if !checkScann {

		fileToCompare := strings.Replace(file, dirComFirst, dirComSecond, -1)

		// Read the first file to compare with dt2
		dt1, err := ioutil.ReadFile(file)
		if err != nil {
			log.Fatal(err)
		}

		// Check if another file exists in the target
		_, err = os.Stat(fileToCompare)
		if err != nil {
			if os.IsNotExist(err) {
				register(fileToCompare, filesDontExists)
				return []byte{}, []byte{}, "", ""
			}
		}

		register(fileToCompare, filesExists)

		dt2, err := ioutil.ReadFile(fileToCompare)
		if err != nil {
			log.Fatal(err)
		}
		newtext := generateSpaces(" " + file)
		scanning.Println(scanningPrint(newtext))

		return dt1, dt2, fileToCompare, file
	}
	return []byte{}, []byte{}, "", ""
}

// execShell("sh", []string{"-c", cmd})
func execShell(command string, args []string) {
	cmd := exec.Command(command, args...)
	var waitStatus syscall.WaitStatus
	if err := cmd.Run(); err != nil {
		if err != nil {
			os.Stderr.WriteString(fmt.Sprintf("Error: %s\n", err.Error()))
		}
		if exitError, ok := err.(*exec.ExitError); ok {
			waitStatus = exitError.Sys().(syscall.WaitStatus)
			fmt.Printf("Output: %s\n", []byte(fmt.Sprintf("%d", waitStatus.ExitStatus())))
		}
	} else {
		// Success
		waitStatus = cmd.ProcessState.Sys().(syscall.WaitStatus)
		fmt.Printf("Output: %s\n", []byte(fmt.Sprintf("%d", waitStatus.ExitStatus())))
	}
}

func (p *Params) GenerateDiffFiles(b1, b2 string) {
	nameFile := strings.Replace(b1, "/", "_", -1)
	newName := nameFile + ".diff"
	cmd := "diff " + b1 + " " + b2

	output, err := exec.Command("sh", "-c", cmd).CombinedOutput()
	openFile, err := os.OpenFile(filepath.Join(nameDirDiffs, newName), os.O_CREATE|os.O_RDWR, 0755)
	if err != nil {
		panic(err)
	}
	defer openFile.Close()
	w := bufio.NewWriter(openFile)
	_, err = w.WriteString(string(output))
	if err != nil {
		panic(err)
	}
	w.Flush()
}

func (p *Params) CompareBetweenTwoFiles(b1, b2 []byte, file, fileToCompare string) {
	if file != "" {
		result := bytes.Compare(b1, b2)
		if result != 0 {
			registerDiffer(file)
			if params.hasExport() {
				go params.GenerateDiffFiles(fileToCompare, file)
			}
		}
	}
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

func readFileDependencie(file, anterior string, signal bool) {
	pathFile := path + "/" + file

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
							readFileDependencie(split[1], pathFile, false)
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
							readFileDependencie(split[1], pathFile, false)
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

func registerDiffer(name string) bool {
	exists, _ := inArray(name, filesDiffers)
	if !exists {
		filesDiffers = append(filesDiffers, name)
	}
	return exists
}

func register(name string, fileRegister []string) bool {
	exists, _ := inArray(name, fileRegister)
	if !exists {
		fileRegister = append(fileRegister, name)
	}
	return exists
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
	index = -1
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
	length := (winsize - 15) - len(str)
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
