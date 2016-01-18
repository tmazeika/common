package security

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"math/big"
	"time"
)

func GeneratePrivKey() (privKey *rsa.PrivateKey, b []byte, err error) {
	const Bits = 4096
	if privKey, err = rsa.GenerateKey(rand.Reader, Bits); err != nil {
		return
	}
	der := x509.MarshalPKCS1PrivateKey(privKey)

	// Pem encode.
	b = pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: der,
	})
	return
}

func CreateCertificate(privKey *rsa.PrivateKey) (b []byte, err error) {
	const SerialNumber = 50977
	const YearsValid = 1000
	template := x509.Certificate{
		SerialNumber:          big.NewInt(SerialNumber),
		SignatureAlgorithm:    x509.SHA512WithRSA,
		PublicKeyAlgorithm:    x509.RSA,
		NotBefore:             time.Now(),
		NotAfter:              time.Now().AddDate(YearsValid, 0, 0),
		BasicConstraintsValid: true,
		IsCA:        true,
		ExtKeyUsage: []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		KeyUsage:    x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
	}
	cert, err := x509.CreateCertificate(rand.Reader, &template, &template, privKey.Public(), privKey)
	if err != nil {
		return
	}

	// Pem encode.
	b = pem.EncodeToMemory(&pem.Block{
		Type:  "CERTIFICATE",
		Bytes: cert,
	})
	return
}
