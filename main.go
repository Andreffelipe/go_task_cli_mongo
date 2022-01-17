package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/urfave/cli/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	collection *mongo.Collection
	ctx        = context.TODO()
)

type Task struct {
	ID        primitive.ObjectID `bson:"_id"`
	CreatedAt time.Time          `bson:"created_at"`
	UpdatedAt time.Time          `bson:"updated_at"`
	Text      string             `bson:"text"`
	Completed bool               `bson:"completed"`
}

func init() {
	db_url := os.Getenv("DATABASE_URL")
	if db_url == "" {
		fmt.Println("Variavel vazia")
		db_url = "mongodb://root:docker@localhost:27017/"
	}
	clientOption := options.Client().ApplyURI(db_url)

	client, err := mongo.Connect(ctx, clientOption)

	if err != nil {
		log.Fatal(err)
	}

	// verifica se a conexão foi bem sucedida
	if err := client.Ping(ctx, nil); err != nil {
		log.Fatal(err)
	}

	// criando o banco de dados e a coleção
	collection = client.Database("tasker").Collection("tasks")
}
func main() {
	app := &cli.App{
		Name:  "Tasker",
		Usage: "A Simple CLI program to manage your task",
		Commands: []*cli.Command{
			{
				Name:    "add",
				Aliases: []string{"a"},
				Usage:   "add a task to the list",
				Action: func(c *cli.Context) error {
					str := c.Args().First()
					if str == "" {
						return errors.New("Cannot ass an empty task")
					}

					task := &Task{
						ID:        primitive.NewObjectID(),
						CreatedAt: time.Now(),
						UpdatedAt: time.Now(),
						Text:      str,
						Completed: false,
					}
					return createTask(task)
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func createTask(task *Task) error {
	_, err := collection.InsertOne(ctx, task)
	return err
}
