package main

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"cmdb-backend/config"
	"cmdb-backend/routes"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/rs/cors"
	"golang.org/x/oauth2"
)

// 配置 OAuth2
var (
	oauth2Config = &oauth2.Config{
		ClientID:     os.Getenv("OAUTH_CLIENT_ID"),
		ClientSecret: os.Getenv("OAUTH_CLIENT_SECRET"),
		RedirectURL:  os.Getenv("RedirectURL"), //"http://localhost:3000/auth/callback",
		Scopes:       []string{"openid", "profile", "email"},
		Endpoint: oauth2.Endpoint{
			AuthURL:  os.Getenv("OAUTH_AUTH_URL"),  // 授权端点
			TokenURL: os.Getenv("OAUTH_TOKEN_URL"), // 令牌端点
		},
	}
	oauthStateString = "random"
	store            = sessions.NewCookieStore([]byte(os.Getenv("SESSION_KEY")))
)

func main() {

	// 初始化数据库
	config.InitDB()

	// 设置路由
	router := mux.NewRouter()
	routes.RegisterRoutes(router)

	// OAuth2 路由
	router.HandleFunc("/auth/login", handleLogin)
	router.HandleFunc("/auth/callback", handleCallback)
	router.HandleFunc("/profile", handleProfile)
	router.HandleFunc("/logout", handleLogout)

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

func handleLogin(w http.ResponseWriter, r *http.Request) {
	url := oauth2Config.AuthCodeURL(oauthStateString)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func handleCallback(w http.ResponseWriter, r *http.Request) {
	if r.FormValue("state") != oauthStateString {
		log.Println("Invalid OAuth state")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	token, err := oauth2Config.Exchange(context.Background(), r.FormValue("code"))
	if err != nil {
		log.Printf("Error while exchanging token: %v\n", err)
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	client := oauth2Config.Client(context.Background(), token)
	userInfo, err := client.Get(os.Getenv("OAUTH_USERINFO_URL")) // 使用用户信息端点
	if err != nil {
		log.Printf("Error getting user info: %v\n", err)
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}
	defer userInfo.Body.Close()

	buf := new(bytes.Buffer)
	buf.ReadFrom(userInfo.Body)
	newStr := buf.String()

	session, _ := store.Get(r, "auth-session")
	session.Values["user-info"] = newStr
	session.Save(r, w)

	http.Redirect(w, r, "/profile", http.StatusTemporaryRedirect)
}

func handleProfile(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "auth-session")
	userInfo := session.Values["user-info"]
	if userInfo == nil {
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	fmt.Fprintf(w, "User Info: %s", userInfo)
}

func handleLogout(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "auth-session")
	session.Options.MaxAge = -1 // 删除会话
	session.Save(r, w)

	http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
}
