package dao

import (
	"EstrategiaConcursos/model"
	"EstrategiaConcursos/mongodb"
	"context"
	"log"
	"reflect"
	"errors"

	"go.mongodb.org/mongo-driver/bson"
)

type RepoDAO struct {
	
}

var DB mongodb.MongoDB
func init() {
	DB = mongodb.MongoDB{}
}

func (d *RepoDAO) InsertRepositories(user string, repositories []*model.Repository) {
	var rps []interface{}
	
	for _, r := range repositories {
		r.User = user
		rps = append(rps, r)
	}
	DB.InsertMany(rps)
}

func (d *RepoDAO) GetRepositoryTags(githubUser string, repoId uint64) []model.Tag {
	repository := d.GetRepository(githubUser, repoId)
	
	return repository.Tags
}

func (d *RepoDAO) GetAllRepositories(user string) []*model.Repository {
	log.Println("[Info] Getting all repositories from MongoDB")
	
	docs := DB.Find(bson.M{"user": user})
	var ret []*model.Repository
	
	for docs.Next(context.Background()) {
		var repo model.Repository
		err := docs.Decode(&repo)
		if err != nil {
			log.Fatal("Error on docoding repository", err)
		}
		ret = append(ret, &repo)
	}
	docs.Close(context.Background())
	
	return ret
}

func (d *RepoDAO) GetRepository(githubUser string, repoId uint64) model.Repository {
	log.Printf("[Info] Fetching [%d] from MongoDB", repoId)
	
	var repo model.Repository
	filter := bson.M{"id": repoId, "user": githubUser}
	docFromDB := DB.FindOne(filter)
	docFromDB.Decode(&repo)
	
	return repo
}

func (d *RepoDAO) AddTag(user string, repoId uint64, tag model.Tag) model.Repository {
	log.Printf("[Info] add tag [%s] to [%d] [%s]", tag, repoId, user)
	
	repo := d.GetRepository(user, repoId)
	
	if !reflect.DeepEqual(model.Repository{}, repo) {
		repo.Tags = append(repo.Tags, tag)
		filter := bson.M{"id": repoId, "user": user}
		update := bson.M{"$set": bson.M{"tags": repo.Tags}}
		DB.UpdateOne(filter, update)
		log.Println("[Info] Tag added to repository: ")
	}
	
	return repo
}

func (d *RepoDAO) SearchByTag(user string, tag model.Tag) []*model.Repository {
	log.Printf("[Info] search tag [%s] in [%s]", tag, user)
	
	arr := []model.Tag{tag}
	docs := DB.Find(bson.M{"tags": bson.M{"$in": arr}})
	var ret []*model.Repository
	
	for docs.Next(context.Background()) {
		var repo model.Repository
		err := docs.Decode(&repo)
		if err != nil {
			log.Fatal("Error on docoding repository", err)
		}
		ret = append(ret, &repo)
	}
	docs.Close(context.Background())
	
	return ret
}

func (d *RepoDAO) DeleteTag(user string, repoId uint64, tag model.Tag) model.Repository {
	log.Printf("[Info] Deleting tag [%s] from [%s] [%d]", tag, user, repoId)
	
	repo := d.GetRepository(user, repoId)
	
	//checking if the repo from MongoDB is not an empty struct
	if !reflect.DeepEqual(model.Repository{}, repo) {
		for i, t := range repo.Tags {
			if t.Name == tag.Name {
				repo.Tags = append(repo.Tags[:i], repo.Tags[i+1:]...)
				break
			}
		}
		DB.UpdateOne(bson.M{"user": user, "id": repoId}, bson.M{"$set": bson.M{"tags": repo.Tags}})
	}
	
	return repo
}

func (d *RepoDAO) UpdateTag(user string, repoId uint64, tagUpdate model.TagUpdate) (model.Repository, error) {
	log.Printf("[Info] Update tag [%s] to [%s] in [%s] [%d]", tagUpdate.OldTag, tagUpdate.NewTag, user, repoId)
	
	repo := d.GetRepository(user, repoId)
	var err error = errors.New("Tag not found in repository")
	if !reflect.DeepEqual(model.Repository{}, repo) {
		for i, t := range repo.Tags {
			if t.Name == tagUpdate.OldTag.Name {
				repo.Tags = append(repo.Tags[:i], repo.Tags[i+1:]...)
				repo.Tags = append(repo.Tags, tagUpdate.NewTag)
				err = nil
				break
			}
		}
		DB.UpdateOne(bson.M{"user": user, "id": repoId}, bson.M{"$set": bson.M{"tags": repo.Tags}})
	}
	
	return repo, err
}

func (d *RepoDAO) UpdateRepositoryCollection(user string, toInsert []*model.Repository, toRemove []*model.Repository) {
	log.Printf("[Info] Refresh [%s] repositories", user)
	
	var idsToRemove []uint64
	for _, r := range toRemove {
		idsToRemove = append(idsToRemove, r.Id)
	}
	
	if len(idsToRemove) > 0 {
		DB.DeleteMany(bson.M{"user": user, "id" : bson.M{"$in": idsToRemove}})
	}

	if len(toInsert) > 0 {
		d.InsertRepositories(user, toInsert)
	}
}