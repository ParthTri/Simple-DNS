package config

import (
	"os"
	"strings"

	"github.com/miekg/dns"
	"gopkg.in/yaml.v3"
)

func ReadConfig(configs []string) (Configs, error) {
	var config Configs = Configs{};

	for _, path := range configs {
		b := Domain{}

		content, err := os.ReadFile(path)
		if err != nil {
			return config, err
		}

		err = yaml.Unmarshal(content, &b)
		if err != nil {
			return config, err
		}

		config[b.Domain + "."] = b
	}

	return config, nil
}

func ReadConfigDir(path string) ([]string, error) {
	files, err := os.ReadDir(path)
	var configFiles []string = []string{};
	if err != nil {
		return nil, err
	}

	for _, file := range files {
		name := strings.Split(file.Name(), ".")
		suffix := name[len(name)-1]
		if suffix == "yaml" || suffix == "yml" {
			configFiles = append(configFiles, path + "/" + file.Name())
		}
	}

	return configFiles, err
}

func (d Domain) GetSubRecord(recordType uint16, subDomain string) Record {
	var records []Record = []Record{};

	switch recordType {
		case dns.TypeA:
			records = d.Records.A;
		case dns.TypeCNAME:
			records = d.Records.CNAME;
		case dns.TypeTXT:
			records = d.Records.TXT;
	}

	for _, record := range records {
		if subDomain == record.Name {
			return record;
		}
	}
	
	return Record{};
}
