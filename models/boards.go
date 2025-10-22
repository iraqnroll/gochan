package models

import (
	"bytes"
	"database/sql"
	"fmt"
	"io"
	"math/rand"
	"mime/multipart"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strconv"
	"strings"
)

type Board struct {
	Id            int
	Uri           string
	Name          string
	Description   string
	Date_created  string
	Date_updated  string
	OwnerId       int
	OwnerUsername string
}

type BoardDto struct {
	Id          int
	Uri         string
	Name        string
	Description string
	Threads     []ThreadDto
}

type ThreadDto struct {
	Id       int
	Posts    []PostDto
	Topic    string
	Locked   bool
	BoardId  int
	BoardUri string
}

type PostDto struct {
	Id            int
	ThreadId      int
	Identifier    string
	Content       string
	PostTimestamp string
	IsOP          bool
}

type BoardService struct {
	DB             *sql.DB
	ImagickService *IMagickService
}

// -==============================[Admin actions]==============================-
func (bs *BoardService) Create(uri, name, description string, ownerId int) (*Board, error) {
	uri = strings.ToLower(uri)

	board := Board{
		Uri:         uri,
		Name:        name,
		Description: description,
		OwnerId:     ownerId,
	}

	row := bs.DB.QueryRow(`
		INSERT INTO boards (uri, name, description, ownerId)
		VALUES ($1, $2, $3, $4) RETURNING id, to_char(date_created, 'YYYY-MM-DD HH24:MI:SS')`, uri, name, description, ownerId)

	err := row.Scan(&board.Id, &board.Date_created)

	if err != nil {
		return nil, fmt.Errorf("BoardService.Create failed : %w", err)
	}

	//TODO: Refactor this into a separate function
	//Create a board folder to store static content.
	path := filepath.Join(".", "static", board.Uri, "banners")
	err = os.MkdirAll(path, 0755)
	if err != nil {
		return &board, fmt.Errorf("BoardService.Create failed :%w", err)
	}

	path = filepath.Join(".", "static", board.Uri, "src")
	err = os.Mkdir(path, 0755)
	if err != nil {
		return &board, fmt.Errorf("BoardService.Create failed :%w", err)
	}

	return &board, nil
}

func (bs *BoardService) Delete(boardId int, boardUri string) error {
	_, err := bs.DB.Exec(`DELETE FROM boards WHERE id = $1`, boardId)
	if err != nil {
		fmt.Println("BoardService.Delete failed : %w", err)
		return err
	}

	path := filepath.Join(".", "static", boardUri)
	err = os.RemoveAll(path)
	if err != nil {
		return fmt.Errorf("BoardService.Create failed :%w", err)
	}

	return nil
}

