package utils

import "mime/multipart"

type FormFile struct {
	File		multipart.File
	FileHeader	*multipart.FileHeader
}
