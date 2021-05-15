package base45

import (
	"bytes"
	"github.com/go-errors/errors"
	gobig "math/big"
)

var qrCharset = []byte("0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ $%*+-./:")
var goCharset = []byte("0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHI")
var qrCharsetLen = 45

var qrToGoLookup, goToQrLookup map[byte]byte

func buildLookups() (qrGo map[byte]byte, goQr map[byte]byte) {
	qrGo = map[byte]byte{}
	goQr = map[byte]byte{}

	for i := 0; i < qrCharsetLen; i++ {
		qrGo[qrCharset[i]] = goCharset[i]
		goQr[goCharset[i]] = qrCharset[i]
	}

	return
}

func ensureLookups() {
	if qrToGoLookup == nil || goToQrLookup == nil {
		qrToGoLookup, goToQrLookup = buildLookups()
	}
}

func Base45Encode(input []byte) ([]byte, error) {
	ensureLookups()

	if len(input) > 0 && input[0] == 0 {
		return nil, errors.Errorf("Could not encode input that starts with a zero-byte")
	}

	goEncoded := []byte(new(gobig.Int).SetBytes(input).Text(45))
	outputLen := len(goEncoded)

	qrEncoded := make([]byte, outputLen)
	for i := 0; i < outputLen; i++ {
		qrEncoded[i] = goToQrLookup[goEncoded[i]]
	}

	return qrEncoded, nil
}

func Base45Decode(input []byte) ([]byte, error) {
	ensureLookups()

	inputLen := len(input)
	goEncodedInput := make([]byte, inputLen)

	for i := 0; i < len(input); i++ {
		goEncodedInput[i] = qrToGoLookup[input[i]]
	}

	decodedInt, ok := new(gobig.Int).SetString(string(goEncodedInput), 45)
	if !ok {
		return nil, errors.Errorf("Could not decode invalid character in input; not alphanumeric")
	}

	return decodedInt.Bytes(), nil
}

// These two functions offer an alternative implementation that is easy to port to other
// languages, because it uses standard big integer primitives. It's quite slow.

var bigQrCharsetLen = gobig.NewInt(int64(qrCharsetLen))

func Base45EncodeAlternative(input []byte) ([]byte, error) {
	inputLen := len(input)
	if inputLen == 0 {
		return []byte{qrCharset[0]}, nil
	}

	if input[0] == 0 {
		return nil, errors.Errorf("Could not encode input that starts with a zero-byte")
	}

	estOutputLen := int(float64(inputLen)*1.4568) + 1
	output := make([]byte, 0, estOutputLen)

	divident, remainder := new(gobig.Int), new(gobig.Int)
	divident.SetBytes(input)

	for len(divident.Bits()) != 0 {
		divident, remainder = divident.QuoRem(divident, bigQrCharsetLen, remainder)
		output = append(output, qrCharset[remainder.Int64()])
	}

	return reverseByteSlice(output), nil
}

func Base45DecodeAlternative(input []byte) ([]byte, error) {
	inputLength := len(input)
	result := gobig.NewInt(0)

	for i, b := range input {
		charsetIndex := bytes.IndexByte(qrCharset, b)
		if charsetIndex == -1 {
			return nil, errors.Errorf("Could not decode invalid character in input; not alphanumeric")
		}

		factor := gobig.NewInt(int64(charsetIndex))

		weight := new(gobig.Int)
		weight.Exp(bigQrCharsetLen, gobig.NewInt(int64(inputLength-i-1)), nil)

		result = result.Add(result, new(gobig.Int).Mul(factor, weight))
	}

	return result.Bytes(), nil
}

func reverseByteSlice(bs []byte) []byte {
	amount := len(bs)
	result := make([]byte, amount)

	for i := 0; i < amount; i++ {
		result[i] = bs[amount-i-1]
	}

	return result
}
