//web weixin client
package main

import (
	"fmt"
	"net/http"
	"net/http/cookiejar"

	// "log"
	"bytes"
	"encoding/json"
	"encoding/xml"
	"io"
	"io/ioutil"

	// "math"
	"math/rand"
	"net/url"
	"os"
	"os/exec"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"time"
	// "strings"
)

var (
	port   = "19090"
	local  = "192.168.99.10"
	remote = "http://192.168.99.10:9090"
)

func debugPrint(content interface{}) {
	if *debug == "on" {
		fmt.Println(content)
	}
}

func init() {
	cfg, err := getConfig("esap")
	if err != nil {
		fmt.Println("未找到配置文件")
	} else {
		port = cfg["port"]
		local = cfg["local"]
		remote = cfg["remote"]
	}
}

type wxweb struct {
	uuid         string
	base_uri     string
	redirect_uri string
	uin          string
	sid          string
	skey         string
	pass_ticket  string
	deviceId     string
	SyncKey      map[string]interface{}
	synckey      string
	User         map[string]interface{}
	BaseRequest  map[string]interface{}
	syncHost     string
	http_client  *http.Client
}

func (self *wxweb) getUuid(args ...interface{}) bool {
	urlstr := "https://login.weixin.qq.com/jslogin"
	urlstr += "?appid=wx782c26e4c19acffb&fun=new&lang=zh_CN&_=" + self._unixStr()
	data, _ := self._get(urlstr, false)
	re := regexp.MustCompile(`"([\S]+)"`)
	find := re.FindStringSubmatch(data)
	if len(find) > 1 {
		self.uuid = find[1]
		return true
	} else {
		return false
	}
}

func (self *wxweb) _run(desc string, f func(...interface{}) bool, args ...interface{}) {
	start := time.Now().UnixNano()
	fmt.Print(desc)
	var result bool
	if len(args) > 1 {
		result = f(args)
	} else if len(args) == 1 {
		result = f(args[0])
	} else {
		result = f()
	}
	useTime := fmt.Sprintf("%.5f", (float64(time.Now().UnixNano()-start) / 1000000000))
	if result {
		fmt.Println("成功,用时" + useTime + "秒")
	} else {
		fmt.Println("失败\n[*] 退出程序")
		os.Exit(1)
	}
}

func (self *wxweb) _post(urlstr string, params map[string]interface{}, jsonFmt bool) ([]byte, error) {
	var err error
	var resp *http.Response
	if jsonFmt == true {
		jsonPost, err := json.Marshal(params)
		if err != nil {
			return []byte(""), Error("json encode fail")
		}

		debugPrint(jsonPost)
		requestBody := bytes.NewBuffer([]byte(jsonPost))
		request, err := http.NewRequest("POST", urlstr, requestBody)
		if err != nil {
			return []byte(""), err
		}
		request.Header.Set("Content-Type", "application/json;charset=utf-8")
		request.Header.Add("Referer", "https://wx.qq.com/")
		request.Header.Add("User-agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_11_3) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/47.0.2526.111 Safari/537.36")
		resp, err = self.http_client.Do(request)
		// resp, err = self.http_client.Post(urlstr, "application/json;charset=utf-8", requestBody)
	} else {
		v := url.Values{}
		for key, value := range params {
			v.Add(key, value.(string))
		}
		resp, err = self.http_client.PostForm(urlstr, v)
	}

	if err != nil || resp == nil {
		fmt.Println(err)
		return []byte(""), err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return []byte(""), err
	} else {
		defer resp.Body.Close()
	}
	return body, nil
}

func (self *wxweb) _get(urlstr string, jsonFmt bool) (string, error) {
	var err error
	res := ""

	request, _ := http.NewRequest("GET", urlstr, nil)
	request.Header.Add("Referer", "https://wx.qq.com/")
	request.Header.Add("User-agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_11_3) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/47.0.2526.111 Safari/537.36")
	resp, err := self.http_client.Do(request)
	if err != nil {
		return res, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return res, err
	}
	return string(body), nil
}

func (self *wxweb) _getpic(urlstr string, jsonFmt bool) (string, error) {
	var err error
	res := ""

	request, _ := http.NewRequest("GET", urlstr, nil)
	request.Header.Add("Referer", "https://wx.qq.com/")
	request.Header.Add("User-agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_11_3) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/47.0.2526.111 Safari/537.36")
	resp, err := self.http_client.Do(request)
	if err != nil {
		return res, err
	}
	defer resp.Body.Close()
	filename := GenerateId()
	filejpg, err := os.Create(filename + ".jpg")
	if err != nil {
		fmt.Println("create jpg err:", err)
	}
	defer time.AfterFunc(30*time.Second, func() { os.Remove(filename + ".jpg") })
	defer filejpg.Close()
	_, err = io.Copy(filejpg, resp.Body)
	if err != nil {
		return "", err
	}
	return filename, nil
}

