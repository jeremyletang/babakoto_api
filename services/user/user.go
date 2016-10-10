package user

import (
	"context"
	"net/http"

	"github.com/jeremyletang/babakoto_api/ctxext"
	"github.com/jeremyletang/babakoto_api/dao"
	"github.com/jeremyletang/babakoto_api/domain"
	"github.com/jeremyletang/babakoto_api/jsend"
	"github.com/jeremyletang/babakoto_api/utils"
	"github.com/jinzhu/gorm"
)

func AddUserInfoToContext(
	f func(context.Context, http.ResponseWriter, *http.Request),
	db *gorm.DB,
) func(context.Context, http.ResponseWriter, *http.Request) {
	return func(ctx context.Context, w http.ResponseWriter, r *http.Request) {
		tokenDao := dao.NewAccessTokenDao(db)
		userDao := dao.NewUserDao(db)
		signupDao := dao.NewUserSignupVerificationDao(db)
		var tokenString string
		var token domain.AccessToken
		var user domain.User
		var err error
		var ok bool

		// get the token string from the context
		if tokenString, ok = ctxext.ExtractAccessTokenString(ctx); !ok {
			utils.WriteJsonResponse(w, http.StatusBadRequest,
				jsend.FailWithName("Missing access token", "access_token"))
			return
		}

		// get the token first
		if token, err = tokenDao.GetById(tokenString); err != nil {
			utils.WriteJsonResponse(w, http.StatusBadRequest,
				jsend.FailWithName("Invalid access token", "access_token"))
		}

		// then get the user from the token userid
		if user, err = userDao.GetById(token.UserId); err != nil {
			utils.WriteJsonResponse(w, http.StatusBadRequest,
				jsend.FailWithName("Invalid access token (no user linked)", "access_token"))
		}

		// then check if the user validated his account, or return failure now
		if _, err = signupDao.GetByUserId(user.Id); err == nil {
			utils.WriteJsonResponse(w, http.StatusBadRequest,
				jsend.FailWithName("Cannot use a non verified user", "access_token"))
			return
		}

		ctx = ctxext.AddUser(ctx, user)
		ctx = ctxext.AddAccessToken(ctx, token)

		// call the final handler
		f(ctx, w, r)
	}
}
