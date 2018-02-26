# ckp - Check PHP files

[![Go Report Card](https://goreportcard.com/badge/github.com/facilitatech/ckp)](https://goreportcard.com/report/github.com/facilitatech/ckp) &nbsp;&nbsp; [![Build Status](https://travis-ci.org/facilitatech/ckp.svg?branch=master)](https://travis-ci.org/facilitatech/ckp)

### Install
```bash
go get github.com/facilitatech/ckp
```

### Options

```bash
--broken-deps    Check all broken dependencies of your project PHP has

--diff           Make diff between two folders recursively

--check          Check dependencies with two folders recursively

--filter-file    Filter files per list, work with --diff and --check

--ignore         Ignore folders

--export         Export diffs file into folder, this only work with --diff

--verbose        Print result

--dep-map        Dependency map, this only work with --check

--excluded-file  Ignore this files on the dependencies check, this only work with --check
```

### Examples

```bash
ckp --broken-deps /var/www/app --verbose

ckp --diff /var/www/app1 /var/www/app2 --verbose
ckp --diff /var/www/app1 /var/www/app2 --ignore vendor,bin --verbose
ckp --diff /var/www/app1 /var/www/app2 --ignore vendor --filter-file files.txt --verbose
ckp --diff /var/www/app1 /var/www/app2 --ignore vendor --filter-file files.txt --export --verbose

ckp --check /var/www/app --verbose
ckp --check /var/www/app --filter-file files.txt --verbose
ckp --check /var/www/app --filter-file files.txt --dep-map --verbose
ckp --check /var/www/app --filter-file files.txt --dep-map --excluded-file ignore-files.txt --verbose
```


### Development with docker
Clone the repository in folder do you prefer
```bash
cd /srv
git clone https://github.com/facilitatech/ckp
```

**Execute the file init.sh for up the docker containers**

The first step executes the option 3, again execute the file with the option 1 when the option 3 is finished!
```bash
https://github.com/facilitatech/ckp for the canonical source repository
Copyright (c) facilita.tech - 2016-2018
(http://facilita.tech)

  __            _ _ _ _         _            _
 / _| __ _  ___(_) (_) |_ __ _ | |_ ___  ___| |__
| |_ / _` |/ __| | | | __/ _` || __/ _ \/ __| '_ \
|  _| (_| | (__| | | | || (_| || ||  __/ (__| | | |
|_|  \__,_|\___|_|_|_|\__\__,_(_)__\___|\___|_| |_|

ckp

DOCKER
Generate new containers ? [ 1 ]
Delete all containers ?   [ 2 ]
Start new build ?         [ 3 ]
```

Example of how to return after executing option number `1`
```bash
Generating new containers ...
Name              Command               State    Ports
--------------------------------------------------------
ckp    reflex -c /var/exec/reflex ...   Up      6060/tcp
ckp is up-to-date
Name              Command               State    Ports
--------------------------------------------------------
ckp    reflex -c /var/exec/reflex ...   Up      6060/tcp
```

Preview the logs.
```bash
docker logs ckp -f
```