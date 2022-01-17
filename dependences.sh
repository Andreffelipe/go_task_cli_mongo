#!/bin/bash
GIT_NAME="meu-git"

echo "inicializando o go.mod"
go mod init github.com/${GIT_NAME}/go_task_cli_mongo

echo "criando variável ambiente"
echo DATABASE_URL= > .env

echo "instalando pacote do mongo"
go get go.mongodb.org/mongo-driver

echo "instalando pacote para cli"
go get github.com/urfave/cli/v2

echo "instalando pacote de cor"
go get gopkg.in/gookit/color.v1