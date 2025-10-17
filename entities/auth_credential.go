package entities

const DB_TABLE_AUTH_CREDENTIAL = "auth_credential"

type AuthCredential struct {
	CredentialId    string `gorm:"column:credential_id;type:VARCHAR(100);primaryKey" json:"credential_id,omitempty" yaml:"credential_id,omitempty" field-id:"credential_id" field-type:"varchar" field-comment:"credential ID(主键)"` // credential ID(主键)
	CredentialName  string `gorm:"column:credential_name;type:VARCHAR(100)" json:"credential_name,omitempty" yaml:"credential_name,omitempty" field-id:"credential_name" field-type:"varchar" field-comment:"名称"`                   // 名称
	CredentialValue string `gorm:"column:credential_value;type:VARCHAR(100)" json:"credential_value,omitempty" yaml:"credential_value,omitempty" field-id:"credential_value" field-type:"varchar" field-comment:"credential"`       // credential
	CredentialType  string `gorm:"column:credential_type;type:VARCHAR(30)" json:"credential_type,omitempty" yaml:"credential_type,omitempty" field-id:"credential_type" field-type:"varchar" field-comment:"credential类型(key,jwt)"` // credential类型(key,jwt)
	ExpireTime      string `gorm:"column:expire_time;type:VARCHAR(30)" json:"expire_time,omitempty" yaml:"expire_time,omitempty" field-id:"expire_time" field-type:"varchar" field-comment:"过期时间"`                                  // 过期时间
	ObjectType      string `gorm:"column:object_type;type:VARCHAR(30)" json:"object_type,omitempty" yaml:"object_type,omitempty" field-id:"object_type" field-type:"varchar" field-comment:"对象类型(user,space,org)"`                  // 对象类型(user,space,org)
	ObjectId        string `gorm:"column:object_id;type:VARCHAR(100)" json:"object_id,omitempty" yaml:"object_id,omitempty" field-id:"object_id" field-type:"varchar" field-comment:"对象ID"`                                         // 对象ID
	RefId           string `gorm:"column:ref_id;type:VARCHAR(100)" json:"ref_id,omitempty" yaml:"ref_id,omitempty" field-id:"ref_id" field-type:"varchar" field-comment:"关联ID"`                                                     // 关联ID
	OrgId           string `gorm:"column:org_id;type:VARCHAR(100)" json:"org_id,omitempty" yaml:"org_id,omitempty" field-id:"org_id" field-type:"varchar" field-comment:"机构ID"`                                                     // 机构ID
	Status          string `gorm:"column:status;type:VARCHAR(30)" json:"status,omitempty" yaml:"status,omitempty" field-id:"status" field-type:"varchar" field-comment:"状态"`                                                        // 状态
	DataUpdDate     string `gorm:"column:data_upd_date;type:VARCHAR(10)" json:"data_upd_date,omitempty" yaml:"data_upd_date,omitempty" field-id:"data_upd_date" field-type:"varchar" field-comment:"数据修改日期(yyyymmdd)"`              // 数据修改日期(yyyymmdd)
	DataUpdTime     string `gorm:"column:data_upd_time;type:VARCHAR(20)" json:"data_upd_time,omitempty" yaml:"data_upd_time,omitempty" field-id:"data_upd_time" field-type:"varchar" field-comment:"数据修改时间(yyyymmddHHMMSS)"`        // 数据修改时间(yyyymmddHHMMSS)
	DataCrtDate     string `gorm:"column:data_crt_date;type:VARCHAR(10)" json:"data_crt_date,omitempty" yaml:"data_crt_date,omitempty" field-id:"data_crt_date" field-type:"varchar" field-comment:"数据创建日期(yyyymmdd)"`              // 数据创建日期(yyyymmdd)
	DataCrtTime     string `gorm:"column:data_crt_time;type:VARCHAR(20)" json:"data_crt_time,omitempty" yaml:"data_crt_time,omitempty" field-id:"data_crt_time" field-type:"varchar" field-comment:"数据创建时间(yyyymmddHHMMSS)"`        // 数据创建时间(yyyymmddHHMMSS)
}

func (AuthCredential) TableName() string { return "auth_credential" }
