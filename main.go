package main

import (
	"encoding/json"
	"flag"
	"github.com/bitly/go-simplejson"
	"github.com/go-ini/ini"
	"github.com/golang/glog"
	"github.com/valyala/fasthttp"
	"io/ioutil"
	"os"
	"strings"
	"time"
)

const (
	CONF_FILE  = `main.conf`
	HttpContentType = "application/json; charset=utf-8"
)

type Conf struct {
	Listen          string `ini:"listen"`
	Api_path_dir    string `ini:"api_path_dir"`
}

type JsonFileData struct {
	Method  		string `json:"method"`
	Parameters  	map[string]string `json:"parameters"`
	SuccessResult   interface{} `json:"success_result"`
}

type JsonResult struct{
	Code  int    `json:"code"`
	Msg   string `json:"message"`
	Data  string `json:"data"`
}

func main()  {
	defer glog.Flush()

	conf := &Conf{}

	initConf(conf) //加载并映射conf文件内容

	conf.ServerStart() //监听http端口并获取
}

func (conf *Conf) ServerStart()  {
	go fasthttp.ListenAndServe(conf.Listen, conf.HandleFastHTTP)
	for {
		time.Sleep(5 * time.Second)
	}
}

func (conf *Conf) HandleFastHTTP(ctx *fasthttp.RequestCtx) {
	url_path := string(ctx.Path())
	if url_path == "/favicon.ico" { //在http包使用的时候，注册了/这个根路径的模式处理，浏览器会自动的请求favicon.ico
		return
	}

	JsonFileData := &JsonFileData{}

	simplejson_obj := conf.initJsonFile(url_path, JsonFileData, ctx) //加载并映射json文件内容
	if simplejson_obj == nil {
		sendBytes(ctx, []byte(`json simplejson error`))
		return
	}

	if JsonFileData.Method == `` {
		sendBytes(ctx, []byte(`json method not nil`))
		return
	}

	method := strings.ToLower(string(ctx.Method()))
	if strings.ToLower(JsonFileData.Method) != method {
		sendBytes(ctx, []byte(`request method by:` + JsonFileData.Method))
		return
	}

	if JsonFileData.SuccessResult == nil {
		sendBytes(ctx, []byte(`json success_result not nil`))
		return
	}

	data_map, err := simplejson_obj.Map();
	if err != nil {
		sendBytes(ctx, []byte(`simplejson to map data_map is nil`))
		return
	}

	success_map, err := simplejson_obj.Get(`success_result`).Map();
	if err != nil {
		sendBytes(ctx, []byte(`simplejson to map data_map is nil`))
		return
	}

	var msg_key string = `message`
	var res_key string = `data`

	for s_k,_ := range success_map  {
		if strings.Contains(`msg,mssage`, string(s_k)) {
			msg_key =  string(s_k)
		} else if strings.Contains(`data,result`, string(s_k)) {
			res_key =  string(s_k)
		}
	}

	if JsonFileData.Parameters != nil { //参数验证
		var val string
		for key, value := range JsonFileData.Parameters {
			if method == `get` {
				val = GetBodyBy(ctx, key)
			}
			if method == `post` {
				val = PostBodyBy(ctx, key)
			}

			if value == `required` && val == `` {
				sendBytes(ctx, []byte(`{"code":5001,"`+msg_key+`":"`+key+` is required","`+res_key+`":[]}`))
				return
			}
		}
	}

	mjson, err := json.Marshal(data_map[`success_result`])
	if err != nil {
		sendBytes(ctx, []byte(`simplejson success_result json.Marshal error`))
		return
	}

	sendBytes(ctx, []byte(string(mjson)))
	return
}

func GetBodyBy(ctx *fasthttp.RequestCtx, field string) string  {
	value := string(ctx.QueryArgs().Peek(field))
	return value
}

func PostBodyBy(ctx *fasthttp.RequestCtx, field string) string  {
	value := string(ctx.FormValue(field))
	return value
}

func (Conf *Conf) initJsonFile(url string, JsonFileData *JsonFileData, ctx *fasthttp.RequestCtx) *simplejson.Json  {
	var simplejson_obj *simplejson.Json
	url = strings.Trim(url, "/")

	json_file := Conf.Api_path_dir+ `/`+ url+ `.json`

	buf, err := ReadFile(json_file)
	if err != nil {
		glog.Warningf(`failed to read %s, %s`,json_file, err.Error())
	} else {
		if err := json.Unmarshal(buf, &JsonFileData); err != nil {
			glog.Warningf(`failed to unmarshal %s, %s`, json_file, err.Error())
		}
		if simplejson_obj, err = simplejson.NewJson(buf); err != nil {
			glog.Warningf(`failed to simplejson %s, %s`, json_file, err.Error())
		}
	}
	return simplejson_obj
}


func sendBytes(ctx *fasthttp.RequestCtx, buf []byte, code ...int) {
	status := fasthttp.StatusOK
	if len(code) > 0 {
		status = code[0]
	}

	ctx.SetBody(buf)
	ctx.SetStatusCode(status)
	ctx.Response.Header.Set("Content-Type", HttpContentType)
	return
}

func ReadFile(file string) (b []byte, err error) {
	b, err = ioutil.ReadFile(file)
	return
}

func initConf(conf interface{}) {
	file_name := flag.String("conf_path", CONF_FILE, "config file path")

	flag.Parse()

	err := ini.MapTo(conf, *file_name)
	if err != nil {
		glog.Errorf(`conf file error: %s`, err.Error())
		os.Exit(1)
	}
}