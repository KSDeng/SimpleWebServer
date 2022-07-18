package main

import (
	"entry-task/config"
	"entry-task/tcpserver/db"
	"fmt"
	"github.com/spf13/viper"
	"golang.org/x/crypto/bcrypt"
)

var configs config.Configurations

func loadConfigs() {
	viper.SetConfigName("config")
	viper.AddConfigPath("./config")
	viper.SetConfigType("yml")
	if err := viper.ReadInConfig(); err != nil {
		fmt.Printf("Error reading config file, %s", err)
	}

	err := viper.Unmarshal(&configs)
	if err != nil {
		fmt.Printf("Unable to decode into struct, %v", err)
	}
}

func HashPassword(password string) string {
	bytes, _ := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes)
}

func main() {
	loadConfigs()
	db.InitDBConnection(configs.Database)

	users := db.SelectAllUsers()
	db.TruncateUserTable()
	for index, user := range users {
		fmt.Println(index)
		db.AddUser(user.Nickname, HashPassword(user.Password), "")
	}

}
