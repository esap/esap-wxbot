package main

import (
	"fmt"
	"io/ioutil"

	// "math/rand"
	"net/http"
	"time"
)

// const (
// 	EMOTICON_QUEST  = "颜文字"
// 	LOVEWORDS_QUEST = "爱你"
// )

// type ApiServer struct {
// 	ApiName string
// }

// func randomEmoticon() string {
// 	happy := []string{"_(┐「ε:)_", "_(:3 」∠)_", "(￣y▽￣)~*捂嘴偷笑",
// 		"・゜・(PД`q｡)・゜・", "(ง •̀_•́)ง", "(•̀ᴗ•́)و ̑̑", "ヽ(•̀ω•́ )ゝ", "(,,• ₃ •,,)",
// 		"(｡˘•ε•˘｡)", "(=ﾟωﾟ)ﾉ", "(○’ω’○)", "(´・ω・`)", "ヽ(･ω･｡)ﾉ", "(。-`ω´-)",
// 		"(´・ω・`)", "(´・ω・)ﾉ", "(ﾉ･ω･)", "  (♥ó㉨ò)ﾉ♡", "(ó㉨ò)", "・㉨・", "( ・◇・)？",
// 		"ヽ(*´Д｀*)ﾉ", "(´°̥̥̥̥̥̥̥̥ω°̥̥̥̥̥̥̥̥｀)", "(╭￣3￣)╭♡", "(☆ﾟ∀ﾟ)", "⁄(⁄ ⁄•⁄ω⁄•⁄ ⁄)⁄.",
// 		"(´-ι_-｀)", "ಠ౪ಠ", "ಥ_ಥ", "(/≥▽≤/)", "ヾ(o◕∀◕)ﾉ ヾ(o◕∀◕)ﾉ ヾ(o◕∀◕)ﾉ", "*★,°*:.☆\\(￣▽￣)/$:*.°★*",
// 		"ヾ (o ° ω ° O ) ノ゙", "╰(*°▽°*)╯ ", " (｡◕ˇ∀ˇ◕）", "o(*≧▽≦)ツ", "≖‿≖✧", ">ㅂ<", "ˋ▽ˊ", "\\(•ㅂ•)/♥",
// 		"✪ε✪", "✪υ✪", "✪ω✪", "눈_눈", ",,Ծ‸Ծ,,", "π__π", "（/TДT)/", "ʅ（´◔౪◔）ʃ", "(｡☉౪ ⊙｡)", "o(*≧▽≦)ツ┏━┓拍桌狂笑",
// 		" (●'◡'●)ﾉ♥", "<(▰˘◡˘▰)>", "｡◕‿◕｡", "(｡・`ω´･)", "(♥◠‿◠)ﾉ", "ʅ(‾◡◝) ", " (≖ ‿ ≖)✧", "（´∀｀*)",
// 		"（＾∀＾）", "(o^∇^o)ﾉ", "ヾ(=^▽^=)ノ", "(*￣∇￣*)", " (*´∇｀*)", "(*ﾟ▽ﾟ*)", "(｡･ω･)ﾉﾞ", "(≡ω≡．)",
// 		"(｀･ω･´)", "(´･ω･｀)", "(●´ω｀●)φ)"}
// 	rand.Seed(time.Now().Unix())
// 	return happy[rand.Intn(len(happy))]
// }

// func randomLoveWord() string {
// 	words := []string{"我唯一害怕的就是失去你，我感觉你快要消失了",
// 		"我养你吧。",
// 		"我能遇见你已经不可思议了",
// 		"我在未来等你",
// 		"我闭着眼睛也看不见自己，但是我却可以看见你",
// 		"You are like everything to me now",
// 		"这世界上，除了你我别无所求",
// 		"有时候我会想你想到无法承受",
// 		"多希望我自己知道怎么放弃你",
// 		"我已经爱上你一个星期了，记得吗",
// 		"送花啊，牵手啊，生日礼物这些我都不擅长，但要是说到结婚，我只希望可以娶你",
// 		"我真是大笨蛋，除了喜欢你什么都不知道",
// 		"该死，我记不起其他没有你的地方了",
// 		"眼睛，是你的眼睛，我知道为什么喜欢你了",
// 		"我喜欢你的一切，特别是你呆呆傻傻的样子",
// 		"为什么喜欢你，我好像也不知道，但我知道没有别人会比我更喜欢你"}
// 	rand.Seed(time.Now().Unix())
// 	return words[rand.Intn(len(words))]
// }

func getAnswer(msg string, uid string, robotName string, pic ...string) (string, error) {
	fmt.Println("[api] 尝试应答=>", msg)
	if len(pic) == 0 {
		pic = append(pic, "")
	}

	httpClient := http.Client{Timeout: 10 * time.Second}
	url := remote + robotName + "?userid=" + uid + "&msg=" + msg + "&pic=" + pic[0]
	resp, err := httpClient.PostForm(url, nil)
	if err != nil {
		fmt.Println("post-err:", err)
		return "", Error("request fail")
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", Error("request fail")
	}

	// if string(body) == "" {
	// 	if strings.Contains(msg, LOVEWORDS_QUEST) {
	// 		e := randomLoveWord()
	// 		return e, nil
	// 	} else {
	// 		e := randomEmoticon()
	// 		return e, nil
	// 	}
	// }

	return string(body), nil
}
