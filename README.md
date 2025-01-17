## xping
xping: a tcp ping programming write with the help of ChatGTP.


### Install

```bash

go install github.com/bopjiang/xping@latest

```

### Usage

```bash

xping -h www.163.com -p 443
PING www.163.com (113.96.128.244) on TCP port 443
From 113.96.128.244: tcp_seq=1 Port 443 open time=43.590 ms
From 183.3.205.208: tcp_seq=2 Port 443 open time=103.197 ms
From 59.36.214.237: tcp_seq=3 Port 443 open time=109.716 ms
From 183.3.205.210: tcp_seq=4 Port 443 open time=339.315 ms

--- www.163.com:443 ping statistics ---
4 packets transmitted, 4 successful, 0 failed

DNS Resolution Statistics:
min/avg/max = 0.822/2.020/2.734 ms

TCP Connection Statistics:
min/avg/max = 42.767/146.935/336.581 ms

Total Time Statistics:
min/avg/max = 43.590/148.954/339.315 ms

```