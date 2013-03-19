package main

import (
	"crypto/sha1"
	"fmt"
	"github.com/msbranco/goconfig"
	"io/ioutil"
	"launchpad.net/goamz/aws"
	"launchpad.net/goamz/s3"
	"log"
	"os"
	"os/user"
	"path/filepath"
	"strings"
)

func assetDir() string {
	usr, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}
	return filepath.Join(usr.HomeDir, ".gitasset")
}

func getConfig() *goconfig.ConfigFile {
	confFile := filepath.Join(assetDir(), "gits3.ini")
	conf, err := goconfig.ReadConfigFile(confFile)
	if err != nil {
		log.Fatal(err)
	}
	return conf
}

func getBucket() *s3.Bucket {
	conf := getConfig()
	awskey, err := conf.GetString("default", "awskey")
	if err != nil {
		log.Fatal(err)
	}
	bucketName, err := conf.GetString("default", "bucket")
	if err != nil {
		log.Fatal(err)
	}

	key_secret := strings.Split(awskey, ":")
	if len(key_secret) != 2 {
		log.Fatal("Bad awskey:" + awskey)
	}
	auth := aws.Auth{key_secret[0], key_secret[1]}
	return s3.New(auth, aws.USEast).Bucket(bucketName)
}

func cachePath(sha1hex string) (dirpath, filename string) {
	dirpath = filepath.Join(assetDir(), "data", string(sha1hex[0:2]), string(sha1hex[2:4]))
	filename = string(sha1hex[4:])
	return
}

func storeToS3(hex string, data []byte) error {
	bucket := getBucket()
	_, err := bucket.GetReader(hex)
	if err == nil {
		log.Println("Already exists in S3: ", hex)
		return err
	}
	log.Println(err)
	return bucket.Put(hex, data, "application/octet-stream", s3.Private)
}

func loadFromS3(hex string) ([]byte, error) {
	bucket := getBucket()
	return bucket.Get(hex)
}

func storeToCache(hex string, data []byte) {
	dirpath, filename := cachePath(hex)
	filePath := filepath.Join(dirpath, filename)
	_, err := os.Lstat(filePath)
	if os.IsExist(err) {
		return
	}
	err = os.MkdirAll(dirpath, os.ModePerm)
	if err != nil {
		log.Fatal(err)
	}
	err = ioutil.WriteFile(filePath, data, os.ModePerm)
	if err != nil {
		log.Fatal(err)
	}
}

func loadFromCache(hex string) ([]byte, error) {
	dirpath, filename := cachePath(hex)
	return ioutil.ReadFile(filepath.Join(dirpath, filename))
}

func store() {
	contents, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		log.Fatal(err)
	}
	sum := sha1.New()
	sum.Write(contents)
	sha1hex := fmt.Sprintf("%x", sum.Sum(nil))
	log.Println("sha1=", sha1hex)
	storeToCache(sha1hex, contents)
	storeToS3(sha1hex, contents)
}

func load() {
	hash, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		log.Fatal(err)
	}
	hex := string(hash)
	contents, err := loadFromCache(hex)
	if os.IsNotExist(err) {
		contents, err = loadFromS3(hex)
		if err != nil {
			log.Fatal(err)
		}
		storeToCache(hex, contents)
	} else if err != nil {
		log.Fatal(err)
	}
	n, err := os.Stdout.Write(contents)
	if n != len(contents) || err != nil {
		log.Fatal(err)
	}
}

func main() {
	log.SetPrefix("gits3:")
	if len(os.Args) < 2 {
		log.Fatal("Invalid argument.")
	}
	switch os.Args[1] {
	case "store":
		store()
	case "load":
		load()
	default:
		log.Fatal("Invalid argument.")
	}
}
