One Punch Man API

![One Punch Man API](https://encrypted-tbn0.gstatic.com/images?q=tbn:ANd9GcTsELxXWj1G2RJo8kTgc0Icu0lG6EpVc7QI7Q&usqp=CAU)

Servidor DB - MySql

------ Configuração MySql -------

1 - Executar o comando: create database onepunchman;

2 - Executar os comandos do arquivo create-tables.sql

3 - Exportar a variável de senha do seu servidor MySql para a variável do sistema (usar o mesmo terminal onde se executa o comando "go run main.go"):
    - Linux, Mac:
        $ export DBPASS=password
    - Windows:
        set DBPASS=password

4 - Ao executar o main.go, a mensagem "Conectado ao Banco de Dados!" deverá aparecer no terminal


------ Executando -------

Endereço: localhost:8080/heroi/novo

body json {

    “nome”: string

    “classe”: “A”|”B”|”C”|”S”
    
    “ranking”: integer
    
}



