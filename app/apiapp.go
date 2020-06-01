package app

import (
	"GIT-Stars/handler"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

const SERVER_PORT = ":8000"

type API struct {
	Router *mux.Router
	Handler *handler.Handler
}

func (a *API) Initialize() {
	a.Router = mux.NewRouter()
	a.InitializeRoutes()
	a.server()
}

func (a *API) InitializeRoutes() {
	/*
		Retorna todos os repositórios do MongDB ou da API GitHub, nessa ordem, de determinado user
	*/
	a.Router.HandleFunc("/starred/{user}", a.Handler.GetAllRepositories).Methods("GET")

	/*
		Update dos repositórios starred para determinado user. Grava no MongoDB os repositórios novos e
		deleta os "unstarred"
	*/
	a.Router.HandleFunc("/starred/{user}", a.Handler.UpdateRepositories).Methods("PUT")

	/*
		Retorna um repósitório {id} de determinado usuario {user} do MongoDB
	*/
	a.Router.HandleFunc("/starred/{user}/{id}", a.Handler.GetRepository).Methods("GET")

	/*
		Procura os repositório de {user} que contenham {tag}
	*/
	a.Router.HandleFunc("/search/{user}/{tag}", a.Handler.TagSearch).Methods("GET")

	/*
		Sugetão de {tag} para determinado repositório {id} de determinado {user}
	*/
	a.Router.HandleFunc("/tagadvisor/{user}/{id}", a.Handler.TagAdvice).Methods("GET")

	/*
		Adiciona uma tag no repositório {id} de determinado {user}. {tag} vai no body
	*/
	a.Router.HandleFunc("/tag/{user}/{id}", a.Handler.TagAdd).Methods("POST")

	/*
		Update de uma tag de um repositório {id} de um determinado {user}. {oldtag} e {tag} vai no body
	*/
	a.Router.HandleFunc("/tag/{user}/{id}", a.Handler.TagUpdate).Methods("PUT")

	/*
		Deleta uma tag de um repositório {id} de um dado {user}
	*/
	a.Router.HandleFunc("/tag/{user}/{id}/{tag}", a.Handler.TagDelete).Methods("DELETE")

	/*
		pega todas as {tags} de um repositório {id} de um dado {user}
	*/
	a.Router.HandleFunc("/tag/{user}/{id}", a.Handler.GetRepositoryTags).Methods("GET")
}

func (a *API) server() {
	log.Fatal(http.ListenAndServe(SERVER_PORT, a.Router))
}
