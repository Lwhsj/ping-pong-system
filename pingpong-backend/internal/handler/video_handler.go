package handler

import (
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	"pingpong-backend/internal/config"

	"github.com/gin-gonic/gin"
)

const maxVideoChunkSize = 1024 * 1024

var unsafeFileNameChars = regexp.MustCompile(`[^a-zA-Z0-9._-]`)

type VideoHandler struct {
	uploadDir   string
	maxUploadMB int64
}

func NewVideoHandler(cfg config.Config) (*VideoHandler, error) {
	uploadDir := filepath.Clean(cfg.UploadDir)
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		return nil, err
	}
	return &VideoHandler{uploadDir: uploadDir, maxUploadMB: cfg.MaxUploadMB}, nil
}

func (h *VideoHandler) UploadVideo(c *gin.Context) {
	if h.maxUploadMB > 0 {
		c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, h.maxUploadMB*1024*1024)
	}

	file, header, err := c.Request.FormFile("file")
	if err != nil {
		c.String(http.StatusBadRequest, "File is empty")
		return
	}
	defer file.Close()

	if isEmptyFile(header, file) {
		c.String(http.StatusBadRequest, "File is empty")
		return
	}

	matchID := c.PostForm("matchId")
	fileName := header.Filename
	if fileName == "" {
		fileName = fmt.Sprintf("video_%s_%d.webm", matchID, time.Now().UnixMilli())
	}
	fileName = unsafeFileNameChars.ReplaceAllString(fileName, "_")
	targetPath := filepath.Join(h.uploadDir, fileName)

	out, err := os.Create(targetPath)
	if err != nil {
		c.String(http.StatusInternalServerError, "Could not upload file: "+err.Error())
		return
	}
	defer out.Close()

	if _, err := io.Copy(out, file); err != nil {
		c.String(http.StatusInternalServerError, "Could not upload file: "+err.Error())
		return
	}
	c.String(http.StatusOK, fileName)
}

func (h *VideoHandler) StreamVideo(c *gin.Context) {
	fileName := c.Param("fileName")
	cleanName := filepath.Base(fileName)
	if cleanName != fileName || cleanName == "." || cleanName == string(filepath.Separator) {
		c.Status(http.StatusBadRequest)
		return
	}

	filePath := filepath.Join(h.uploadDir, cleanName)
	info, err := os.Stat(filePath)
	if err != nil || info.IsDir() {
		c.Status(http.StatusNotFound)
		return
	}

	file, err := os.Open(filePath)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}
	defer file.Close()

	contentLength := info.Size()
	contentType := videoContentType(cleanName)
	c.Header("Accept-Ranges", "bytes")
	c.Header("Content-Type", contentType)

	rangeHeader := c.GetHeader("Range")
	if rangeHeader == "" {
		c.Header("Content-Length", strconv.FormatInt(contentLength, 10))
		c.Status(http.StatusOK)
		_, _ = io.Copy(c.Writer, file)
		return
	}

	start, end, ok := parseRange(rangeHeader, contentLength)
	if !ok {
		c.Header("Content-Range", fmt.Sprintf("bytes */%d", contentLength))
		c.Status(http.StatusRequestedRangeNotSatisfiable)
		return
	}

	chunkLength := end - start + 1
	if chunkLength > maxVideoChunkSize {
		chunkLength = maxVideoChunkSize
		end = start + chunkLength - 1
	}

	c.Header("Content-Length", strconv.FormatInt(chunkLength, 10))
	c.Header("Content-Range", fmt.Sprintf("bytes %d-%d/%d", start, end, contentLength))
	c.Status(http.StatusPartialContent)
	_, _ = file.Seek(start, io.SeekStart)
	_, _ = io.CopyN(c.Writer, file, chunkLength)
}

func isEmptyFile(header *multipart.FileHeader, file multipart.File) bool {
	if header.Size == 0 {
		return true
	}
	if seeker, ok := file.(io.Seeker); ok {
		current, _ := seeker.Seek(0, io.SeekCurrent)
		end, err := seeker.Seek(0, io.SeekEnd)
		_, _ = seeker.Seek(current, io.SeekStart)
		return err == nil && end == 0
	}
	return false
}

func videoContentType(fileName string) string {
	lower := strings.ToLower(fileName)
	switch {
	case strings.HasSuffix(lower, ".webm"):
		return "video/webm"
	case strings.HasSuffix(lower, ".mp4"):
		return "video/mp4"
	default:
		return "application/octet-stream"
	}
}

func parseRange(header string, contentLength int64) (int64, int64, bool) {
	if !strings.HasPrefix(header, "bytes=") || contentLength <= 0 {
		return 0, 0, false
	}
	value := strings.TrimPrefix(header, "bytes=")
	parts := strings.SplitN(value, "-", 2)
	if len(parts) != 2 {
		return 0, 0, false
	}

	var start, end int64
	var err error
	if parts[0] == "" {
		suffix, err := strconv.ParseInt(parts[1], 10, 64)
		if err != nil || suffix <= 0 {
			return 0, 0, false
		}
		if suffix > contentLength {
			suffix = contentLength
		}
		return contentLength - suffix, contentLength - 1, true
	}

	start, err = strconv.ParseInt(parts[0], 10, 64)
	if err != nil || start < 0 || start >= contentLength {
		return 0, 0, false
	}
	if parts[1] == "" {
		end = contentLength - 1
	} else {
		end, err = strconv.ParseInt(parts[1], 10, 64)
		if err != nil || end < start {
			return 0, 0, false
		}
		if end >= contentLength {
			end = contentLength - 1
		}
	}
	return start, end, true
}
