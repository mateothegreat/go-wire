# msgpack over tcp

```bash
cd receiver && go test
cd sender && go test
```

```bash
../wire/sender 🌱 main ➜ go test
2024/06/01 22:57:54 Connected to server
2024/06/01 22:57:54 Sending data length: 1321879
PASS
ok      github.com/mateothegreat/go-wire/sender 0.008s

../wire/receiver 🌱 main ➜ go test
2024/06/01 22:57:48 Listening on port 15000
2024/06/01 22:57:54 Accepted connection from 127.0.0.1:56466
2024/06/01 22:57:54 Reading data of length: 1321879
2024/06/01 22:57:54 Received 1321879 bytes
2024/06/01 22:57:54 Connection closed by client
2024/06/01 22:57:54 Received image from test, len=1321856
```