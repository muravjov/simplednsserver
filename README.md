# Simple DNS Server with Go

To run the program you need to install https://golang.org/doc/install first.

```bash
    $ go get github.com/muravjov/simplednsserver
    $ simplednsserver --a-record example.com:1.2.3.4
    2019/10/27 22:53:30 simplednsserver listening at :8053...

    $ dig @localhost -p 8053 +short example.com
    1.2.3.4
```
