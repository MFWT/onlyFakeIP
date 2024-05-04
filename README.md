# onlyFakeIP

一个**只会**把域名映射到特定fakeIP的DNS服务器

A DNS server that **ONLY** maps domain names to a specific fakeIP

**本程序只是作者随手写的一个小工具，除了程序可以运行之外，作者不对其性能和安全性做任何担保，用于实际环境前请三思！**

**Except that the program can be run, the author does not make any guarantees about its performance and security. Use it in actual environments at your own risk!**

# How?

原理其实非常简单：**对域名执行cityHash32算法之后右移8位并加上一个前缀，从而组成了一个有效的IP地址，来作为该域名的A记录应答**

The principle is actually very simple: **After performing the cityHash32 algorithm on the domain name, shift it to the right by 8 bits and add a prefix, thereby forming a valid IP address as the A record response for the domain name.**

# For?

如果你用过SNI代理，你会明白我在说什么：

现代浏览器，在特定情况下（某两个域名解析到的IP一样，且前者使用的证书CN中包括了后面的域名）会重用之前使用过但还未断开的连接。对于SNI代理的场景，可能会出现流量被错误路由的情况，导致网页无法打开。

虽然正确的实现应该是：服务器在收到类似的请求后，返回HTTP状态码421，此时客户端应该放弃之前的连接，新开连接再次请求。但是考虑到互联网上还有鬼知道多少的服务器没有正确实现这个代码，与其被动规避，倒不如一开始就避免额外的连接重用，导致不必要的麻烦。

本程序单纯用于实现**不同域名可以解析出大致不同的IP**这个目标，因为他的目标场景是SNI代理而不是普通的透明代理（只要求解析出来的IP可以路由到SNI代理服务器），没有从IP反向解析回域名这个需求，所以一方面没有设计和代理程序的联动功能，另一方面也允许『两个不同的域名解析出来同一个IP』的情况出现（但是概率应该不大）。

对于同一个二层局域网内的，运行Linux的软路由，要想让他对一个IP段都做出反应，**不妨试试内核支持的AnyIP功能**：

```shell
ip route add local 11.0.0.0/8 dev eth0 #eth0是你的内网网卡
```

If you've used SNI agents, you'll know what I'm talking about:

Modern browsers, under certain circumstances (the IPs resolved to the same two domain names, and the certificate CN used by the former includes the latter domain name) will reuse previously used but not yet disconnected connections. For SNI proxy scenarios, traffic may be incorrectly routed, causing the web page to fail to open.

Although the correct implementation should be: after receiving a similar request, the server returns HTTP status code 421. At this time, the client should abandon the previous connection, open a new connection and request again. But considering that there are still countless servers on the Internet that have not implemented this code correctly, instead of passively avoiding it, it is better to avoid unnecessary trouble caused by additional connection reuse in the first place.

This program is purely used to achieve the goal of **different domain names can be resolved to roughly different IPs**, because its target scenario is an SNI proxy rather than an ordinary transparent proxy (it only requires that the resolved IP can be routed to the SNI proxy server) , there is no need to reversely resolve the domain name from the IP, so on the one hand, there is no linkage function with the agent program, and on the other hand, it also allows the situation of "two different domain names to resolve to the same IP" (but the probability should not be high).

For a soft router running Linux, if you want it to respond to an IP segment, you might as well try the AnyIP function supported by the kernel:

```shell
ip route add local 11.0.0.0/8 dev eth0 # eth0 is your LAN interfere
```

# Install

## Pre-built binary packages

到releases页面下载即可

Go to releases page and download it

## Compiling

`go build -o onlyFakeIP main.go`

Windows下运行build.bat也行

or run the `build.bat` on Windows

# Usage
```shell
$ ./onlyFakeIP -h

Usage of onlyFakeIP:
  -b string
        Set the listening address (default is dual-stack)
  -h    show this help
  -p port
        Set the listening port (default is 53) (default "53")
  -prefix prefix
        Set the prefix of the fakeIP (default is 11 for 11.x.x.x) (default "11")
```

程序默认使用的前缀是11.0.0.0，这个段隶属于美国国防部（也是很经常被『公网私用』的一个段），如果你担心隐私问题，不妨修改为10.0.0.0，或者某个你可能一辈子都用不上的IP段

The default prefix used by the program is 11.0.0.0. This segment belongs to the U.S. Department of Defense (it is also a segment that is often used for "public IP using on private networks"). If you don’t really want to be **accidentally** spied on by the DoD , you may wish to change it to 10.0.0.0, or someone you An IP segment that may never be used in a lifetime

**推荐先使用SmartDNS等DNS缓存与转发程序**，把确实需要被代理的域名分流发到本程序，以减轻发生意外的可能性

It is recommended to first use a DNS caching and forwarding program such as SmartDNS to offload the domain names that really need to be proxied to this program to reduce the possibility of accidents.

# License

**MIT**
