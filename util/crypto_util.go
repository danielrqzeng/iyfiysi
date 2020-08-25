package util

import (
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
)

func Md5sum(data []byte) string {
	h := md5.New()
	h.Write(data)
	d := hex.EncodeToString(h.Sum(nil))
	return d
}

func Sha1sum(data []byte) string {
	h := sha1.New()
	h.Write(data)
	d := hex.EncodeToString(h.Sum(nil))
	return d
}
func Sha256sum(data []byte) string {
	h := sha256.New()
	h.Write(data)
	d := hex.EncodeToString(h.Sum(nil))
	return d
}

func HmacSha1(data []byte, key string) []byte {
	m := hmac.New(sha1.New, []byte(key))
	m.Write([]byte(data))
	return m.Sum(nil)
}

func HmacSha256(data []byte, key string) []byte {
	m := hmac.New(sha256.New, []byte(key))
	m.Write([]byte(data))
	return m.Sum(nil)
}

func Base64Encode(data []byte) (str string) {
	str = base64.StdEncoding.EncodeToString(data)
	return
}
