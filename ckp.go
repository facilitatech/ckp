// ckp - Check PHP dependencies and diff whatever files
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
	scanning      *log.Logger
	found         *log.Logger
	notFound      *log.Logger
	result        *log.Logger
	empty         *log.Logger
	info          *log.Logger

	scanningPrint func(v ...interface{}) string
	foundPrint    func(v ...interface{}) string
	notFoundPrint func(v ...interface{}) string
	resultPrint   func(v ...interface{}) string
	infoPrint     func(v ...interface{}) string

	winsize       int

	brokenDependencyLogger,
	dependencyLogger,
	dependencyMapLogger,
	directoryLogger,
	files,
	filesDep,
	filesDiffers,
	ignoreFolders []string

	params        = new(Params)
	puts          = fmt.Println
)

const (
	nameDirDiffs = "diffs"
)

// Winsize have the sizes of the window terminal
// this is used for configure the printed colors
type Winsize struct {
	Row    uint16
	Col    uint16
	Xpixel uint16
	Ypixel uint16
}

// Params used for manager the parameters who a
// user has passed to the program
type Params struct {
	Options       []string
	IndexOf       map[string]int
	Position      map[int]string
	Deps          chan string
	ExcludedFiles []string
	FilterFiles   []string
	Path          string
}

func main() {

	// Set arguments passed to ckp
	params.Set(os.Args)

	// Get size of terminal window where the ckp is running.
	winsize = getWidth()

	params.Help()
	params.FilterFiles = params.SetFilesParams("--filter-file")
	params.ExcludedFiles = params.SetFilesParams("--excluded-file")

	scanningColor := gocolorize.NewColor("green+h:black")
	resultColor := gocolorize.NewColor("white+h:black")
	infoColor := gocolorize.NewColor("white+h:black")
	foundColor := gocolorize.NewColor("black+i:yellow")
	notFoundColor := gocolorize.NewColor("black+i:red")

	scanningPrint = scanningColor.Paint
	foundPrint = foundColor.Paint
	notFoundPrint = notFoundColor.Paint
	resultPrint = resultColor.Paint
	infoPrint = infoColor.Paint

	scanning = log.New(os.Stdout, scanningPrint("SCANNING  -->  "), 0)
	found = log.New(os.Stdout, foundPrint("FOUND          "), 0)
	notFound = log.New(os.Stdout, notFoundPrint("NOT FOUND      "), 0)
	result = log.New(os.Stdout, resultPrint("RESULT    -->  "), 0)
	info  = log.New(os.Stdout, infoPrint("INFO      -->  "), 0)
	empty = log.New(os.Stdout, resultPrint("               "), 0)

	// using this only for analysis of the dependencies of
	// the programs PHP at the moment
	if params.Count() >= 2 {
		params.BrokenDeps()
	}

	// This session initialize diff analysis and your options
	// --ignore Ignore folders who which are not part of the process
	// --export Export the data obtained from the diffs
	if params.Count() >= 4 {
		params.Check()
		params.Diff()
	}
}

/////////////////////////////////////////////////////
// Helper functions
/////////////////////////////////////////////////////

// GetAll return all parameters passed to the program
func (p *Params) GetAll() []string {
	return p.Options
}

// Get return one specific parameter per position in the slice
func (p *Params) Get(name int) string {
	return p.Options[name]
}

// Count return the total of parameters passed to program
func (p *Params) Count() int {
	return len(p.Options)
}

// IndexOf return the index of the parameter in the map per name
func (p *Params) GetIndexOf(name string) int {
	return p.IndexOf[name]
}

// Position return the name of the parameter in the map per index
func (p *Params) GetPosition(index int) string {
	return p.Position[index]
}

// Has check if the parameter exists in the map passed to program
func (p *Params) Has(option string) bool {
	for i := range p.Options {
		if p.Options[i] == option {
			return true
		}
	}
	return false
}

// Set is the setter function for parameters
func (p *Params) Set(params []string) {
	if len(params) > 1 {

		p.IndexOf = make(map[string]int, len(params))
		p.Position = make(map[int]string, len(params))

		for i := range params {
			p.Options = append(p.Options, params[i])
			p.IndexOf[params[i]] = int(i)
			p.Position[i] = params[i]
		}
	}
}