func (bs *BoardService) GetAdminBoardList() ([]Board, error) {
	var result []Board

	rows, err := bs.DB.Query(`SELECT
		b.id,
		b.uri,
		b.name,
		b.description,
		COALESCE(to_char(b.date_created, 'YYYY-MM-DD HH24:MI:SS'), 'Never') AS date_created,
		COALESCE(to_char(b.date_updated, 'YYYY-MM-DD HH24:MI:SS'), 'Never') AS date_updated,
		usr.username AS ownerUsername
		FROM boards AS b
		INNER JOIN users AS usr ON usr.id = b.ownerId`)

	if err != nil {
		return nil, fmt.Errorf("BoardService.GetAdminBoardList failed : %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var board Board
		err := rows.Scan(&board.Id, &board.Uri, &board.Name, &board.Description, &board.Date_created, &board.Date_updated, &board.OwnerUsername)
		if err != nil {
			fmt.Println("BoardService.GetAdminBoardList loop failed : %w", err)
		}
		result = append(result, board)
	}

	return result, nil
}

// -==============================[Global actions]==============================-
func (bs *BoardService) GetBoardList() ([]BoardDto, error) {
	var result []BoardDto

	rows, err := bs.DB.Query(`SELECT id, uri, name, description FROM boards`)

	if err != nil {
		return nil, fmt.Errorf("BoardService.GetBoardList failed : %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var board BoardDto
		err := rows.Scan(&board.Id, &board.Uri, &board.Name, &board.Description)
		if err != nil {
			fmt.Println("BoardService.GetBoardList loop failed : %w", err)
		}
		result = append(result, board)
	}

	return result, nil
}

func (bs *BoardService) GetBoard(uri string) (*BoardDto, error) {
	//Get upper-board data, then move on to threads/posts.
	result, err := bs.GetBoardQuery(uri, nil)
	if err != nil {
		return nil, err
	}

	fmt.Printf("Board id %d \n", result.Id)

	result.Threads, err = bs.GetThreadsQuery(result.Id, nil)
	if err != nil {
		return nil, err
	}

	fmt.Printf("First thread id %d\n", result.Threads[0].Id)

	posts, err := bs.GetPostsQuery(result.Threads)
	if err != nil {
		return nil, err
	}

	bs.SortPostsIntoThreads(result.Threads, posts)

	return result, nil
}

func (bs *BoardService) GetThread(thread_id int, board_uri string) (*BoardDto, error) {
	result, err := bs.GetBoardQuery(board_uri, nil)
	if err != nil {
		return nil, err
	}

	fmt.Printf("Board id %d \n", result.Id)

	result.Threads, err = bs.GetThreadsQuery(result.Id, &thread_id)
	if err != nil {
		return nil, err
	}

	fmt.Printf("First thread id %d\n", result.Threads[0].Id)

	posts, err := bs.GetPostsQuery(result.Threads)
	if err != nil {
		return nil, err
	}

	bs.SortPostsIntoThreads(result.Threads, posts)

	return result, nil
}

func (bs *BoardService) CreateThread(board_id int, topic, identifier, content string) (*ThreadDto, error) {
	post := PostDto{
		Identifier: identifier,
		Content:    content,
		IsOP:       true,
	}

	result := ThreadDto{
		BoardId: board_id,
		Locked:  false,
		Topic:   topic,
	}

	row := bs.DB.QueryRow(`
		INSERT INTO threads(board_id, topic)
		VALUES ($1, $2) RETURNING id, to_char(date_created, 'YYYY-MM-DD HH24:MI:SS')`, board_id, topic)

	err := row.Scan(&result.Id, &post.PostTimestamp)

	if err != nil {
		return nil, fmt.Errorf("BoardService.CreateThread failed while creating a new thread : %w", err)
	}

	//We successfully created a new thread, attach the OP's post to it before returning the new thread
	row = bs.DB.QueryRow(`
		INSERT INTO posts(thread_id, identifier, content, post_timestamp, is_op)
		VALUES ($1, $2, $3, $4, $5) RETURNING id`, result.Id, post.Identifier, post.Content, post.PostTimestamp, post.IsOP)

	err = row.Scan(&post.Id)

	if err != nil {
		return nil, fmt.Errorf("BoardService.CreateThread failed while attaching OP post to new thread : %w", err)
	}

	result.Posts = append(result.Posts, post)

	return &result, nil
}

func (bs *BoardService) CreateReply(thread_id int, identifier, content string) error {
	var result int

	row := bs.DB.QueryRow(`
		INSERT INTO posts(thread_id, identifier, content)
		VALUES ($1, $2, $3) RETURNING id`, thread_id, identifier, content)

	err := row.Scan(&result)
	if err != nil {
		return fmt.Errorf("BoardService.CreateReply failed : %w", err)
	}

	return nil
}

//-==============================[Utility functions]==============================-

func (bs *BoardService) GetBoardQuery(uri string, board_id *int) (*BoardDto, error) {
	var result BoardDto
	var rows *sql.Row

	if board_id != nil {
		rows = bs.DB.QueryRow(`SELECT id, uri, name, description FROM boards WHERE board_id = $1`, *board_id)
	} else {
		rows = bs.DB.QueryRow(`SELECT id, uri, name, description FROM boards WHERE uri = $1`, uri)
	}

	err := rows.Scan(&result.Id, &result.Uri, &result.Name, &result.Description)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, err
		}

		return nil, fmt.Errorf("GetBoardQuery failed : %w", err)
	}

	return &result, err
}

// We pass thread_id as a pointer for it's ability to get here as nil, if we pass thread_id, we get a single thread, if its nil we fetch threads for the whole board.
// TODO: Refactor this garbage, optimally we'd fetch all our board posts/threads with a single query....
func (bs *BoardService) GetThreadsQuery(board_id int, thread_id *int) ([]ThreadDto, error) {
	var result []ThreadDto

	//Retrieve specific thread for a specific board
	if thread_id != nil {
		var thread ThreadDto
		row := bs.DB.QueryRow(`
		SELECT id,
			locked,
			board_id,
			topic
		FROM threads
		WHERE board_id = $1 AND id = $2`, board_id, *thread_id)

		err := row.Scan(&thread.Id, &thread.Locked, &thread.BoardId, &thread.Topic)
		if err != nil {
			return nil, fmt.Errorf("GetThreadsQuery failed : %w", err)
		}

		result = append(result, thread)
		return result, nil
	} else {
		threads, err := bs.DB.Query(`
		SELECT id,
			locked,
			board_id,
			topic
		FROM threads
		WHERE board_id = $1`, board_id)

		if err != nil {
			if err == sql.ErrNoRows { //Why the fuck is 0 rows result considered an error ?????
				return result, nil
			}
			return nil, fmt.Errorf("GetThreadsQuery failed : %w", err)
		}
		defer threads.Close()

		for threads.Next() {
			var thread ThreadDto
			err = threads.Scan(&thread.Id, &thread.Locked, &thread.BoardId, &thread.Topic)

			if err != nil {
				fmt.Println("BoardService.GetBoard thread loop failed : %w", err)
			}

			result = append(result, thread)
		}

		return result, nil
	}
}

