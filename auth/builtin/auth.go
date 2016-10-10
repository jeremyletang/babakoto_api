package builtinauth

import (
	"context"
	"net/http"
	"time"

	log "github.com/cihub/seelog"
	"github.com/gorilla/mux"
	"github.com/jeremyletang/babakoto_api/ctxext"
	"github.com/jeremyletang/babakoto_api/dao"
	"github.com/jeremyletang/babakoto_api/domain"
	"github.com/jeremyletang/babakoto_api/jsend"
	"github.com/jeremyletang/babakoto_api/utils"
	"github.com/jeremyletang/babakoto_api/utils/errmsg"
	"github.com/jinzhu/gorm"
	"github.com/satori/go.uuid"
	"golang.org/x/crypto/bcrypt"
)

const (
	// seconds for a day
	defaultUserSignupVerificationTtl = 86400
	// accesstoken = two days
	defaultAccessTokenTtl = 172800
)

type BuiltinAuth struct {
	db *gorm.DB
}

func NewBuiltinAuth(db *gorm.DB) BuiltinAuth {
	return BuiltinAuth{db: db}
}

type LoginRequest struct {
	Identifier string `json:"identifier"`
	Password   string `json:"password"`
}

type SignupRequest struct {
	Email    string `json:"email"`
	Username string `json:"username"`
	Password string `json:"password"`
}

func loginValidator(l *LoginRequest) map[string]interface{} {
	errors := map[string]interface{}{}
	if l.Identifier == "" {
		errors["identifier"] = errmsg.MissingFieldError
	}
	if l.Password == "" {
		errors["password"] = errmsg.MissingFieldError
	}
	return errors
}

func generateAccessToken(db *gorm.DB, u domain.User) (domain.AccessToken, error) {
	// if access token exist just delete it and return a new one
	tokenDao := dao.NewAccessTokenDao(db)
	if _, err := tokenDao.GetByUserId(u.Id); err == nil {
		log.Infof("[builtinauth.generateAccessToken] access token for user [id=%s] already exist, delete it", u.Id)
	}

	newToken := domain.AccessToken{
		Id:        uuid.NewV4().String(),
		UserId:    u.Id,
		Ttl:       defaultAccessTokenTtl,
		CreatedAt: time.Now(),
	}

	if err := tokenDao.Create(newToken); err != nil {
		log.Errorf("[builtinauth.generateAccessToken] unable to save token: %s", err.Error())
		return newToken, err
	}

	return newToken, nil
}

func (ba *BuiltinAuth) Login(w http.ResponseWriter, r *http.Request) {
	var login LoginRequest
	if err := utils.ReadRequestBody(r, &login); err != nil {
		log.Errorf("[builtinauth.Login] invalid request body: %s", err.Error())
		utils.WriteJsonResponse(w, http.StatusBadRequest, jsend.Fail("invalid json"))
	} else {
		// validation error
		if err := loginValidator(&login); len(err) != 0 {
			log.Errorf("[builtinauth.Login] validation error: %#v", err)
			utils.WriteJsonResponse(w, http.StatusBadRequest, jsend.Fail(err))
			return
		}

		// request is good let's process it
		// first get the user
		userDao := dao.NewUserDao(ba.db)
		u, err := userDao.GetByEmailOrUsername(login.Identifier)
		if err != nil {
			log.Errorf("[builtinauth.Login] unknow identifier: %s", login.Identifier)
			utils.WriteJsonResponse(w, http.StatusBadRequest,
				jsend.FailWithName("unable to login", "login"))
			return
		}

		// try to match the password
		if err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(login.Password)); err != nil {
			log.Errorf("[builtinauth.Login] invalid password for identifier: %s", login.Identifier)
			utils.WriteJsonResponse(w, http.StatusBadRequest,
				jsend.FailWithName("unable to login", "login"))
			return
		}

		// password have matched, generate or regenerate the access_token
		at, err := generateAccessToken(ba.db, u)
		if err != nil {
			utils.WriteJsonResponse(w, http.StatusInternalServerError,
				jsend.Error("internal error"))
		}

		// hide password for now
		u.Password = ""
		res := map[string]interface{}{}
		res["access_token"] = at
		res["user"] = u

		utils.WriteJsonResponse(w, http.StatusOK,
			jsend.New(res))
	}
}
func (ba *BuiltinAuth) Logout(
	ctx context.Context,
	w http.ResponseWriter,
	r *http.Request,
) {
	if token, ok := ctxext.ExtractAccessToken(ctx); ok {
		tokenDao := dao.NewAccessTokenDao(ba.db)
		if err := tokenDao.Delete(token.Id); err != nil && err != gorm.ErrRecordNotFound {
			utils.WriteJsonResponse(w, http.StatusInternalServerError,
				jsend.Error("database error"))
			return
		} else {
			utils.WriteJsonResponse(w, http.StatusOK, jsend.New(nil))
		}
	} else {
		utils.WriteJsonResponse(w, http.StatusOK, jsend.Error("internal error"))
	}
}