func (self *wxweb) _unixStr() string {
	return strconv.Itoa(int(time.Now().Unix()))
}

func (self *wxweb) genQRcode(args ...interface{}) bool {
	urlstr := "https://login.weixin.qq.com/qrcode/" + self.uuid
	urlstr += "?t=webwx"
	urlstr += "&_=" + self._unixStr()
	path := "qrcode.jpg"
	out, err := os.Create(path)
	resp, err := self._get(urlstr, false)
	_, err = io.Copy(out, bytes.NewReader([]byte(resp)))
	if err != nil {
		return false
	} else {
		if runtime.GOOS == "darwin" {
			exec.Command("open", path).Run()
		} else {
			go func() {
				http.HandleFunc("/qrcode", func(w http.ResponseWriter, req *http.Request) {
					http.ServeFile(w, req, "qrcode.jpg")
					return
				})
				http.HandleFunc("/p", func(w http.ResponseWriter, req *http.Request) {
					pi := req.FormValue("id")
					picfile := pi + ".jpg"
					http.ServeFile(w, req, picfile)
					return
				})
				http.ListenAndServe(fmt.Sprint(":", port), nil)
			}()
		}
		return true
	}
}

func (self *wxweb) waitForLogin(tip int) bool {
	time.Sleep(time.Duration(tip) * time.Second)
	url := "https://login.weixin.qq.com/cgi-bin/mmwebwx-bin/login"
	url += "?tip=" + strconv.Itoa(tip) + "&uuid=" + self.uuid + "&_=" + self._unixStr()
	data, _ := self._get(url, false)
	re := regexp.MustCompile(`window.code=(\d+);`)
	find := re.FindStringSubmatch(data)
	if len(find) > 1 {
		code := find[1]
		if code == "201" {
			return true
		} else if code == "200" {
			re := regexp.MustCompile(`window.redirect_uri="(\S+?)";`)
			find := re.FindStringSubmatch(data)
			if len(find) > 1 {
				r_uri := find[1] + "&fun=new"
				self.redirect_uri = r_uri
				re = regexp.MustCompile(`/`)
				finded := re.FindAllStringIndex(r_uri, -1)
				self.base_uri = r_uri[:finded[len(finded)-1][0]]
				return true
			}
			return false
		} else if code == "408" {
			fmt.Println("[登陆超时]")
		} else {
			fmt.Println("[登陆异常]")
		}
	}
	return false
}

func (self *wxweb) login(args ...interface{}) bool {
	data, _ := self._get(self.redirect_uri, false)
	type Result struct {
		Skey        string `xml:"skey"`
		Wxsid       string `xml:"wxsid"`
		Wxuin       string `xml:"wxuin"`
		Pass_ticket string `xml:"pass_ticket"`
	}
	v := Result{}
	err := xml.Unmarshal([]byte(data), &v)
	if err != nil {
		fmt.Printf("error: %v", err)
		return false
	}
	self.skey = v.Skey
	self.sid = v.Wxsid
	self.uin = v.Wxuin
	self.pass_ticket = v.Pass_ticket
	self.BaseRequest = make(map[string]interface{})
	self.BaseRequest["Uin"], _ = strconv.Atoi(v.Wxuin)
	self.BaseRequest["Sid"] = v.Wxsid
	self.BaseRequest["Skey"] = v.Skey
	self.BaseRequest["DeviceID"] = self.deviceId
	return true
}

func (self *wxweb) webwxinit(args ...interface{}) bool {
	url := fmt.Sprintf("%s/webwxinit?pass_ticket=%s&skey=%s&r=%s", self.base_uri, self.pass_ticket, self.skey, self._unixStr())
	params := make(map[string]interface{})
	params["BaseRequest"] = self.BaseRequest
	res, err := self._post(url, params, true)
	if err != nil {
		return false
	}
	ioutil.WriteFile("tmp.txt", []byte(res), 777)
	data := make(map[string]interface{})
	err = json.Unmarshal(res, &data)
	if err != nil {
		return false
	}
	self.User = data["User"].(map[string]interface{})
	self.SyncKey = data["SyncKey"].(map[string]interface{})
	self._setsynckey()

	retCode := data["BaseResponse"].(map[string]interface{})["Ret"].(float64)
	return retCode == 0
}

