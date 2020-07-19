package zhsig

import (
	"bytes"
	"crypto/sha512"
	"io"
)

// CalcHashBytes はバイト配列からハッシュを計算します。
func CalcHashBytes(data []byte) []byte {
	h, _ := CalcHashStream(bytes.NewReader(data))
	return h
}

// CalcHashFile はファイルからハッシュを計算します。
func CalcHashFile(filename string) ([]byte, error) {
	f, err := di.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	fileHash, err := CalcHashStream(f)
	if err != nil {
		return nil, err
	}
	return fileHash, nil
}

// CalcHashStream はストリームからハッシュを計算します。
func CalcHashStream(reader io.Reader) ([]byte, error) {
	buf := make([]byte, 1024)
	h := sha512.New()
	for {
		c, err := reader.Read(buf)
		if err == io.EOF {
			return h.Sum(nil), nil
		}
		if err != nil {
			return nil, err
		}
		h.Write(buf[:c])
	}
}
