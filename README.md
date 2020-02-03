# Construindo uma API Rest para o desafio Estratégia Concursos

Link do enunciado do desafio: https://github.com/estrategiahq/desafio-desenvolvimento

A solução apresentada foi elaborada em Go usando [Gorila Mux](https://github.com/gorilla/mux) para rotear os endpoints e o [MongoDB](https://www.mongodb.com/) como banco de dados de documentos e seguinto a metodologia TDD ([test-driven methodology](https://pt.wikipedia.org/wiki/Test-driven_development)). Foi pensada para ser fault tolerant e atuar como uma solução de backend.

### Objetivos
* Aprender e ficar proficiente em Go
* Estudo e uso de um banco de dados NoSQL
* Aprender a suar MongoDB e Go implemnentando uma API Rest
* Cumprir os requisitos de avaliação técnica do processo seletivo

### Requisitos
* Binários Go instralados e configurados. 
    * Windows: [download](https://golang.org/dl/)
    * Linux: 
        ```console
        ubuntu@usr$ sudo add-apt-repository ppa:longsleep/golang-backports
        ubuntu@usr$ sudo apt-get update
        ubuntu@usr$ sudo apt-get install golang-go
        ```
* MongoDB instalado e rodando como serviço na porta default *27017*
    * Windows: [download](https://www.mongodb.com/download-center/community)
    * Linux: 
        ```console
        ubuntu@usr$ sudo apt update
        ubuntu@usr$ sudo apt install -y mongodb
        ```
### Dependências Go
Antes de rodar a aplicação, tendo todos os requisitos, devemos baixar as dependências do projeto
* gorillaz mux: 
    ```console
    $ go get github.com/gorilla/mux
    ```
* mongoDB driver:
    ```console
    $ go get go.mongodb.org/mongo-driver/mongo
    ```

### Ferramentas de desenvolvimento
A solução foi desenvolvida em ambiente Windows usando as seguintes ferramentas:
* VSCode: https://code.visualstudio.com/
    * Extensão [Go](https://github.com/Microsoft/vscode-go)
    * Extensão [openapi-lint](https://github.com/Mermade/openapi-lint-vscode)
* Mongo Compass: https://www.mongodb.com/products/compass
* PostMan: https://www.getpostman.com/

### OpenAPI
A API foi descrita usando a especificação [OpenAPI 3](http://spec.openapis.org/oas/v3.0.2) no formato [YAML](https://yaml.org/). 
O arquivo [openapi.yaml](openapi.yaml) pode ser validado com a ferramenta [Swagger](https://swagger.io) por exemplo.

### Rodando os testes
Os Testes foram feitos apenas para as funções handler, que servem os requests dos endpoints. Os testes são simples, levando apenas em consideração os código de retorno http.

```console
$ go test -timeout 30 EstrategiaConcursos/handler
```

### Rodando o BackEnd
Na pasta raiz do projeto:
```console
$ go build
$ ./EstrategiaConcursos
```
Para testar basta fazer as requests para os endpoints em http://localhost:8000 via brower/curl ou com a ferramenta **PostMan** 




