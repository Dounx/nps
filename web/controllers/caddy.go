package controllers

import (
	"os"
	"path/filepath"

	"ehang.io/nps/lib/caddy"
	"ehang.io/nps/lib/common"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
)

var caddyClient *caddy.Client

type CaddyController struct {
	BaseController
}

func init() {
	err := beego.LoadAppConfig("ini", filepath.Join(common.GetRunPath(), "conf", "nps.conf"))
	if err != nil {
		logs.Error("load config file error " + err.Error())
		os.Exit(0)
	}

	caddyClient = &caddy.Client{
		Host: "localhost",
		Port: beego.AppConfig.String("caddy_admin_port"),
	}
}

func (s *CaddyController) List() {
	if !s.Data["isAdmin"].(bool) {
		s.Abort("401")
	}

	if s.Ctx.Request.Method == "GET" {
		s.Data["menu"] = "caddy"
		s.SetInfo("caddy list")
		s.display("caddy/list")
		return
	}

	list, err := caddyClient.GetReverseProxyList()
	if err != nil {
		logs.Error(err)
		s.Abort("500")
	}

	cnt := len(list)
	s.AjaxTable(list, cnt, cnt, nil)
}

func (s *CaddyController) Add() {
	if !s.Data["isAdmin"].(bool) {
		s.Abort("401")
	}

	if s.Ctx.Request.Method == "GET" {
		s.Data["menu"] = "caddy"
		s.SetInfo("add caddy config")
		s.display()
		return
	}

	lastID, err := caddyClient.GetReverseProxyLastID()
	if err != nil {
		logs.Error(err)
		s.Abort("500")
	}

	c := &caddy.ReverseConfig{
		ID:           lastID + 1,
		UpstreamPath: s.getEscapeString("upstream_path"),
		MatchHost:    s.getEscapeString("match_host"),
		MatchPath:    s.getEscapeString("match_path"),
	}

	_ = caddyClient.AddReverseProxy(c)

	s.AjaxOk("add success")
}

func (s *CaddyController) Edit() {
	id := int64(s.GetIntNoErr("id"))
	if s.Ctx.Request.Method == "GET" {
		c, err := caddyClient.GetReverseProxy(id)
		if err != nil {
			s.error()
		} else {
			s.Data["c"] = c
		}

		s.Data["menu"] = "caddy"
		s.SetInfo("edit client")
		s.display()
		return
	}

	c, err := caddyClient.GetReverseProxy(id)
	if err != nil {
		s.error()
		s.AjaxErr("caddy config id not found")
		return
	} else {
		c.UpstreamPath = s.getEscapeString("upstream_path")
		c.MatchHost = s.getEscapeString("match_host")
		c.MatchPath = s.getEscapeString("match_path")

		_ = caddyClient.UpdateReverseProxy(c)
	}
	s.AjaxOk("save success")
}

func (s *CaddyController) Del() {
	id := int64(s.GetIntNoErr("id"))
	_ = caddyClient.DeleteReverseProxy(id)

	s.AjaxOk("delete success")
}