// SetFilesParams faz a leitura do arquivo passado como parâmetro
// --filter-file, --excluded-file ou qualquer outro parâmetro
// que passe uma lista de arquivos.
// usage: params.SetFilesParams("--filter-file")
func (p *Params) SetFilesParams(param string) (rdata []string) {
	if p.Has(param) {
		list := p.GetPosition(p.GetIndexOf(param) + 1)
		if list == "" {
			puts("Not found parameters from " + param)
			puts("Usage:")
			puts("    Help: ckp --help")
			os.Exit(2)
		}

		file, err := os.Open(list)
		if err != nil {
			panic(err)
		}
		defer file.Close()

		fs := bufio.NewScanner(file)
		for fs.Scan() {
			rdata = append(rdata, fs.Text())
		}

		if err := fs.Err(); err != nil {
			panic(err)
		}
	}
	return rdata
}

/////////////////////////////////////////////////////
// Parameters
/////////////////////////////////////////////////////

// Help this are the instructions for the users
func (p *Params) Help() {
	if p.Has("--help") {
		c := exec.Command("clear")
		c.Stdout = os.Stdout
		c.Run()

		puts("NAME")
		puts("      ckp - Check PHP files")
		puts("")
		puts("SYNOPSIS")
		puts("      ckp [OPTIONS]...")
		puts("")
		puts("DESCRIPTION")
		puts("      Check PHP dependencies and diff whatever files")
		puts(" ")
		puts("OPTIONS")
		puts("      --broken-deps    Check all broken dependencies of your project PHP has")
		puts(" ")
		puts("      --diff           Make diff between two folders recursively")
		puts(" ")
		puts("      --check          Check dependencies with two folders recursively")
		puts(" ")
		puts("      --filter-file    Filter files per list, work with --diff and --check")
		puts(" ")
		puts("      --ignore         Ignore folders")
		puts(" ")
		puts("      --export         Export diffs file into folder, this only work with --diff")
		puts(" ")
		puts("      --verbose        Print result")
		puts(" ")
		puts("      --dep-map        Dependency map, this only work with --check")
		puts(" ")
		puts("      --excluded-file  Ignore this files on the dependencies check, this only work with --check")
		puts(" ")
		puts("EXAMPLES")
		puts("      ckp --broken-deps /var/www/app --verbose")
		puts(" ")
		puts("      ckp --diff /var/www/app1 /var/www/app2 --verbose")
		puts("      ckp --diff /var/www/app1 /var/www/app2 --ignore vendor,bin --verbose")
		puts("      ckp --diff /var/www/app1 /var/www/app2 --ignore vendor --filter-file files.txt --verbose")
		puts("      ckp --diff /var/www/app1 /var/www/app2 --ignore vendor --filter-file files.txt --export --verbose")
		puts(" ")
		puts("      ckp --check /var/www/app --verbose")
		puts("      ckp --check /var/www/app --filter-file files.txt --verbose")
		puts("      ckp --check /var/www/app --filter-file files.txt --dep-map --verbose")
		puts("      ckp --check /var/www/app --filter-file files.txt --dep-map --excluded-file ignore-files.txt --verbose")
		puts(" ")
		puts("AUTHOR")
		puts("      Lucas Alves")
		puts(" ")
		puts("REPORTING BUGS")
		puts("      Report bugs on <https://github.com/facilitatech/ckp/issues>")
		puts(" ")
		puts("COPYRIGHT")
		puts("      Copyright (c) 2017 Facilita.tech Author.")
		puts("      BSD ")
		puts("      See the LICENSE: <https://github.com/facilitatech/ckp/blob/master/LICENSE>")
		puts(" ")
		os.Exit(2)
	}
}

// BrokenDeps
func (p *Params) BrokenDeps() {
	if p.Has("--broken-deps") {
		dirDependencies := p.GetPosition(p.GetIndexOf("--broken-deps") + 1)
		if dirDependencies == "" {
			puts("Not found folders for analysis!")
			puts("Usage:")
			puts("    Help: ckp --help")
			os.Exit(2)
		}
		p.Path = p.GetPosition(p.GetIndexOf("--broken-deps") + 1)

		if !p.IsFolderExists(p.Path) {
			puts("Not found folders for analysis!")
			puts("Usage:")
			puts("    Help: ckp --help")
			os.Exit(2)
		}
		// initiate read directories
		p.ReadDir(p.Path, false, "php")
		p.ResultDisplay()
	}
}

