package toolkit

import (
	"errors"
	"fmt"
	"net/http"
)

// UploadedFile is a struct used to save information about an uploaded file
type UploadedFile struct {
	NewFileName      string
	OriginalFileName string
	FileSize         int64
}

func (t *Tools) UploadOneFile(r *http.Request, uploadDir string, rename ...bool) (*UploadedFile, error) {
	renameFile := true

	if len(rename) > 0 {
		renameFile = rename[0]
	}

	files, err := t.UploadFiles(r, uploadDir, renameFile)
	if err != nil {
		return nil, err
	}

	return files[0], nil
}

func (t *Tools) UploadFiles(r *http.Request, uploadDir string, rename ...bool) ([]*UploadedFile, error) {
	renameFile := true
	if len(rename) > 0 {
		renameFile = rename[0]
	}

	var uploadedFiles []*UploadedFile

	if t.MaxFileSize == 0 {
		t.MaxFileSize = 1024 * 1024 * 1024
	}

	err := t.CreateDirIfNotExist(uploadDir)
	if err != nil {
		return nil, err
	}

	if err := r.ParseMultipartForm(int64(t.MaxFileSize)); err != nil {
		return nil, errors.New("the uploaded file is too big")
	}

	fmt.Println("I am coming here")

	for _, fHeaders := range r.MultipartForm.File {
		for _, header := range fHeaders {
			uploadedFiles, err := func(UploadedFiles []*UploadedFile) ([]*UploadedFile, error) {
				var uploadFile UploadedFile
				infile, err := header.Open()
				if err != nil {
					return nil, err
				}
				defer infile.Close()

				buff := make([]byte, 512)
				_, err = infile.Read(buff)
				if err != nil {
					return nil, err
				}

				// Check to see if file type is permitted
				fileType := http.DetectContentType(buff)
				allowed := IsFileContentTypeAllowed(t.AllowedFileTypes, fileType)

				if !allowed {
					return nil, errors.New("the uploaded file type is not allowed")
				}

				_, err = infile.Seek(0, 0)
				if err != nil {
					return nil, err
				}

				uploadFile.OriginalFileName = header.Filename
				uploadFile.NewFileName = t.GetFileName(renameFile, header.Filename)

				fileSize, err := UploadFileToTargetFolder(infile, uploadDir, uploadFile.NewFileName)
				if err != nil {
					return nil, err
				}
				uploadFile.FileSize = fileSize

				uploadedFiles = append(uploadedFiles, &uploadFile)

				return uploadedFiles, nil

			}(uploadedFiles)

			if err != nil {
				return uploadedFiles, err
			}
		}
	}
	return uploadedFiles, nil
}
