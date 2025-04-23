package main

import (
	"fmt"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID         int    `json:"id"`
	Name       string `json:"name"`
	Email      string `json:"email"`
	Password   string `json:"password"`
	Created_at string `json:"created_at"`
}

func GetUsers() ([]User, error) {
	rows, err := db.Query("SELECT id,name,email,password,created_at FROM Users")
	if err != nil {
		return nil, err
	}

	var users []User

	for rows.Next() {
		var u User
		err := rows.Scan(&u.ID, &u.Name, &u.Email, &u.Password, &u.Created_at)
		if err != nil {
			return nil, err
		}
		users = append(users, u)
	}
	if err = rows.Err(); err != nil {
		fmt.Println("Error")
		return nil, err
	}
	return users, nil
}

func GetUserByID(id int) (User, error) {
	row := db.QueryRow("SELECT id,name,email,password,created_at FROM Users WHERE id = $1", id)

	var u User

	err := row.Scan(&u.ID, &u.Name, &u.Email, &u.Password, &u.Created_at)

	if err != nil {
		return User{}, err
	}
	return u, nil

}

func CreateUser(user User) (User, error) {
	var u User

	//Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return User{}, err
	}

	err = db.QueryRow("INSERT INTO Users (name, email, password,created_at) VALUES ($1, $2, $3,$4) RETURNING id, name, email, password, created_at", user.Name, user.Email, hashedPassword, time.Now()).Scan(&u.ID, &u.Name, &u.Email, &u.Password, &u.Created_at)
	if err != nil {
		return User{}, err
	}
	return u, nil
}

func UpdateUser(id int, user *User) (User, error) {
	var u User

	//Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return User{}, err
	}

	err = db.QueryRow("UPDATE public.Users SET name=$1, email=$2, password=$3	WHERE id=$4 RETURNING id,name,email,password;", user.Name, user.Email, hashedPassword, id).Scan(&u.ID, &u.Name, &u.Email, &u.Password)
	if err != nil {
		return User{}, err
	}
	return u, err
}

func DeleteUser(id int) error {
	_, err := db.Exec("DELETE FROM public.Users WHERE id=$1;", id)
	return err
}

func CheckEmail(email string) (int, error) {
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM Users WHERE email = $1", email).Scan(&count)
	if err != nil {
		return count, err
	}
	return count, nil
}

func CheckId(id int) (int, error) {
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM Users WHERE id = $1", id).Scan(&count)
	if err != nil {
		return count, err
	}
	return count, nil
}
