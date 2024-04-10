package main

import (
	"flag"
	"log"
	"github.com/miekg/dns"
)

func main() {
	// 定义命令行参数
	host := flag.String("host", "", "IP address to listen on")
	port := flag.String("port", "53", "Port to listen on")
	flag.Parse()

	// 创建服务器地址
	serverAddr := *host + ":" + *port

	// 创建 DNS 服务器
	server := &dns.Server{Addr: serverAddr, Net: "udp"}

	// 处理 DNS 请求
	dns.HandleFunc(".", func(w dns.ResponseWriter, r *dns.Msg) {
		m := new(dns.Msg)
		m.SetReply(r)

		// 输出请求的域名
		for _, q := range r.Question {
			log.Printf("Received DNS query for domain: %s\n", q.Name)
			rr, err := dns.NewRR(q.Name + " IN A 0.0.0.0")
			if err == nil {
				m.Answer = append(m.Answer, rr)
			}
		}

		// 将回复发送回客户端
		if err := w.WriteMsg(m); err != nil {
			log.Printf("Failed to send response: %v", err)
		}
	})

	// 启动 DNS 服务器
	log.Printf("DNS server listening on %s...\n", serverAddr)
	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("Failed to start DNS server: %v", err)
	}
}
