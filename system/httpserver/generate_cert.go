package httpserver

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"github.com/gazercloud/gazernode/system/settings"
	"github.com/gazercloud/gazernode/utilities/logger"
	"io/ioutil"
	"math/big"
	"os"
	"time"
)

func certPublic(ss *settings.Settings) []byte {
	publicKeyPath := ss.ServerDataPath() + "/tls_public.key"
	//privateKeyPath := ss.ServerDataPath() + "/tls_private.key"

	bs, err := ioutil.ReadFile(publicKeyPath)
	if err != nil {
		logger.Println("certPublic error", err)
		return make([]byte, 0)
	}
	return bs
}

func certPrivate(ss *settings.Settings) []byte {
	//publicKeyPath := ss.ServerDataPath() + "/tls_public.key"
	privateKeyPath := ss.ServerDataPath() + "/tls_private.key"

	bs, err := ioutil.ReadFile(privateKeyPath)
	if err != nil {
		logger.Println("certPrivate error", err)
		return make([]byte, 0)
	}
	return bs
}

func generateTLS(ss *settings.Settings) {
	logger.Println("generateTLS")
	publicKeyPath := ss.ServerDataPath() + "/tls_public.key"
	privateKeyPath := ss.ServerDataPath() + "/tls_private.key"

	filesExist := true

	privateFile, err := os.Stat(privateKeyPath)
	if err != nil {
		filesExist = false
		logger.Println("generateTLS private key not found")
	} else {
		logger.Println("generateTLS private key found. size:", privateFile.Size())
	}

	publicFile, err := os.Stat(publicKeyPath)
	if err != nil {
		filesExist = false
		logger.Println("generateTLS public key not found")
	} else {
		logger.Println("generateTLS public key found. size:", publicFile.Size())
	}

	if filesExist {
		logger.Println("generateTLS keys found")
		return
	}

	logger.Println("generateTLS keys not found. generating ...")

	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		logger.Println("generateTLS generating error (ecdsa.GenerateKey)", err)
		return
	}

	template := x509.Certificate{
		SerialNumber: big.NewInt(1),
		Subject: pkix.Name{
			Organization: []string{"GazerNode"},
		},
		NotBefore: time.Now(),
		NotAfter:  time.Now().Add(time.Hour * 24 * 365),

		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
	}

	certificateBytes, err := x509.CreateCertificate(rand.Reader, &template, &template, publicKey(privateKey), privateKey)
	if err != nil {
		logger.Println("generateTLS generating error (x509.CreateCertificate)", err)
		return
	}
	out := &bytes.Buffer{}
	pem.Encode(out, &pem.Block{Type: "CERTIFICATE", Bytes: certificateBytes})
	ioutil.WriteFile(publicKeyPath, out.Bytes(), 0600)
	out.Reset()
	pem.Encode(out, pemBlockForKey(privateKey))
	ioutil.WriteFile(privateKeyPath, out.Bytes(), 0600)

	logger.Println("generateTLS OK")
}

func publicKey(priv interface{}) interface{} {
	switch k := priv.(type) {
	case *rsa.PrivateKey:
		return &k.PublicKey
	case *ecdsa.PrivateKey:
		return &k.PublicKey
	default:
		return nil
	}
}

func pemBlockForKey(priv interface{}) *pem.Block {
	switch k := priv.(type) {
	case *rsa.PrivateKey:
		return &pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(k)}
	case *ecdsa.PrivateKey:
		b, err := x509.MarshalECPrivateKey(k)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Unable to marshal ECDSA private key: %v", err)
			os.Exit(2)
		}
		return &pem.Block{Type: "EC PRIVATE KEY", Bytes: b}
	default:
		return nil
	}
}
