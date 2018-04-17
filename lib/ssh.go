package lib

import (
	"errors"
	"crypto/rsa"
	"crypto/rand"
	"encoding/pem"
	"crypto/x509"
	"os"
	"bytes"
	"golang.org/x/crypto/ssh"
	"io/ioutil"
	"time"
	"fmt"
	"log"
)

type SSH struct {
	Username string          // 验证用户信息的用户名
	Password string          // 验证用户信息的密码
	Host string              // 服务器地址
	Port int                 // 服务器端口
	PrivateKeyFile string    // 设置用以验证的私钥地址
	KeyPairSavePath string   // 密钥对保存路径，请务必在路劲最后不要带上斜杠
	PrivateKey string        // 私钥，生成密钥对是保存的私钥
	PublicKey string         // 公钥，生成密钥对时保存的公钥
}
/******************************************
*              通用方法                   *
*                                         *
*******************************************/
// 通过rsa 生成支持ssh访问的密钥对
// 感谢 https://github.com/nanobox-io/golang-ssh/blob/master/key.go 对方法写成的帮助
func (c *SSH) GenKeyPair() error {
	// 私钥
	privateKey,err := rsa.GenerateKey(rand.Reader,2048)
	if err != nil{
		return err
	}
	privateKeyPEM:=&pem.Block{
		Type:"RSA PRIVATE KEY",
		Bytes:x509.MarshalPKCS1PrivateKey(privateKey),
	}
	var private bytes.Buffer
	if err := pem.Encode(&private,privateKeyPEM);err != nil{
		return err
	}
	// 公钥
	pub,err := ssh.NewPublicKey(&privateKey.PublicKey)
	if err != nil {
		return err
	}
	public := ssh.MarshalAuthorizedKey(pub)

	// 保存key pair
	if c.KeyPairSavePath != ""{
		// 保存私钥
		privateFile,err := os.Create(c.KeyPairSavePath+"/id_rsa")
		if err != nil{
			return err
		}
		defer privateFile.Close()
		pem.Encode(privateFile,privateKeyPEM)

		// 保存公钥

		pubFile,err := os.Create(c.KeyPairSavePath+"/id_rsa.pub")
		if err != nil{
			return err
		}
		defer pubFile.Close()
		_,err = pubFile.Write(public)
		if err != nil{
			return err
		}
	}
	c.PrivateKeyFile = private.String()
	c.PublicKey = string(public)
	return nil
}


// 生成rsa
func GenerateSSHKeyPair(savepath string) (string,string,error) {
	// 生成私钥文件
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return "","",err
	}
	derStream := x509.MarshalPKCS1PrivateKey(privateKey)
	block := &pem.Block{
		Type:  "私钥",
		Bytes: derStream,
	}
	file, err := os.Create("tmp/private.pem")
	if err != nil {
		return "","",err
	}
	err = pem.Encode(file, block)
	if err != nil {
		return "","",err
	}
	// 生成公钥文件
	publicKey := &privateKey.PublicKey
	derPkix, err := x509.MarshalPKIXPublicKey(publicKey)
	if err != nil {
		return "","",err
	}
	block = &pem.Block{
		Type:  "公钥",
		Bytes: derPkix,
	}
	file, err = os.Create("tmp/public.pem")
	if err != nil {
		return "","",err
	}
	err = pem.Encode(file, block)
	if err != nil {
		return "","",err
	}

	return "","",nil
}

// TODO 密钥对验证



/******************************************
*            getter setter  方法区        *
*                                         *
*******************************************/
// 设置密钥对保存路径
func (c *SSH) SetKeyPairSavePath(path string) *SSH {
	c.KeyPairSavePath = path
	return c
}

// 设置用户名
func (c *SSH)SetUsername(username string) (*SSH){
	c.Username = username
	return c
}

//设置密码
func (c *SSH) SetPassword(password string) (*SSH) {
	c.Password = password
	return c
}

//设置服务器地址
func (c *SSH) SetHost(host string) (*SSH) {
	c.Host = host
	return c
}

//设置端口
func (c *SSH)SetPort(port int) (*SSH) {
	c.Port = port
	return c
}

// 设置用以验证的私钥文件
func (c *SSH)SetPrivateKeyFile(privateKeyFilePath string) (*SSH) {
	c.PrivateKeyFile = privateKeyFilePath
	return c
}

/******************************************
*                                         *
*                                         *
*******************************************/
// 检测ssh配置
func (c *SSH) CheckConfig() (error) {
	if c.Host ==""{
		return errors.New("请设置主机")
	}
	if c.Port == 0 {
		c.Port = 22
	}
	if c.Username == ""{
		return errors.New("请设置用户名")
	}
	if c.PrivateKeyFile == "" || c.Password == ""{
		return errors.New("请设置使用秘钥还是使用密码验证，请务必设置其中一项")
	}
	return nil
}

// ssh链接设置，感谢：http://blog.ralch.com/tutorial/golang-ssh-connection/ 的样例
func (c *SSH)Connect() (*ssh.Session,error) {
	var (
		auth []ssh.AuthMethod
		addr string
		clientConfig *ssh.ClientConfig
		client *ssh.Client
		session *ssh.Session
		err error
	)
	auth = make([]ssh.AuthMethod,0)
	if c.Password != ""{
		auth = append(auth,ssh.Password(c.Password))
	}else {
		buffer,err := ioutil.ReadFile(c.PrivateKeyFile)
		if err != nil{
			return nil,err
		}
		key,err := ssh.ParsePrivateKey(buffer)
		if err != nil{
			return nil,err
		}
		auth = append(auth,ssh.PublicKeys(key))
	}

	clientConfig = &ssh.ClientConfig{
		User:c.Username,
		Auth:auth,
		Timeout:30 * time.Second,
	}

	addr = fmt.Sprintf("%s:%d",c.Host,c.Port)
	if client,err = ssh.Dial("tcp",addr,clientConfig);err != nil{
		return nil,err
	}
	if session,err =  client.NewSession();err != nil{
		return nil,err
	}
	return session,nil
}

func (c *SSH) Run(cmd string) (string) {
	session,err := c.Connect()
	if err != nil{
		log.Fatal(err)
	}
	defer session.Close()
	var stuoutBuf bytes.Buffer
	session.Stdout = &stuoutBuf
	session.Run(cmd)
	fmt.Println(stuoutBuf.String())
	return stuoutBuf.String()
}