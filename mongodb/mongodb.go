package mongodb

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const MONGO_BACKEND = "mongodb://localhost:27017"
const DB_NAME = "github_tags"
const COLLECTION_REPOS = "collection_repos"

type MongoDB struct {

}

var db *mongo.Database
var collectionRepos *mongo.Collection

/*
	Inicializando o MongoDB
	TODO: Remover os 'Fatal' para rodar o serviço graciosamente
*/
func init() {
	//configuração de conexão do MongoDB. Aqui configurado para conectar na maquina local e com a configuração default de instalação
	//rodando como serviço.
	var clientOptions = options.Client().ApplyURI(MONGO_BACKEND)

	client, err := mongo.NewClient(clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	err = client.Connect(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("[info] Connected to MongoDB")

	db = client.Database(DB_NAME)
	collectionRepos = db.Collection(COLLECTION_REPOS)

}
/*
	Funções "façade" para desacoplar as funções do DB de outras partes do código.
*/
func (m *MongoDB) GetCollection() *mongo.Collection {
	return collectionRepos
}

func (m *MongoDB) InsertMany(data []interface{}) {
	_, err := collectionRepos.InsertMany(context.TODO(), data)
	
	if err != nil {
		log.Fatal(err)
	}
}

func (m *MongoDB) Find(filter interface{}) *mongo.Cursor {
	docs, err := collectionRepos.Find(context.Background(), filter)
	
	if err != nil {
		log.Fatal("Error finding", err)
	}
	return docs
}

func (m *MongoDB) FindOne(filter interface{}) *mongo.SingleResult {
	return collectionRepos.FindOne(context.Background(), filter)
}

func (m *MongoDB) UpdateOne(filter interface{}, update interface{}) {
	_, err := collectionRepos.UpdateOne(context.Background(), filter, update)
	
	if err != nil {
		log.Fatal("Error updating", err)
	}
}

func (m *MongoDB) DeleteMany(data interface{}) {
	_, err := collectionRepos.DeleteMany(context.Background(), data)
	
	if err != nil {
		log.Fatal("Error deleting", err)
	}
}
