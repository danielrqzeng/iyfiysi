package comm

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha1"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/asn1"
	"encoding/pem"
	"errors"
	"fmt"
	"io/ioutil"
	"math/big"
	"os"
	"time"
)

// rsaPublicKey reflects the ASN.1 structure of a PKCS#1 public key.
type rsaPublicKey struct {
	N *big.Int
	E int
}

// GenerateSubjectKeyID generates SubjectKeyId used in Certificate
// Id is 160-bit SHA-1 hash of the value of the BIT STRING subjectPublicKey
func GenerateSubjectKeyID(pub crypto.PublicKey) ([]byte, error) {
	var pubBytes []byte
	var err error
	switch pub := pub.(type) {
	case *rsa.PublicKey:
		pubBytes, err = asn1.Marshal(rsaPublicKey{
			N: pub.N,
			E: pub.E,
		})
		if err != nil {
			return nil, err
		}
	default:
		return nil, errors.New("only RSA public key is supported")
	}

	hash := sha1.Sum(pubBytes)
	return hash[:], nil
}

//项目的签发证书
func CreateCA(c, o, ou, cn string, caPemFile, caKeyFile string) (err error) {
	max := new(big.Int).Lsh(big.NewInt(1), 128)   //把 1 左移 128 位，返回给 big.Int
	serialNumber, _ := rand.Int(rand.Reader, max) //返回在 [0, max) 区间均匀随机分布的一个随机值
	subject := pkix.Name{                         //Name代表一个X.509识别名。只包含识别名的公共属性，额外的属性被忽略。
		Country:            []string{c},
		Organization:       []string{o},
		OrganizationalUnit: []string{ou},
		CommonName:         cn,
	}
	issuer := pkix.Name{ //Name代表一个X.509识别名。只包含识别名的公共属性，额外的属性被忽略。
		Country:            []string{c},
		Organization:       []string{o},
		OrganizationalUnit: []string{ou},
		CommonName:         cn,
	}
	template := x509.Certificate{
		SerialNumber: serialNumber, // SerialNumber 是 CA 颁布的唯一序列号，在此使用一个大随机数来代表它
		Subject:      subject,
		Issuer:       issuer,
		// activate CA
		BasicConstraintsValid: true,
		IsCA:                  true,
		// Not allow any non-self-issued intermediate CA, sets MaxPathLen=0
		MaxPathLenZero: true,
		NotBefore:      time.Now().Add(-600).UTC(),
		NotAfter:       time.Now().Add(100 * 365 * 24 * time.Hour),     //一百年不变
		KeyUsage:       x509.KeyUsageCertSign | x509.KeyUsageCRLSign,   //ca
		ExtKeyUsage:    []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth}, // 密钥扩展用途的序列
		//IPAddresses:  []net.IP{net.ParseIP("127.0.0.1")},
	}

	pk, err := rsa.GenerateKey(rand.Reader, 4096) //生成一对具有指定字位数的RSA密钥
	if err != nil {
		return
	}
	subjectKeyID, err := GenerateSubjectKeyID(&pk.PublicKey)
	if err != nil {
		return
	}
	template.SubjectKeyId = subjectKeyID

	derBytes, err := x509.CreateCertificate(rand.Reader, &template, &template, &pk.PublicKey, pk)
	if err != nil {
		return
	}
	certOut, _ := os.Create(caPemFile)
	pem.Encode(certOut, &pem.Block{Type: "CERTIFICATE", Bytes: derBytes, Headers: nil})
	certOut.Close()
	keyOut, _ := os.Create(caKeyFile)
	pem.Encode(keyOut, &pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(pk), Headers: nil})
	keyOut.Close()
	return
}

//项目的csr
func CreateCSR(c, o, ou, cn string, csrFile, keyFile string) (err error) {
	subject := pkix.Name{ //Name代表一个X.509识别名。只包含识别名的公共属性，额外的属性被忽略。
		Country:            []string{c},
		Organization:       []string{o},
		OrganizationalUnit: []string{ou},
		CommonName:         cn,
	}

	template := x509.CertificateRequest{
		Subject: subject,
	}

	pk, err := rsa.GenerateKey(rand.Reader, 4096) //生成一对具有指定字位数的RSA密钥
	if err != nil {
		fmt.Println(err)
		return
	}

	csrBytes, _ := x509.CreateCertificateRequest(rand.Reader, &template, pk)
	certOut, _ := os.Create(csrFile)
	pem.Encode(certOut, &pem.Block{Type: "CERTIFICATE REQUEST", Bytes: csrBytes, Headers: nil})
	certOut.Close()
	keyOut, _ := os.Create(keyFile)
	pem.Encode(keyOut, &pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(pk), Headers: nil})
	keyOut.Close()
	return
}

