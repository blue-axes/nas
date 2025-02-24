package simple_upload

import (
	"fmt"
	api "github.com/blue-axes/tmpl/http/api"
	"github.com/blue-axes/tmpl/pkg/constants"
	"github.com/blue-axes/tmpl/pkg/context"
	"github.com/blue-axes/tmpl/pkg/errors"
	"github.com/blue-axes/tmpl/service"
	"github.com/blue-axes/tmpl/types"
	"github.com/blue-axes/tmpl/types/api_schema"
	"github.com/labstack/echo/v4"
	"io"
	"net/http"
	"net/url"
	"path"
	"strings"
	"time"
)

type (
	FileObjectHandler struct {
		*api.Handler
		svc *service.Service
	}
)

func New(svc *service.Service) *FileObjectHandler {
	h := &FileObjectHandler{
		Handler: api.New(svc),
		svc:     svc,
	}
	return h
}

func (h FileObjectHandler) Schema(c echo.Context) error {
	ctx, _ := c.Get(constants.CtxKeyContext).(*context.Context)
	// 获取文件元信息
	var (
		req struct {
			api_schema.Filename
		}
		err error
	)
	if err := c.Bind(&req); err != nil {
		return err
	}
	req.Name, err = h.validFilename(req.Name)
	if err != nil {
		return h.ToHttpError(err)
	}
	// 查询文件是否存在
	info, err := h.svc.SimpleGetFileInfo(ctx, req.Name)
	if err != nil {
		return h.ToHttpError(err)
	}
	header := c.Response().Header()
	header.Set(echo.HeaderContentLength, fmt.Sprintf("%d", info.Size))
	header.Set(echo.HeaderContentDisposition, fmt.Sprintf("attachment; filename=%s", path.Dir(path.Clean(info.Name))))
	header.Set(echo.HeaderContentType, echo.MIMEOctetStream)
	return nil
}

func (h FileObjectHandler) Download(c echo.Context) error {
	ctx, _ := c.Get(constants.CtxKeyContext).(*context.Context)
	var (
		req struct {
			api_schema.Filename
			Download bool `query:"Download"`
		}
		err error
	)
	if err := c.Bind(&req); err != nil {
		return err
	}
	req.Name, err = h.validFilename(req.Name)
	if err != nil {
		return h.ToHttpError(err)
	}
	rc, info, err := h.svc.SimpleGetFileContent(ctx, req.Name)
	if err != nil {
		return h.ToHttpError(err)
	}
	defer rc.Close()
	if req.Download {
		// 设置header  attachment
		header := c.Response().Header()
		header.Set(echo.HeaderContentLength, fmt.Sprintf("%d", info.Size))
		header.Set(echo.HeaderContentDisposition, fmt.Sprintf("attachment; filename=%s", path.Base(path.Clean(info.Name))))
		header.Set(echo.HeaderContentType, echo.MIMEOctetStream)
		http.ServeContent(c.Response(), c.Request(), path.Base(req.Name), time.Now(), rc)
		return nil
	}
	// 显示文件
	header := c.Response().Header()
	header.Set(echo.HeaderContentLength, fmt.Sprintf("%d", info.Size))
	header.Set(echo.HeaderContentType, types.Ext2MimeType(info.Ext))
	_, err = io.Copy(c.Response(), rc)
	return err
}

func (h FileObjectHandler) Upload(c echo.Context) error {
	ctx, _ := c.Get(constants.CtxKeyContext).(*context.Context)
	var (
		req struct {
			api_schema.Filename
			Overwrite bool `form:"Overwrite"`
		}
	)
	upFile, err := c.FormFile("File")
	if err != nil {
		return err
	}
	if err := c.Bind(&req); err != nil {
		return err
	}
	req.Name, err = h.validFilename(req.Name)
	if err != nil {
		return h.ToHttpError(err)
	}
	upf, err := upFile.Open()
	if err != nil {
		return err
	}
	err = h.svc.SimpleSaveFile(ctx, req.Name, upf, req.Overwrite)
	_ = upf.Close()
	return h.RespJson(c, nil, err)
}

func (h FileObjectHandler) Delete(c echo.Context) error {
	ctx, _ := c.Get(constants.CtxKeyContext).(*context.Context)
	var (
		req struct {
			api_schema.Filename
		}
		err error
	)
	if err := c.Bind(&req); err != nil {
		return err
	}
	req.Name, err = h.validFilename(req.Name)
	if err != nil {
		return h.ToHttpError(err)
	}
	// 查询文件是否存在
	err = h.svc.SimpleDeleteFile(ctx, req.Name)
	return h.RespJson(c, nil, err)
}

