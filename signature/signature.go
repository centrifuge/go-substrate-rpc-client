package signature

import (
	"errors"
	"regexp"

	"golang.org/x/crypto/ed25519"
)

const DEV_PHRASE = "bottom drive obey lake curtain smoke basket hold race lonely fit walk"

type SupportedKeyType int

const (
	ED25519 SupportedKeyType = iota + 1
)

type Keyring struct {
	Type  SupportedKeyType
	Pairs map[string]KeyringPair
}

func (kr *Keyring) AddFromURI(SURI string, meta map[string]interface{}, tp SupportedKeyType) {

}

type KeyringPair interface {
	Type() SupportedKeyType
	Address() string
	DecodePkcs8(passphrase string, encoded []byte)
	EncodePkcs8(passphrase string) []byte
	Meta() map[string]interface{}
	IsLocked() bool
	Lock()
	PublicKey() []byte
	SetMeta(meta map[string]interface{})
	Sign(message []byte) []byte
	// TODO
	Json(passphrase string)
	Verify(message []byte, signature []byte) bool
}

func sign(privKey []byte, msg []byte) {
	ed25519.Sign(privKey, msg)
}

var reCapture = regexp.MustCompile("^(\\w+( \\w+)*)((//?[^/]+)*)(///(.*))?$")
var reJunction = regexp.MustCompile("//(/?)([^/]+)/g")

// extractKey Extracts the phrase, path and password from a SURI format for specifying secret keys
// `<secret>/<soft-key>//<hard-key>///<password>` (the `///password` may be omitted, and
// `/<soft-key>` and `//<hard-key>` maybe repeated and mixed). The secret can be a hex string,
// mnemonic phrase or a string (to be padded)
func extractKey(suri string) (password, phrase string, path []DeriveJunction, err error) {
	matches := reCapture.FindStringSubmatch(suri)
	if len(matches) < 6 {
		return "", "", nil, errors.New("suri is not correct")
	}

	// TODO
	//path, err = extractKeyPath(matches[3])
	//if err != nil {
	//	return "", "", nil, err
	//}
	//fmt.Println(matches[6], matches[3], matches[1])
	// TODO process the path
	return matches[6], matches[1], path, nil
}

type DeriveJunction struct {
	chainCode []byte
	isHard    bool
}

//func extractKeyPath(s string) ([]DeriveJunction, error) {
//	parts := reJunction.FindStringSubmatch(s)
//	path := make([]DeriveJunction, 1)
//	constructed := ""
//
//	if len(parts) > 0 {
//		constructed = strings.Join(parts, "")
//		for _, p := range parts {
//			append(path, )
//		}
//	}
//}

/**
  /**
   * @description Generate a payload and pplies the signature from a keypair

sign (method: Method, account: KeyringPair, { blockHash, era, nonce, version }: SignatureOptions): ExtrinsicSignature {
	const signer = new Address(account.publicKey());
	const signingPayload = new SignaturePayload({
	nonce,
	method,
	era: era || this.era || IMMORTAL_ERA,
	blockHash
	});
	const signature = new Signature(signingPayload.sign(account, version as RuntimeVersion));

return this.injectSignature(signature, signer, signingPayload.nonce, signingPayload.era);
}
*/