//项目的证书
func CreateCert(c, o, ou, cn string,
	caPemFile, caKeyFile, csrFile, crtFile, keyFile string,
	dnsName []string,
	expireDay int) (err error) {
	//先读取ca
	caData, err := ioutil.ReadFile(caPemFile)
	if err != nil {
		fmt.Println(err)
		return
	}
	caPemBlock, _ := pem.Decode(caData)
	if caPemBlock == nil {
		err = errors.New("cannot find the next PEM formatted block")
		fmt.Println(err)
		return
	}
	crts, err := x509.ParseCertificates(caPemBlock.Bytes)
	if len(crts) != 1 {
		fmt.Println("len(crts) != 1")
		return
	}
	caPem := crts[0]
	if !caPem.IsCA {
		fmt.Println("!caPem.IsCA")
		return
	}
	caData, err = ioutil.ReadFile(caKeyFile)
	if err != nil {
		fmt.Println(err)
		return
	}
	caKeyBlock, _ := pem.Decode(caData)
	if caKeyBlock == nil {
		err = errors.New("cannot find the next PEM formatted block")
		return
	}
	/*	password := ""
		cakeyByte, err := x509.DecryptPEMBlock(caKeyBlock, []byte(password))
		if err != nil {
			fmt.Println("line254", err)
			return
		}
		caPK, err := x509.ParsePKCS1PrivateKey(cakeyByte)
		if err != nil {
			fmt.Println("line2595", err)
			return
		}
	*/
	caPK, err := x509.ParsePKCS1PrivateKey(caKeyBlock.Bytes)
	if err != nil {
		return
	}

	//再读取csr
	csrData, err := ioutil.ReadFile(csrFile)
	if err != nil {
		fmt.Println(err)
		return
	}
	pemBlock, _ := pem.Decode(csrData)
	if pemBlock == nil {
		return
	}
	csrRaw, err := x509.ParseCertificateRequest(pemBlock.Bytes)
	if err != nil {
		fmt.Println(err)
		return
	}

	SubjectKeyId, err := GenerateSubjectKeyID(csrRaw.PublicKey)
	if err != nil {
		return
	}

	max := new(big.Int).Lsh(big.NewInt(1), 128)   //把 1 左移 128 位，返回给 big.Int
	serialNumber, _ := rand.Int(rand.Reader, max) //返回在 [0, max) 区间均匀随机分布的一个随机值
	subject := pkix.Name{                         //Name代表一个X.509识别名。只包含识别名的公共属性，额外的属性被忽略。
		Country:            []string{c},
		Organization:       []string{o},
		OrganizationalUnit: []string{ou},
		CommonName:         cn,
	}
	template := x509.Certificate{
		SubjectKeyId: SubjectKeyId,
		SerialNumber: serialNumber, // SerialNumber 是 CA 颁布的唯一序列号，在此使用一个大随机数来代表它
		Subject:      subject,
		NotBefore:    time.Now(),
		NotAfter:     time.Now().Add(time.Duration(expireDay*24) * time.Hour),
		KeyUsage:     x509.KeyUsageKeyEncipherment | x509.KeyUsageDataEncipherment | x509.KeyUsageDigitalSignature | x509.KeyUsageKeyAgreement,
		ExtKeyUsage: []x509.ExtKeyUsage{
			x509.ExtKeyUsageServerAuth,
			x509.ExtKeyUsageClientAuth,
		}, // 密钥扩展用途的序列
		RawSubject: csrRaw.RawSubject,
		//IPAddresses:  []net.IP{net.ParseIP("127.0.0.1")},
		DNSNames: dnsName,
		//EmailAddresses:[]string{"loongtime@gmail.com"},
	}

	derBytes, err := x509.CreateCertificate(rand.Reader, &template, caPem, csrRaw.PublicKey, caPK) //DER 格式
	if err != nil {
		return
	}
	certOut, _ := os.Create(crtFile)
	pem.Encode(certOut, &pem.Block{Type: "CERTIFICATE", Bytes: derBytes, Headers: nil})
	certOut.Close()
	return
}
