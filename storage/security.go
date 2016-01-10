package storage

import (
	"crypto/tls"
	"crypto/rsa"
	"crypto/rand"
	"crypto/x509"
	"math/big"
	"time"
	"encoding/pem"
	"io/ioutil"
)

const (
	KeyMode  = 0600
	CertMode = 0600
)

func (s storage) Certificate(keyName, certName string) (*tls.Certificate, error) {
	keyPath := s.filePath(keyName)
	certPath := s.filePath(certName)

	if ! fileExists(keyPath) || ! fileExists(certPath) {
		if err := saveNewCertificate(keyPath, certPath); err != nil {
			return nil, err
		}
	}

	cert, err := tls.LoadX509KeyPair(certPath, keyPath)

	if err != nil {
		return nil, err
	}

	return &cert, nil
}

func saveNewCertificate(keyPath, certPath string) error {
	keyData, certData, err := createCertificate()

	if err != nil {
		return err
	}

	err = ioutil.WriteFile(keyPath, keyData, KeyMode)

	if err != nil {
		return err
	}

	return ioutil.WriteFile(certPath, certData, CertMode)
}

func createCertificate() (keyData []byte, certData []byte, err error) {
	const RSABits = 4096

	key, err := rsa.GenerateKey(rand.Reader, RSABits)

	if err != nil {
		return
	}

	keyData = x509.MarshalPKCS1PrivateKey(key)
	ca := &x509.Certificate{
		SerialNumber:          big.NewInt(50977),
		SignatureAlgorithm:    x509.SHA512WithRSA,
		PublicKeyAlgorithm:    x509.RSA,
		NotBefore:             time.Now(),
		NotAfter:              time.Now().AddDate(1000, 0, 0),
		BasicConstraintsValid: true,
		IsCA:                  true,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		KeyUsage: x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
	}
	certData, err = x509.CreateCertificate(rand.Reader, ca, ca, &key.PublicKey, key)

	if err != nil {
		return
	}

	keyData = pem.EncodeToMemory(&pem.Block{
		Type: "RSA PRIVATE KEY",
		Bytes: keyData,
	})

	certData = pem.EncodeToMemory(&pem.Block{
		Type: "CERTIFICATE",
		Bytes: certData,
	})

	return
}
