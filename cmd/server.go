package main

import (
	"DNS-Resolver/lib/config"
	"fmt"
	"log"
	"net"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/miekg/dns"
	"github.com/patrickmn/go-cache"
	"golang.org/x/net/publicsuffix"
)

var CONFIG config.Configs
var BASE_CONFIG config.Config
var CACHE cache.Cache

func extractSubdomainAndRoot(url string) (string, string) {
	// Regex to extract the domain part from a URL
	regex := regexp.MustCompile(`(?:https?:\/\/)?(?:www\.)?([^\/\s]+)`)
	matches := regex.FindStringSubmatch(url)

	if len(matches) < 2 {
		return "@", "" // No valid match found
	}

	domain := matches[1]
	domain = strings.TrimSuffix(domain, ".")

	// Get the root domain using publicsuffix package
	root, err := publicsuffix.EffectiveTLDPlusOne(domain)
	if err != nil {
		return "@", domain // Fallback if publicsuffix fails
	}

	// Extract subdomain by removing the root domain
	if domain == root {
		return "@", root // No subdomain
	}

	subdomain := strings.TrimSuffix(domain, "."+root)
	return subdomain, root
}

func createResourceRecord(q *dns.Question, record config.Record) dns.RR {
	var response dns.RR

	switch q.Qtype {
	case dns.TypeA:
		{
			response = &dns.A{
				Hdr: dns.RR_Header{Name: q.Name, Rrtype: dns.TypeA, Class: dns.ClassINET, Ttl: 60},
				A:   net.ParseIP(record.Value),
			}
		}
	case dns.TypeCNAME:
		{
			response = &dns.CNAME{
				Hdr: dns.RR_Header{
					Name:   q.Name,
					Rrtype: dns.TypeCNAME,
					Ttl:    uint32(record.TTL),
				},
				Target: record.Value,
			}
		}
	}

	return response
}

func handleQuery(q *dns.Question, msg *dns.Msg) []dns.RR {
	client := new(dns.Client)
	subDomain, rootDomain := extractSubdomainAndRoot(q.Name)
	domain := CONFIG[rootDomain+"."]

	record := domain.GetSubRecord(q.Qtype, subDomain)

	// * Look up the domain from another DNS server
	if domain.Domain == "" || len(domain.Domain) == 0 {
		fmt.Println("Not Found")
		var in *dns.Msg
		var err error

		msg := new(dns.Msg)
		msg.Question = append(msg.Question, *q)
		fmt.Println(msg.Question)
		for _, server := range BASE_CONFIG.DNS_Resolvers {
			fmt.Println(server + ":53")

			if !BASE_CONFIG.DNS_Over_HTTPS && !BASE_CONFIG.DNS_Over_TLS {
				in, _, err = client.Exchange(msg, server+":53")
			}

			fmt.Println(in)
			fmt.Println(in.Response)

			if err != nil && strings.Contains(err.Error(), "connection refused") {
				continue
			} else if err != nil || len(in.Answer) == 0 {
				continue
			}

			break
		}

		CACHE.Set(q.Name, in.Answer, cache.DefaultExpiration)
		return in.Answer
	}

	results := []dns.RR{}
	results = append(results, createResourceRecord(q, record))
	return results
}

// handleDNSRequest processes incoming DNS queries
func handleDNSRequest(w dns.ResponseWriter, r *dns.Msg) {
	msg := dns.Msg{}
	msg.SetReply(r)

	for _, q := range r.Question {
		var res []dns.RR
		if cacheRes, found := CACHE.Get(q.Name); found {
			fmt.Println("Hit cache")
			if len(cacheRes.([]dns.RR)) == 0 {
				CACHE.Delete(q.Name)
				res = handleQuery(&q, &msg)
			} else {
				res = cacheRes.([]dns.RR)
			}
		} else {
			res = handleQuery(&q, &msg)
		}

		msg.Answer = append(msg.Answer, res...)
	}

	err := w.WriteMsg(&msg)
	if err != nil {
		log.Printf("Failed to write DNS response: %v", err)
	}
}

// TODO: Create a cache for recently searched addresses
func main() {
	configFiles, err := config.ReadConfigDir("test-configs")
	if err != nil {
		fmt.Println(err)
	}

	configs, err := config.ReadConfig(configFiles)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	baseConfig, err := config.ReadBaseConfig("test-configs/config.yml")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	} else if len(baseConfig.DNS_Resolvers) == 0 {
		fmt.Println("Need to provide DNS resolvers")
		os.Exit(1)
	}

	BASE_CONFIG = baseConfig
	CONFIG = configs

	fmt.Println(configs)
	fmt.Println(baseConfig)

	// * SETUP Cache
	cache := cache.New(10*time.Minute, 20*time.Minute)
	CACHE = *cache

	// Create a DNS server
	dns.HandleFunc(".", handleDNSRequest)

	server := &dns.Server{
		Addr: "0.0.0.0:8080", // Listen on all interfaces, port 8080
		Net:  "udp",          // UDP-based DNS
	}

	fmt.Println("DNS server listening on port 8080...")
	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
