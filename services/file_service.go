package services

import (
	"bytes"
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

type FileService struct {
}

func NewFileService() *FileService {
	return &FileService{}
}

// Create a directory to store board-specific static content
func (fs *FileService) CreateBoardStatic(board_uri string) error {
	parent_path := filepath.Join(".", "static", board_uri)

	//Create a directory for board specific banners
	path := filepath.Join(parent_path, "banners")
	err := os.MkdirAll(path, 0755)
	if err != nil {
		return fmt.Errorf("FileService.CreateBoardStatic failed :%w", err)
	}

	//Create a directory for thread content
	path = filepath.Join(parent_path, "src")
	err = os.Mkdir(path, 0755)
	if err != nil {
		return fmt.Errorf("FileService.CreateBoardStatic failed :%w", err)
	}

	return nil
}

// Deletes a board-specific directory for static content (recursive)
func (fs *FileService) RemoveBoardStatic(board_uri string) error {
	path := filepath.Join(".", "static", board_uri)
	err := os.RemoveAll(path)
	if err != nil {
		return fmt.Errorf("FileService.RemoveBoardStatic failed :%w", err)
	}

	return nil
}

// Returns a uri for a randomly-picked board-specific banner
func (fs *FileService) GetBoardBannerUri(board_uri string) (string, error) {
	parent_path := filepath.Join(".", "static", board_uri, "banners")
	banners, err := os.ReadDir(parent_path)
	if err != nil {
		return "", fmt.Errorf("FileService.GetBoardBannerUri failed : %w", err)
	}

	if len(banners) > 0 {
		result := banners[rand.Intn(len(banners))]

		return filepath.Join(parent_path, result.Name()), nil
	}

	return "", nil
}

// Handles writing uploaded files to disk and generating thumbnails.
// TODO: Refactor this garbage, make imagemagick conversion parameters configurable from config.
func (bs *FileService) HandleFileUploads(files []*multipart.FileHeader, board_uri string, thread_id, post_id int) (string, error) {
	var result string
	for i := range files {
		//fmt.Printf("Processing uploaded file : %s\n", files[i].Filename)
		file, err := files[i].Open()
		if err != nil {
			return "", err
		}
		defer file.Close()

		f_ext := path.Ext(files[i].Filename)
		src_fn, thumb_fn := createThreadStaticDirectories(board_uri, f_ext, thread_id, post_id, i)

		dst, err := os.Create(src_fn)
		if err != nil {
			return "", err
		}
		if _, err := io.Copy(dst, file); err != nil {
			return "", err
		}
		dst.Close()

		//attempt to generate a thumbnail from uploaded src image.
		//TODO: this will probably get changed once i containerize all dependencies
		m_cmd := exec.Command("magick", src_fn, "-thumbnail", "200x200", thumb_fn)
		var stderr bytes.Buffer
		m_cmd.Stderr = &stderr
		err = m_cmd.Run()
		if err != nil {
			fmt.Println(fmt.Sprint(err) + ": " + stderr.String())
			return "", err
		}

		result = result + fmt.Sprintf("%d-%d%s;", post_id, i, f_ext)
	}
	return result, nil
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
