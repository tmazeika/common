package security

import (
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"github.com/transhift/appdir"
	"os"
)

func Certificate(keyName, certName string, dir *appdir.AppDir) (cert tls.Certificate, err error) {
	const FileMode = 0600

	var privKey *rsa.PrivateKey
	// Read or generate private key.
	if err = dir.IfExistsOrOtherwise(keyName, func(b []byte) (err error) {
		p, _ := pem.Decode(b)
		privKey, err = x509.ParsePKCS1PrivateKey(p.Bytes)
		return
	}, func(file *os.File) (b []byte, err error) {
		// Set file mode.
		if err = file.Chmod(FileMode); err != nil {
			return
		}
		privKey, b, err = GeneratePrivKey()
		return
	}); err != nil {
		return
	}

	// Generate certificate if not exists.
	if err = dir.IfNExists(certName, func(file *os.File) (b []byte, err error) {
		// Set file mode.
		if err = file.Chmod(FileMode); err != nil {
			return
		}
		return CreateCertificate(privKey)
	}); err != nil {
		return
	}
	return tls.LoadX509KeyPair(dir.FilePath(certName), dir.FilePath(keyName))
}