// Diff
func (p *Params) Diff() {
	if p.Has("--diff") {
		positionDiff := p.GetIndexOf("--diff") + 1
		if (p.Count() - positionDiff) < 2 {
			puts("Not found folders for analysis!")
			puts("Usage:")
			puts("    Help: ckp --help")
			os.Exit(2)
		}

		firstPath := p.GetPosition(p.GetIndexOf("--diff") + 1)
		secondPath := p.GetPosition(p.GetIndexOf("--diff") + 2)

		if !p.IsFolderExists(firstPath) || !p.IsFolderExists(secondPath) {
			puts("Not found folders for analysis!")
			puts("Usage:")
			puts("    Help: ckp --help")
			os.Exit(2)
		}
		if p.Has("--ignore") {
			ignore := p.GetPosition(p.GetIndexOf("--ignore") + 1)
			if ignore == "" {
				puts("Not found parameters from --ignore!")
				puts("Usage:")
				puts("    Help: ckp --help")
				os.Exit(2)
			}
			if strings.Contains(ignore, "--") {
				puts("Be careful, this may not work.")
				puts("--ignore ", ignore)
				puts("Usage:")
				puts("    Help: ckp --help")
			}
			split := strings.Split(ignore, ",")
			for i := range split {
				removeSpace := strings.Trim(split[i], " ")
				ignoreFolders = append(ignoreFolders, removeSpace)
			}
		}
		// initiate read directories
		pwd, err := os.Getwd()
		if err != nil {
			panic(err)
		}
		firstPathFQDN := pwd + "/" + firstPath

		p.Export(nameDirDiffs)
		p.FilterFile(firstPathFQDN, firstPath, secondPath)
		p.ReadRecursiveDir(firstPathFQDN, firstPath, secondPath)
		p.ResultDisplay()
	}
}

// Check é o pararâmetro para analisar dependências de um determinado alvo, no caso
// algum diretório que é especificado pelo usuário, acessado recursivamente
// cada arquivo/diretório e resgatado as entradas de "require" e "include" que
// caracterizam uma dependência do programa.
// a opção --dep-map faz uma analise profunda de todos os arquivos que usam essas dependências
func (p *Params) Check() {
	if p.Has("--check") {
		positionDiff := p.GetIndexOf("--check") + 1
		if (p.Count() - positionDiff) < 2 {
			puts("Not found folders for analysis!")
			puts("Usage:")
			puts("    Help: ckp --help")
			os.Exit(2)
		}

		p.Path = p.GetPosition(p.GetIndexOf("--check") + 1)

		if !p.IsFolderExists(p.Path) {
			puts("Not found folders for analysis!")
			puts("Usage:")
			puts("    Help: ckp --help")
			os.Exit(2)
		}
		// initiate read directories
		pwd, err := os.Getwd()
		if err != nil {
			panic(err)
		}
		pathFQDN := pwd + "/" + p.Path

		p.FilterFileCheck(pathFQDN, p.Path)

		// quando o parâmetro --dep-map é utilizado que será responsável por
		// executar o preview dos resultados finais será a última função a ser executada
		// Deep()
		if !p.Has("--dep-map") {
			p.ResultDisplay()
		}
		params.MapDeps(p.Path, pathFQDN)
	}
}

// Deep inicializa uma analise profunda das dependências
// em todos os arquivos do alvo identificando onde as dependências estão sendo
// usadas em demais rotinas do sistema ou arquivos.
func (p *Params) MapDeps(path, pathFQDN string) {
	if p.Has("--dep-map") {
		br := generateSpaces(" ")
		empty.Println(resultPrint(br))

		newtext := generateSpaces(" Initiated dependency map build, waiting ...")
		info.Println(infoPrint(newtext))
		p.InitiateDeepReport(path, pathFQDN)
	}
}

// FilterFile
func (p *Params) FilterFile(directory, dirComFirst, dirComSecond string) {
	if p.Has("--filter-file") {
		filter := p.GetPosition(p.GetIndexOf("--filter-file") + 1)
		if filter == "" {
			puts("Not found parameters from --filter-file!")
			puts("Usage:")
			puts("    Help: ckp --help")
			os.Exit(2)
		}
		p.ReadListFiles(directory, dirComFirst, dirComSecond)

		p.ResultDisplay()
		os.Exit(2)
	}
}

// Export() gera o diretório onde será exportado
// os arquivos de log quando usado esse parâmetro --export
func (p *Params) Export(name string) {
	if p.Has("--export") {
		pwd, err := os.Getwd()
		if err != nil {
			panic(err)
		}
		err = os.MkdirAll(pwd+"/"+name, 0755)
		if err != nil {
			panic(err)
		}
	}
}

