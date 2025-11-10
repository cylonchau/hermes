package model

// 表名设置
func (Zone) TableName() string        { return "dns_zones" }
func (Record) TableName() string      { return "dns_records" }
func (ARecord) TableName() string     { return "dns_a_records" }
func (AAAARecord) TableName() string  { return "dns_aaaa_records" }
func (CNAMERecord) TableName() string { return "dns_cname_records" }
func (MXRecord) TableName() string    { return "dns_mx_records" }
func (TXTRecord) TableName() string   { return "dns_txt_records" }
func (SRVRecord) TableName() string   { return "dns_srv_records" }
func (SOARecord) TableName() string   { return "dns_soa_records" }
func (NSRecord) TableName() string    { return "dns_ns_records" }