func (self *wxweb) _setsynckey() {
	keys := []string{}
	for _, keyVal := range self.SyncKey["List"].([]interface{}) {
		key := strconv.Itoa(int(keyVal.(map[string]interface{})["Key"].(float64)))
		value := strconv.Itoa(int(keyVal.(map[string]interface{})["Val"].(float64)))
		keys = append(keys, key+"_"+value)
	}
	self.synckey = strings.Join(keys, "|")
	debugPrint(self.synckey)
}

func (self *wxweb) synccheck() (string, string) {
	urlstr := fmt.Sprintf("https://%s/cgi-bin/mmwebwx-bin/synccheck", self.syncHost)
	v := url.Values{}
	v.Add("r", self._unixStr())
	v.Add("sid", self.sid)
	v.Add("uin", self.uin)
	v.Add("skey", self.skey)
	v.Add("deviceid", self.deviceId)
	v.Add("synckey", self.synckey)
	v.Add("_", self._unixStr())
	urlstr = urlstr + "?" + v.Encode()
	data, _ := self._get(urlstr, false)
	re := regexp.MustCompile(`window.synccheck={retcode:"(\d+)",selector:"(\d+)"}`)
	find := re.FindStringSubmatch(data)
	if len(find) > 2 {
		retcode := find[1]
		selector := find[2]
		debugPrint(fmt.Sprintf("retcode:%s,selector,selector%s", find[1], find[2]))
		return retcode, selector
	} else {
		return "9999", "0"
	}
}

func (self *wxweb) testsynccheck(args ...interface{}) bool {
	SyncHost := []string{
		"webpush.wx.qq.com",
		"webpush2.wx.qq.com",
		"webpush.wechat.com",
		"webpush1.wechat.com",
		"webpush2.wechat.com",
		"webpush1.wechatapp.com",
		//"webpush.wechatapp.com"
	}
	for _, host := range SyncHost {
		self.syncHost = host
		retcode, _ := self.synccheck()
		if retcode == "0" {
			return true
		}
	}
	return false
}

func (self *wxweb) webwxstatusnotify(args ...interface{}) bool {
	urlstr := fmt.Sprintf("%s/webwxstatusnotify?lang=zh_CN&pass_ticket=%s", self.base_uri, self.pass_ticket)
	params := make(map[string]interface{})
	params["BaseRequest"] = self.BaseRequest
	params["Code"] = 3
	params["FromUserName"] = self.User["UserName"]
	params["ToUserName"] = self.User["UserName"]
	params["ClientMsgId"] = int(time.Now().Unix())
	res, err := self._post(urlstr, params, true)
	if err != nil {
		return false
	}
	data := make(map[string]interface{})
	err = json.Unmarshal(res, &data)
	if err != nil {
		return false
	}
	retCode := data["BaseResponse"].(map[string]interface{})["Ret"].(float64)
	return retCode == 0
}

func (self *wxweb) webwxsync() interface{} {
	urlstr := fmt.Sprintf("%s/webwxsync?sid=%s&skey=%s&pass_ticket=%s", self.base_uri, self.sid, self.skey, self.pass_ticket)
	params := make(map[string]interface{})
	params["BaseRequest"] = self.BaseRequest
	params["SyncKey"] = self.SyncKey
	params["rr"] = ^int(time.Now().Unix())
	res, err := self._post(urlstr, params, true)
	if err != nil {
		return false
	}
	data := make(map[string]interface{})
	err = json.Unmarshal(res, &data)
	if err != nil {
		return false
	}
	retCode := data["BaseResponse"].(map[string]interface{})["Ret"].(float64)
	if retCode == 0 {
		self.SyncKey = data["SyncKey"].(map[string]interface{})
		self._setsynckey()
	}
	return data
}

