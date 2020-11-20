package sysInit

import (
	"bufio"
	"go-admin/core/utils/log"
	"os"
	"strings"
)

var (
	config map[string]string
)
func init(){
	conf := make(map[string]string)

	f, err := os.Open("./core/conf/conf.conf")   //因为bufio需要的是一个*os.File类型，所以我们换个方式读取，稍后再介绍一下
	if err != nil {
		log.Error.Println(err)
	}

	defer func() {
		if err = f.Close(); err != nil {
			log.Error.Println(err)
		}
	}()

	s := bufio.NewScanner(f)
	for s.Scan() {
		keyval := s.Text()

		if keyval != "" {
			arr := strings.Split(keyval, "=")
			conf[strings.TrimSpace(arr[0])] = strings.TrimSpace(arr[1])
		}
	}

	err = s.Err()
	if err != nil {
		log.Error.Println(err)
	}

	config = conf

	log.Debug.Println("init")
}

func GetConfigValue(key string) string{
	return config[key]
}
