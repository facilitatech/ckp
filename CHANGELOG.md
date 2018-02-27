# CHANGELOG

## v0.0.10

_Release: 2018-02-27_

- Refactoring ReadFileDependencie() function
- Removed FilterFileCheck(), is not used more
- Removed ReadListFiles() and ReadListFilesCheck()
- Refactoring FilterFile() to CheckFilterFile()


## v0.0.9

_Release: 2018-02-25_

- Fixed bug on SetFilesParams.
- Included new install command to travis.
- Generated build to darwin too.
- Fixed bug -> Many Open Files
- Changed build.sh to work with binary to linux.
- Change name of main.go to ckp.go.
- Included new feature --filter-file


## v0.0.8

_Release: 2018-02-08_

- Removed condition of count the number of params
   - Only work when the param is passed, ignore the total
     of parameters passed
- Removed function register() and improved the docs
- Fixed the error on the title of docker
- Changed the instruction and improved the README
- Using gofmt for format the code
- Create rules and validations for some parameters
   - Rules for --check-dependencies and --ignore
   - The program only beginning when the parameters passed:
     is not empty, the name/path of the folder is valid, if the
     parameters passed has the options: --check-dependencies or --ignore
- Fixed name of LICENSE
- Changed option --check to --check-dependencies
- Implementing new feature --diff
   - Option --diff receive three parameters
     target1 target2 and one option --ignore, that's possible
     ignore some folder
   - Changed the log report, generated the log file with the affected
   - Changed the function display to receive more fields.
- Ignore txt files
- Created readRecursiveDir function
   - Created function "readRecursiveDir()"
     This function analyzes the entire folder of the set of
     the parameters and returns execute scan recursively
- Changed package name and descriptions.
- Changed description of the project
- Ignored binary files.
- Changed the name of package for ckp
- Changed the build directory
- Generated two binary files, darwin and linux
- Change directory of main file and README name package


## v0.0.7

_Release: 2017-12-02_

- Criado licence para o projeto
- Alterado versão do docker-compose para a 3
- Refatorado o Dockerfile
   - Enviado as dependências para o GOROOT onde ficam isoladas
     do diretório onde é feito o volume
   - Alterado versão do Go para 1.9
- Alterado nome da empresa que mantém projeto
- Ignorado arquivos e atualizado README


## v0.0.6

_Release: 2017-03-21_

- Refatoração e implementação de novas features

    - registerFile(name string) bool, retornando um
      booleano, usado na func de readFile() para verificar
      os arquivos que já foram escaneados, sendo ignorados esses
      arquivos em uma segunda leitura
    - Criado func registerLog para separar trecho de código da func
      generateLog, passado o nome do arquivo e o próprio slice que
      armazena os nomes dos arquivos escaneados
    - Criado func resultDisplay que monta o resultado dos arquivos e
      pastas afetadas pelo escaneamento e printa os arquivos
      que possuem dependências quebradas
    - Criado func writeLog que escreve as dependências quebradas
      em um arquivo txt


## v0.0.5

_Release: 2017-03-20_

- Correção de inconsistência no search das inclusões nos
  arquivos .php


## v0.0.4

_Release: 2017-03-17_

- Refatoração, novas features, alteração na estrutura

    - Alterado layout, removido log gerando mais
      espaço para o nome dos arquivos grandes
    - Criado um preview do resultado da execução para mostrar
      os arquivos afetados, diretórios que foram escaneados
      e total de inclusões quebradas
    - Criado registrador para os arquivos abertos e pastas
      escaneadas gerando o total no final da execução
    - Criado condições para checar se os arquivos ou pastas
      que já foram analisados não entrem novamente no slice
      de registro
    - Inserido comentários @todo indicando onde deve ser melhorado
      o código, implementação de escaneamento para
      arquivos que usam namespaces [ use ] ou
      arquivos que usam inclusão fora do escopo, retornando
      diretórios [ ../../../ ]


## v0.0.3

_Release: 2017-03-16_

- Correção de erro index out of range


## v0.0.2

_Release: 2017-03-16_

- Adicionado na hora do build nova dependência

    - github.com/agtorre/gocolorize para geração dos logs com opção de cores no stdout
    - Alterado o ponto de montagem do volume no docker-compose,
      enviando somente o arquivo main.go -> ./src/app/:/go/src/app


## v0.0.1

_Release: 2017-03-16_

- First commit

    - Adicionado estrutura inicial
    - Criado provisionamento com docker para o ambiente
    - Adicionado README e CHANGELOG
    - Adicionado arquivo init.sh para automatização de tarefas