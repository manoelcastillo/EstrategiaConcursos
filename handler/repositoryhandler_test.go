package handler

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
)

var router *mux.Router
var handler Handler
func init() {
	router = mux.NewRouter()
	handler = Handler{}
}
/*
	Nesta suite de testes só estou testando os códigos de resposta da API, mais especificamente os códigos esperados (200 - OK)
	TODO: extender os testes para analisar os conteúdos recebidos e os códigos de respostas inesperados
*/

func TestGetAllRepositories(t *testing.T) {
	router.HandleFunc("/starred/{user}", handler.GetAllRepositories).Methods("GET")

	req, err := http.NewRequest("GET", "/starred/manoelcastillo", nil)
	if err != nil {
		t.Fatal(err)
	}

	response := request(req)
	checkResponseCode(t, 200, response.Code)

}

func TestUpdateRepositories(t *testing.T) {
	router.HandleFunc("/starred/{user}", handler.UpdateRepositories).Methods("PUT")

	req, err := http.NewRequest("PUT", "/starred/manoelcastillo", nil)
	if(err != nil) {
		t.Fatal(err)
	}

	response := request(req)
	checkResponseCode(t, 202, response.Code) //quando não houve alteração dos repositórios com estrela no gihub
}

func TestGetRepository(t *testing.T) {
	router.HandleFunc("/starred/{user}/{id}", handler.GetRepository).Methods("GET")

	req, err := http.NewRequest("GET", "/starred/manoelcastillo/101190593", nil)
	if(err != nil) {
		t.Fatal(err)
	}

	response := request(req)
	checkResponseCode(t, 200, response.Code)
}

func TestGetRepositoryTags(t *testing.T) {
	router.HandleFunc("/search/{user}/{tag}", handler.TagSearch).Methods("GET")

	req, err := http.NewRequest("GET", "/search/manoelcastillo/javascript", nil)
	if err != nil {
		t.Fatal(err)
	}

	response := request(req)
	checkResponseCode(t, 200, response.Code)

}

func TestTagAdvice(t *testing.T) {
	router.HandleFunc("/tagadvisor/{user}/{id}", handler.TagAdvice).Methods("GET")

	req, err := http.NewRequest("GET", "/tagadvisor/manoelcastillo/101190593", nil)
	if err != nil {
		t.Fatal(err)
	}

	response := request(req)
	checkResponseCode(t, 200, response.Code)
}

func TestTagAdd(t *testing.T) {
	router.HandleFunc("/tag/{user}/{id}", handler.TagAdd).Methods("POST")

	payload := []byte(`{"name":"javascript"}`)
	req, err := http.NewRequest("POST", "/tag/manoelcastillo/101190593", bytes.NewBuffer(payload))
	if err != nil {
		t.Fatal(err)
	}

	response := request(req)
	checkResponseCode(t, 202, response.Code)
}

func TestTagUpdate(t *testing.T) {
	router.HandleFunc("/tag/{user}/{id}", handler.TagUpdate).Methods("PUT")

	payload := []byte(`{"oldtag":{"name":"javascript"},"newtag":{"name":"scriptjava"}}`)
	req, err := http.NewRequest("PUT", "/tag/manoelcastillo/101190593", bytes.NewBuffer(payload))
	if err != nil {
		t.Fatal(err)
	}

	response := request(req)
	checkResponseCode(t, 200, response.Code)
}

func TestTagDelete(t *testing.T) {
	router.HandleFunc("/tag/{user}/{id}/{tag}", handler.TagDelete).Methods("DELETE")

	req, err := http.NewRequest("DELETE", "/tag/manoelcastillo/101190593/scriptjava", nil)
	if err != nil {
		t.Fatal(err)
	}

	response := request(req)
	checkResponseCode(t, 200, response.Code)
}

func TestGetRespositoryTags(t *testing.T) {
	router.HandleFunc("/tag/{user}/{id}", handler.GetRepositoryTags).Methods("GET")

	req, err := http.NewRequest("GET", "/tag/manoelcastillo/101190593", nil)
	if err != nil {
		t.Fatal(err)
	}
	
	response := request(req)
	checkResponseCode(t, 200, response.Code)
}

func checkResponseCode(t *testing.T, expected int, received int) {
	if expected != received {
		t.Errorf("Código esperado [%d] recebido [%d]", expected, received)
	}
}

func request(req *http.Request) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)
	return rr
}
