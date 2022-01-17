package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/urfave/cli/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gopkg.in/gookit/color.v1"
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
			{
				Name:    "all",
				Aliases: []string{"l"},
				Usage:   "List all tasks",
				Action: func(_ *cli.Context) error {
					tasks, err := getAll()
					if err != nil {
						if err == mongo.ErrNoDocuments {
							fmt.Println("Nada para ver aqui.\nExecute `add 'task'` para adicionar uma tarefa")
							return nil
						}
						return err
					}
					printTask(tasks)
					return nil
				},
			},
			{
				Name:    "done",
				Aliases: []string{"d"},
				Usage:   "complete a task on the list",
				Action: func(c *cli.Context) error {
					text := c.Args().First()
					return completeTask(text)
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func printTask(task []*Task) {
	for i, t := range task {
		if t.Completed {
			color.Green.Printf("%d: %s\n", i+1, t.Text)
		}
		color.Yellow.Printf("%d: %s\n", i+1, t.Text)
	}
}

func createTask(task *Task) error {
	_, err := collection.InsertOne(ctx, task)
	return err
}

func getAll() ([]*Task, error) {
	filter := bson.D{}
	return filterTask(filter)
}

func completeTask(text string) error {
	filter := bson.D{primitive.E{Key: "text", Value: text}}

	update := bson.D{primitive.E{Key: "$set", Value: bson.D{
		primitive.E{Key: "completed", Value: true},
	}}}

	t := &Task{}

	return collection.FindOneAndUpdate(ctx, filter, update).Decode(&t)
}

func filterTask(filter interface{}) ([]*Task, error) {
	var tasks []*Task

	cur, err := collection.Find(ctx, filter)
	if err != nil {
		return tasks, err
	}

	defer cur.Close(ctx)

	for cur.Next(ctx) {
		var t Task
		err := cur.Decode(&t)

		if err != nil {
			return tasks, err
		}
		tasks = append(tasks, &t)
	}

	if err := cur.Err(); err != nil {
		if err != nil {
			return tasks, err
		}
	}

	if len(tasks) == 0 {
		return tasks, mongo.ErrNoDocuments
	}
	return tasks, nil
}
