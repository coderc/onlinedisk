package file_store

import (
	"os"

	"github.com/coderc/onlinedisk-util/mapper"
)

func Download(fileUUID int64) ([]byte, error) {
	fileModel, err := mapper.QueryFileByUUID(fileUUID)
	if err != nil {
		return nil, err
	}

	file, err := os.Open(fileModel.Path)
	if err != nil {
		return nil, err
	}

	fileBytes := make([]byte, fileModel.Size)
	_, err = file.Read(fileBytes)
	if err != nil {
		return nil, err
	}

	return fileBytes, nil
}
