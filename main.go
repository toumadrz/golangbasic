package main

import (
	"database/sql"
	"fmt"
	"log"
	"strconv"

	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

var db *sql.DB

func main() {
	// โหลด .env มาใช้
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	// ผูก .env
	host := os.Getenv("DB_HOST")
	port, _ := strconv.Atoi(os.Getenv("DB_PORT")) // Convert port to int
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")

	//เชื่อมต่อ
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	// Open a connection
	sdb, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatal(err)
	}
	db = sdb
	defer db.Close()

	// Check the connection
	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Successfully connected!")

	app := fiber.New()
	app.Get("/users", getUsers)
	app.Get("/user/:id", getUser)
	app.Post("/user", createUser)
	app.Put("/user/:id", updateUser)
	app.Delete("/user/:id", deleteUser)
	app.Listen(":8000")
}

// Fiber get user ทั้งหมด
func getUsers(c *fiber.Ctx) error {
	user, err := GetUsers()
	if err != nil {
		log.Println("GetUsers error:", err)
		return c.SendStatus(fiber.StatusBadRequest)
	}
	return c.JSON(user)
}

// Fiber get user จาก Id
func getUser(c *fiber.Ctx) error {
	Id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}
	user, err := GetUserByID(Id)
	if err == sql.ErrNoRows {
		return c.SendString("Id not found")
	}
	if err != nil {
		log.Println("GetUsers error:", err)
		return c.SendStatus(fiber.StatusBadRequest)
	}
	return c.JSON(user)
}

// Fiber create user
func createUser(c *fiber.Ctx) error {
	u := new(User)
	if err := c.BodyParser(u); err != nil {
		log.Println("BodyParser error:", err)
		return c.SendStatus(fiber.StatusBadRequest)
	}

	// เช็ค email ซ้ำ
	count, err := CheckEmail(u.Email)
	if err != nil {
		log.Println("Check email duplicate error:", err)
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	// ถ้า count มากกว่า 0 แปลว่า email ซ้ำ
	if count > 0 {
		return c.SendString("duplicate email")
	}

	// email ไม่ซ้ำ จึงcreateuser
	user, err := CreateUser(*u)
	if err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	return c.JSON(user)
}

// Fiber update user
func updateUser(c *fiber.Ctx) error {
	Id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	u := new(User)
	if err := c.BodyParser(u); err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	user, err := UpdateUser(Id, u)

	if err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}
	return c.JSON(user)

}

// Fiber delete user
func deleteUser(c *fiber.Ctx) error {
	Id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	// เช็คว่ามีidไหม
	count, err := CheckId(Id)
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	// ถ้า count <= 0 แปลว่า ไม่มี id นี้
	if count <= 0 {
		return c.SendString("no id")
	}
	//มีidค่อยลบ
	err = DeleteUser(Id)

	if err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}
	return c.SendString("Delete user successfully.")
}
