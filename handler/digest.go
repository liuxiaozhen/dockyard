package handler

import (
	"bytes"
	"crypto"
	"fmt"
	"hash"
	"io"
	"strings"

	"github.com/docker/libtrust"
)

// Algorithm identifies and implementation of a digester by an identifier.
// Note the that this defines both the hash algorithm used and the string
// encoding.
type Algorithm string

// supported digest types
const (
	SHA256         Algorithm = "sha256"           // sha256 with hex encoding
	SHA384         Algorithm = "sha384"           // sha384 with hex encoding
	SHA512         Algorithm = "sha512"           // sha512 with hex encoding
	TarsumV1SHA256 Algorithm = "tarsum+v1+sha256" // supported tarsum version, verification only

	// Canonical is the primary digest algorithm used with the distribution
	// project. Other digests may be used but this one is the primary storage
	// digest.
	Canonical = SHA256
)

var (
	// TODO(stevvooe): Follow the pattern of the standard crypto package for
	// registration of digests. Effectively, we are a registerable set and
	// common symbol access.

	// algorithms maps values to hash.Hash implementations. Other algorithms
	// may be available but they cannot be calculated by the digest package.
	algorithms = map[Algorithm]crypto.Hash{
		SHA256: crypto.SHA256,
		SHA384: crypto.SHA384,
		SHA512: crypto.SHA512,
	}
)

// Available returns true if the digest type is available for use. If this
// returns false, New and Hash will return nil.
func (a Algorithm) Available() bool {
	h, ok := algorithms[a]
	if !ok {
		return false
	}

	// check availability of the hash, as well
	return h.Available()
}

func (a Algorithm) New() Digester {
	return &digester{
		alg:  a,
		hash: a.Hash(),
	}
}

func (a Algorithm) Hash() hash.Hash {
	if !a.Available() {
		return nil
	}

	return algorithms[a].New()
}

type Digester interface {
	Hash() hash.Hash // provides direct access to underlying hash instance.
	Digest() string
}

// digester provides a simple digester definition that embeds a hasher.
type digester struct {
	alg  Algorithm
	hash hash.Hash
}

func (d *digester) Hash() hash.Hash {
	return d.hash
}

func (d *digester) Digest() string {
	return string(fmt.Sprintf("%s:%x", d.alg, d.hash.Sum(nil)))
}

func FromReader(rd io.Reader) (string, error) {
	digester := Canonical.New()

	if _, err := io.Copy(digester.Hash(), rd); err != nil {
		return "", err
	}

	return digester.Digest(), nil
}

func Payload(data []byte) ([]byte, error) {
	jsig, err := libtrust.ParsePrettySignature(data, "signatures")
	if err != nil {
		return nil, err
	}

	// Resolve the payload in the manifest.
	return jsig.Payload()
}

func DigestManifest(data []byte) (string, error) {
	p, err := Payload(data)
	if err != nil {
		if !strings.Contains(err.Error(), "missing signature key") {
			return "", err
		}

		p = data
	}

	digest, err := FromReader(bytes.NewReader(p))
	if err != nil {
		return "", err
	}

	return digest, err
}
