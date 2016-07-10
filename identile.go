package main

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
)

type IdentileCodeAlgoType int

const (
	MD5 IdentileCodeAlgoType = iota
	SHA1
	SHA256
	SHA512
)

func GetIdentileCode(data string, salt string) uint32 {
	return GetIdentileCodeByAlgo(data, salt, MD5)

}

func GetIdentileCodeByAlgo(data string, salt string, algo IdentileCodeAlgoType) uint32 {
	data = data + "@" + salt
	var result uint32 = 0
	switch algo {
	case SHA1:
		sum := sha1.Sum([]byte(data))
		result = (uint32(sum[0]&0xFF) << 24) | (uint32(sum[6]&0xFF) << 16) |
			(uint32(sum[11]&0xFF) << 8) | uint32(sum[19]&0xFF)
		break
	case SHA256:
		sum := sha256.Sum256([]byte(data))
		result = (uint32(sum[0]&0xFF) << 24) | (uint32(sum[9]&0xFF) << 16) |
			(uint32(sum[16]&0xFF) << 8) | uint32(sum[31]&0xFF)
		break
	case SHA512:
		sum := sha512.Sum512([]byte(data))
		result = (uint32(sum[0]&0xFF) << 24) | (uint32(sum[17]&0xFF) << 16) |
			(uint32(sum[49]&0xFF) << 8) | uint32(sum[63]&0xFF)
		break
	case MD5:
		fallthrough
	default:
		sum := md5.Sum([]byte(data))
		result = (uint32(sum[0]&0xFF) << 24) | (uint32(sum[5]&0xFF) << 16) |
			(uint32(sum[8]&0xFF) << 8) | uint32(sum[15]&0xFF)
		break
	}
	return result
}

func GetIdentileAlgoByString(str string) IdentileCodeAlgoType {
    switch str{
        case "sha1":
            return SHA1
        case "sha256":
            return SHA256
        case "sha512":
            return SHA512
        case "md5":
            fallthrough
        default:
            return MD5
    }
}
