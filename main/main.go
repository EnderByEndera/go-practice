package main

import "log"

func main() {
	initRoute()
	err := router.RunTLS(":8000", "1_www.italktoyou.cn_bundle.crt", "2_www.italktoyou.cn.key")
	//err := router.Run(":8000")
	if err != nil {
		log.Fatal(err)
	}
}
