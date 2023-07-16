package file_store

import "github.com/coderc/onlinedisk-util/mapper"

func Delete(userId, fileId int64) error {
	return mapper.DeleteUserFile(userId, fileId)
}