func (self *wxweb) handleMsg(r interface{}) {
	// fmt.Printf("[*] raw: %v \n", r)
	myNickName := self.User["NickName"].(string)
	for _, msg := range r.(map[string]interface{})["AddMsgList"].([]interface{}) {
		// fmt.Printf("[*] message: %v \n", msg)
		// msg = msg.(map[string]interface{})
		msgType := msg.(map[string]interface{})["MsgType"].(float64)
		// fmt.Println("msgtype=", msgType)
		fromUserName := msg.(map[string]interface{})["FromUserName"].(string)
		// name = self.getUserRemarkName(msg['FromUserName'])
		content := msg.(map[string]interface{})["Content"].(string)
		content = strings.Replace(content, "&lt;", "<", -1)
		content = strings.Replace(content, "&gt;", ">", -1)
		content = strings.Replace(content, " ", " ", 1)
		msgid := msg.(map[string]interface{})["MsgId"]
		if msgType == 1 || msgType == 3 {
			var ans string
			var err error

			if msgType == 3 { // 处理图片
				u := fmt.Sprintf("https://wx.qq.com/cgi-bin/mmwebwx-bin/webwxgetmsgimg?&MsgID=%v&skey=%s", msgid, strings.Replace(self.skey, "@", "%40", -1))
				picurl, err := self._getpic(u, false)
				if err != nil {
					fmt.Println("getpic err:", err)
				}
				ans, err = self.getReplyByApi("图片", fromUserName, "http://"+local+":"+port+"/p?id="+picurl)
			} else if fromUserName[:2] == "@@" {
				debugPrint(content + "|0045")
				contentSlice := strings.Split(content, ":<br/>")
				// people := contentSlice[0]
				content = contentSlice[1]
				if strings.Contains(content, "@"+myNickName) {
					realcontent := strings.TrimSpace(strings.Replace(content, "@"+myNickName, "", 1))
					debugPrint(realcontent + "|0046")
					ans, err = self.getReplyByApi(realcontent, fromUserName)
				}
			} else {
				ans, err = self.getReplyByApi(content, fromUserName)
			}
			debugPrint(ans)
			debugPrint(content)
			if err != nil {
				debugPrint(err)
			} else if ans != "" {
				go self.webwxsendmsg(ans, fromUserName)
			}
		}
	}
}

func (self *wxweb) getReplyByApi(realcontent string, fromUserName string, pic ...string) (string, error) {
	return getAnswer(realcontent, fromUserName, self.User["NickName"].(string), pic...)
}

func (self *wxweb) webwxsendmsg(message string, toUseNname string) bool {
	urlstr := fmt.Sprintf("%s/webwxsendmsg?pass_ticket=%s", self.base_uri, self.pass_ticket)
	clientMsgId := self._unixStr() + "0" + strconv.Itoa(rand.Int())[3:6]
	params := make(map[string]interface{})
	params["BaseRequest"] = self.BaseRequest
	msg := make(map[string]interface{})
	msg["Type"] = 1
	msg["Content"] = message
	msg["FromUserName"] = self.User["UserName"]
	msg["ToUserName"] = toUseNname
	msg["LocalID"] = clientMsgId
	msg["ClientMsgId"] = clientMsgId
	params["Msg"] = msg
	data, err := self._post(urlstr, params, true)
	if err != nil {
		debugPrint(err)
		return false
	} else {
		debugPrint(data)
		return true
	}
}

func (self *wxweb) _init() {
	gCookieJar, _ := cookiejar.New(nil)
	httpclient := http.Client{
		CheckRedirect: nil,
		Jar:           gCookieJar,
	}
	self.http_client = &httpclient
	rand.Seed(time.Now().Unix())
	str := strconv.Itoa(rand.Int())
	self.deviceId = "e" + str[2:17]
}

func (self *wxweb) start() {
	fmt.Println("[*] 微信网页版 ... 开动")
	self._init()
	self._run("[*] 正在获取 uuid ... ", self.getUuid)
	self._run("[*] 正在获取 二维码 ... ", self.genQRcode)
	if runtime.GOOS == "darwin" {
		fmt.Println("[*] 请使用微信扫描二维码以登录 ... ")
	} else {
		fmt.Println("[*] 打开链接扫码登录 http://" + local + ":" + port + "/qrcode")
	}
	for {
		if self.waitForLogin(1) == false {
			continue
		}
		fmt.Println("[*] 请在手机上点击确认以登录 ... ")
		if self.waitForLogin(0) == false {
			continue
		}
		break
	}
	self._run("[*] 正在登录 ... ", self.login)
	self._run("[*] 微信初始化 ... ", self.webwxinit)
	self._run("[*] 开启状态通知 ... ", self.webwxstatusnotify)
	self._run("[*] 进行同步线路测试 ... ", self.testsynccheck)
	for {
		retcode, selector := self.synccheck()
		if retcode == "1100" {
			fmt.Println("[*] 你在手机上登出了微信，再见")
			break
		} else if retcode == "1101" {
			fmt.Println("[*] 你在其他地方登录了 WEB 版微信，再见")
			break
		} else if retcode == "0" {
			if selector == "2" {
				r := self.webwxsync()
				debugPrint(r)
				switch r.(type) {
				case bool:
				default:
					self.handleMsg(r)
				}
			} else if selector == "0" {
				time.Sleep(1)
			} else if selector == "6" || selector == "4" {
				self.webwxsync()
				time.Sleep(1)
			}
		}
	}

}
