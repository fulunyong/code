package file

import (
	"encoding/json"
	"fmt"
	"github.com/fulunyong/code/common"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

//rootPath 上传跟目录
//maxSize 总文件大小限制
//singleSize 单个文件大小限制
func UploadFileHandler(rootPath string, maxSize, singleSize int64) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response := new(common.BaseResponse)
		response.Code = common.ResponseError
		response.Msg = "上传失败"
		fmt.Println(response)
		dataMap := make(map[string]string)
		fmt.Println(dataMap)
		// 校验总文件大小
		r.Body = http.MaxBytesReader(w, r.Body, maxSize)
		if err := r.ParseMultipartForm(maxSize); err != nil {
			response.Msg = fmt.Sprintf("%s%d%s", "总文件大小超出限制，最大上传限制:", maxSize/1024, "KB")
			Render(w, *response)
			return
		}
		//判断有没有上传有文件
		fileForm := r.MultipartForm
		if fileForm == nil {
			response.Msg = "当前请求体没有上传文件！"
			Render(w, *response)
			return
		}
		//前缀路径
		dir := "/"
		//模块名称
		module := r.PostFormValue("module")
		if module != "" {
			dir = fmt.Sprintf("%s/%s", dir, module)
		}
		//用户id
		userId := r.PostFormValue("userId")
		if userId != "" {
			dir = fmt.Sprintf("%s/%s", dir, userId)
		}
		//用户id
		path := r.PostFormValue("path")
		if path != "" {
			dir = fmt.Sprintf("%s/%s", dir, path)
		}
		//创建文件夹
		_, dirError := MkdirOnNotExist(fmt.Sprintf("%s%s", rootPath, dir))
		if dirError != nil {
			response.Msg = fmt.Sprintf("%s path: %s err:%s", "文件夹创建失败, ", fmt.Sprintf("%s%s", rootPath, dir), dirError.Error())
			Render(w, *response)
			return
		}
		//判断文件数量是否为0
		fileHeaders := fileForm.File
		if 0 == len(fileHeaders) {
			response.Msg = "当前上传文件数量为0"
			Render(w, *response)
			return
		}
		//处理文件
		for k := range fileHeaders {
			file, fileHeader, dirError := r.FormFile(k)
			if dirError != nil {
				response.Msg = fmt.Sprintf("%s err:%s", "读取上传的文件失败, ", dirError.Error())
				Render(w, *response)
				return
			}
			if file == nil || fileHeader == nil {
				response.Msg = "读取文件流失败！"
				Render(w, *response)
				return
			}
			//文件名称
			fileName := fileHeader.Filename
			fmt.Println("fileName:", fileName)
			//校验单个文件大小
			if singleSize < fileHeader.Size {
				response.Msg = fmt.Sprintf("key:%s 单文件大小超出限制 %d%s", k, singleSize/1024, "KB")
				Render(w, *response)
				return
			}
			filePath := fmt.Sprintf("%s/%s", dir, fileName)
			//判断文件是否存在
			exists, dirError := PathExists(fmt.Sprintf("%s%s", rootPath, filePath))
			if dirError != nil {
				response.Msg = fmt.Sprintf("读取文件失败:%s", dirError.Error())
				Render(w, *response)
				return
			}
			if exists {
				//存在同名文件处理 拼接时间戳
				filePath = fmt.Sprintf("%s/%d%s", dir, time.Now().Unix(), fileName)
			}
			newFile, dirError := os.Create(fmt.Sprintf("%s%s", rootPath, filePath))
			if dirError != nil {
				response.Msg = fmt.Sprintf("创建文件失败:%s %s", fmt.Sprintf("%s%s", rootPath, filePath), dirError.Error())
				Render(w, *response)
				return
			}
			bytes, _ := ioutil.ReadAll(file)
			_, dirError = newFile.Write(bytes)
			if dirError != nil {
				response.Msg = fmt.Sprintf("文件保存失败:%s %s", fmt.Sprintf("%s%s", rootPath, filePath), dirError.Error())
				Render(w, *response)
				return
			}
			_ = newFile.Close()
			_ = file.Close()
			//返回处理
			s := dataMap[k]
			if s != "" {
				dataMap[k] = fmt.Sprintf("%s,%s", s, filePath)
			} else {
				dataMap[k] = filePath
			}
		}

		response.Code = common.ResponseOK
		response.Data = dataMap
		response.Msg = "文件上传成功！"
		Render(w, *response)
	})
}

//判断文件是否存在  不存在则创建
func MkdirOnNotExist(path string) (bool, error) {
	//判断是否存在
	exists, err := PathExists(path)
	if err == nil {
		if exists {
			return true, nil
		} else {
			err := os.Mkdir(path, os.ModePerm)
			if err == nil {
				return true, nil
			}
			return false, err
		}
	}
	return false, err
}

//判断是否存在
//如果返回的错误为nil,说明文件或文件夹存在
//如果返回的错误类型使用os.IsNotExist()判断为true,说明文件或文件夹不存在
//如果返回的错误为其它类型,则不确定是否在存在
func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

//统一返回处理
func Render(w http.ResponseWriter, response common.BaseResponse) {
	w.WriteHeader(http.StatusOK)
	bytes, _ := json.Marshal(response)
	_, _ = w.Write(bytes)
}
