// +build debug

package main

import "log"

func dLogd(v ...interface{}) {
	for _, v := range v {
		log.Printf("%+#v", v)
	}
}

func dLogf(f, v ...interface{}) {
	log.Printf(f, v...)
}

func dLog(v ...interface{}) {
	log.Print(v...)
}

func dLogln(v ...interface{}) {
	log.Println(v...)
}