func (h FileObjectHandler) ReadDir(c echo.Context) error {
	ctx, _ := c.Get(constants.CtxKeyContext).(*context.Context)
	var (
		req struct {
			api_schema.Filename
		}
		resp struct {
			List []api_schema.FileInfo `json:"List,omitempty"`
		}
	)
	if err := c.Bind(&req); err != nil {
		return err
	}
	req.Name = h.correctName(req.Name)
	req.Name = path.Clean(req.Name + "/")
	if req.Name == "/" || req.Name == "" {
		req.Name = ""
	} else {
		req.Name += "/"
	}

	entry, err := h.svc.SimpleListFiles(ctx, req.Name)
	if err != nil {
		return err
	}
	var distinctMap = map[string]bool{}
	for _, item := range entry {
		var (
			fileType = "file"
			size     = item.Size
		)
		if item.Name != req.Name && path.Dir(item.Name) != path.Dir(req.Name) {
			fileType = "dir"
			size = 0
		}

		shortName := h.pickupFilename(item.Name, req.Name)
		if distinctMap[shortName] {
			continue
		}
		distinctMap[shortName] = true

		resp.List = append(resp.List, api_schema.FileInfo{
			Name:     shortName,
			Size:     size,
			FileType: fileType,
		})
	}
	return h.RespJson(c, resp, nil)
}

func (h FileObjectHandler) MultiUpload(c echo.Context) error {
	ctx, _ := c.Get(constants.CtxKeyContext).(*context.Context)
	form, err := c.MultipartForm()
	if err != nil {
		return err
	}
	var (
		value   = url.Values(form.Value)
		dirPath = value.Get("Dir")
		force   = strings.ToLower(value.Get("Overwrite")) == "true"
	)
	dirPath, err = h.validFilename(dirPath)
	if err != nil {
		return err
	}
	for _, fArr := range form.File {
		for _, upFile := range fArr {
			upf, err := upFile.Open()
			if err != nil {
				return err
			}
			//写入文件
			filename := path.Join(dirPath, upFile.Filename)
			err = h.svc.SimpleSaveFile(ctx, filename, upf, force)
			_ = upf.Close()
			if err != nil {
				return err
			}
		}
	}
	return h.RespJson(c, nil, nil)
}

func (h FileObjectHandler) validFilename(name string) (string, error) {
	/**
	./     不支持
	../    不支持
	../../  不支持
	/..		不支持
	/.		不支持
	././    不支持
	xxx/../../zzz  不支持
	xxx/././ff 不支持
	xx/./../zz 不支持
	ss./cc. => 不变还是不支持？   本来想支持最终不支持 ，测试发现windows上 a. 和 a同名
	xx../ww.. => 不变还是不支持？ 本来想支持最终不支持 ，测试发现windows上 a.. 和 a同名

	总结: 不支持路径中带有 ./ 和 ../这中遍历目录的操作
	*/
	name = strings.TrimSpace(name)
	name = h.correctName(name)
	//if strings.HasPrefix(name, "./") || strings.HasPrefix(name, "../") {
	//	return "", errors.WithCode(constants.ErrCodeInvalidArgs, "filename is invalid")
	//}
	if strings.HasSuffix(name, ".") || strings.HasSuffix(name, "..") {
		return "", errors.WithCode(constants.ErrCodeInvalidArgs, "filename is invalid")
	}
	if strings.Contains(name, "./") || strings.Contains(name, "../") {
		return "", errors.WithCode(constants.ErrCodeInvalidArgs, "filename is invalid")
	}

	return name, nil
}

func (h FileObjectHandler) correctName(name string) string {
	// 纠正文件名中的错误
	//  路径中的 \ 替换成 /
	//  /开头的文件名修改为相对路径
	//  以 / 结尾的文件名改去掉结尾的/
	name = strings.TrimSpace(name)
	name = strings.ReplaceAll(name, "\\", "/")
	name = strings.Trim(name, "/")
	return name
}

func (h FileObjectHandler) pickupFilename(filename string, parentDir string) string {
	filename = strings.Replace(filename, path.Clean(parentDir)+"/", "", 1)
	name := ""
	for _, c := range []rune(filename) {
		if c == '/' {
			break
		}
		name += string(c)
	}
	return name
}
