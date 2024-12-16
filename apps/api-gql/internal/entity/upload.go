package entity

import (
	"io"
)

type Upload struct {
	File        io.ReadSeeker
	Filename    string
	Size        int64
	ContentType string
}
