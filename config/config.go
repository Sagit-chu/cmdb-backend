package config

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

var DB *sql.DB

func InitDB() {
	var err error
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")

	// 生成数据源名称（DSN）
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", dbUser, dbPassword, dbHost, dbPort, dbName)

	// 定义最大重试次数和重试间隔
	maxRetries := 10
	retryInterval := 5 * time.Second

	// 循环尝试连接数据库，直到成功或达到最大重试次数
	for i := 0; i < maxRetries; i++ {
		DB, err = sql.Open("mysql", dsn)
		if err != nil {
			log.Printf("Failed to connect to database (attempt %d/%d): %v", i+1, maxRetries, err)
		} else {
			err = DB.Ping()
			if err == nil {
				fmt.Println("Database connection established")
				return
			}
			log.Printf("Failed to ping database (attempt %d/%d): %v", i+1, maxRetries, err)
		}

		// 等待一段时间再重试
		log.Printf("Retrying in %v...", retryInterval)
		time.Sleep(retryInterval)
	}

	// 如果达到最大重试次数，程序退出
	log.Fatalf("Could not connect to database after %d attempts: %v", maxRetries, err)
}
