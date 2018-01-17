
这是一个golang爬虫demo 爬去一个美女图片网站的首页所有图片
采用golang 多线程的方式爬取图片 将爬到的图片保存到本地
代码中有用到goquery 网页数据解析框架 chan 控制goroutine 进行下载

http://www.umei.cc/
一个妹子图片网站  请求的 header 必须带着 Referer 否则404 （比较简单的一种反爬虫策略）
用wireshark 抓取浏览器请求图片的数据就可以得到 Referer

//代码不复杂，适合新手学习

goquery 传送门 https://godoc.org/github.com/PuerkitoBio/goquery



