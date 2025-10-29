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

// type BoardService struct {
// 	boardRepo      BoardRepository
// 	threadRepo     ThreadRepository
// 	postRepo       PostRepository
// 	ImagickService *IMagickService
// }

// type BoardRepository interface {
// 	GetAll() ([]BoardDto, error)
// 	GetById(id int) (BoardDto, error)
// 	GetByUri(uri string) (BoardDto, error)
// 	GetAllForAdmin() ([]Board, error)

// 	CreateNew(uri, name, description string, owner_id int) (Board, error)
// 	Delete(id int) error
// }

// type ThreadRepository interface {
// 	CreateNew(board_id int, topic string) (ThreadDto, error)
// 	GetById(thread_id int) (ThreadDto, error)
// 	GetAllByBoard(board_id int) ([]ThreadDto, error)
// }

// type PostRepository interface {
// 	GetAllByThread(thread_id int) ([]PostDto, error)
// 	CreateNew(thread_id int, identifier, content string, is_op bool) (PostDto, error)
// }

// -==============================[Admin actions]==============================-
// func (bs *BoardService) Create(uri, name, description string, ownerId int) (*Board, error) {
// 	uri = strings.ToLower(uri)

// 	board, err := bs.boardRepo.CreateNew(uri, name, description, ownerId)

// 	if err != nil {
// 		return nil, fmt.Errorf("BoardService.Create failed : %w", err)
// 	}

// 	//TODO: Refactor this into a separate function
// 	//Create a board folder to store static content.
// 	path := filepath.Join(".", "static", board.Uri, "banners")
// 	err = os.MkdirAll(path, 0755)
// 	if err != nil {
// 		return &board, fmt.Errorf("BoardService.Create failed :%w", err)
// 	}

// 	path = filepath.Join(".", "static", board.Uri, "src")
// 	err = os.Mkdir(path, 0755)
// 	if err != nil {
// 		return &board, fmt.Errorf("BoardService.Create failed :%w", err)
// 	}

// 	return &board, nil
// }

// func (bs *BoardService) Delete(boardId int, boardUri string) error {
// 	err := bs.boardRepo.Delete(boardId)
// 	if err != nil {
// 		fmt.Println("BoardService.Delete failed : %w", err)
// 		return err
// 	}

// 	path := filepath.Join(".", "static", boardUri)
// 	err = os.RemoveAll(path)
// 	if err != nil {
// 		return fmt.Errorf("BoardService.Create failed :%w", err)
// 	}

// 	return nil
// }

// func (bs *BoardService) GetAdminBoardList() ([]Board, error) {
// 	boards, err := bs.boardRepo.GetAllForAdmin()
// 	if err != nil {
// 		return nil, fmt.Errorf("BoardService.GetAdminBoardList failed : %w", err)
// 	}

// 	return boards, nil
// }

// -==============================[Global actions]==============================-
// func (bs *BoardService) GetBoardList() ([]BoardDto, error) {
// 	boards, err := bs.boardRepo.GetAll()
// 	if err != nil {
// 		return nil, fmt.Errorf("BoardService.GetBoardList failed : %w", err)
// 	}

// 	return boards, nil
// }

// func (bs *BoardService) GetBoard(uri string) (*BoardDto, error) {
// 	//Get upper-board data, then move on to threads/posts.
// 	result, err := bs.boardRepo.GetByUri(uri)
// 	if err != nil {
// 		return nil, err
// 	}

// 	result.Threads, err = bs.threadRepo.GetAllByBoard(result.Id)
// 	if err != nil {
// 		return nil, err
// 	}

// 	posts, err := bs.GetPostsQuery(result.Threads)
// 	if err != nil {
// 		return nil, err
// 	}

// 	bs.SortPostsIntoThreads(result.Threads, posts)

// 	return &result, nil
// }

// func (bs *BoardService) GetThread(thread_id int, board_uri string) (*BoardDto, error) {
// 	result, err := bs.boardRepo.GetByUri(board_uri)
// 	if err != nil {
// 		return nil, err
// 	}

// 	thread, err := bs.threadRepo.GetById(thread_id)
// 	if err != nil {
// 		return nil, err
// 	}

// 	result.Threads = append(result.Threads, thread)

// 	posts, err := bs.postRepo.GetAllByThread(thread.Id)
// 	if err != nil {
// 		return nil, err
// 	}

// 	bs.SortPostsIntoThreads(result.Threads, posts)

// 	return &result, nil
// }

// func (bs *BoardService) CreateThread(board_id int, topic, identifier, content string) (*ThreadDto, error) {
// 	result, err := bs.threadRepo.CreateNew(board_id, topic)
// 	if err != nil {
// 		return nil, fmt.Errorf("BoardService.CreateThread failed while creating a new thread : %w", err)
// 	}

// 	//We successfully created a new thread, attach the OP's post to it before returning the new thread
// 	post, err := bs.postRepo.CreateNew(result.Id, identifier, content, true)
// 	if err != nil {
// 		return nil, fmt.Errorf("BoardService.CreateThread failed while attaching OP post to new thread : %w", err)
// 	}

// 	result.Posts = append(result.Posts, post)
// 	return &result, nil
// }

// func (bs *BoardService) CreateReply(thread_id int, identifier, content string) error {
// 	_, err := bs.postRepo.CreateNew(thread_id, identifier, content, false)
// 	if err != nil {
// 		return fmt.Errorf("BoardService.CreateReply failed : %w", err)
// 	}

// 	return nil
// }

// -==============================[Utility functions]==============================-
// We pass thread_id as a pointer for it's ability to get here as nil, if we pass thread_id, we get a single thread, if its nil we fetch threads for the whole board.

// func (bs *BoardService) GetPostsQuery(threads []ThreadDto) ([]PostDto, error) {
// 	var result []PostDto

// 	for _, thread := range threads {
// 		posts, err := bs.postRepo.GetAllByThread(thread.Id)
// 		if err != nil {
// 			return nil, fmt.Errorf("GetPostsQuery failed : %w", err)
// 		}
// 		result = append(result, posts...)
// 	}

// 	return result, nil
// }

// func (bs *BoardService) SortPostsIntoThreads(threads []ThreadDto, posts []PostDto) {
// 	postHashMap := make(map[int][]PostDto)

// 	for _, post := range posts {
// 		postHashMap[post.ThreadId] = append(postHashMap[post.ThreadId], post)
// 	}

// 	for i := range threads {
// 		threadId := threads[i].Id
// 		threads[i].Posts = postHashMap[threadId]
// 	}
// }

// TODO: refactor to use board repository and split into two functions.
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
