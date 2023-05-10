package iso8859

import "golang.org/x/text/encoding/charmap"

func EncodeAsBytes(str string) ([]byte, error) {
	encoder := charmap.ISO8859_1.NewEncoder()
	return encoder.Bytes([]byte(str))
}

func Encode(str string) (string, error) {
	encoder := charmap.ISO8859_1.NewEncoder()
	return encoder.String(str)
}

func DecodeBytes(b []byte) (string, error) {
	decoder := charmap.ISO8859_1.NewDecoder()
	return decoder.String(string(b))
}
func Decode(str string) (string, error) {
	decoder := charmap.ISO8859_1.NewDecoder()
	return decoder.String(str)
}
