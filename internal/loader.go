package internal

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"os"
	"path/filepath"
)

func (s *Server) upload(input multipart.File, handler *multipart.FileHeader) (string, error) {

	log.Printf("Uploaded File: %+v\r\nFile Size: %+v\r\nMIME Header: %+v\r\n", handler.Filename, handler.Size, handler.Header)
	hash := md5.New()

	if _, err := io.Copy(hash, input); err != nil {
		return "", err
	}

	md5 := hex.EncodeToString(hash.Sum(nil))

	subDir := md5[:2]

	rootDir := s.config.Directory + s.config.Route
	fileDir := rootDir + subDir

	outPath := filepath.Join(fileDir, filepath.Base(md5))

	if _, err := os.Stat(outPath); !os.IsNotExist(err) {
		// file exists

		return "", fmt.Errorf("This file exist")
	}

	if _, err := os.Stat(rootDir); os.IsNotExist(err) {
		//path does not exist
		log.Printf("Path %s doesn't exsist, creating path...\r\n", rootDir)
		_ = os.Mkdir(rootDir, os.ModePerm)
	}

	if _, err := os.Stat(fileDir); os.IsNotExist(err) {
		//path does not exist
		log.Printf("Path %s doesn't exsist, creating path...\r\n", fileDir)
		_ = os.Mkdir(fileDir, os.ModePerm)
	}

	out, err := os.OpenFile(outPath, os.O_CREATE|os.O_WRONLY, 0600)
	if err != nil {
		return "", err
	}

	defer out.Close()

	if _, err := input.Seek(0, io.SeekStart); err != nil {
		return "", err
	}

	if _, err := io.Copy(out, input); err != nil {
		return "", err
	}

	return md5, nil
}
