package utils

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// 签名算法，目前为固定值
const algorithm = "TC3-HMAC-SHA256"

func sha256hex(v string) string {
	b := sha256.Sum256([]byte(v))
	return hex.EncodeToString(b[:])
}

func tencentDNSHmasha256(v, key string) string {
	hashed := hmac.New(sha256.New, []byte(key))
	hashed.Write([]byte(v))
	return string(hashed.Sum(nil))
}

// TencentDNSSigner 腾讯云 DNSPod 签名方法
func TencentDNSSigner(r *http.Request, secretId, secretKey, action, payload string) {
	timestamp := time.Now().Unix()
	timestampStr := strconv.FormatInt(timestamp, 10)

	// 请求字符串
	canonicalHeaders := "content-type:application/json\nhost:dnspod.tencentcloudapi.com\nx-tc-action:" + strings.ToLower(action) + "\n"
	signedHeaders := "content-type;host;x-tc-action"
	hashedRequestPayload := sha256hex(payload)
	canonicalRequest := "POST\n/\n\n" + canonicalHeaders + "\n" + signedHeaders + "\n" + hashedRequestPayload

	// 待签名字符串
	date := time.Unix(timestamp, 0).UTC().Format("2006-01-02")
	credentialScope := date + "/dnspod/tc3_request"
	hashedCanonicalRequest := sha256hex(canonicalRequest)
	string2sign := algorithm + "\n" + timestampStr + "\n" + credentialScope + "\n" + hashedCanonicalRequest

	// 签名
	secretDate := tencentDNSHmasha256(date, "TC3"+secretKey)
	secretService := tencentDNSHmasha256("dnspod", secretDate)
	secretSigning := tencentDNSHmasha256("tc3_request", secretService)
	signature := hex.EncodeToString([]byte(tencentDNSHmasha256(string2sign, secretSigning)))

	// 拼接 Authorization
	authorization := algorithm + " Credential=" + secretId + "/" + credentialScope + ", SignedHeaders=" + signedHeaders + ", Signature=" + signature

	r.Header.Set("Authorization", authorization)
	r.Header.Set("Host", "dnspod.tencentcloudapi.com")
	r.Header.Set("X-TC-Action", action)
	r.Header.Add("X-TC-Timestamp", timestampStr)
}
