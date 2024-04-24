package main

import (
	"log"
	"nas/internal/dns"
	"nas/utils"
	"os"
	"time"
)

var lastIP string

var interName string

var secretId string
var secretKey string

var domain string
var subDomain string
var recordType string

func init() {
	secretId = os.Getenv("SECRET_ID")
	secretKey = os.Getenv("SECRET_KEY")
	domain = os.Getenv("DOMAIN")
	subDomain = os.Getenv("SUB_DOMAIN")
	recordType = os.Getenv("RECORD_TYPE")
	interName = os.Getenv("INTERFACE")

	if secretId == "" || secretKey == "" {
		panic("SECRET_ID and SECRET_KEY environment variable required")
	}

	if domain == "" {
		panic("DOMAIN environment variable required")
	}
}

func main() {
	client := dns.NewTencentDNSClient(secretId, secretKey)

	records, err := client.DescribeRecordList(domain, subDomain, "")
	if err != nil {
		log.Fatalf("未找到相关记录 DOMAIN=%s SUB_DOMAIN=%s RECORD_TYPE=%s 请先建立一条随机记录", domain, subDomain, recordType)
	}

	targetRecord := records[0]

	lastIP = targetRecord.Value

	log.Printf("主循环开始，最新记录为: %s\n", lastIP)

	for {
		_, targetIP, err := utils.GetIPByInterface(interName)
		if err != nil || targetIP == "" {
			log.Fatalf("获取 IP 地址错误: %v\n", err)
		}

		if targetIP != lastIP {
			log.Printf("更新 IP 地址: %s\n", targetIP)
			newId, err := client.ModifyDNS(domain, subDomain, recordType, targetRecord.RecordId, targetIP)
			if err != nil {
				log.Printf("更新 IP 解析记录错误: %v\n", err)
			} else {
				targetRecord.RecordId = newId
				lastIP = targetIP
			}
		}

		time.Sleep(time.Minute)
	}

}
