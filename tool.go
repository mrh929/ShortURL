package main

import (
	"crypto/md5"
	"encoding/binary"
	"io"
	"math/rand"
	"net/url"
)

func parseRawURL(rawurl string) (corr_url string, err error) {
	u, err := url.ParseRequestURI(rawurl)
	if err != nil || u.Host == "" {
		u, err = url.ParseRequestURI("https://" + rawurl)
		if err != nil {
			return
		}
	}
	corr_url = u.String()
	// re := regexp.MustCompile(`(http|ftp|https):\/\/[\w\-_]+(\.[\w\-_]+)+([\w\-\.,@?^=%&:/~\+#]*[\w\-\@?^=%&/~\+#])?`)
	// result := re.FindAllStringSubmatch(corr_url, -1)
	// if result == nil {
	// 	fmt.Println("illegal URL")
	// }
	return
}

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

func randSeq(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func keyGenerate(url string) (rand_key string, needInsert bool, err error) {
	h := md5.New()
	io.WriteString(h, url)
	var seed uint64 = binary.BigEndian.Uint64(h.Sum(nil))
	rand.Seed(int64(seed))

	for {
		rand_key = randSeq(5) // generate one key based on url
		r_url, err_ := urlSelect(rand_key)
		if err_ != nil { // handler error
			err = err_
			return
		}

		if r_url == url { // if the same url is in sql, ignore it
			needInsert = false
			return
		}

		if r_url == "" {
			break
		}
	}
	needInsert = true
	return
}
