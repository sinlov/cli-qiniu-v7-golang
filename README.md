[TOC]

# this is qiniu golang CLI utils

- os OSX 10.11.6
- golang version 1.9.2

# build

```sh
./mutil_build.sh
```
- this script will build `OS X Linux Windows` at default path at `${ProjectRoot}/build`

# use

add you qiniu key at `main.go`

```golang
	// set your qiniu set
	downloadURLHead string = "http://xx.clouddn.com/"
	// set accessKey
	accessKeyDefault string = ""
	// set secretKey
	secretKeyDefault string = ""
	// set your want use bucket
	bucketDefault string = ""
```

```
go build main main.go
./main -h
./main -l [file for upload]
```

or use by main -h

###License

---

Copyright 2017 sinlovgm@gmail.com

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

   http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
