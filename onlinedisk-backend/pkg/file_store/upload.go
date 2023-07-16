package file_store

import (
	"fmt"

	"github.com/coderc/onlinedisk-util/mapper"
	"github.com/coderc/onlinedisk-util/model"
)

const (
	errorFileNotExist = "文件不存在"
)

func Chunk(fileSHA1, fileName string, userId int64) error {
	fileModel := hasFile(fileSHA1)
	if fileModel != nil && fileModel.SHA1 == fileSHA1 { // 秒传成功
		userFileModel := &model.UserFileModel{
			UserId:   userId,
			FileId:   fileModel.UUID,
			FileName: fileName,
		}
		err := mapper.InsertUserFile(userFileModel)
		if err != nil {
			return err
		} else {
			return nil
		}
	}
	// 秒传失败，没有此文件
	return fmt.Errorf(errorFileNotExist)
}

func Upload(fileModel *model.FileModel, userFileModel *model.UserFileModel) error {
	err := mapper.InsertFile(fileModel) // table_file
	if err != nil {
		return err
	}
	return mapper.InsertUserFile(userFileModel) // table_user_file
}

func hasFile(sha1 string) *model.FileModel {
	fileModel, err := mapper.QueryFileBySHA1(sha1)
	if err != nil {
		return nil
	}
	return fileModel
}