func signupValidator(l *SignupRequest) map[string]interface{} {
	errors := map[string]interface{}{}
	if l.Email == "" {
		errors["mail"] = errmsg.MissingFieldError
	}
	if l.Username == "" {
		errors["username"] = errmsg.MissingFieldError
	}
	if l.Password == "" {
		errors["password"] = errmsg.MissingFieldError
	}
	return errors
}

func checkExistsByEmailOrUsername(db *gorm.DB, email, username string) map[string]interface{} {
	errors := map[string]interface{}{}
	userDao := dao.NewUserDao(db)
	if _, err := userDao.GetByMail(email); err == nil {
		errors[email] = errmsg.MailAlreadyUsed
	}
	if _, err := userDao.GetByUsername(username); err == nil {
		errors[username] = errmsg.UsernameAlreadyUsed
	}
	return errors
}

func (ba *BuiltinAuth) Signup(w http.ResponseWriter, r *http.Request) {
	var signup SignupRequest
	if err := utils.ReadRequestBody(r, &signup); err != nil {
		log.Errorf("[builtinauth.Signup] invalid request body: %s", err.Error())
		utils.WriteJsonResponse(w, http.StatusBadRequest, jsend.Fail("invalid json"))
	} else {
		// validation error
		if err := signupValidator(&signup); len(err) != 0 {
			log.Errorf("[builtinauth.Signup] validation error: %#v", err)
			utils.WriteJsonResponse(w, http.StatusBadRequest, jsend.Fail(err))
			return
		}

		// request is good let's process it
		// first check if a user with this username or email already exist
		if err := checkExistsByEmailOrUsername(ba.db, signup.Email, signup.Username); len(err) != 0 {
			log.Errorf("[builtinauth.Signup] user data already exists: %#v", err)
			utils.WriteJsonResponse(w, http.StatusBadRequest, jsend.Fail(err))
			return
		}

		cryptedPassword, _ := bcrypt.GenerateFromPassword([]byte(signup.Password), bcrypt.DefaultCost)

		// user don't exists so create it
		newUser := domain.User{
			Id:        uuid.NewV4().String(),
			Username:  signup.Username,
			Email:     signup.Email,
			Password:  string(cryptedPassword),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		// save user
		userDao := dao.NewUserDao(ba.db)
		if err := userDao.Create(newUser); err != nil {
			log.Errorf("[builtinauth.Signup] unable to create a new user: %s", error.Error)
			utils.WriteJsonResponse(w, http.StatusInternalServerError,
				jsend.Error("database error"))
			return
		}

		// create and save user signup verification request
		userSignupVerif := domain.UserSignupVerification{
			Id:        uuid.NewV4().String(),
			UserId:    newUser.Id,
			Ttl:       defaultUserSignupVerificationTtl,
			CreatedAt: time.Now(),
		}

		// save verif
		signupDao := dao.NewUserSignupVerificationDao(ba.db)
		if err := signupDao.Create(userSignupVerif); err != nil {
			log.Errorf("[builtinauth.Signup] unable to create user verification: %s", error.Error)
			utils.WriteJsonResponse(w, http.StatusInternalServerError,
				jsend.Error("database error"))
			return
		}

		// create response
		res := map[string]interface{}{}
		res["user"] = newUser
		res["signup_verification"] = userSignupVerif
		utils.WriteJsonResponse(w, http.StatusOK, jsend.New(res))
	}
}

func (ba *BuiltinAuth) Verify(w http.ResponseWriter, r *http.Request) {
	// get path parameters
	vars := mux.Vars(r)
	verifId := vars["id"]

	// try to get the verif from the id
	verifDao := dao.NewUserSignupVerificationDao(ba.db)
	if _, err := verifDao.GetById(verifId); err != nil {
		utils.WriteJsonResponse(w, http.StatusBadRequest,
			jsend.FailWithName("Invalid user signup verification id", "id"))
		return
	}

	// remove it to validate the user
	if err := verifDao.Delete(verifId); err != nil {
		utils.WriteJsonResponse(w, http.StatusInternalServerError,
			jsend.Error("database error"))
		return
	}

	utils.WriteJsonResponse(w, http.StatusOK, jsend.New(nil))
}

func (ba *BuiltinAuth) TokenInfos(
	ctx context.Context,
	w http.ResponseWriter,
	r *http.Request) {
	infos := map[string]interface{}{}
	accessToken, _ := ctxext.ExtractAccessToken(ctx)
	user, _ := ctxext.ExtractUser(ctx)
	infos["access_token"] = accessToken
	infos["user"] = user
	utils.WriteJsonResponse(w, http.StatusOK, jsend.New(infos))
}
