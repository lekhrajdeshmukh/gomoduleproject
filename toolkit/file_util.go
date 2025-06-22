package toolkit

import (
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
)

func IsFileContentTypeAllowed(files []string, fileType string) bool {
	if len(files) > 0 {
		for _, x := range files {
			if strings.EqualFold(fileType, x) {
				return true
			}
		}
		return false
	} else {
		return true
	}
}

func (t *Tools) GetFileName(isRename bool, fileName string) string {
	if isRename {
		return fmt.Sprintf("%s%s", t.RandomStringUsingRandInt(25), filepath.Ext(fileName))
	} else {
		return fileName
	}
}

func UploadFileToTargetFolder(inFile multipart.File, targetDirectory, fileName string) (fileSize int64, err error) {
	var outFile *os.File
	defer outFile.Close()

	if outFile, err = os.Create(filepath.Join(targetDirectory, fileName)); err != nil {
		return 0, err
	} else {
		fileSize, err = io.Copy(outFile, inFile)
		if err != nil {
			return 0, err
		}
	}
	return fileSize, nil

}
