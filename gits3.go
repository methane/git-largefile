package main

import (
	"crypto/sha1"
	"flag"
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

var (
	configSection string
	assetDirPath  string
	localMode     bool
)

func init() {
	flag.StringVar(&configSection, "section", "default", "Section name of config file.")
	flag.StringVar(&assetDirPath, "assetpath", "~/.gitasset", "Asset directory")
	flag.BoolVar(&localMode, "local", false, "Local mode (not use S3)")
}

func assetDir() string {
	p := assetDirPath
	if p[0:2] == "~/" {
		usr, err := user.Current()
		if err != nil {
			log.Fatal(err)
		}
		p = filepath.Join(usr.HomeDir, p[2:])
	}
	return p
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
	awskey, err := conf.GetString(configSection, "awskey")
	if err != nil {
		log.Fatal(err)
	}
	bucketName, err := conf.GetString(configSection, "bucket")
	if err != nil {
		log.Fatal(err)
	}

	key_secret := strings.Split(awskey, ":")
	if len(key_secret) != 2 {
		log.Fatal("Bad awskey:" + awskey)
	}
	auth := aws.Auth{key_secret[0], key_secret[1]}
	return s3.New(auth, aws.APNortheast).Bucket(bucketName)
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

func calcSha1String(data []byte) string {
	sum := sha1.New()
	sum.Write(data)
	return fmt.Sprintf("%x", sum.Sum(nil))
}

func store() {
	data, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		log.Fatal(err)
	}
	sha1hex := calcSha1String(data)
	storeToCache(sha1hex, data)
	if !localMode {
		if err = storeToS3(sha1hex, data); err != nil {
			log.Fatal(err)
		}
	}
	writeStdout([]byte(sha1hex))
}

func isValidHash(hex string) bool {
	if len(hex) != 40 {
		log.Println("warn: hash length is ", len(hex))
		return false
	}
	for _, c := range hex {
		if '0' <= c && c <= '9' {
			continue
		}
		if 'a' <= c && c <= 'f' {
			continue
		}
		log.Println("warn: hash contains ", c, hex)
		return false
	}
	return true
}

func writeStdout(contents []byte) {
	_, err := os.Stdout.Write(contents)
	if err != nil {
		log.Fatal(err)
	}
}

func load() {
	hash, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		log.Fatal(err)
	}
	hex := string(hash)
	if !isValidHash(hex) {
		writeStdout(hash)
		return
	}
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
	writeStdout(contents)
}

func upload() {
	datadir := filepath.Join(assetDir(), "data")
	store := func(path string, info os.FileInfo, err error) error {
		if err != nil {
			log.Println("Skip:", path, err)
			return nil // Skip this directory.
		}
		if info.IsDir() {
			return nil
		}
		data, err := ioutil.ReadFile(path)
		if err != nil {
			log.Println("Skip:", path, err)
			return nil // Skip this file.
		}
		sha1hex := calcSha1String(data)
		return storeToS3(sha1hex, data)
	}
	filepath.Walk(datadir, store)
}

const usageStr = `Usage:
  gits3 [options] store < input-file > shafile
  gits3 [options] load < shafile > output-file
  gits3 [options] upload

Options:`

func usage() {
	fmt.Fprintln(os.Stderr, usageStr)
	flag.PrintDefaults()
}

func main() {
	log.SetPrefix("gits3:")
	flag.Parse()
	if flag.NArg() != 1 {
		usage()
		os.Exit(1)
	}
	switch flag.Args()[0] {
	case "store":
		store()
	case "load":
		load()
	case "upload":
		upload()
	default:
		log.Fatal("Invalid argument.")
	}
}
