# ICE (Image Compression Engine)

## Installation

1. If you have glide install - `$ glide update` if not, `$ go get`
2. `$ go run main.go` or `$ go build && ./Ice`
3. Build for Linux (amd64) ` $ GOARCH=amd64 GOOS=linux go build && ./Ice`
4. Add the following to your profile

```
export AWS_S3_BUCKET=""
export AWS_ACCESS_KEY_ID=""
export AWS_SECRET_ACCESS_KEY=""
```

## Api Docs

Multiple dimension image cropping.

```
POST /resize

BODY
file (file)
width[] (integer)
height[] (integer)
width[]
height[]
```
This will crop the same file with several different sets of dimensions.

'On the fly' image resizing.

```
GET /resize?file=my-file.jpg&width=60&height=60

```
