# ICE (Image Compression Engine)

## Installation

1. If you have glide install - `$ glide update` if not, `$ go get`
2. `$ go run main.go` or `$ go build && ./Ice -port 2000 -bucket myBucket`
3. Build for Linux (amd64) ` $ GOARCH=amd64 GOOS=linux go build && ./Ice`
4. Add the following to your profile

```
export AWS_BUCKET_NAME=""
export AWS_ACCESS_KEY_ID=""
export AWS_SECRET_ACCESS_KEY=""
```

## Api Docs

Multiple dimension image cropping.

```
POST /resize

BODY (json)
{
    "files": [
        { 
            "filename": "121.jpg", 
            "dimensions": [
                { "height": 50, "width": 50 },
                { "height": 100, "width": 80 },
                { "height": 104, "width": 80 },
                { "height": 106, "width": 80 },
                { "height": 107, "width": 80 },
                { "height": 108, "width": 80 },
                { "height": 109, "width": 80 },
                { "height": 134, "width": 80 },
                { "height": 145, "width": 80 }
            ]
        }
    ]
}
```

This will crop the same file with several different sets of dimensions.

'On the fly' image resizing.

```
GET /resize?file=my-file.jpg&width=60&height=60

```
