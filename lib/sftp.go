package lib

import (
	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
	"io/ioutil"
	"time"
	"fmt"
	"os"
	"path"
)

type SFTP struct {
	Username string
	Password string
	Host string
	Port int
	PrivateKeyFile string
	PutGet string                  // put,get
	Localfile string
	Remotefile string
}

func (c *SFTP)Connect() (*sftp.Client,error) {
	var (
		auth []ssh.AuthMethod
		addr string
		clientConfig *ssh.ClientConfig
		sshClient *ssh.Client
		sftpClient *sftp.Client
		err error
	)
	auth = make([]ssh.AuthMethod,0)
	if c.Password == ""{
		buffer,err := ioutil.ReadFile(c.PrivateKeyFile)
		if err != nil{
			return nil,err
		}
		key,err := ssh.ParsePrivateKey(buffer)
		if err != nil{
			return nil,err
		}
		auth = append(auth,ssh.PublicKeys(key))
	}else {
		auth = append(auth,ssh.Password(c.Password))
	}
	clientConfig = &ssh.ClientConfig{
		User:c.Username,
		Auth:auth,
		Timeout:30 * time.Second,
	}
	addr = fmt.Sprintf("%s:%d",c.Host,c.Port)
	fmt.Println(addr)
	if sshClient,err = ssh.Dial("tcp",addr,clientConfig);err != nil{
		return nil,err
	}
	if sftpClient,err =  sftp.NewClient(sshClient);err!= nil{
		return nil,err
	}
	return sftpClient,nil
}

func (c *SFTP) Put(localFile string) *SFTP {
	c.PutGet = "put"
	c.Localfile = localFile
	return c
}

func (c *SFTP) GetTo(localFile string) *SFTP {
	c.PutGet = "get"
	c.Localfile = localFile
	return c
}

func (c *SFTP) To(RemoteOrLocalPath string)  {
	sftpclient,err := c.Connect()
	if err != nil{
		fmt.Println(err)
	}
	defer sftpclient.Close()
	srcFile,err := os.Open(c.Localfile)
	if err != nil{
		fmt.Println(err)
	}
	defer srcFile.Close()
	var remoteFileName = path.Base(c.Localfile)
	dstFile,err := sftpclient.Create(path.Join(RemoteOrLocalPath,remoteFileName))
	if err != nil{
		fmt.Println(err)
	}
	defer dstFile.Close()
	buf := make([]byte,1024)
	for {
		n,_ := srcFile.Read(buf)
		if n == 0 {
			break
		}
		dstFile.Write(buf)
	}

}

func (c *SFTP) From(remoteFilePath string)  {
	sftpclient,err := c.Connect()
	if err != nil{
		fmt.Println(err)
		os.Exit(1)
	}
	defer sftpclient.Close()
	fmt.Println(sftpclient)
	srcFile,err := sftpclient.Open(remoteFilePath)
	if err != nil{
		fmt.Println(err)
		os.Exit(1)
	}
	defer srcFile.Close()
	var localFileName = path.Base(remoteFilePath)
	localfilePath := path.Join(c.Localfile,localFileName)
	dstFile,err := os.Create(localfilePath)
	if err != nil{
		fmt.Println(err)
	}
	defer dstFile.Close()
	_,err = srcFile.WriteTo(dstFile)
	if err != nil{
		fmt.Println(err)
	}
}