package dao

import (
	"crypto/md5"
	"encoding/hex"
)

func mD5(data []byte) string {
	_md5 := md5.New()
	_md5.Write(data)
	return hex.EncodeToString(_md5.Sum([]byte("")))
}
