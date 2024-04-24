package dns

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"nas/models"
	"nas/utils"
	"net/http"
)

const tencentApiVersion = "2021-03-23"
const tencentApiEndPoint = "https://dnspod.tencentcloudapi.com"

// TencentDNSClient 腾讯云 DNSPod 客户端
type TencentDNSClient struct {
	secretId  string
	secretKey string
	instance  *http.Client
}

// NewTencentDNSClient 新建腾讯云 DNSPod 客户端
func NewTencentDNSClient(secretId, secretKey string) *TencentDNSClient {
	return &TencentDNSClient{
		secretId:  secretId,
		secretKey: secretKey,
		instance:  &http.Client{},
	}
}

// DescribeRecordList 获取域名的解析记录列表
func (client *TencentDNSClient) DescribeRecordList(domain, subDomain, recordType string) ([]models.DNSTencentDescribeRecordItem, error) {
	payload := make(map[string]any, 3)
	payload["Domain"] = domain
	payload["Subdomain"] = subDomain
	payload["RecordType"] = recordType

	var recordList models.DNSTencentDescribeRecordList
	if err := client.send("DescribeRecordList", payload, &recordList); err != nil {
		return nil, err
	}

	if recordList.Response.Error.Message != "" {
		return nil, errors.New(recordList.Response.Error.Message)
	}

	return recordList.Response.RecordList, nil
}

// CreateRecord 添加记录
func (client *TencentDNSClient) CreateRecord(domain, subDomain, recordType, value string) (int, error) {
	payload := make(map[string]any, 5)
	payload["Domain"] = domain
	payload["RecordType"] = recordType
	payload["RecordLine"] = "默认"
	payload["Value"] = value
	payload["SubDomain"] = subDomain

	var createRecord models.DNSTencentCreateRecord
	if err := client.send("CreateRecord", payload, &createRecord); err != nil {
		return 0, err
	}

	if createRecord.Response.Error.Message != "" {
		return 0, errors.New(createRecord.Response.Error.Message)
	}

	return createRecord.Response.RecordId, nil
}

// DeleteRecord 删除记录
func (client *TencentDNSClient) DeleteRecord(domain string, recordId int) error {
	payload := make(map[string]any, 2)
	payload["Domain"] = domain
	payload["RecordId"] = recordId

	var deleteRecord models.DNSTencentDeleteRecord
	if err := client.send("DeleteRecord", payload, &deleteRecord); err != nil {
		return err
	}

	if deleteRecord.Response.Error.Message != "" {
		return errors.New(deleteRecord.Response.Error.Message)
	}

	return nil
}

// ModifyDNS 修改记录
func (client *TencentDNSClient) ModifyDNS(domain, subDomain, recordType string, recordId int, value string) (int, error) {
	payload := make(map[string]any, 6)
	payload["Domain"] = domain
	payload["RecordType"] = recordType
	payload["RecordId"] = recordId
	payload["RecordLine"] = "默认"
	payload["Value"] = value
	payload["SubDomain"] = subDomain

	var modifyRecord models.DNSTencentModifyRecord
	if err := client.send("ModifyRecord", payload, &modifyRecord); err != nil {
		return 0, err
	}

	if modifyRecord.Response.Error.Message != "" {
		return 0, errors.New(modifyRecord.Response.Error.Message)
	}

	return modifyRecord.Response.RecordId, nil
}

// send 发送请求
func (client *TencentDNSClient) send(action string, payload map[string]any, v any) error {
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	r, _ := http.NewRequest(http.MethodPost, tencentApiEndPoint, bytes.NewReader(payloadBytes))

	r.Header.Set("Content-Type", "application/json")
	r.Header.Set("X-TC-Version", tencentApiVersion)

	utils.TencentDNSSigner(r, client.secretId, client.secretKey, action, string(payloadBytes))

	resp, err := client.instance.Do(r)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	if v != nil {
		return json.Unmarshal(body, v)
	}

	return nil
}
