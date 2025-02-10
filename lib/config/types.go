package config

type Record struct {
	Name string
	Value string
	TTL int64
}


type Records struct {
	A  []Record `yaml:"A"`
	CNAME []Record `yaml:"CNAME"`
	TXT []Record `yaml:"TXT"`
}

type Domain struct {
	Domain string `yaml:"domain"`
	Records Records
}

type Configs map[string]Domain

