package main

import (
	"context"
	"crypto/md5"
	"fmt"
	"io"
	"os"

	"github.com/mkideal/cli"
	"github.com/qiniu/api.v7/auth/qbox"
	"github.com/qiniu/api.v7/storage"
)

const (
	versionName string = "1.0.0"
	commInfo    string = "qiniu sdk v7 utils"
	// set your qiniu set
	downloadURLHead string = "http://xx.clouddn.com/"
	// set accessKey
	accessKeyDefault string = ""
	// set secretKey
	secretKeyDefault string = ""
	// set your want use bucket
	bucketDefault string = ""
)

type config struct {
	A string
	B int
	C bool
}

type myPutRet struct {
	Key    string
	Hash   string
	Fsize  int
	Bucket string
	Name   string
}

type filterCLI struct {
	cli.Helper
	Version   bool   `cli:"version" usage:"version"`
	Verbose   bool   `cli:"verbose" usage:"see Verbose of utils"`
	AccessKey string `cli:"a,accessKey" usage:"qiniu AccessKey"`
	SecretKey string `cli:"s,secretKey" usage:"qiniu SecretKey"`
	Bucket    string `cli:"b,bucket" usage:"qiniu Bucket"`
	Expires   uint32 `cli:"e,expires" usage:"Expires as key" dft:"3600"`
	Overwrite string `cli:"o,overwrite" usage:"qiniu keyToOverwrite"`
	LocalFile string `cli:"l,localFile" usage:"want upload localFile"`
	FileKey   string `cli:"k,fileKey" usage:"want save fileKey"`
	FileName  string `cli:"n,fileName" usage:"want show fileName"`
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Warning you input is error pleae use -h to see help")
		os.Exit(-1)
	}
	cli.Run(new(filterCLI), func(ctx *cli.Context) error {
		argv := ctx.Argv().(*filterCLI)
		if argv.Version {
			ctx.String(commInfo + "\n\tversion: " + versionName)
			os.Exit(0)
		}
		var accessKey string
		if argv.AccessKey != "" {
			accessKey = argv.AccessKey
		} else {
			fmt.Printf("use default accessKeyDefault: %s\n", accessKeyDefault)
			accessKey = accessKeyDefault
		}
		var secretKey string
		if argv.SecretKey != "" {
			secretKey = argv.SecretKey
		} else {
			fmt.Printf("use default secretKeyDefault: %s\n", secretKeyDefault)
			secretKey = secretKeyDefault
		}

		var bucket string
		if argv.Bucket != "" {
			bucket = argv.Bucket
		} else {
			fmt.Printf("use default bucketDefault: %s\n", bucketDefault)
			bucket = bucketDefault
		}

		var localFile string
		if argv.LocalFile != "" {
			localFile = argv.LocalFile
		} else {
			ctx.String("not set LocalFile exit 1")
			os.Exit(1)
		}

		var fileKey string
		if argv.FileKey != "" {
			fileKey = argv.FileKey
		} else {
			fmt.Printf("localFile: %s\n", localFile)
			hashKey, err := readFileHash(localFile)
			if err != nil {
				return err
			} else {
				fileKey = hashKey
			}

		}

		var fileName string
		if argv.FileName != "" {
			fileName = argv.FileName
		} else {
			fmt.Printf("use fileName to fileKey: %s\n", fileKey)
			fileName = fileKey
		}

		var expires uint32
		expires = 3600
		if argv.Expires > 0 {
			expires = argv.Expires
		}

		var overwrite string
		if argv.Overwrite != "" {
			overwrite = argv.Overwrite
		}

		//check set
		fmt.Printf("set \n-> accessKey: %v\n-> secretKey: %v\n-> bucket: %v\n-> expires: %v\n", accessKey, secretKey, bucket, expires)
		fmt.Printf("local \n-> fileKey: %v\n-> fileName: %v\n", fileKey, fileName)
		fmt.Printf("-> update file: %v\n", localFile)
		fmt.Printf("change \n-> overwrite: %v\n", overwrite)

		//os.Exit(0)
		// start post
		var putPolicy storage.PutPolicy
		if overwrite != "" {
			putPolicy = storage.PutPolicy{
				Scope:      fmt.Sprintf("%s:%s", bucket, overwrite),
				ReturnBody: `{"key":"$(key)","hash":"$(etag)","fsize":$(fsize),"bucket":"$(bucket)","name":"$(x:name)"}`,
			}
		} else {
			putPolicy = storage.PutPolicy{
				Scope:      bucket,
				ReturnBody: `{"key":"$(key)","hash":"$(etag)","fsize":$(fsize),"bucket":"$(bucket)","name":"$(x:name)"}`,
			}
		}

		putPolicy.Expires = expires
		mac := qbox.NewMac(accessKey, secretKey)
		upToken := putPolicy.UploadToken(mac)
		cfg := storage.Config{}
		formUploader := storage.NewFormUploader(&cfg)
		ret := myPutRet{}
		putExtra := storage.PutExtra{
			Params: map[string]string{
				"x:name": fileName,
			},
		}
		err := formUploader.PutFile(context.Background(), &ret, upToken, fileKey, localFile, &putExtra)
		if err != nil {
			fmt.Println(err)
			return err
		}
		fmt.Printf("upload file success\n")
		fmt.Printf("-> Bucket: %s\n-> Key: %s\n-> \n-> FileName: %sFileSize: %d\n-> Hash: %v\n-> Download: %v\n", ret.Bucket, ret.Key, ret.Name, ret.Fsize, ret.Hash, fmt.Sprintf("%s%s", downloadURLHead, ret.Key))
		fmt.Println(ret.Bucket, ret.Key, ret.Fsize, ret.Hash, ret.Name)
		return nil
	})
}

func readFileHash(pathFile string) (string, error) {
	f, err := os.Open(pathFile)
	if err != nil {
		fmt.Println("Open", err)
		return "", err
	}
	defer f.Close()
	md5hash := md5.New()
	if _, err := io.Copy(md5hash, f); err != nil {
		fmt.Println("Copy", err)
		return "", err
	}
	md5hash.Sum(nil)
	return fmt.Sprintf("%x", md5hash.Sum(nil)), nil
}
