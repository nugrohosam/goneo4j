package main

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/neo4j/neo4j-go-driver/neo4j"
)

func main() {
	app := fiber.New()

	app.Get("/", func(c *fiber.Ctx) error {
		fmt.Println("Ready to connect")
		data, err := helloWorld("bolt://127.0.0.1:7687", "neo4j", "root")
		fmt.Println("After connect")
		if err != nil {
			return err
		}

		fmt.Println("results :", data)
		return nil
	})

	app.Listen(":3000")
}

func helloWorld(uri, username, password string) (interface{}, error) {
	driver, err := neo4j.NewDriver(uri, neo4j.BasicAuth(username, password, ""))
	fmt.Println("Got driver of neo4j")
	if err != nil {
		return nil, err
	}

	fmt.Println("Prepare defer disconnect driver of neo4j")
	defer driver.Close()

	fmt.Println("Prepare session of driver neo4j")
	session, _ := driver.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})

	fmt.Println("Prepare defer disconnect session of driver neo4j")
	defer session.Close()

	greeting, err := session.ReadTransaction(func(transaction neo4j.Transaction) (interface{}, error) {
		result, err := transaction.Run(
			"match (a:User) return a.name as name",
			map[string]interface{}{"name": "Rojak"})

		if err != nil {
			return nil, err
		}

		data := []interface{}{}
		for {
			if result.Next() {
				value, _ := result.Record().Get("name")
				data = append(data, value)
			} else {
				break
			}
		}

		return data, result.Err()
	})

	if err != nil {
		return nil, err
	}

	return greeting, nil
}
