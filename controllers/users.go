package controllers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/iraqnroll/gochan/context"
	"github.com/iraqnroll/gochan/models"
	"github.com/iraqnroll/gochan/views"
)

type Users struct {
	Templates struct {
		Login  Template
		Create Template
		Admin  Template
		Users  Template
		Boards Template
	}
	PageData views.BasePageData

	UserService    *models.UserService
	SessionService *models.SessionService
	BoardService   *models.BoardService
}

type UserMiddleware struct {
	SessionService *models.SessionService
}

func (umw UserMiddleware) SetUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenCookie, err := readCookie(r, CookieSession)
		if err != nil {
			fmt.Println(err)
			next.ServeHTTP(w, r)
			return
		}

		user, err := umw.SessionService.GetUser(tokenCookie)
		if err != nil {
			fmt.Println(err)
			next.ServeHTTP(w, r)
			return
		}

		//If we get to this point, we have an active user - store it in the context.
		ctx := r.Context()
		ctx = context.WithUser(ctx, user)
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
}

func (umw UserMiddleware) RequireUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := context.User(r.Context())
		if user == nil {
			http.Redirect(w, r, "/", http.StatusFound)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func (u Users) LoginForm(w http.ResponseWriter, r *http.Request) {
	u.Templates.Login.Execute(w, r, u.PageData)
}

func (u Users) CreateForm(w http.ResponseWriter, r *http.Request) {
	u.Templates.Create.Execute(w, r, u.PageData)
}

func (u Users) AdminForm(w http.ResponseWriter, r *http.Request) {
	//user := context.User(r.Context())
	u.Templates.Admin.Execute(w, r, u.PageData)
}

func (u Users) UsersForm(w http.ResponseWriter, r *http.Request) {
	data, err := u.UserService.GetUserList()
	if err != nil {
		http.Error(w, "Unable to fetch user list"+err.Error(), http.StatusInternalServerError)
	}

	u.PageData.PageData = data
	u.Templates.Users.Execute(w, r, u.PageData)
}

func (u Users) BoardsForm(w http.ResponseWriter, r *http.Request) {
	data, err := u.BoardService.GetAdminBoardList()
	if err != nil {
		http.Error(w, "Unable to fetch board list"+err.Error(), http.StatusInternalServerError)
	}

	u.PageData.PageData = data
	u.Templates.Boards.Execute(w, r, u.PageData)
}

func (u Users) Login(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Unable to parse form submission.", http.StatusBadRequest)
	}

	var data struct {
		username string
		password string
	}
	data.username = r.FormValue("username")
	data.password = r.FormValue("password")

	//Try to authenticate
	user, err := u.UserService.Authenticate(data.username, data.password)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Something went wrong.", http.StatusInternalServerError)
		return
	}

	//Create a new active session on login
	session, err := u.SessionService.CreateSession(user.Id)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Something went horribly wrong...", http.StatusInternalServerError)
		return
	}

	//Authentication successful, add the Set-Cookie header before writting the response
	setCookie(w, CookieSession, session.Token)
	http.Redirect(w, r, "/admin", http.StatusFound)
}

func (u Users) Logout(w http.ResponseWriter, r *http.Request) {
	token, err := readCookie(r, CookieSession)
	if err != nil {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	err = u.SessionService.DeleteSession(token)
	if err != nil {
		fmt.Println("Logout failed : %w", err)
		http.Error(w, "Something went horribly wrong...", http.StatusInternalServerError)
		return
	}

	//Delete the session cookie and redirect to homepage
	deleteCookie(w, CookieSession)
	http.Redirect(w, r, "/", http.StatusFound)
}

func (u Users) Create(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")
	email := r.FormValue("email")
	password := r.FormValue("password")
	user_type, usert_err := strconv.Atoi(r.FormValue("user_type"))
	if usert_err != nil {
		http.Error(w, "Invalid user_type entry, only numeric values are allowed.", http.StatusBadRequest)
		return
	}

	user, err := u.UserService.Create(username, password, email, user_type)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Something went horribly wrong....", http.StatusInternalServerError)
		return
	}
	fmt.Println(w, "User created: %+v", user)

	http.Redirect(w, r, "/admin/users", http.StatusFound)
}

func (u Users) Delete(w http.ResponseWriter, r *http.Request) {
	userId, err := strconv.Atoi(r.FormValue("userId"))
	if err != nil {
		http.Error(w, "Invalid userId provided, only numeric values are allowed.", http.StatusBadRequest)
		return
	}

	err = u.UserService.Delete(userId)
	if err != nil {
		http.Error(w, "Failed to delete user", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/admin/users", http.StatusFound)

}

func (u Users) CreateBoard(w http.ResponseWriter, r *http.Request) {
	//TODO: Restrict board and user management only to specific user type.
	user := context.User(r.Context())
	if user == nil {
		http.Redirect(w, r, "/admin/boards", http.StatusFound)
		return
	}

	uri := r.FormValue("uri")
	name := r.FormValue("name")
	description := r.FormValue("description")

	board, err := u.BoardService.Create(uri, name, description, user.Id)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Something went horribly wrong....", http.StatusInternalServerError)
		return
	}
	fmt.Println(w, "New board created : %s", board.Uri)
	http.Redirect(w, r, "/admin/boards", http.StatusFound)
}

func (u Users) DeleteBoard(w http.ResponseWriter, r *http.Request) {
	boardId, err := strconv.Atoi(r.FormValue("boardId"))
	if err != nil {
		http.Error(w, "Invalid boardId provided, only numeric values are allowed.", http.StatusBadRequest)
		return
	}

	boardUri := r.FormValue("boardUri")

	err = u.BoardService.Delete(boardId, boardUri)
	if err != nil {
		http.Error(w, "Failed to delete board", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/admin/boards", http.StatusFound)

}
