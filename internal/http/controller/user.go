package controller

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/riskiramdan/evermos/internal/appcontext"
	"github.com/riskiramdan/evermos/internal/data"
	"github.com/riskiramdan/evermos/internal/http/response"
	"github.com/riskiramdan/evermos/internal/types"
	"github.com/riskiramdan/evermos/internal/user"
)

// UserController represents the user controller
// swagger:ignore
type UserController struct {
	userService user.ServiceInterface
	dataManager *data.Manager
}

// UserList user list and count
type UserList struct {
	Data  []*user.User `json:"data"`
	Count int          `json:"count"`
}

// Login swagger:operation POST /v1/login Users Login
//
// Login.
//
// ---
// parameters:
// - name: body
//   in: body
//   required: true
//   schema:
//     $ref: "#/definitions/LoginParams"
// responses:
//   200:
//     description: "Ok"
//     schema:
//       $ref: "#/definitions/LoginResponse"
//   default:
//     description: "Error"
//     schema:
//       $ref: "#/definitions/ErrorResponse"
//
func (a *UserController) Login(w http.ResponseWriter, r *http.Request) {
	var err *types.Error

	decoder := json.NewDecoder(r.Body)

	var params user.LoginParams
	errDecode := decoder.Decode(&params)
	if errDecode != nil {
		err = &types.Error{
			Path:    ".UserController->Login()",
			Message: errDecode.Error(),
			Error:   errDecode,
			Type:    "golang-error",
		}
		response.Error(w, "Bad Request", http.StatusBadRequest, *err)
		return
	}

	var sess *user.LoginResponse
	errTransaction := a.dataManager.RunInTransaction(r.Context(), func(ctx context.Context) error {
		sess, err = a.userService.Login(r.Context(), params.Email, params.Password)
		if err != nil {
			return err.Error
		}
		return nil
	})
	if errTransaction != nil {
		err.Path = ".UserController->Login()" + err.Path
		if err.Error == user.ErrWrongPassword || err.Error == data.ErrNotFound {
			response.Error(w, "Email / password is wrong", http.StatusBadRequest, *err)
		} else {
			response.Error(w, "Internal Server Error", http.StatusInternalServerError, *err)
		}
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:  "sessionId",
		Value: sess.SessionID,
	})

	response.JSON(w, http.StatusOK, sess)
}

// ChangePassword swagger:operation POST /v1/login Users ChangePassword
//
// Change the user password.
//
// ---
// parameters:
// - name: body
//   in: body
//   required: true
//   schema:
//     $ref: "#/definitions/ChangePasswordParams"
// responses:
//   200:
//     description: "Ok"
//   default:
//     description: "Error"
//     schema:
//       $ref: "#/definitions/ErrorResponse"
//
func (a *UserController) ChangePassword(w http.ResponseWriter, r *http.Request) {
	var err *types.Error

	decoder := json.NewDecoder(r.Body)
	var params user.ChangePasswordParams
	errDecode := decoder.Decode(&params)
	if errDecode != nil {
		err = &types.Error{
			Path:    ".UserController->ChangePassword()",
			Message: errDecode.Error(),
			Error:   errDecode,
			Type:    "golang-error",
		}
		response.Error(w, "Bad Request", http.StatusBadRequest, *err)
		return
	}

	userID := appcontext.UserID(r.Context())

	err = a.userService.ChangePassword(r.Context(), userID, params.OldPassword, params.NewPassword)
	if err != nil {
		err.Path = ".UserController->ChangePassword()" + err.Path
		if err.Error == user.ErrWrongPassword {
			response.Error(w, "Wrong old password", http.StatusBadRequest, *err)
		} else {
			response.Error(w, "Internal Server Error", http.StatusInternalServerError, *err)
		}
		return
	}

	response.JSON(w, http.StatusNoContent, "")
}

