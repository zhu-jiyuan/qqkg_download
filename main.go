package main

import (
	"encoding/json"
	"fmt"
	"github.com/schollz/progressbar/v3"
	"io"
	"net/http"
	"os"
	"regexp"
	"sync"
	"time"
)

const DOWNLOAD_URL = "http://cgi.kg.qq.com/fcgi-bin/fcg_get_play_url?shareid="
const DOWNLOAD_DIR = "./downloads/"

func downloadMusic(shareid string, title string, nickname string, wg *sync.WaitGroup) {
	url := DOWNLOAD_URL + shareid
	fileName := title + "-" + nickname + ".m4a"
	filePath := DOWNLOAD_DIR + fileName
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("http.Get error: ", err)
		return
	}
	defer resp.Body.Close()

	file_m4a, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY, 0644)
	defer file_m4a.Close()

	// content, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("OpenFile error: ", err)
		return
	}

	bar := progressbar.DefaultBytes(
		resp.ContentLength,
		fileName,
	)
	io.Copy(io.MultiWriter(file_m4a, bar), resp.Body)

	wg.Done()
}

func appendUgc(uid string, ugclist chan Ugc, startCount int, wg *sync.WaitGroup) {
	ugcList, _, _ := getMusicList(uid, startCount)
	for _, ugc := range ugcList {
		ugclist <- ugc
	}
	wg.Done()
}

func startAddUgc(uid string, ugcList chan Ugc, ugcTotalCount int) {
	var requestCount int = ugcTotalCount / 10
	if ugcTotalCount%10 != 0 {
		requestCount += 1
	}

	defer close(ugcList)
	wg := sync.WaitGroup{}
	wg.Add(requestCount)
	for i := 1; i <= requestCount; i++ {
		go appendUgc(uid, ugcList, i, &wg)
	}
	wg.Wait()
}

func getMusicList(uid string, startCount int) (ugcList []Ugc, nickname string, ugcTotalCount int) {
	tsMs := time.Now().UnixMilli()
	url := fmt.Sprintf("https://node.kg.qq.com/fcgi-bin/kg_ugc_get_homepage?outCharset=utf-8&from=1&nocache=%d&format=json&type=get_uinfo&start=%d&num=10&share_uid=%s&g_tk=1164660242&g_tk_openkey=1164660242", tsMs, startCount, uid)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println("http.NewRequest error: ", err)
		return
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (iPhone; CPU iPhone OS 16_6 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/16.6 Mobile/15E148 Safari/604.1")
	req.Header.Set("Referer", "https://static-play.kg.qq.com/node/personal_v2/?uid="+uid)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("client.Do error: ", err)
		return
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("ioutil.ReadAll error: ", err)
		return
	}

	var RespPersonHomepage PersonHomepage
	err = json.Unmarshal(body, &RespPersonHomepage)
	if err != nil {
		fmt.Println("json.Unmarshal error: ", err)
		return
	}

	nickname = RespPersonHomepage.Data.Nickname
	ugcTotalCount = RespPersonHomepage.Data.UgcTotalCount

	return RespPersonHomepage.Data.Ugclist, nickname, ugcTotalCount

}

func parseUidFromUrl(url string) (uid string) {
	if url == "" {
		fmt.Println("url is empty.")
		return ""
	}
	re := regexp.MustCompile(`uid=([^&]+)`)
	match := re.FindStringSubmatch(url)

	if len(match) < 1 {
		fmt.Println("uid not found in url.")
		return ""
	}

	uid = match[1]
	fmt.Println("uid: ", uid)
	return uid
}

func main() {
	args := os.Args

	// var uid string = "639a9f8c2c2b3388"
	if len(args) < 2 {
		fmt.Println("Please provide a valid url.\ne.g. https://node.kg.qq.com/personal?uid=1223342")
		return
	}

	var uid string = parseUidFromUrl(args[1])
	if uid == "" {
		return
	}

	if _, err := os.Stat(DOWNLOAD_DIR); os.IsNotExist(err) {
		os.Mkdir(DOWNLOAD_DIR, 0755)
	}

	_, nickname, ugcTotalCount := getMusicList(uid, 1)
	fmt.Println("nickname: ", nickname, "ugcList: ", ugcTotalCount)
	var ugcChan chan Ugc = make(chan Ugc, ugcTotalCount)
	go startAddUgc(uid, ugcChan, ugcTotalCount)

	wg := sync.WaitGroup{}
	wg.Add(ugcTotalCount)
	for ugc := range ugcChan {
		go downloadMusic(ugc.Shareid, ugc.Title, nickname, &wg)
	}
	wg.Wait()
}
