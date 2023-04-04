# :punch: One Punch Man API :punch:

`Servidor DB - MySql`

## Configuração MySq

1. Executar o comando: `create database onepunchman;`

2. Executar os comandos do arquivo `create-tables.sql`

3. Exportar a variável de senha do seu servidor MySql para a variável do sistema (usar o mesmo terminal onde se executa o comando "go run main.go"):
    - Linux, Mac: `$ export DBPASS=password`
    - Windows: `set DBPASS=password`

4. Ao executar o main.go, a mensagem `"Conectado ao Banco de Dados!"` deverá aparecer no terminal


## Executando

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

404 (Não encontrado)  {“message”: “herói não encontrado”}
``` 
### Cadastro de novo Herói
`localhost:8080/heroi/novo (POST)`
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
409 (Conflict) - Já Existente | Body: {“message”: “herói já cadastrado”}
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
