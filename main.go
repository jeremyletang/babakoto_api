package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	log "github.com/cihub/seelog"
	"github.com/gorilla/mux"
	"github.com/jeremyletang/babakoto_api/auth/builtin"
	"github.com/jeremyletang/babakoto_api/ctxext"
	"github.com/jeremyletang/babakoto_api/jsend"
	"github.com/jeremyletang/babakoto_api/services/user"
	"github.com/jeremyletang/babakoto_api/utils"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/rs/cors"
)

var db *gorm.DB

type MysqlConfig struct {
	User     string `json:"user"`
	Password string `json:"password"`
	Ip       string `json:"ip"`
	Port     string `json:"port"`
	Database string `json:"database"`
}

type Config struct {
	Mysql MysqlConfig `json:"mysql"`
}

func init() {
	logger, err := log.LoggerFromConfigAsString(seelogConfig)
	if err != nil {
		panic("unable to find logger configuration")
	}
	log.ReplaceLogger(logger)
}

// "root:root@tcp(192.168.99.100:3307)/babakoto?parseTime=True"
func getConfig() Config {
	c, err := ioutil.ReadFile(".babakoto.config.json")
	if err != nil {
		panic(fmt.Sprintf("[getConfig] unable to read config file: %s", err.Error()))
	}

	var config Config
	if err := json.Unmarshal(c, &config); err != nil {
		panic(fmt.Sprintf("[getConfig] invalid config format: %s", err.Error()))
	}

	return config
}

func main() {
	// read config
	config := getConfig()
	// init db
	var err error
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?parseTime=True",
		config.Mysql.User,
		config.Mysql.Password,
		config.Mysql.Ip,
		config.Mysql.Port,
		config.Mysql.Database)

	if db, err = gorm.Open("mysql", dsn); err != nil {
		panic(fmt.Sprintf("[main] unable to initialize gorm: %s", err.Error()))
	}
	defer db.Close()

	r := makeRoutes()
	handler := cors.New(cors.Options{AllowedHeaders: []string{"*"}, AllowedMethods: []string{"GET", "POST", "PUT", "DELETE"}}).Handler(r)
	log.Info("Starting http server")
	log.Critical(http.ListenAndServe(fmt.Sprintf(":%v", 9992), handler))

}

func makeRoutes() *mux.Router {
	r := mux.NewRouter()

	// builtin auth routes
	builtinAuth := builtinauth.NewBuiltinAuth(db)
	r.HandleFunc("/api/v1/user/login",
		builtinAuth.Login).Methods("POST")
	r.HandleFunc("/api/v1/user/signup",
		builtinAuth.Signup).Methods("POST")
	r.HandleFunc("/api/v1/user/verify/{id}",
		builtinAuth.Verify).Methods("GET")
	// need login
	r.HandleFunc("/api/v1/user/token-infos",
		addContext(addUserInfo(builtinAuth.TokenInfos, db))).Methods("GET")
	r.HandleFunc("/api/v1/user/logout",
		addContext(addUserInfo(builtinAuth.Logout, db))).Methods("GET")

	return r
}

func getToken(r *http.Request) (string, error) {
	var token string
	if token = r.Header.Get("authorization"); token == "" {
		if token = r.Header.Get("Authorization"); token == "" {
			if token = r.Form.Get("access_token"); token == "" {
				return token, errors.New("Missing access token")
			}
		}
	}

	return token, nil
}

func addUserInfo(
	f func(context.Context, http.ResponseWriter, *http.Request),
	db *gorm.DB,
) func(context.Context, http.ResponseWriter, *http.Request) {
	return user.AddUserInfoToContext(f, db)
}

func addContext(
	f func(context.Context, http.ResponseWriter, *http.Request),
) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		token, err := getToken(r)

		if err != nil {
			utils.WriteJsonResponse(w, http.StatusUnauthorized,
				jsend.FailWithName(err.Error(), "access_token"))
			return
		}

		ctx := ctxext.AddAccessTokenString(context.Background(), token)
		// finally call the handler
		f(ctx, w, r)
	}
}