// ResultDisplay print the information collected when finished the process
func (p *Params) ResultDisplay() {
	// scan result
	for j := 0; j < 2; j++ {
		line := generateSpaces(" ")
		empty.Println(resultPrint(line))
	}

	if len(brokenDependencyLogger) != 0 {
		p.WriteLog("broken_dependencies.log", brokenDependencyLogger)
		brokenDependencies := generateSpaces("Broken dependencies: " + strconv.Itoa(len(brokenDependencyLogger)))
		empty.Println(resultPrint(brokenDependencies))

		filesOpened1 := generateSpaces("Files opened: " + strconv.Itoa(len(files)))
		empty.Println(resultPrint(filesOpened1))
	}

	if len(dependencyLogger) != 0 {
		p.WriteLog("dependencies.log", dependencyLogger)
	}

	if len(dependencyMapLogger) != 0 {
		p.WriteLog("dependencies_map.log", dependencyMapLogger)
	}

	if len(filesDiffers) != 0 {
		p.WriteLog("differences.log", filesDiffers)
		if params.Has("--export") {
			folderExported := generateSpaces("Folder exported: " + nameDirDiffs)
			empty.Println(resultPrint(folderExported))
		}
		filesDiffers := generateSpaces("Files differs: " + strconv.Itoa(len(filesDiffers)))
		empty.Println(resultPrint(filesDiffers))
	}

	if len(directoryLogger) != 0 {
		directoriesScanned := generateSpaces("Directories scanned: " + strconv.Itoa(len(directoryLogger)))
		empty.Println(resultPrint(directoriesScanned))
	}

	if len(dependencyMapLogger) != 0 {
		filesDeep := generateSpaces("Dependency map: " + strconv.Itoa(len(dependencyMapLogger)))
		empty.Println(resultPrint(filesDeep))
	}

	footer := generateSpaces(" ")
	empty.Println(resultPrint(footer))
}

// IsFolderExists return bool true if the folder exists and false if not
func (p *Params) IsFolderExists(d string) bool {
	_, err := os.Stat(d)
	if err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}

// InitiateDeepReport faz a leitura dos arquivos que foram gerados por ReadFileDependencieCheck()
// armazenados os arquivos que são dependências no slice dependencyLogger no momento da execução,
// cada arquivo desse é checado em todos os arquivos do sistema alvo para identificar se ele está sendo
// usado em algum trecho do sistema, ignorando somente os arquivos que já estão na lista de filtros
// passada como parâmetro --filter-file e os que estão na lista de exclusão que é passada para o parâmetro
// --excluded-file, devido muitos dos arquivos que são dependências as vezes não serem importantes para a
// pesquisa, ou são arquivos de conexão, sessão ou semelhantes.
func (p *Params) InitiateDeepReport(path, pathFQDN string) {
	// dependencyLogger possui todos os arquivos encontrados como dependências
	// dentro de cada arquivo analisado a partir da lista passada como parâmetro
	// ckp --check directory --filter-file ecidademarica/pessoal.out
	for i := range dependencyLogger {
		// inArray retorna "true" quando o arquivo estiver na lista passada,
		// p.ExcludedFiles é um slice que contém a lista passada pelo usuário
		// com os arquivos que devem ser ignorados.
		exists, _ := inArray(dependencyLogger[i], p.ExcludedFiles)
		if !exists {
			if p.Has("--verbose") {
				newtext := generateSpaces(" " + dependencyLogger[i])
				scanning.Println(scanningPrint(newtext))
			}
			p.ReadRecursiveDir(pathFQDN, path, dependencyLogger[i])
		}
	}
	p.ResultDisplay()
}

// FilterFileCheck
func (p *Params) FilterFileCheck(directory, dirComFirst string) {
	if p.Has("--filter-file") {
		filter := p.GetPosition(p.GetIndexOf("--filter-file") + 1)
		if filter == "" {
			puts("Not found parameters from --filter-file!")
			puts("Usage:")
			puts("    Help: ckp --help")
			os.Exit(2)
		}
		p.ReadListFilesCheck(directory, dirComFirst)
	}
}

// ReadListFilesCheck
func (p *Params) ReadListFilesCheck(directory, dirComFirst string) {
	filterFile := p.GetPosition(p.GetIndexOf("--filter-file") + 1)
	file, err := os.Open(filterFile)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	fs := bufio.NewScanner(file)
	for fs.Scan() {
		p.ReadFileDependencieCheck(fs.Text(), directory, dirComFirst,"", false)
	}

	if err := fs.Err(); err != nil {
		panic(err)
	}
}