func (bs *BoardService) GetPostsQuery(threads []ThreadDto) ([]PostDto, error) {
	var result []PostDto

	for _, thread := range threads {
		threadPostRows, err := bs.DB.Query(`SELECT p.id,
			p.thread_id,
			p.identifier,
			p.content,
			COALESCE(to_char(p.post_timestamp, 'YYYY-MM-DD HH24:MI:SS'), 'Never') AS post_timestamp,
			p.is_op
		FROM posts AS p
		INNER JOIN threads AS th ON th.id = p.thread_id
		WHERE th.id = $1`, thread.Id)

		if err != nil {
			return nil, fmt.Errorf("GetPostsQuery failed : %w", err)
		}
		defer threadPostRows.Close()

		for threadPostRows.Next() {
			var post PostDto
			err = threadPostRows.Scan(&post.Id, &post.ThreadId, &post.Identifier, &post.Content, &post.PostTimestamp, &post.IsOP)
			if err != nil {
				fmt.Println("GetPostsQuery failed : %w", err)
			}

			result = append(result, post)
		}
	}

	return result, nil
}

func (bs *BoardService) SortPostsIntoThreads(threads []ThreadDto, posts []PostDto) {
	postHashMap := make(map[int][]PostDto)

	for _, post := range posts {
		postHashMap[post.ThreadId] = append(postHashMap[post.ThreadId], post)
	}

	for i := range threads {
		threadId := threads[i].Id
		threads[i].Posts = postHashMap[threadId]
	}
}

func (bs *BoardService) CheckBoard(uri string, board_id *int) (int, error) {
	var result int
	var row *sql.Row

	if board_id != nil {
		row = bs.DB.QueryRow(`SELECT id FROM boards WHERE id = $1`, *board_id)
	} else {
		row = bs.DB.QueryRow(`SELECT id FROM boards WHERE uri = $1`, uri)
	}
	err := row.Scan(&board_id)

	if err != nil {
		return -1, fmt.Errorf("BoardService.CheckBoard failed : %w", err)
	}

	return result, nil
}

func (bs *BoardService) GetBoardBannerUri(boardUri string) (string, error) {
	path := filepath.Join(".", "static", boardUri, "banners")

	banners, err := os.ReadDir(path)
	if err != nil {
		return "", fmt.Errorf("failed to retrieve board banners %w", err)
	}

	if len(banners) > 0 {
		result := banners[rand.Intn(len(banners))]

		return filepath.Join(".", "static", boardUri, "banners", result.Name()), nil
	}

	return "", nil
}

// Handles writing uploaded files to disk and generating thumbnails.
// TODO: Refactor this garbage, make imagemagick conversion parameters configurable from config.
func (bs *BoardService) HandleFileUploads(files []*multipart.FileHeader, board_uri string, thread_id, post_id int) error {
	fmt.Printf("Received %d files for processing...\n", len(files))

	for i := range files {
		fmt.Printf("Processing uploaded file : %s\n", files[i].Filename)
		file, err := files[i].Open()

		if err != nil {
			return err
		}
		defer file.Close()

		src_fn, thumb_fn := createThreadStaticDirectories(board_uri, path.Ext(files[i].Filename), thread_id, post_id, i)

		dst, err := os.Create(src_fn)
		if err != nil {
			return err
		}
		if _, err := io.Copy(dst, file); err != nil {
			return err
		}

		dst.Close()

		//fmt.Printf("Successfully uploaded src file : %s\n thumbnail path : %s\n, src path : %s\n", files[i].Filename, thumb_fn, src_fn)

		//attempt to generate a thumbnail from uploaded src image.
		m_cmd := exec.Command("magick", src_fn, "-thumbnail", "200x200", thumb_fn)
		//fmt.Printf("Command : %s\n", m_cmd.String())
		var stderr bytes.Buffer
		m_cmd.Stderr = &stderr
		err = m_cmd.Run()
		if err != nil {
			fmt.Println(fmt.Sprint(err) + ": " + stderr.String())
			return err
		}
	}
	return nil
}

// Creates directories & returns filepaths for thumbnails and source images that were uploaded for a specific post
// TODO: only allow whitelisted and validated extensions to reach this part of code...
func createThreadStaticDirectories(board_uri, file_ext string, thread_id, post_id, f_idx int) (string, string) {
	threadSrcPath := filepath.Join(".", "static", board_uri, "src", strconv.Itoa(thread_id))
	threadThbPath := filepath.Join(".", "static", board_uri, "thumbs", strconv.Itoa(thread_id))

	err := os.MkdirAll(threadSrcPath, 0755)
	if err != nil {
		panic(err)
	}

	err = os.MkdirAll(threadThbPath, 0755)
	if err != nil {
		panic(err)
	}

	filename := strconv.Itoa(post_id) + "-" + strconv.Itoa(f_idx) + file_ext

	return filepath.Join(threadSrcPath, filename), filepath.Join(threadThbPath, filename)
}
