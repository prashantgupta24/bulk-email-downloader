# Bulk-email-downloader
A bulk email downloader from an IMAP server written in go. It downloads bulk emails from an `IMAP` server using configurable number of go-routines. It saves them in an Array called `Messages`. 

### Config

Create a `.env` file in the same directory as the `main.go` file, and fill in the following properties:
```
IMAP_SERVER=
EMAIL=
PASSWORD=
```

### Example output

```
go run main.go

2018/07/11 16:24:35 Connecting to server...
2018/07/11 16:24:37 Logged in
2018/07/11 16:24:37 Fetching emails from number 501:1000
2018/07/11 16:24:37 Fetching emails from number 1:500
2018/07/11 16:24:37 Fetching emails from number 1001:1500
2018/07/11 16:24:42 Logging out
2018/07/11 16:24:42 Read 1216 messages
2018/07/11 16:24:42 Time taken 6.814280051s
```
