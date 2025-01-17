package chainhash

import (
	"encoding/hex"
	"fmt"
)

// HashSize of array used to store hashes.  See Hash.
const HashSize = 32

// MaxHashStringSize is the maximum length of a Hash hash string.
const MaxHashStringSize = HashSize * 2

// ErrHashStrSize describes an error that indicates the caller specified a hash string that has too many characters.
var ErrHashStrSize = fmt.Errorf(
	"max hash string length is %v bytes",
	MaxHashStringSize,
)

// Hash is used in several of the bitcoin messages and common structures. It typically represents the double sha256 of
// data.
type Hash [HashSize]byte

// String returns the Hash as the hexadecimal string of the byte-reversed hash.
func (hash Hash) String() string {
	for i := 0; i < HashSize/2; i++ {
		hash[i], hash[HashSize-1-i] = hash[HashSize-1-i], hash[i]
	}
	return hex.EncodeToString(hash[:])
}

// CloneBytes returns a copy of the bytes which represent the hash as a byte slice. NOTE: It is generally cheaper to
// just slice the hash directly thereby reusing the same bytes rather than calling this method.
func (hash *Hash) CloneBytes() []byte {
	newHash := make([]byte, HashSize)
	copy(newHash, hash[:])
	return newHash
}

// SetBytes sets the bytes which represent the hash. An error is returned if the number of bytes passed in is not
// HashSize.
func (hash *Hash) SetBytes(newHash []byte) (e error) {
	nhlen := len(newHash)
	if nhlen != HashSize {
		return fmt.Errorf(
			"invalid hash length of %v, want %v", nhlen,
			HashSize,
		)
	}
	copy(hash[:], newHash)
	return nil
}

// IsEqual returns true if target is the same as hash.
func (hash *Hash) IsEqual(target *Hash) bool {
	if hash == nil && target == nil {
		return true
	}
	if hash == nil || target == nil {
		return false
	}
	return *hash == *target
}

// NewHash returns a new Hash from a byte slice. An error is returned if the number of bytes passed in is not HashSize.
func NewHash(newHash []byte) (*Hash, error) {
	var sh Hash
	e := sh.SetBytes(newHash)
	if e != nil {
		E.Ln(e)
		return nil, e
	}
	return &sh, e
}

// NewHashFromStr creates a Hash from a hash string. The string should be the hexadecimal string of a byte-reversed
// hash, but any missing characters result in zero padding at the end of the Hash.
func NewHashFromStr(hash string) (ret *Hash, e error) {
	ret = new(Hash)
	if e = Decode(ret, hash); E.Chk(e) {
		return
	}
	return
}

// Decode decodes the byte-reversed hexadecimal string encoding of a Hash to a destination.
func Decode(dst *Hash, src string) (e error) {
	// Return error if hash string is too long.
	if len(src) > MaxHashStringSize {
		return ErrHashStrSize
	}
	// Hex decoder expects the hash to be a multiple of two.  When not, pad with a leading zero.
	var srcBytes []byte
	if len(src)%2 == 0 {
		srcBytes = []byte(src)
	} else {
		srcBytes = make([]byte, 1+len(src))
		srcBytes[0] = '0'
		copy(srcBytes[1:], src)
	}
	// Hex decode the source bytes to a temporary destination.
	var reversedHash Hash
	if _, e = hex.Decode(
		reversedHash[HashSize-hex.DecodedLen(len(srcBytes)):],
		srcBytes,
	); E.Chk(e) {
		return e
	}
	// Reverse copy from the temporary hash to destination. Because the temporary was zeroed, the written result will be
	// correctly padded.
	for i, b := range reversedHash[:HashSize/2] {
		dst[i], dst[HashSize-1-i] = reversedHash[HashSize-1-i], b
	}
	return
}
