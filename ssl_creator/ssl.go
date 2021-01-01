package ssl_creator

import (
	"bufio"
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"math/big"
	"net"
	"time"
)

func CreateSimpleSSLCert(endTime string, ipAddr string, organizations, goods []string, webName string) (cert *bytes.Buffer, key *bytes.Buffer, err error) {
	cert = new(bytes.Buffer)
	key = new(bytes.Buffer)

	certWriter := bufio.NewWriter(cert)
	keyWriter := bufio.NewWriter(key)

	max := new(big.Int).Lsh(big.NewInt(1), 128)
	serialNumber, _ := rand.Int(rand.Reader, max)

	subject := pkix.Name{
		Organization:       organizations,
		OrganizationalUnit: goods,
		CommonName:         webName,
	}
	//timeDuration := years * 365 * 24 * time.Hour
	desiredTime, err := time.Parse("2006-Jan-02", endTime)
	if err != nil {
		return nil, nil, err
	}
	template := x509.Certificate{SerialNumber: serialNumber,
		Subject:     subject,
		NotBefore:   time.Now(),
		NotAfter:    desiredTime,
		KeyUsage:    x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage: []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		IPAddresses: []net.IP{net.ParseIP(ipAddr)},
	}

	pk, _ := rsa.GenerateKey(rand.Reader, 2048)
	derBytes, _ := x509.CreateCertificate(rand.Reader, &template, &template, &pk.PublicKey, pk)

	_ = pem.Encode(certWriter, &pem.Block{Type: "CERTIFICATE", Bytes: derBytes})
	_ = pem.Encode(keyWriter, &pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(pk)})
	certWriter.Flush()
	keyWriter.Flush()
	return cert, key, nil
}
