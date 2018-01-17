package main

import (
	"fmt"
	"net/http"
	"io"
	"github.com/PuerkitoBio/goquery"
	"log"
	"strings"
	"os"
	//"time"
	"time"
	"runtime"
)

//一个妹子图片网站 请求的 header 必须带着 Referer 否则404 （比较简单的一种反爬虫策略）
var url = "http://www.umei.cc/"

var c chan int

func main() {
	runtime.GOMAXPROCS(4)
	spider()
	//testDownLoad()
}

//url->Document->所有图片url->开启多线程进行下载->保存到本地
func spider() {
	doc, err := goquery.NewDocument(url)
	if err != nil {
		log.Fatal(err)
	}
	urls := ImageRule(doc, match);
	fmt.Println("共解析到", len(urls), "图片地址")
	c = make(chan int)
	for _, s := range urls {
		fmt.Println(s)
		go downloadImage(s)
	}
	//可以等待一会儿，留时间给子goroutine 执行
	//但是这种方式不怎么靠谱 //直接采用chan 的方式
	//time.Sleep(1e9*10)
	for i := 0; i < len(urls); i++ {
		<-c
	}
}

// 单独测试了以下 下载方法
func testDownLoad() {
	var url_img = "http://i1.umei.cc/uploads/tu/201608/164/hanguomeinv.jpg";
	//var  url_img  = "http://t1.mmonly.cc/uploads/tu/sm/201601/19/005.jpg";
	downloadImage(url_img)
}

func match(image string) {
	fmt.Println(image);
}

func getData(url string) (eader io.Reader, err error) {
	req := buildRequest(url)
	resp, err := http.DefaultClient.Do(req)
	defer resp.Body.Close()
	return io.Reader(resp.Body), err
}

// 得到一个网页中所有 ImageUrl
func parseImageUrl(reader io.Reader) (res []string, err error) {

	doc, err := goquery.NewDocumentFromReader(reader)
	if err != nil {
		return nil, err
	}
	fmt.Println(doc.Url)
	ImageRule(doc, func(image string) {
		res = append(res, image)
	})
	return res, nil
}

func ImageRule(doc *goquery.Document, f func(image string)) (urls []string) {
	str := make([]string, 0)
	//直接找以img 开头的标签 过滤掉不符合规则的url 即可
	doc.Find("img").Each(func(i int, s *goquery.Selection) {
		url, result := s.Attr("src")
		if result {
			if strings.HasSuffix(url, ".jpg") {
				str = append(str, url)
			}
		}
	})
	return str
}

//根据url 创建http 请求的 request
//网站有反爬虫策略 wireshark 不解释
func buildRequest(url string) *http.Request {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		panic(err)
	}
	//	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/60.0.3112.78 Safari/537.36")
	//	req.Header.Set("Cookie", "Hm_lvt_c605a31292b623d214d012ec2a737685=1516111586; Hm_lpvt_c605a31292b623d214d012ec2a737685=1516111613")
	//req.Header.Set("If-None-Match", "5a309bab-26057")
	req.Header.Set("Referer", "http://www.umei.cc/")
	//req.Header.Set("If-Modified-Since", "Wed, 13 Dec 2017 03:16:59 GMT")
	return req
}
// 下载图片
func downloadImage(url string) {
	fileName := getNameFromUrl(url)
	req := buildRequest(url)
	http.DefaultClient.Timeout = 10 * time.Second;
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println("failed download ")
		panic(err)
	}
	if resp.StatusCode != http.StatusOK {
		fmt.Println("failed download " + url)
		return
	}
	defer func() {
		resp.Body.Close()
		if r := recover(); r != nil {
			fmt.Println(r)
		}
		c <- 0
	}()

	fmt.Println("begin download " + fileName)
	os.MkdirAll("./images/", 0777)
	localFile, _ := os.OpenFile("./images/"+fileName, os.O_CREATE|os.O_RDWR, 0777)

	if _, err := io.Copy(localFile, resp.Body); err != nil {
		panic("failed save " + fileName)
	}

	fmt.Println("success download " + fileName)
}
// 判读文件夹是否存在
func isExist(dir string) bool {
	_, err := os.Stat(dir)
	if err == nil {
		return true
	}
	return os.IsExist(err)
}

// 通过url 得到图片名字
func getNameFromUrl(url string) string {
	arr := strings.Split(url, "/")
	return arr[len(arr)-1]
}
