package main

import (
	"cmdb-backend/config"
	"cmdb-backend/routes"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

func main() {

	// 初始化数据库
	config.InitDB()

	// 设置路由
	router := mux.NewRouter()
	routes.RegisterRoutes(router)

	// 配置 CORS 中间件
	corsHandler := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"}, // 允许所有来源（你可以指定具体的来源）
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
	}).Handler

	// 启动服务器
	port := "3000"
	fmt.Println("Server is running on port", port)
	log.Fatal(http.ListenAndServe(":"+port, corsHandler(router)))
}
