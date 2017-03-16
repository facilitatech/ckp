# CHANGELOG

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