// ReadFileDependencieCheck
func (p *Params) ReadFileDependencieCheck(file, directory, dirComFirst, anterior string, signal bool) {
	pathFile := p.Path + "/" + file
	if signal {
		pathFile = file
	}

	nFile, err := os.Open(pathFile)
	if err == nil {
		if len(filesDep) == 0 {
			filesDep = append(filesDep, pathFile)
		}
		checkScann := registerFileDep(pathFile)
		if !checkScann {
			newtext := generateSpaces(" " + pathFile)

			if p.Has("--verbose") {
				scanning.Println(scanningPrint(newtext))
			}

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
						registerDeep(split[1])
						newtext = generateSpaces(" [ require ] found: " + split[1] + " in file -> " + pathFile)

						if p.Has("--verbose") {
							found.Println(foundPrint(newtext))
						}
						if strings.Contains(split[1], ".php") {
							p.ReadFileDependencieCheck(split[1], pathFile, directory, dirComFirst, false)
						}
					}
				}

				indexInclude := strings.Index(text, "include") // include or include_once
				if indexInclude != -1 {
					split := strings.Split(text, "\"")
					if len(split) == 3 {
						registerDeep(split[1])
						newtext = generateSpaces(" [ include ] found: " + split[1] + " in file -> " + pathFile)

						if p.Has("--verbose") {
							found.Println(foundPrint(newtext))
						}
						if strings.Contains(split[1], ".php") {
							p.ReadFileDependencieCheck(split[1], pathFile, directory, dirComFirst, false)
						}
					}
				}
			}
		}
	}
	return
}

// ReadListFiles
func (p *Params) ReadListFiles(directory, dirComFirst, dirComSecond string) {
	filterFile := p.GetPosition(p.GetIndexOf("--filter-file") + 1)
	file, err := os.Open(filterFile)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	fs := bufio.NewScanner(file)
	for fs.Scan() {
		p.CompareBetweenTwoFiles(p.OpenTwoFiles(directory + "/" + fs.Text(), dirComFirst, dirComSecond))
	}

	if err := fs.Err(); err != nil {
		panic(err)
	}
}

// WriteLog
func (p *Params) WriteLog(fileToWrite string, data []string) error {
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

// ReadDir
func (p *Params) ReadDir(directory string, signal bool, extension string) {
	dirs, err := ioutil.ReadDir(directory)
	if err != nil {
		panic(err)
	}
	for _, file := range dirs {
		// check for files with extension .php
		if strings.Contains(file.Name(), "."+extension) {
			filename := file.Name()
			if signal {
				filename = directory + "/" + file.Name()
			}
			p.ReadFileDependencie(filename, "null", signal)
		} else if file.IsDir() {
			if len(directoryLogger) == 0 {
				directoryLogger = append(directoryLogger, directory+"/"+file.Name())
			}
			registerDirectory(directory + "/" + file.Name())
			p.ReadDir(directory+"/"+file.Name(), true, extension)
		}
	}
	return
}

// readRecursiveDir
func (p *Params) ReadRecursiveDir(directory, dirComFirst, dirComSecond string) {
	dirs, err := ioutil.ReadDir(directory)
	if err != nil {
		panic(err)
	}
	for _, file := range dirs {
		fileOrDirName := directory + "/" + file.Name()
		if file.IsDir() {
			// Check if file aren't in the []ignoreFolders
			ignore, _ := inArray(file.Name(), ignoreFolders)
			if !ignore {
				registerDirectory(fileOrDirName)
				p.ReadRecursiveDir(fileOrDirName, dirComFirst, dirComSecond)
			}
			// if is continue to another record
			continue
		}
		if p.Has("--diff") {
			// if is not a folder put on into  -> openTwoFiles -> compareBetweenTwoFiles
			p.CompareBetweenTwoFiles(p.OpenTwoFiles(fileOrDirName, dirComFirst, dirComSecond))
		}
		if p.Has("--dep-map") {
			// Nenhum arquivo de dependência pode estar na lista de filtro
			// porque qualquer arquivo de dependência já veio provido dessa lista de filtros.
			exists, _ := inArray(file.Name(), p.FilterFiles)
			if !exists {
				// Faz uma leitura em cada arquivo de cada pasta
				scanned := p.ScanFile(fileOrDirName)
				// Faz uma busca do arquivo procurado no texto retornado por scanFile()
				exist := p.SearchOnScanned(scanned, dirComSecond)
				if exist {
					// quando encontrado registra a ocorrência no slice dependencyMapLogger
					registerDependencyMap(fileOrDirName)
				}
			}
		}
	}
}

// SearchOnScanned
func (p *Params) SearchOnScanned(data []string, search string) bool {
	for i := range data {
		indexRequire := strings.Index(data[i], search)
		if indexRequire != -1 {
			return true
		}
	}
	return false
}

// ScanFile
func (p *Params) ScanFile(file string) []string {
	nFile, err := os.Open(file)
	if err != nil {
		panic(err)
	}
	defer nFile.Close()
	var tx []string

	scanner := bufio.NewScanner(nFile)
	scanner.Split(bufio.ScanLines)

	for scanner.Scan() {
		tx = append(tx, scanner.Text())
	}
	return tx
}

// OpenTwoFiles
func (p *Params) OpenTwoFiles(file, dirComFirst, dirComSecond string) ([]byte, []byte, string, string) {
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
				return []byte{}, []byte{}, "", ""
			}
		}
		dt2, err := ioutil.ReadFile(fileToCompare)
		if err != nil {
			log.Fatal(err)
		}
		if p.Has("--verbose") {
			newtext := generateSpaces(" " + file)
			scanning.Println(scanningPrint(newtext))
		}
		return dt1, dt2, fileToCompare, file
	}
	return []byte{}, []byte{}, "", ""
}

