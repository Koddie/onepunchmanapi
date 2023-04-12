# :punch: One Punch Man API :punch:

`Servidor DB - MySql`

## Configuração MySql

1. Executar o comando: `create database onepunchman;`

2. Executar os comandos do arquivo `create-tables.sql`

3. Exportar a variável de senha do seu servidor MySql para a variável do sistema (usar o mesmo terminal onde se executa o comando "go run main.go"):
    - Linux, Mac: `$ export DBPASS=password`
    - Windows: `set DBPASS=password`

4. Ao executar o main.go, a mensagem `"Conectado ao Banco de Dados!"` deverá aparecer no terminal


## Executando

Navegue até a raiz do projeto no Terminal ou CMD e execute o comando: `go run main.go`

## End Points

### Listar Heróis
`localhost:8080/heroi (GET)`
```
Body do json do request:
{
    [“nome”: string]
    [“classe”: “A”|”B”|”C”|”S”]
    [“ranking”: integer"]
} 
```
```
retornos: 
Código 200 (OK)
Body: [{
“id”: integer
“nome”: string
“classe”: “A”|”B”|”C”|”S”
“ranking”: integer
},
...
{
“id”: integer
“nome”: string
“classe”: “A”|”B”|”C”|”S”
“ranking”: integer
}]

404 (Não encontrado) Body: {“message”: “herói não encontrado”}
``` 
### Listar Herói
`localhost:8080/heroi/{id} (GET)`

```
retornos: 
Código 200 (OK)
Body: {
“id”: integer
“nome”: string
“classe”: “A”|”B”|”C”|”S”
“ranking”: integer
}

404 (Não encontrado) Body: {“message”: “herói não encontrado”}
``` 
### Cadastro de novo Herói
`localhost:8080/heroi/novo (PUT)`
```
Body do json do request:
{
    “nome”: string
    “classe”: “A”|”B”|”C”|”S”
    “ranking”: integer 
} 
```
```
Retornos: 
200 (OK) Body: {id: integer}
409 (Conflict) Body: {“message”: “herói já cadastrado”}
```
### Alterar Herói
`localhost:8080/heroi/{id} (POST)`
```
Body do json do request:
{
    [“classe”: “A”|”B”|”C”|”S”]
    [“ranking”: integer ]
} 
```
```
retorno:
204 (No Content)
```
### Excluir Herói
`localhost:8080/heroi/{id} (DELETE)`
```
retorno:
204 (No Content)
```
### Erros Diversos
```
Retorno:
500 (Internal Server Error) Body: {"Mensagem..."}
```
