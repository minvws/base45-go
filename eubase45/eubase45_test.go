package eubase45

import (
	"github.com/go-errors/errors"
	"math/rand"
	"testing"
)

var testPairs = map[string][]byte {
	"QED8WEX0": []byte("ietf!"),
	"%69 VD92EX0": []byte("Hello!!"),
	"UJCLQE7W581": []byte("base-45"),
}

func TestEUBase45Decode(t *testing.T) {
	for encoded, decoded := range testPairs {
		shouldBeDecoded, err := EUBase45Decode([]byte(encoded))
		if err != nil {
			t.Fatal(errors.WrapPrefix(err, "Could not decode", 0))
		}

		err = bytesMatch(shouldBeDecoded, decoded)
		if err != nil {
			t.Fatal(errors.WrapPrefix(err, "Decoded bytes don't match", 0))
		}
	}
}

func TestEUBase45Encode(t *testing.T) {
	for encoded, decoded := range testPairs {
		shouldBeEncoded := EUBase45Encode(decoded)

		err := bytesMatch(shouldBeEncoded, []byte(encoded))
		if err != nil {
			t.Fatal(errors.WrapPrefix(err, "Encoded bytes don't match", 0))
		}
	}
}

func TestEUBase45EncodeDecode(t *testing.T) {
	for i := 0; i < 128; i++ {
		input := genBytes(rand.Intn(2048))
		encoded := EUBase45Encode(input)

		decoded, err := EUBase45Decode(encoded)
		if err != nil {
			t.Fatal(errors.WrapPrefix(err, "Could not decode", 0))
		}

		err = bytesMatch(input, decoded)
		if err != nil {
			t.Fatal(errors.WrapPrefix(err, "Decoded bytes don't match", 0))
		}
	}
}

func genBytes(amount int) []byte {
	b := make([]byte, amount)
	rand.Read(b)

	return b
}

func bytesMatch(b1, b2 []byte) error {
	if len(b1) != len(b2) {
		return errors.Errorf("Lengths don't match")
	}

	for i := 0; i < len(b1); i++ {
		if b1[i] != b2[i] {
			return errors.Errorf("Contents don't match")
		}
	}

	return nil
}

