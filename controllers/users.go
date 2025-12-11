

// func (u Users) AdminForm(w http.ResponseWriter, r *http.Request) {
// 	//user := context.User(r.Context())
// 	u.Templates.Admin.Execute(w, r, u.PageData)
// }

// func (u Users) UsersForm(w http.ResponseWriter, r *http.Request) {
// 	data, err := u.UserService.GetUserList()
// 	if err != nil {
// 		http.Error(w, "Unable to fetch user list"+err.Error(), http.StatusInternalServerError)
// 	}

// 	u.PageData.PageData = data
// 	u.Templates.Users.Execute(w, r, u.PageData)
// }

// func (u Users) BoardsForm(w http.ResponseWriter, r *http.Request) {
// 	data, err := u.BoardService.GetAdminBoardList()
// 	if err != nil {
// 		http.Error(w, "Unable to fetch board list"+err.Error(), http.StatusInternalServerError)
// 	}

// 	u.PageData.PageData = data
// 	u.Templates.Boards.Execute(w, r, u.PageData)
// }

// func (u Users) Logout(w http.ResponseWriter, r *http.Request) {
// 	token, err := readCookie(r, CookieSession)
// 	if err != nil {
// 		http.Redirect(w, r, "/", http.StatusFound)
// 		return
// 	}

// 	err = u.SessionService.DeleteSession(token)
// 	if err != nil {
// 		fmt.Println("Logout failed : %w", err)
// 		http.Error(w, "Something went horribly wrong...", http.StatusInternalServerError)
// 		return
// 	}

// 	//Delete the session cookie and redirect to homepage
// 	deleteCookie(w, CookieSession)
// 	http.Redirect(w, r, "/", http.StatusFound)
// }

// func (u Users) Create(w http.ResponseWriter, r *http.Request) {
// 	username := r.FormValue("username")
// 	email := r.FormValue("email")
// 	password := r.FormValue("password")
// 	user_type, usert_err := strconv.Atoi(r.FormValue("user_type"))
// 	if usert_err != nil {
// 		http.Error(w, "Invalid user_type entry, only numeric values are allowed.", http.StatusBadRequest)
// 		return
// 	}

// 	user, err := u.UserService.Create(username, password, email, user_type)
// 	if err != nil {
// 		fmt.Println(err)
// 		http.Error(w, "Something went horribly wrong....", http.StatusInternalServerError)
// 		return
// 	}
// 	fmt.Println(w, "User created: %+v", user)

// 	http.Redirect(w, r, "/admin/users", http.StatusFound)
// }

// func (u Users) Delete(w http.ResponseWriter, r *http.Request) {
// 	userId, err := strconv.Atoi(r.FormValue("userId"))
// 	if err != nil {
// 		http.Error(w, "Invalid userId provided, only numeric values are allowed.", http.StatusBadRequest)
// 		return
// 	}

// 	err = u.UserService.Delete(userId)
// 	if err != nil {
// 		http.Error(w, "Failed to delete user", http.StatusInternalServerError)
// 		return
// 	}

// 	http.Redirect(w, r, "/admin/users", http.StatusFound)

// }

// func (u Users) CreateBoard(w http.ResponseWriter, r *http.Request) {
// 	//TODO: Restrict board and user management only to specific user type.
// 	user := context.User(r.Context())
// 	if user == nil {
// 		http.Redirect(w, r, "/admin/boards", http.StatusFound)
// 		return
// 	}

// 	uri := r.FormValue("uri")
// 	name := r.FormValue("name")
// 	description := r.FormValue("description")

// 	board, err := u.BoardService.Create(uri, name, description, user.Id)
// 	if err != nil {
// 		fmt.Println(err)
// 		http.Error(w, "Something went horribly wrong....", http.StatusInternalServerError)
// 		return
// 	}
// 	fmt.Println(w, "New board created : %s", board.Uri)
// 	http.Redirect(w, r, "/admin/boards", http.StatusFound)
// }

// func (u Users) DeleteBoard(w http.ResponseWriter, r *http.Request) {
// 	boardId, err := strconv.Atoi(r.FormValue("boardId"))
// 	if err != nil {
// 		http.Error(w, "Invalid boardId provided, only numeric values are allowed.", http.StatusBadRequest)
// 		return
// 	}

// 	boardUri := r.FormValue("boardUri")

// 	err = u.BoardService.Delete(boardId, boardUri)
// 	if err != nil {
// 		http.Error(w, "Failed to delete board", http.StatusInternalServerError)
// 		return
// 	}

// 	http.Redirect(w, r, "/admin/boards", http.StatusFound)

// }