// UpdateUser swagger:operation PUT /v1/users Users UpdateUser
//
// Update user.
//
// ---
// parameters:
// - name: body
//   in: body
//   required: true
//   schema:
//     $ref: "#/definitions/ChangePasswordParams"
// responses:
//   200:
//     description: "Ok"
//     schema:
//       $ref: "#/definitions/Profile"
//   default:
//     description: "Error"
//     schema:
//       $ref: "#/definitions/ErrorResponse"
// security:
//   - bearerAuth: []
//
func (a *UserController) UpdateUser(w http.ResponseWriter, r *http.Request) {
	var err *types.Error

	decoder := json.NewDecoder(r.Body)

	var params *user.UpdateUserParams
	errDecode := decoder.Decode(&params)
	if errDecode != nil {
		err = &types.Error{
			Path:    ".UserController->UpdateUser()",
			Message: errDecode.Error(),
			Error:   errDecode,
			Type:    "golang-error",
		}
		response.Error(w, "Bad Request", http.StatusBadRequest, *err)
		return
	}
	var sUserID = chi.URLParam(r, "userId")
	userID, errConversion := strconv.Atoi(sUserID)
	if errConversion != nil {
		err = &types.Error{
			Path:    ".UserController->UpdateUser()",
			Message: errConversion.Error(),
			Error:   errConversion,
			Type:    "golang-error",
		}
		response.Error(w, "Bad Request", http.StatusBadRequest, *err)
		return
	}

	var singleUser *user.User
	errTransaction := a.dataManager.RunInTransaction(r.Context(), func(ctx context.Context) error {
		singleUser, err = a.userService.UpdateUser(ctx, userID, params)
		if err != nil {
			return err.Error
		}
		return nil
	})
	if errTransaction != nil {
		err.Path = ".UserController->UpdateUser()" + err.Path
		if errTransaction == user.ErrEmailAlreadyExists {
			response.Error(w, "Alamat Email Sudah Terdaftar", http.StatusUnprocessableEntity, *err)
		} else {
			response.Error(w, "Internal Server Error", http.StatusInternalServerError, *err)
		}
		return
	}
	response.JSON(w, http.StatusOK, singleUser)

}

// CreateUser swagger:operation POST /v1/users Users CreateUser
//
// Create user.
//
// ---
// parameters:
// - name: body
//   in: body
//   required: true
//   schema:
//     $ref: "#/definitions/ChangePasswordParams"
// responses:
//   200:
//     description: "Ok"
//     schema:
//       $ref: "#/definitions/Profile"
//   default:
//     description: "Error"
//     schema:
//       $ref: "#/definitions/ErrorResponse"
// security:
//   - bearerAuth: []
//
func (a *UserController) CreateUser(w http.ResponseWriter, r *http.Request) {
	var err *types.Error

	decoder := json.NewDecoder(r.Body)

	var params *user.CreateUserParams
	errDecode := decoder.Decode(&params)
	if errDecode != nil {
		err = &types.Error{
			Path:    ".UserController->CreateUser()",
			Message: errDecode.Error(),
			Error:   errDecode,
			Type:    "golang-error",
		}
		response.Error(w, "Bad Request", http.StatusBadRequest, *err)
		return
	}

	errTransaction := a.dataManager.RunInTransaction(r.Context(), func(ctx context.Context) error {
		_, err = a.userService.CreateUser(ctx, params)
		if err != nil {
			return err.Error
		}
		return nil
	})
	if errTransaction != nil {
		err.Path = ".UserController->CreateUser()" + err.Path
		if errTransaction == user.ErrEmailAlreadyExists {
			response.Error(w, "Alamat Email Sudah Terdaftar", http.StatusUnprocessableEntity, *err)
		} else {
			response.Error(w, "Internal Server Error", http.StatusInternalServerError, *err)
		}

		return
	}

	response.JSON(w, http.StatusOK, "User Created Successfully")
}

// DeleteUser swagger:operation POST /v1/users Users DeleteUser
//
// Delete user.
//
// ---
// parameters:
// - name: body
//   in: body
//   required: true
//   schema:
//     $ref: "#/definitions/ChangePasswordParams"
// responses:
//   200:
//     description: "Ok"
//     schema:
//       $ref: "#/definitions/Profile"
//   default:
//     description: "Error"
//     schema:
//       $ref: "#/definitions/ErrorResponse"
// security:
//   - bearerAuth: []
//
func (a *UserController) DeleteUser(w http.ResponseWriter, r *http.Request) {
	var err *types.Error
	var sUserID = chi.URLParam(r, "userId")
	userID, errConversion := strconv.Atoi(sUserID)
	if errConversion != nil {
		err = &types.Error{
			Path:    ".UserController->DeleteUser()",
			Message: errConversion.Error(),
			Error:   errConversion,
			Type:    "golang-error",
		}
		response.Error(w, "Bad Request", http.StatusBadRequest, *err)
		return
	}

	errTransaction := a.dataManager.RunInTransaction(r.Context(), func(ctx context.Context) error {
		err = a.userService.DeleteUser(ctx, userID)
		if err != nil {
			return err.Error
		}
		return nil
	})
	if errTransaction != nil {
		err.Path = ".USerController->DeleteUser()" + err.Path
		response.Error(w, "Internal Server Error", http.StatusInternalServerError, *err)
		return
	}
	response.JSON(w, http.StatusNoContent, "")

}

