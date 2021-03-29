package structure

type CasbinContent struct {
	Ptype           string `json:"ptype" gorm:"column:p_type"`
	AuthorityObject string `json:"role_name" gorm:"column:v0"`
	Path            string `json:"path" gorm:"column:v1"`
	Method          string `json:"method" gorm:"column:v2"`
}

// Casbin structure for input parameters
type CasbinRelation struct {
	AuthorityObject string          `json:"authority_object"`
	CasbinBucket    []CasbinContent `json:"casbin_bucket"`
}
