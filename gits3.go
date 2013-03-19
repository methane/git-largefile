package main

import (
	"crypto/sha1"
	"fmt"
	"io/ioutil"
	//"launchpad.net/goamz/s3"
	"log"
	"os"
	"os/user"
	"path/filepath"
)

func cachePath(sha1hex string) (dirpath, filename string) {
	usr, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}
	dirpath = filepath.Join(usr.HomeDir, ".gitasset", "data", string(sha1hex[0:2]), string(sha1hex[2:4]))
	filename = string(sha1hex[4:])
	return
}

func store() {
	contents, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		log.Fatal(err)
	}
	sum := sha1.New()
	sum.Write(contents)
	sha1hex := fmt.Sprintf("%x", sum.Sum(nil))
	//log.Println("sha1 =", sha1hex)
	dirpath, filename := cachePath(sha1hex)
	err = os.MkdirAll(dirpath, os.ModePerm)
	if err != nil {
		log.Fatal(err)
	}
	err = ioutil.WriteFile(filepath.Join(dirpath, filename), contents, os.ModePerm)
	if err != nil {
		log.Fatal(err)
	}
}

func load() {
	hash, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		log.Fatal(err)
	}
	dirpath, filename := cachePath(string(hash))
	contents, err := ioutil.ReadFile(filepath.Join(dirpath, filename))
	if err != nil {
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
