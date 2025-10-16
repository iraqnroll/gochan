package models

import (
	"database/sql"
	"fmt"
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
	DB *sql.DB
}

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

	return &board, nil
}

func (bs *BoardService) Delete(boardId int) error {
	_, err := bs.DB.Exec(`DELETE FROM boards WHERE id = $1`, boardId)
	if err != nil {
		fmt.Println("BoardService.Delete failed : %w", err)
		return err
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
	var result BoardDto
	rows := bs.DB.QueryRow(`SELECT id, uri, name, description FROM boards WHERE uri = $1`, uri)
	err := rows.Scan(&result.Id, &result.Uri, &result.Name, &result.Description)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, err
		}

		return nil, fmt.Errorf("BoardService.GetBoard failed : %w", err)
	}

	//Since we found a valid registered board, we populate it with content :
	//TODO: Refactor this garbage, optimally we'd fetch all our board posts/threads with a single query....
	threads, err := bs.DB.Query(`
		SELECT id,
			locked,
			board_id,
			topic
		FROM threads
		WHERE board_id = $1`, result.Id)

	if err != nil {
		if err == sql.ErrNoRows {
			return &result, nil
		}
		return &result, fmt.Errorf("BoardService.GetBoard threads failed : %w", err)
	}
	defer threads.Close()

	for threads.Next() {
		var thread ThreadDto
		err = threads.Scan(&thread.Id, &thread.Locked, &thread.BoardId, &thread.Topic)

		if err != nil {
			fmt.Println("BoardService.GetBoard thread loop failed : %w", err)
		}
		result.Threads = append(result.Threads, thread)
	}

	postRows, err := bs.DB.Query(`SELECT p.id,
			p.thread_id,
			p.identifier,
			p.content,
			COALESCE(to_char(p.post_timestamp, 'YYYY-MM-DD HH24:MI:SS'), 'Never') AS post_timestamp,
			p.is_op
		FROM posts AS p
		INNER JOIN threads AS th ON th.Id = p.thread_id
		WHERE th.board_id = $1`, result.Id)

	if err != nil {
		return nil, fmt.Errorf("BoardService.GetBoard failed : %w", err)
	}
	defer postRows.Close()

	postHashMap := make(map[int][]PostDto)

	for postRows.Next() {
		var post PostDto
		err = postRows.Scan(&post.Id, &post.ThreadId, &post.Identifier, &post.Content, &post.PostTimestamp, &post.IsOP)
		if err != nil {
			fmt.Println("BoardService.GetBoard post loop failed : %w", err)
		}
		postHashMap[post.ThreadId] = append(postHashMap[post.ThreadId], post)
	}

	for i := range result.Threads {
		threadId := result.Threads[i].Id
		result.Threads[i].Posts = postHashMap[threadId]
	}

	return &result, nil
}

func (bs *BoardService) GetThread(id int) (*ThreadDto, error) {
	var result ThreadDto

	row := bs.DB.QueryRow(`
		SELECT th.id,
			th.locked,
			th.board_id,
			th.topic,
			brd.uri AS board_uri
		FROM threads AS th
		INNER JOIN boards brd ON brd.id = th.board_id
		WHERE th.id = $1`, id)

	err := row.Scan(&result.Id, &result.Locked, &result.BoardId, &result.Topic, &result.BoardUri)
	if err != nil {
		return nil, fmt.Errorf("BoardService.GetBoard threads failed : %w", err)
	}

	postRows, err := bs.DB.Query(`SELECT id,
			thread_id,
			identifier,
			content,
			COALESCE(to_char(post_timestamp, 'YYYY-MM-DD HH24:MI:SS'), 'Never') AS post_timestamp,
			is_op
		FROM posts
		WHERE thread_id= $1`, result.Id)

	if err != nil {
		return nil, fmt.Errorf("BoardService.GetThread failed : %s", err)
	}
	defer postRows.Close()

	for postRows.Next() {
		var post PostDto
		err = postRows.Scan(&post.Id, &post.ThreadId, &post.Identifier, &post.Content, &post.PostTimestamp, &post.IsOP)
		if err != nil {
			fmt.Printf("BoardService.GetThread failed : %s", err)
			continue
		}
		result.Posts = append(result.Posts, post)
	}

	return &result, nil
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

func (bs *BoardService) CheckBoard(uri string) (int, error) {
	var board_id int

	rows := bs.DB.QueryRow(`SELECT id FROM boards WHERE uri = $1`, uri)
	err := rows.Scan(&board_id)

	if err != nil {
		return -1, fmt.Errorf("BoardService.CheckBoard failed : %w", err)
	}

	return board_id, nil
}
