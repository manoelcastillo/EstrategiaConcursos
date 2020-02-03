package handler

import (
	"EstrategiaConcursos/dao"
	"EstrategiaConcursos/model"
	"encoding/json"
	"io/ioutil"
	"log"
	"fmt"
	"net/http"
	"strconv"
	"time"
	"reflect"

	"github.com/gorilla/mux"
)

const STARRED_URI = "https://api.github.com/users/%s/starred"

type Handler struct {

}

var RepoDao dao.RepoDAO
func init() {
	RepoDao = dao.RepoDAO{}
}

/*
	Primeiro tenta pegar os repositórios de determinado usuário no MongoBD, caso esteja vazia a collection, pega da 
	API do GitHub e insere na collection.
*/
func (h *Handler) GetAllRepositories(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(405) //method not allowed
		return
	}

	gitHubUser := getRequestUser(r)

	if gitHubUser == "" {
		w.WriteHeader(400) //bad request
		return 
	}

	var reposData []*model.Repository = RepoDao.GetAllRepositories(gitHubUser)
	
	if reposData == nil {
		log.Println("[Info] Fetching from Github ", gitHubUser)
		var err error
		reposData, err = h.FetchGithub(gitHubUser)

		if err != nil {
			w.WriteHeader(500) //error on fetching from github
			return
		}

		if len(reposData) > 0 {
			RepoDao.InsertRepositories(gitHubUser, reposData)
		}
	}

	err := json.NewEncoder(w).Encode(reposData)
	if err != nil {
		log.Println("[Error] ", err)
	}
}

/*
	Função para atualizar os repositórios de determinado usuário. Remove os des-etrelados e insere os estrelados
*/
func (h *Handler) UpdateRepositories(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		w.WriteHeader(405) //method not allowed
		return
	}

	gitHubUser := getRequestUser(r)	
	if gitHubUser == "" {
		w.WriteHeader(400) //bad request
		return
	}
	
	dataFromGithub, err := h.FetchGithub(gitHubUser)
	if err != nil {
		w.WriteHeader(400)
		return
	}

	var responseCode = 202

	dataFromMongo := RepoDao.GetAllRepositories(gitHubUser)
	if len(dataFromGithub) > 0 {
		if dataFromMongo == nil {
			//something wrong is not right
			json.NewEncoder(w).Encode(dataFromGithub)
			return;
		}
		var toRemove []*model.Repository
		var toInsert []*model.Repository
		//checking repositories unstarred
		for _, dbRep := range dataFromMongo {
			var found = false
			for _, ghRep := range dataFromGithub {
				if dbRep.Id == ghRep.Id {
					found = true
					break
				}
			}
			if !found {
				toRemove = append(toRemove, dbRep)
			}
		}
		//checking new repositories starred
		for _, ghRep := range dataFromGithub {
			var found = false
			for _, dbRep := range dataFromMongo {
				if ghRep.Id == dbRep.Id {
					found = true
					break
				}
			}
			if !found {
				toInsert = append(toInsert, ghRep)
			}
		}

		if len(toRemove) > 0 || len(toInsert) > 0 {
			responseCode = 200
		}

		RepoDao.UpdateRepositoryCollection(gitHubUser, toInsert, toRemove)
	}
	
	w.WriteHeader(responseCode) // 202 accepted (nada alterado), 200 ok (houve update)

}

func (h *Handler) FetchGithub(user string) ([]*model.Repository, error) {
	url := fmt.Sprintf(STARRED_URI, user)
	client := http.Client{
		Timeout: time.Second * 2,
	}

	githubAPIReq, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	githubAPIReq.Header.Set("User-Agent", "Golang")
	res, getErr := client.Do(githubAPIReq)
	if getErr != nil {
		return nil, getErr
	}

	body, readErr := ioutil.ReadAll(res.Body)
	if readErr != nil {
		log.Println("[Error] body ", readErr)
		return nil, readErr
	}

	repos := make([]*model.Repository, 0)
	jsonErr := json.Unmarshal(body, &repos)
	if jsonErr != nil {
		log.Println("[Error] unmarshal ", jsonErr)
		return nil, jsonErr
	}

	return repos, nil

}

