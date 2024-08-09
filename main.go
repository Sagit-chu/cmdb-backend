package main

import (
	"cmdb-backend/config"
	"cmdb-backend/routes"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	// 初始化数据库
	config.InitDB()

	// 设置路由
	router := mux.NewRouter()
	routes.RegisterRoutes(router)

	// 启动服务器
	port := "3000"
	fmt.Println("Server is running on port", port)
	log.Fatal(http.ListenAndServe(":"+port, router))
}
