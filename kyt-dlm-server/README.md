# Go API Server for kyt-dlm-server

This is the KYT API Server

## Overview

### Running the server

To run the server, follow these simple steps:

```
go run main.go
```

To run the server in a docker container
```
docker build --network=host -t  .
```

Once the image is built, just run
```
docker run --rm -it kyt-dlm-server
```

### Known Issue

Endpoints sharing a common path may result in issues. For example, `/v2/pet/findByTags` and `/v2/pet/:petId` will result in an issue with the Gin framework. For more information about this known limitation, please refer to [gin-gonic/gin#388](https://github.com/gin-gonic/gin/issues/388) for more information.

A workaround is to manually update the path and handler. Please refer to [gin-gonic/gin/issues/205#issuecomment-296155497](https://github.com/gin-gonic/gin/issues/205#issuecomment-296155497) for more information.
