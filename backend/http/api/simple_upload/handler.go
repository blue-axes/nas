package simple_upload

import (
	"fmt"
	api "github.com/blue-axes/tmpl/http/api"
	"github.com/blue-axes/tmpl/pkg/constants"
	"github.com/blue-axes/tmpl/pkg/errors"
	"github.com/blue-axes/tmpl/service"
	"github.com/blue-axes/tmpl/types/api_schema"
	"github.com/labstack/echo/v4"
	"io"
	"net/http"
	"net/url"
	"os"
	"path"
	"strings"
	"time"
)

type (
	FileObjectHandler struct {
		*api.Handler
		svc     *service.Service
		rootDir string
	}
)

func New(svc *service.Service) *FileObjectHandler {
	h := &FileObjectHandler{
		Handler: api.New(svc),
		svc:     svc,
		rootDir: path.Clean(svc.Config().Nas.SimpleUploadRoot),
	}
	return h
}

func (h FileObjectHandler) Schema(c echo.Context) error {
	var (
		req struct {
			api_schema.Filename
		}
		resp struct {
			api_schema.FileInfo
		}
	)
	if err := c.Bind(&req); err != nil {
		return err
	}
	info, err := h.svc.FsStat(path.Join(h.rootDir, req.Name))
	if os.IsNotExist(err) {
		return h.RespJson(c, nil, echo.NewHTTPError(http.StatusNotFound))
	}
	fileType := "file"
	if info.IsDir() {
		fileType = "dir"
	}
	if c.Request().Method == http.MethodHead {
		header := c.Response().Header()
		header.Set("File-Name", info.Name())
		header.Set("File-Size", fmt.Sprintf("%d", info.Size()))
		header.Set("File-Type", fileType)
		return h.RespJson(c, nil, nil)
	}
	resp.Name = info.Name()
	resp.Size = uint64(info.Size())
	resp.FileType = fileType
	return h.RespJson(c, resp, nil)
}

func (h FileObjectHandler) ReadDir(c echo.Context) error {
	var (
		req struct {
			api_schema.Filename
		}
		resp struct {
			List []api_schema.FileInfo `json:"List"`
		}
	)
	if err := c.Bind(&req); err != nil {
		return err
	}
	filename := path.Join(h.rootDir, req.Name)
	info, err := h.svc.FsStat(filename)
	if os.IsNotExist(err) {
		return h.RespJson(c, nil, echo.NewHTTPError(http.StatusNotFound))
	}
	if !info.IsDir() {
		return h.RespJson(c, nil, echo.NewHTTPError(http.StatusBadRequest, "invalid dir"))
	}
	// 目录
	entry, err := h.svc.FsReadDir(filename)
	if err != nil {
		return err
	}
	for _, item := range entry {
		var (
			fileType        = "file"
			size     uint64 = 0
		)
		if item.IsDir() {
			fileType = "dir"
		}
		info, err := item.Info()
		if err == nil {
			size = uint64(info.Size())
		}
		resp.List = append(resp.List, api_schema.FileInfo{
			Name:     item.Name(),
			Size:     size,
			FileType: fileType,
		})
	}
	return h.RespJson(c, resp, nil)
}

func (h FileObjectHandler) Download(c echo.Context) error {
	var (
		req struct {
			api_schema.Filename
		}
	)
	if err := c.Bind(&req); err != nil {
		return err
	}
	var (
		filename = path.Join(h.rootDir, req.Name)
	)
	_, err := h.svc.FsStat(filename)
	if os.IsNotExist(err) {
		return h.RespJson(c, nil, echo.NewHTTPError(http.StatusNotFound))
	}
	f, err := h.svc.FsOpenFile(filename, os.O_RDONLY)
	if err != nil {
		return err
	}
	// 设置header  attachment
	c.Response().Header().Set(echo.HeaderContentDisposition, fmt.Sprintf("attachment; filename=%s", path.Base(req.Name)))
	c.Response().Header().Set(echo.HeaderContentType, "application/octet-stream")
	http.ServeContent(c.Response(), c.Request(), path.Base(req.Name), time.Now(), f)
	return nil
}

func (h FileObjectHandler) Upload(c echo.Context) error {
	filename := c.Param("Name")
	upFile, err := c.FormFile("File")
	if err != nil {
		return err
	}
	upf, err := upFile.Open()
	if err != nil {
		return err
	}
	filename = path.Join(h.rootDir, filename)
	err = h.svc.FsMkdirAll(path.Dir(filename))
	if err != nil {
		_ = upf.Close()
		return err
	}
	f, err := h.svc.FsOpenFile(filename, os.O_CREATE|os.O_TRUNC|os.O_WRONLY)
	if err != nil {
		_ = upf.Close()
		return err
	}
	_, err = io.Copy(f, upf)
	_ = upf.Close()
	_ = f.Close()
	if err != nil {
		return err
	}

	return h.RespJson(c, nil, nil)
}

func (h FileObjectHandler) Delete(c echo.Context) error {
	var (
		req struct {
			api_schema.Filename
		}
	)
	if err := c.Bind(&req); err != nil {
		return err
	}
	var (
		filename = path.Join(h.rootDir, req.Name)
	)
	_, err := h.svc.FsStat(filename)
	if os.IsNotExist(err) {
		return h.RespJson(c, nil, echo.NewHTTPError(http.StatusNotFound))
	}
	err = h.svc.FsRemove(filename)
	return h.RespJson(c, nil, err)
}

func (h FileObjectHandler) MultiUpload(c echo.Context) error {
	form, err := c.MultipartForm()
	if err != nil {
		return err
	}
	var (
		value   = url.Values(form.Value)
		dirPath = value.Get("Dir")
		force   = strings.ToLower(value.Get("Force")) == "true"
	)
	dirPath = path.Clean(dirPath)
	if dirPath == "" {
		dirPath = "/"
	}
	for _, fArr := range form.File {
		for _, upFile := range fArr {
			upf, err := upFile.Open()
			if err != nil {
				return err
			}

			//写入文件
			filename := path.Join(h.rootDir, dirPath, upFile.Filename)
			err = h.svc.FsMkdirAll(path.Dir(filename))
			if err != nil {
				_ = upf.Close()
				return err
			}
			if force {
				f, err := h.svc.FsOpenFile(filename, os.O_CREATE|os.O_TRUNC|os.O_WRONLY)
				if err != nil {
					_ = upf.Close()
					return err
				}
				_, err = io.Copy(f, upf)
				_ = upf.Close()
				_ = f.Close()
				if err != nil {
					return err
				}
			} else {
				// 判断文件是否存在
				info, err := h.svc.FsStat(filename)
				if err == nil && info.Name() != "" {
					_ = upf.Close()
					return errors.WithCode(constants.ErrCodeFileExists)
				}
				f, err := h.svc.FsOpenFile(filename, os.O_CREATE|os.O_TRUNC|os.O_WRONLY)
				if err != nil {
					_ = upf.Close()
					return err
				}
				_, err = io.Copy(f, upf)
				_ = upf.Close()
				_ = f.Close()
				if err != nil {
					return err
				}
			}
		}
	}
	return h.RespJson(c, nil, nil)
}
