package db

import (
	"database/sql"
	"entry-task/config"
	"fmt"
	"github.com/go-sql-driver/mysql"
	"log"
)

var db_driver *sql.DB

type User_DB struct {
	Id        int64
	Nickname  string
	Password  string
	PhotoPath string
}

func InitDBConnection(configs config.DatabaseConfigs) {
	cfg := mysql.Config{
		User:   configs.DB_User,
		Passwd: configs.DB_Password,
		Net:    configs.DB_Net,
		Addr:   configs.DB_Addr,
		DBName: configs.DB_Name,
	}
	var err error
	db_driver, err = sql.Open(configs.DB_Driver, cfg.FormatDSN())
	if err != nil {
		log.Fatal(err)
		return
	}

	db_driver.SetMaxOpenConns(configs.DB_Max_Conn)
	db_driver.SetMaxIdleConns(configs.DB_Max_Idle_Conn)
	db_driver.Exec("use user_info_db")
}

func PhotoPathById(id int64) (string, error) {
	var photoPath string

	row := db_driver.QueryRow("SELECT profile_photo_path from user_tab where id = ?", id)
	if err := row.Scan(&photoPath); err != nil {
		if err == sql.ErrNoRows {
			return photoPath, fmt.Errorf("PhotoPathById(), no such user, id = %d", id)
		}
		return photoPath, fmt.Errorf("PhotoPathById(), %d: %v", id, err)
	}
	return photoPath, nil
}

func UserInfoByNickname(nickname string) (User_DB, error) {
	var user User_DB

	row := db_driver.QueryRow("SELECT * FROM user_tab WHERE user_nickname = ?", nickname)
	if err := row.Scan(&user.Id, &user.Nickname, &user.Password, &user.PhotoPath); err != nil {
		return user, fmt.Errorf("UserInfoByNickname %s: %v", nickname, err)
	}
	return user, nil
}

func UserInfoById(id int64) (User_DB, error) {
	var user User_DB

	row := db_driver.QueryRow("SELECT * FROM user_tab WHERE id = ?", id)
	if err := row.Scan(&user.Id, &user.Nickname, &user.Password, &user.PhotoPath); err != nil {
		return user, fmt.Errorf("UserInfoById %d: %v", id, err)
	}
	return user, nil
}

func UpdatePhotoPathById(id int64, photoPath string) error {
	_, err := db_driver.Exec("UPDATE user_tab SET profile_photo_path = ? WHERE id = ?", photoPath, id)
	if err != nil {
		return err
	}
	return nil
}

func UpdateNicknameById(id int64, nickname string) error {
	_, err := db_driver.Exec("UPDATE user_tab SET user_nickname = ? WHERE id = ?", nickname, id)
	if err != nil {
		return err
	}
	return nil
}

func AddUser(nickname string, password string, photo_path string) error {
	_, err := db_driver.Exec("INSERT INTO user_tab (user_nickname, user_password, profile_photo_path) VALUES (?, ?, ?);", nickname, password, photo_path)
	if err != nil {
		return fmt.Errorf("addUser: %v", err)
	}
	return nil
}

func SelectAllUsers() []User_DB {
	var users []User_DB

	rows, err := db_driver.Query("SELECT * FROM user_tab")
	if err != nil {
		log.Fatal(err)
		return nil
	}

	defer rows.Close()

	for rows.Next() {
		var user User_DB
		if err := rows.Scan(&user.Id, &user.Nickname, &user.Password, &user.PhotoPath); err != nil {
			log.Fatal("Scan error")
			return nil
		}
		users = append(users, user)
	}
	return users
}

func TruncateUserTable() {
	db_driver.Exec("TRUNCATE user_tab")
}
