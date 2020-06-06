package util

import (
	"crypto"
	"crypto/des"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
)

// base解码
func BASE64DecodeStr(src string) string {
	a, err := base64.StdEncoding.DecodeString(src)
	if err != nil {
		return ""
	}
	return string(a)
}

//Des解密
func decrypt(crypted, key []byte) ([]byte, error) {
	if len(crypted) < 1 || len(key) < 1 {
		return nil, errors.New("wrong data or key")
	}
	block, err := des.NewCipher(key)
	if err != nil {
		return nil, err
	}
	out := make([]byte, len(crypted))
	dst := out
	bs := block.BlockSize()
	if len(crypted)%bs != 0 {
		return nil, errors.New("wrong crypted size")
	}

	for len(crypted) > 0 {
		block.Decrypt(dst, crypted[:bs])
		crypted = crypted[bs:]
		dst = dst[bs:]
	}

	return out, nil
}

//[golang ECB 3DES Decrypt]
func TripleEcbDesDecrypt(crypted, key []byte) ([]byte, error) {
	tkey := make([]byte, 24, 24)
	copy(tkey, key)
	k1 := tkey[:8]
	k2 := tkey[8:16]
	k3 := tkey[16:]
	buf1, err := decrypt(crypted, k3)
	if err != nil {
		return nil, err
	}
	buf2, err := encrypt(buf1, k2)
	if err != nil {
		return nil, err
	}
	out, err := decrypt(buf2, k1)
	if err != nil {
		return nil, err
	}
	out = JDUnPadding(out)
	return out, nil
}

func VerifyPKCS1v15(msg, sign, publicKey []byte, hashType crypto.Hash) bool {
	block, _ := pem.Decode(publicKey)
	if block == nil {
		return false
	}
	pub, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		panic(err)
	}

	err = rsa.VerifyPKCS1v15(pub.(*rsa.PublicKey), hashType, msg, sign)
	return err == nil
}
