package config

// * A type for the application for base DNS server and other data
type Config struct {
	DNS_Resolvers  []string
	DNS_Over_HTTPS bool
	DNS_Over_TLS   bool
}

type Record struct {
	Name  string
	Value string
	TTL   int64
}

type Records struct {
	A     []Record `yaml:"A"`
	CNAME []Record `yaml:"CNAME"`
	TXT   []Record `yaml:"TXT"`
}

type Domain struct {
	Domain  string `yaml:"domain"`
	Records Records
}

type Configs map[string]Domain
