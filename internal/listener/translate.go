package listener

import "encoding/base64"

// implement your message decode function
func messageDecode(data string) ([]byte, error) {
	result, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		return nil, err
	}

	return result, nil
}

// implement your message encode function
func messageEncode(data []byte) string {
	return base64.StdEncoding.EncodeToString(data)
}
