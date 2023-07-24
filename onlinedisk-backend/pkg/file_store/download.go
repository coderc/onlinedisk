package file_store

import (
	oss "onlinedisk-backend/pkg/oss_client"
	"os"

	"github.com/coderc/onlinedisk-util/mapper"
)

func Download(fileUUID int64) ([]byte, error) {
	fileModel, err := mapper.QueryFileByUUID(fileUUID)
	if err != nil {
		return nil, err
	}

	localPath, err := oss.GetClient().Download(fileUUID)
	if err != nil {
		return nil, err
	}

	file, err := os.Open(localPath)
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
