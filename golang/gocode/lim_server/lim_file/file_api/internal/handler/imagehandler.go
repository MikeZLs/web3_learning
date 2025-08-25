package handler

import (
	"errors"
	"fim_server/common/response"
	"fim_server/lim_file/file_api/internal/logic"
	"fim_server/lim_file/file_api/internal/svc"
	"fim_server/lim_file/file_api/internal/types"
	"fim_server/utils"
	"fim_server/utils/random"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"strings"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func ImageHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.ImageRequest
		if err := httpx.Parse(r, &req); err != nil {
			response.Response(r, w, nil, err)
			return
		}
		imageType := r.FormValue("imageType")
		switch imageType {
		case "avatar", "group_avatar", "chat":
		default:
			response.Response(r, w, nil, errors.New("imageType只能为 avatar,group_avatar,chat"))
			return
		}

		file, fileHead, err := r.FormFile("image")
		if err != nil {
			response.Response(r, w, nil, err)
			return
		}

		// 文件大小限制
		//mSize := float64(fileHead.Size) / float64(1024) / float64(1024)

		if fileHead.Size > svcCtx.Config.FileSize {
			response.Response(r, w, nil, fmt.Errorf("图片大小超过限制，最大只能上传%.1fMB大小的图片", float64(svcCtx.Config.FileSize)/1024/1024))
			return
		}

		// 文件后缀白名单
		nameList := strings.Split(fileHead.Filename, ".") // name.png  1.lance.png 1.lance.png
		var suffix string
		if len(nameList) > 1 {
			suffix = nameList[len(nameList)-1]
		}

		if !utils.InList(svcCtx.Config.WhiteList, suffix) {
			response.Response(r, w, nil, errors.New("图片非法"))
			return
		}

		// 文件重名
		// 在保存文件之前，去读一些文件列表  如果有重名的，算一下它们两个的hash值，一样的就不用写了
		// 它们的hash如果是不一样的，就把最新的这个重命名一下 {old_name}_xxxx.{suffix}

		dirPath := path.Join("lim_file/file_api/", svcCtx.Config.UploadDir, imageType)
		fmt.Println("dirPath", dirPath)
		dir, err := os.ReadDir(dirPath)
		if err != nil {
			os.MkdirAll(dirPath, 0666)
		}

		filePath := path.Join(svcCtx.Config.UploadDir, imageType, fileHead.Filename)
		imageData, _ := io.ReadAll(file)

		fileName := fileHead.Filename

		l := logic.NewImageLogic(r.Context(), svcCtx)
		resp, err := l.Image(&req)
		resp.Url = "/" + filePath
		if utils.InDir(dir, fileHead.Filename) {
			// 重名了
			// 先读之前的文件，去算一下它的hash
			byteData, _ := os.ReadFile(filePath)
			oldFileHash := utils.MD5(byteData)
			newFileHash := utils.MD5(imageData)
			if oldFileHash == newFileHash {
				// 两个文件是一样的
				fmt.Println("两个文件是一样的")
				response.Response(r, w, resp, nil)
				return
			}

			// 两个文件是不一样的
			// 改名的操作 name.png
			var prefix = utils.GetFilePrefix(fileName)
			newPath := fmt.Sprintf("%s_%s.%s", prefix, random.RandStr(4), suffix)
			filePath = path.Join(svcCtx.Config.UploadDir, imageType, newPath)
			// 改了的名字，还是重名了  这个地方就得递归判断了

		}
		err = os.WriteFile("lim_file/file_api/"+filePath, imageData, 0666)
		if err != nil {
			response.Response(r, w, nil, err)
			return
		}
		resp.Url = "/" + filePath
		fmt.Println("resp.Url", resp.Url)
		response.Response(r, w, resp, err)

	}
}