// ListUser swagger:operation POST /v1/users Users ListUser
//
// List user.
//
// ---
// parameters:
// - name: body
//   in: body
//   required: true
//   schema:
//     $ref: "#/definitions/ChangePasswordParams"
// responses:
//   200:
//     description: "Ok"
//     schema:
//       $ref: "#/definitions/Profile"
//   default:
//     description: "Error"
//     schema:
//       $ref: "#/definitions/ErrorResponse"
// security:
//   - bearerAuth: []
//
func (a *UserController) ListUser(w http.ResponseWriter, r *http.Request) {
	var err *types.Error

	queryValues := r.URL.Query()
	var limit = 10
	var errConversion error
	if queryValues.Get("limit") != "" {
		limit, errConversion = strconv.Atoi(queryValues.Get("limit"))
		if errConversion != nil {
			err = &types.Error{
				Path:    ".UserController->ListUser()",
				Message: errConversion.Error(),
				Error:   errConversion,
				Type:    "golang-error",
			}
			response.Error(w, "Bad Request", http.StatusBadRequest, *err)
			return
		}
	}

	var page = 1
	if queryValues.Get("page") != "" {
		page, errConversion = strconv.Atoi(queryValues.Get("page"))
		if errConversion != nil {
			err = &types.Error{
				Path:    ".UserController->ListUser()",
				Message: errConversion.Error(),
				Error:   errConversion,
				Type:    "golang-error",
			}
			response.Error(w, "Bad Request", http.StatusBadRequest, *err)
			return
		}
	}

	var search = queryValues.Get("search")

	if limit < 0 {
		limit = 10
	}
	if page < 0 {
		page = 1
	}
	userList, count, err := a.userService.ListUsers(r.Context(), &user.FindAllUsersParams{
		Limit:  limit,
		Search: search,
		Page:   page,
	})
	if err != nil {
		err.Path = ".UserController->ListUser()" + err.Path
		if err.Error != data.ErrNotFound {
			response.Error(w, "Internal Server Error", http.StatusInternalServerError, *err)
			return
		}
	}
	if userList == nil {
		userList = []*user.User{}
	}

	response.JSON(w, http.StatusOK, UserList{
		Data:  userList,
		Count: count,
	})
}

// Logout swagger:operation GET /v1/logout Users Logout
//
// Logout from the system.
//
// ---
// responses:
//   200:
//     description: "Ok"
//   default:
//     description: "Error"
//     schema:
//       $ref: "#/definitions/ErrorResponse"
// security:
//   - bearerAuth: []
//
func (a *UserController) Logout(w http.ResponseWriter, r *http.Request) {
	var err *types.Error

	// get token from the context
	// log it out!
	loginToken, ok := r.Context().Value(appcontext.KeySessionID).(string)
	if !ok {
		errUserID := errors.New("failed to get user id from request context")
		response.Error(w, "Internal Server Error", http.StatusInternalServerError, types.Error{
			Path:    ".UserController->Logout()",
			Message: errUserID.Error(),
			Error:   errUserID,
			Type:    "golang-error",
		})
		return
	}

	err = a.userService.Logout(r.Context(), loginToken)
	if err != nil {
		err.Path = ".UserController->Logout()" + err.Path
		response.Error(w, "Internal Server Error", http.StatusInternalServerError, *err)
		return
	}

	response.JSON(w, http.StatusNoContent, "")
}

// NewUserController creates a new user controller
func NewUserController(
	userService user.ServiceInterface,
	dataManager *data.Manager,
) *UserController {
	return &UserController{
		userService: userService,
		dataManager: dataManager,
	}
}
