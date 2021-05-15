package base45

import (
	"fmt"
	"github.com/go-errors/errors"
	"math/rand"
	"testing"
)

var benchmarkBytes = genBytes(1024)

func TestBase45EncodeDecode(t *testing.T) {
	for i := 0; i < 32; i++ {
		input := genBytes(rand.Intn(2048))

		encoded, err := Base45Encode(input)
		if err != nil {
			t.Fatal(errors.WrapPrefix(err, "Could not encode", 0))
		}

		decoded, err := Base45Decode(encoded)
		if err != nil {
			t.Fatal(errors.WrapPrefix(err, "Could not decode", 0))
		}

		err = bytesMatch(input, decoded)
		if err != nil {
			fmt.Println(input)
			fmt.Println(string(encoded))
			fmt.Println(decoded)
			t.Fatal(errors.WrapPrefix(err, "Decoded bytes don't match", 0))
		}
	}
}

func TestBase45EncodeDecodeAlternative(t *testing.T) {
	for i := 0; i < 32; i++ {
		input := genBytes(rand.Intn(2048))

		encoded, err := Base45EncodeAlternative(input)
		if err != nil {
			t.Fatal(errors.WrapPrefix(err, "Could not encode", 0))
		}

		decoded, err := Base45DecodeAlternative(encoded)
		if err != nil {
			t.Fatal(errors.WrapPrefix(err, "Could not decode", 0))
		}

		err = bytesMatch(input, decoded)
		if err != nil {
			t.Fatal(errors.WrapPrefix(err, "Decoded bytes don't match", 0))
		}
	}
}

func TestBase45ImplementationCorrespondence(t *testing.T) {
	for i := 0; i < 32; i++ {
		input := genBytes(rand.Intn(2048))
		e1, _ := Base45Encode(input)
		e2, _ := Base45EncodeAlternative(input)

		err := bytesMatch(e1, e2)
		if err != nil {
			fmt.Println(e1)
			fmt.Println(e2)
			t.Fatal(errors.WrapPrefix(err, "Implementations don't match", 0))
		}
	}
}

func BenchmarkBase45Encode(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Base45Encode(benchmarkBytes)
	}
}

func BenchmarkBase45Decode(b *testing.B) {
	encoded, _ := Base45Encode(benchmarkBytes)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, err := Base45Decode(encoded)
		if err != nil {
			b.FailNow()
		}
	}
}

func BenchmarkBase45EncodeAlternative(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Base45EncodeAlternative(benchmarkBytes)
	}
}

func BenchmarkBase45DecodeAlternative(b *testing.B) {
	encoded, _ := Base45Encode(benchmarkBytes)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, err := Base45DecodeAlternative(encoded)
		if err != nil {
			b.FailNow()
		}
	}
}

func genBytes(amount int) []byte {
	// Don't generate values that start with a zero-byte
	var b []byte
	for {
		b = make([]byte, amount)
		rand.Read(b)

		if amount == 0 || b[0] != 0 {
			break
		}
	}

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