package main

import (
"context"
"database/sql"
"fmt"
"os"
"github.com/gofiber/fiber/v2"
_ "github.com/go-sql-driver/mysql"
)
func getEnv(env string, defaultVal string) (key string) {
	if key = os.Getenv(env); key == "" {
		key = defaultVal
	}
	return key
}
var (
	ctx = context.Background()
)
func main() {
	// test mysql
	host := getEnv("MYSQL_SERVICE_HOST", "localhost")
	port := getEnv("MYSQL_SERVICE_PORT", "3306")
	password := getEnv("MYSQL_ROOT_PASSWORD", "root")
	dsn := fmt.Sprintf("root:%s@tcp(%s:%s)/red_envelope", password, host, port)
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		panic(err)
	}
	if err = db.Ping(); err != nil {
		panic(err)
	}
	// http server
	app := fiber.New()
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!!!")
	})
	panic(app.Listen(":8080"))
}
