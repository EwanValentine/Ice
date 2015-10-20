# ICE (Image Compression Engine)

## Installation

1. If you have glide install - `$ glide update` if not, `$ go get`
2. `$ go run main.go` or `$ go build && ./Ice`

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
