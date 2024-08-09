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
				// 检查并创建表
				err = ensureTablesExist()
				if err != nil {
					log.Fatal("Failed to ensure tables exist:", err)
				}
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

// ensureTablesExist 检查并创建表
func ensureTablesExist() error {
	// 检查表是否存在
	query := `
		SELECT COUNT(*)
		FROM information_schema.tables 
		WHERE table_schema = DATABASE() 
		AND table_name = 'cmdb_assets'
	`
	var count int
	err := DB.QueryRow(query).Scan(&count)
	if err != nil {
		return fmt.Errorf("failed to check if table exists: %v", err)
	}

	// 如果表不存在，则创建
	if count == 0 {
		fmt.Println("Table cmdb_assets does not exist. Creating...")
		createTableQuery := `
		CREATE TABLE cmdb_assets (
			id INT AUTO_INCREMENT PRIMARY KEY,
			ip VARCHAR(255) NOT NULL,
			application_system VARCHAR(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci,
			application_manager VARCHAR(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci,
			overall_manager VARCHAR(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci,
			is_virtual_machine BOOLEAN,
			resource_pool VARCHAR(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci,
			data_center VARCHAR(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci,
			rack_location VARCHAR(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci,
			sn_number VARCHAR(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci,
			out_of_band_ip VARCHAR(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
		) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci;
		`
		_, err := DB.Exec(createTableQuery)
		if err != nil {
			return fmt.Errorf("failed to create table: %v", err)
		}
		fmt.Println("Table cmdb_assets created successfully.")
	} else {
		fmt.Println("Table cmdb_assets already exists.")
	}

	return nil
}
