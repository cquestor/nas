package models

// DNSTencentError 腾讯云 DNSPod 错误响应
type DNSTencentError struct {
	Error struct {
		Code    string `json:"Code"`
		Message string `json:"Message"`
	} `json:"Error,omitempty"`
}

// DNSTencentDescribeRecordList 腾讯云 DNSPod 记录列表
type DNSTencentDescribeRecordList struct {
	Response struct {
		DNSTencentError
		RequestId       string `json:"RequestId"`
		RecordCountInfo struct {
			SubdomainCount int `json:"SubdomainCount"`
			TotalCount     int `json:"TotalCount"`
			ListCount      int `json:"ListCount"`
		} `json:"RecordCountInfo"`
		RecordList []DNSTencentDescribeRecordItem `json:"RecordList"`
	} `json:"Response"`
}

// DNSTencentDescribeRecordItem 腾讯云 DNSPod 记录信息
type DNSTencentDescribeRecordItem struct {
	Value         string `json:"Value"`
	Status        string `json:"Status"`
	UpdatedOn     string `json:"UpdatedOn"`
	Name          string `json:"Name"`
	Line          string `json:"Line"`
	LineId        string `json:"LineId"`
	Type          string `json:"Type"`
	MonitorStatus string `json:"MonitorStatus"`
	Remark        string `json:"Remark"`
	RecordId      int    `json:"RecordId"`
	TTL           int    `json:"TTL"`
	MX            int    `json:"MX"`
}

// DNSTencentCreateRecord 添加腾讯云 DNSPod 记录
type DNSTencentCreateRecord struct {
	Response struct {
		DNSTencentError
		RequestId string `json:"RequestId"`
		RecordId  int    `json:"RecordId"`
	} `json:"Response"`
}

// DNSTencentDeleteRecord 删除腾讯云 DNSPod 记录
type DNSTencentDeleteRecord struct {
	Response struct {
		DNSTencentError
		RequestId string `json:"RequestId"`
	} `json:"Response"`
}

// DNSTencentModifyRecord 修改腾讯云 DNSPod 记录
type DNSTencentModifyRecord struct {
	Response struct {
		DNSTencentError
		RequestId string `json:"RequestId"`
		RecordId  int    `json:"RecordId"`
	} `json:"Response"`
}
