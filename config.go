package main

import (
	"strconv"
	"fmt"
	"os"
	"errors"
	"encoding/json"
	"io/ioutil"
	)

var conf_fn = "conf.json"

const app_name	string="App Name: dogu"

const app_descr	string="Descr: 道具はHTTPでものを転送する"

type Arg_func func(*configuration, []string) (error)

// Command line switch elements
type cmd_line_items struct {
	id		int
	Switch		string
	Help_srt	string
	Has_arg		bool
	Needed		bool
	Func		Arg_func
}

// Represents the application configuration
type configuration struct {
	Host		string
	Target		string
	Request_dom	string
	Daemon		bool
	Port		int
	cmdlineNeeds    map[string] bool
}

// Instance of default configuration values
var	Default_config  configuration = configuration{
	Host:		"https://localtunnel.me",
	Target:		"localhost",
	Request_dom:	"",
	Daemon:		false,
	Port:		8080,
	cmdlineNeeds:   map[string] bool{},
	}

// Inserts a commandline item item, which is composed by:
// * switch string
// * switch descriptio
// * if the switch requires an additiona argument
// * a pointer to the function that manages the switch
// * the configuration that gets updated
func push_cmd_line_item(Switch string, Help_str string, Has_arg bool, Needed bool, Func Arg_func, cmd_line *[]cmd_line_items){
	*cmd_line = append(*cmd_line, cmd_line_items{id: len(*cmd_line)+1, Switch: Switch, Help_srt: Help_str, Has_arg: Has_arg, Needed: Needed, Func: Func})
}

// This function initializes configuration parser subsystem
// Inserts all the commandline switches suppported by the application
func cmd_line_item_init() ([]cmd_line_items){
	var res	[]cmd_line_items

	push_cmd_line_item("-h", "Specifies localtunnel host",			true,  false,	func_host,	&res)
	push_cmd_line_item("-r", "Specifies request domain",			true,  true,	func_dom,	&res)
	push_cmd_line_item("-p", "Specifies http service local port",		true,  false,	func_port,	&res)
	push_cmd_line_item("-j", "Specifies config file",			true,  false,	func_jconf,	&res)
	push_cmd_line_item("-d", "Demonize (linux)",				false,  false,	func_daemon,	&res)

	return res
}

func func_daemon	(conf *configuration,fn []string)		(error){
	(*conf).Daemon=true
	return nil
}

func func_help		(conf *configuration,fn []string)		(error){
	return errors.New("Command Help")
}

func func_host(conf *configuration, host []string)			(error){
	(*conf).Host=host[0]
	return nil
}
func func_dom(conf *configuration, dom []string)			(error){
	(*conf).Request_dom=dom[0]
	return nil
}
func func_jconf		(conf *configuration,fn []string)		(error){
	jsonFile, err := os.Open(fn[0])
	if err != nil {
		return err
		}
	byteValue, _ := ioutil.ReadAll(jsonFile)
	jsonFile.Close()
	err=json.Unmarshal(byteValue, conf)
	if err != nil {
		return err
		}
	return nil
}
func func_port	(conf *configuration, port []string)	(error){
	s, err := strconv.Atoi(port[0])
	if err!=nil {
		return err
		}
	(*conf).Port=s
	return nil
}

// Uses commandline args to generate the help string
func print_help(lines []cmd_line_items){

	fmt.Println(app_name)
	fmt.Println(app_descr)
	for _,item := range lines{
		fmt.Printf(
			"\t%s\t%s\t%s\n",
			item.Switch,
			func (a bool)(string){
				if a {
					return "<v>"
					}
				return ""
			}(item.Has_arg),
			item.Help_srt,
			)
		}
}

// Used to parse the command line and generate the command line
func args_parse(lines []cmd_line_items)(configuration, error){
	var	extra		bool=false;
	var	conf		configuration=Default_config
	var 	f		Arg_func

	for _, item := range lines{
		if item.Needed {
			conf.cmdlineNeeds[item.Switch]=false
			}
		}

	for _, os_arg := range os.Args[1:] {
		if !extra {
			for _, arg := range lines{
				if arg.Switch==os_arg {
					if arg.Needed {
						conf.cmdlineNeeds[arg.Switch]=true
						}
					if arg.Has_arg{
						f=arg.Func
						extra=true
						break
						}
					err := arg.Func(&conf, []string{})
					if err != nil {
						return Default_config, err
						}
					}
				}
			continue
			}
		if extra{
			err := f(&conf,[]string{os_arg})
			if err != nil {
				return Default_config, err
				}
			extra=false
			}

		}
	if extra {
		 return  Default_config, errors.New("Missing switch arg")
		}

	res:=true
	for _, element := range conf.cmdlineNeeds {
		res = res && element
		}
	if res {
		return	conf, nil
		}
	return Default_config, errors.New("Missing needed arg")
}
