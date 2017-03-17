# CHANGELOG

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