package utils

import (
	"time"
	"os"
	"fmt"
	"io"
	"net/http"
	"path/filepath"
	"mime/multipart"
	"bytes"

	"github.com/google/uuid"

)

func SaveFile(file multipart.File, header *multipart.FileHeader) (string, error) {
	if err := os.MkdirAll("uploads", 0755); err != nil {
		return "", err
	}

	buf := make([]byte, 512)
	n, err := file.Read(buf)
	if err != nil && err != io.EOF {
		return "", err
	}

	contentType := http.DetectContentType(buf[:n])

	if contentType != "image/jpeg" &&
		contentType != "image/png" &&
		contentType != "image/gif" &&
		contentType != "image/webp" {
		return "", fmt.Errorf("invalid image type: %s", contentType)
	}

	ext := filepath.Ext(header.Filename)

	filename := fmt.Sprintf("%d_%s%s",
		time.Now().UnixNano(),
		uuid.NewString(),
		ext,
	)

	path := filepath.Join("uploads", filename)

	out, err := os.Create(path)
	if err != nil {
		return "", err
	}
	defer out.Close()

	reader := io.MultiReader(bytes.NewReader(buf[:n]), file)

	if _, err := io.Copy(out, reader); err != nil {
		return "", err
	}

	return path, nil
}