// GenerateDiffFiles generate diffs files into the file system on the folder "diffs"
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

// CompareBetweenTwoFiles compare two files when the option --export is passed to program
func (p *Params) CompareBetweenTwoFiles(b1, b2 []byte, file, fileToCompare string) {
	result := bytes.Compare(b1, b2)
	if result != 0 && file != "" {
		registerDiffer(file)
		if params.Has("--export") {
			params.GenerateDiffFiles(fileToCompare, file)
		}
	}
}

// GenerateLog
func (p *Params) GenerateLog(dependencia, fileorigem string) {
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
	if len(brokenDependencyLogger) == 0 {
		brokenDependencyLogger = append(brokenDependencyLogger, text)
	}
	brokenDependencyLogger = p.RegisterLog(text, brokenDependencyLogger)
}

func (p *Params) ReadFileDependencie(file, anterior string, signal bool) {
	pathFile := p.Path + "/" + file

	if signal {
		pathFile = file
	}

	nFile, err := os.Open(pathFile)
	if err != nil {
		p.GenerateLog(pathFile, anterior)
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
							p.ReadFileDependencie(split[1], pathFile, false)
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
							p.ReadFileDependencie(split[1], pathFile, false)
						}
					}
				}
			}
		}
	}
	return
}

// RegisterLog
func (p *Params) RegisterLog(text string, brokenDependencyLogger []string) []string {
	// check if logger exists
	exists, _ := inArray(text, brokenDependencyLogger)
	if !exists {
		// @todo improvement for scanning this routes ../
		index := strings.Index(text, "../")
		if index == -1 {
			brokenDependencyLogger = append(brokenDependencyLogger, text)
		}
	}
	return brokenDependencyLogger
}

func registerDiffer(name string) bool {
	exists, _ := inArray(name, filesDiffers)
	if !exists {
		filesDiffers = append(filesDiffers, name)
	}
	return exists
}

func registerDeep(name string) bool {
	exists, _ := inArray(name, dependencyLogger)
	if !exists {
		dependencyLogger = append(dependencyLogger, name)
	}
	return exists
}

func registerDependencyMap(name string) bool {
	exists, _ := inArray(name, dependencyMapLogger)
	if !exists {
		dependencyMapLogger = append(dependencyMapLogger, name)
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

func registerFileDep(name string) bool {
	exists, _ := inArray(name, filesDep)
	if !exists {
		filesDep = append(filesDep, name)
	}
	return exists
}

func registerDirectory(name string) {
	exists, _ := inArray(name, directoryLogger)
	if !exists {
		directoryLogger = append(directoryLogger, name)
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
	data := []byte(str)
	for i := 0; i < length; i++ {
		data = append(data, '\u0020')
	}
	return string(data)
}

func getWidth() int {
	ws := &Winsize{}
	retCode, _, err := syscall.Syscall(syscall.SYS_IOCTL,
		uintptr(syscall.Stdin),
		uintptr(syscall.TIOCGWINSZ),
		uintptr(unsafe.Pointer(ws)))

	if int(retCode) == -1 {
		panic(err)
	}
	return int(ws.Col)
}