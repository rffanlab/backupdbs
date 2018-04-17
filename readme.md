# 由来
由于在windows 下使用sftp 很麻烦，所以自己开发了这么一波。用来定时从远程Linux服务器上拉取备份数据库文件  
另外，由于是独立的ssh key文件，因此，即便是key被截获了，~~也不会影响程序所在服务器的安全，事实上scp也没有影响！所以这点没鸟用。~~。
# 是需要添加的库
go get -u github.com/pkg/sftp  
go get -u golang.org/x/crypto/ssh  
# 使用
一切的精华都在配置文件里面。因为配置文件里面。  
将setting.conf.example文件修改为setting.conf 然后填写一下里面的配置文件就OK了。