# clist
一个简单的常用命令console UI工具, 启动后显示一列常用的命令，选择需要执行的，Enter 键选择执行。

工作中由于需要ssh到很多台服务器上，每次都只能用history | grep ssh 来查找，略微繁琐

下载下来build之后，将clist和build出来的可执行文件放到/bin目录下，chmod a+x /bin/clist, 然后执行 clist 即可

由于用的Go的CUI库不怎么会用，目前只能通过手工输入命令到 HOME 目录下的 .clist 文件夹里的 clist 文件, 将你需要经常使用的命令写入到这个文件，一行写一条即可。

目前还支持注释说明功能， 譬如 `ssh shahuwang@localhost|登录到服务器`, 以 | 做分割
