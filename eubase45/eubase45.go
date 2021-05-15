package eubase45

import (
	"bytes"
	"github.com/go-errors/errors"
)

var qrCharset = []byte("0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ $%*+-./:")
var qrCharsetLen = 45
var qrCharsetLenSquared = 2025

func EUBase45Decode(input []byte) ([]byte, error) {
	inputLen := len(input)
	lastChunkSize := inputLen % 3
	tripletInputLen := inputLen - lastChunkSize

	if lastChunkSize == 1 {
		return nil, errors.Errorf("Could not decode input with an invalid length")
	}

	var val int

	charIndexes := make([]int, 0, inputLen)
	for i := 0; i < inputLen; i++ {
		val = bytes.IndexByte(qrCharset, input[i])
		if val == -1 {
			return nil, errors.Errorf("Invalid character in EUBase45; not alphanumeric")
		}

		charIndexes = append(charIndexes, val)
	}

	decodedLen := int(float32(tripletInputLen) / 1.5 + float32(lastChunkSize) * 0.5)
	decoded := make([]byte, 0, decodedLen)

	for i := 0; i < tripletInputLen; i += 3 {
		val = charIndexes[i] + charIndexes[i+1] * qrCharsetLen + charIndexes[i+2] * qrCharsetLenSquared
		decoded = append(decoded, byte(val / 256))
		decoded = append(decoded, byte(val % 256))
	}

	if lastChunkSize == 2 {
		val = charIndexes[tripletInputLen] + charIndexes[tripletInputLen + 1] * qrCharsetLen
		decoded = append(decoded, byte(val))
	}

	return decoded, nil
}

func EUBase45Encode(input []byte) ([]byte) {
	inputLen := len(input)
	lastChunkSize := inputLen % 2
	pairInputLen := inputLen - lastChunkSize
	encodedLen := int(float32(pairInputLen) * 1.5) + lastChunkSize * 2

	encoded := make([]byte, 0, encodedLen)

	var val int
	for i := 0; i < pairInputLen; i += 2 {
		val = int(input[i]) * 256 + int(input[i+1])
		encoded = append(encoded, qrCharset[val % qrCharsetLen])
		encoded = append(encoded, qrCharset[val / qrCharsetLen % qrCharsetLen])
		encoded = append(encoded, qrCharset[val / qrCharsetLenSquared % qrCharsetLen])
	}

	if lastChunkSize == 1 {
		lastByte := int(input[pairInputLen])
		encoded = append(encoded, qrCharset[lastByte % qrCharsetLen])

		lastSymbol := qrCharset[0]
		if lastByte >= qrCharsetLen {
			lastSymbol = qrCharset[lastByte / qrCharsetLen]
		}

		encoded = append(encoded, lastSymbol)
	}

	return encoded
}