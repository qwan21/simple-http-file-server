package internal

import (
	"crypto/md5"
	"encoding/hex"
	"io"
	"log"
	"mime/multipart"
	"os"
	"path/filepath"
)

func (s *Server) upload(input multipart.File, handler *multipart.FileHeader) (string, error) {

	log.Printf("Uploaded File: %+v\nFile Size: %+v\nMIME Header: %+v\n", handler.Filename, handler.Size, handler.Header)
	hash := md5.New()

	if _, err := io.Copy(hash, input); err != nil {
		return "", err
	}

	md5 := hex.EncodeToString(hash.Sum(nil))

	subDir := md5[:2]

	_ = os.Mkdir(s.config.Directory+s.config.Route, os.ModeDir)
	_ = os.Mkdir(s.config.Directory+s.config.Route+subDir, os.ModeDir)

	outPath := filepath.Join(s.config.Directory+s.config.Route+subDir, filepath.Base(md5))

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
