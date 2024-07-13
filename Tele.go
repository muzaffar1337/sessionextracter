package main

import (
	"net/url"
	"github.com/valyala/fasthttp"
)

type Bot struct {
	Token string
	id    string
	Msgid [1]string
}

func (b *Bot) EditMessage(s, msgid string) []byte {
	_,res,_ := fasthttp.Post(nil,"https://api.telegram.org/bot"+b.Token+"/editmessagetext?chat_id="+b.id+"&message_id="+msgid+"&text="+url.QueryEscape(string(s)),nil)
	return res
}

func (b *Bot) SendMessage(s string) []byte {
	_,res,_ := fasthttp.Post(nil,"https://api.telegram.org/bot"+b.Token+"/sendMessage?chat_id="+b.id+"&text="+url.QueryEscape(s)+"&parse_mode=HTML",nil)
	return res
}