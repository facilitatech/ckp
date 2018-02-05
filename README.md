# ckp - Check PHP files

### Instalação

Clone o repositório na pasta de preferência
```bash
cd /srv
git clone https://github.com/facilitatech/dependency-check-php
```

**Executar o arquivo `init.sh` para iniciar os containers**
Execute a opção número `3` para efetuar o `build` logo em seguida executar novamente o arquivo
com a opção número `1`

```bash
https://github.com/facilitatech/dependency-check-php for the canonical source repository 
Copyright (c) facilita.tech - 2016-2017
(http://facilita.tech)  
 
  __            _ _ _ _         _            _     
 / _| __ _  ___(_) (_) |_ __ _ | |_ ___  ___| |__  
| |_ / _` |/ __| | | | __/ _` || __/ _ \/ __| '_ \ 
|  _| (_| | (__| | | | || (_| || ||  __/ (__| | | |
|_|  \__,_|\___|_|_|_|\__\__,_(_)__\___|\___|_| |_|
                                                   
dependency-check-php 

DOCKER
Generate new containers ? [ 1 ] 
Delete all containers ?   [ 2 ] 
Start new build ?         [ 3 ]
```

Exemplo de como deve ser o retorno depois da execução da opção número `1`

```bash
Generating new containers ...
Name              Command               State    Ports
--------------------------------------------------------
app    reflex -c /var/exec/reflex ...   Up      6060/tcp
app is up-to-date
Name              Command               State    Ports
--------------------------------------------------------
app    reflex -c /var/exec/reflex ...   Up      6060/tcp
```

**Acessando o container**

```bash
docker exec -it app bash
```

**Executando o binário**

somente executar o comando abaixo passando dois parâmetros para o programa;
o `--check` inicia o processo de analise dos arquivos e o segundo parâmetro
deve ser a pasta onde encontra-se os arquivos com extensão `.php`

```bash
app --check meus_arquivos
```