func (h *Handler) GetRepository(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(405) //method not allowed
		return
	}

	gitHubUser, repoId, err := getUserAndRepoId(r)
	if err != nil || gitHubUser == "" {
		w.WriteHeader(400) //bad request
		return
	}
	
	repo := RepoDao.GetRepository(gitHubUser, repoId)
	if reflect.DeepEqual(model.Repository{}, repo) {
		w.WriteHeader(204) //no content
		return
	}
	
	json.NewEncoder(w).Encode(repo)
}

func (h *Handler) GetRepositoryTags(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(405) //method not allowed
		return
	}

	gitHubUser, repoId, err := getUserAndRepoId(r)
	if err != nil || gitHubUser == "" {
		w.WriteHeader(400) //bad request
		return
	}

	tags := RepoDao.GetRepositoryTags(gitHubUser, repoId)
	json.NewEncoder(w).Encode(tags)
}

func (h *Handler) TagAdd(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(405) //method not allowed
		return
	}

	gitHubUser, repoId, err1 := getUserAndRepoId(r)
	t, err2 := getTagFromBody(r)
	if err1 != nil || err2 != nil || gitHubUser == "" {
		w.WriteHeader(400) //bad request
		return
	}
	repository := RepoDao.AddTag(gitHubUser, repoId, t)
	
	w.WriteHeader(202)
	json.NewEncoder(w).Encode(repository)
}

func (h *Handler) TagSearch(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(405) //method not allowed
		return
	}

	gitHubUser := getRequestUser(r)
	t := mux.Vars(r)["tag"]
	if t == "" || gitHubUser == "" {
		w.WriteHeader(400) //bad request
		return
	}

	repos := RepoDao.SearchByTag(gitHubUser, model.Tag {Name: t})
	json.NewEncoder(w).Encode(repos)
}

func (h *Handler) TagAdvice(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(405) //method not allowed
		return
	}

	gitHubUser, repoId, err := getUserAndRepoId(r)
	if err != nil || gitHubUser == "" {
		w.WriteHeader(400) //bad request
		return
	}

	repo := RepoDao.GetRepository(gitHubUser, repoId)
	possibleTags := GetPossibleTags(repo)
	json.NewEncoder(w).Encode(possibleTags)
}

func (h *Handler) TagDelete(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		w.WriteHeader(405) //method not allowed
		return
	}

	gitHubUser, repoId, err := getUserAndRepoId(r)
	t := mux.Vars(r)["tag"]
	if err != nil || t == "" || gitHubUser == "" {
		w.WriteHeader(400) //bad request
		return
	}

	repo := RepoDao.DeleteTag(gitHubUser, repoId, model.Tag {Name: t})
	json.NewEncoder(w).Encode(repo)
}

func (h *Handler) TagUpdate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		w.WriteHeader(405) //method not allowed
		return
	}

	gitHubUser, repoId, err1 := getUserAndRepoId(r)
	if err1 != nil || gitHubUser == "" {
		w.WriteHeader(400) //bad request
		return
	}

	var tagUpdate model.TagUpdate
	err2 := json.NewDecoder(r.Body).Decode(&tagUpdate)
	if err2 != nil {
		log.Fatal("Error decoding body", err2)
	}
	
	repo, err := RepoDao.UpdateTag(gitHubUser, repoId, tagUpdate)
	if err != nil {
		w.WriteHeader(400) //bad request, tag not in repository
		return
	}	
	json.NewEncoder(w).Encode(repo)
}

func getRequestUser(r *http.Request) string {
	return mux.Vars(r)["user"]
}

func getRequestRepoId(r *http.Request) (uint64, error) {
	return strconv.ParseUint(mux.Vars(r)["id"], 10, 64)
}

func getUserAndRepoId(r *http.Request) (string, uint64, error) {
	user := getRequestUser(r)
	repoId, err := getRequestRepoId(r)
	return user, repoId, err
}

func getTagFromBody(r *http.Request) (model.Tag, error) {
	var t model.Tag
	err := json.NewDecoder(r.Body).Decode(&t)
	
	return t, err
}
