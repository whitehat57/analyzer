package main

import (
	"bufio"
	"crypto/tls"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/fatih/color"
)

// Pengecekan jika header mengandung string tertentu
func contains(header http.Header, key, value string) bool {
	return strings.Contains(header.Get(key), value)
}

func checkCDN(headers http.Header) bool {
	return contains(headers, "Server", "cloudflare") || contains(headers, "X-CDN", "akamai")
}

func checkLoadBalancer(headers http.Header) bool {
	return contains(headers, "X-Forwarded-For", "") || contains(headers, "Via", "")
}

func checkFirewall(headers http.Header) bool {
	return contains(headers, "CF-RAY", "") || contains(headers, "X-Sucuri-ID", "")
}

func checkRateLimiting(headers http.Header) bool {
	return contains(headers, "Retry-After", "") || contains(headers, "X-RateLimit-Limit", "")
}

func checkAutoscaling(headers http.Header) bool {
	return contains(headers, "X-Amzn-Trace-Id", "") || contains(headers, "X-Request-Id", "")
}

func checkWAF(headers http.Header) bool {
	return contains(headers, "X-Mod-Security", "") || contains(headers, "X-Sucuri-ID", "")
}

// Cek latensi sebagai indikator CDN atau autoscaling
func checkLatency(url string) {
	fmt.Println("🌍 Memeriksa latensi...")
	start := time.Now()
	_, err := http.Get(url)
	if err != nil {
		fmt.Printf("Error fetching the URL: %v\n", err)
		return
	}
	duration := time.Since(start)

	fmt.Printf("⏱ Latensi: %v\n", duration)
	if duration < 100*time.Millisecond {
		color.Green("✅ Kemungkinan CDN Detected (latensi rendah)")
	} else {
		color.Red("❌ CDN Tidak Terdeteksi")
	}
}

// Inspeksi sertifikat SSL untuk mengetahui penerbit
func checkSSL(url string) {
	fmt.Println("🔐 Memeriksa sertifikat SSL...")
	if !strings.HasPrefix(url, "https") {
		color.Red("❌ Tidak menggunakan HTTPS, SSL tidak bisa diperiksa")
		return
	}

	u := strings.Split(url, "//")[1]
	conn, err := tls.Dial("tcp", u+":443", nil)
	if err != nil {
		fmt.Printf("Error fetching SSL certificate: %v\n", err)
		return
	}
	defer conn.Close()

	cert := conn.ConnectionState().PeerCertificates[0]
	issuer := cert.Issuer.Organization

	fmt.Printf("🔑 Sertifikat SSL diterbitkan oleh: %s\n", issuer)
	if strings.Contains(strings.Join(issuer, ""), "Cloudflare") || strings.Contains(strings.Join(issuer, ""), "Akamai") {
		color.Green("✅ Kemungkinan CDN/WAF Detected (SSL)")
	} else {
		color.Red("❌ CDN/WAF Tidak Terdeteksi berdasarkan SSL")
	}
}

// Teknik fuzzing sederhana untuk mendeteksi WAF dan Rate Limiting
func fuzzingCheck(url string) {
	fmt.Println("🧪 Melakukan fuzzing sederhana untuk memeriksa WAF dan Rate Limiting...")

	client := &http.Client{}
	req, _ := http.NewRequest("GET", url+"/../../../../etc/passwd", nil)
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("Error during fuzzing: %v\n", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusForbidden || resp.StatusCode == http.StatusNotFound {
		color.Green("✅ WAF Detected (respons 403 atau 404)")
	} else {
		color.Red("❌ WAF Tidak Terdeteksi")
	}

	for i := 0; i < 10; i++ {
		_, err = http.Get(url)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			return
		}
		time.Sleep(200 * time.Millisecond)
	}

	if resp.StatusCode == 429 {
		color.Green("✅ Rate Limiting Detected (status 429)")
	} else {
		color.Red("❌ Rate Limiting Tidak Terdeteksi")
	}
}

func analyzeWebsite(url string) {
	resp, err := http.Get(url)
	if err != nil {
		fmt.Printf("Error fetching the URL: %v\n", err)
		return
	}
	defer resp.Body.Close()

	headers := resp.Header

	fmt.Println("======================================")
	fmt.Printf("🔍 Analyzing: %s\n", url)
	fmt.Println("======================================")

	if checkCDN(headers) {
		color.Green("✅ CDN Detected")
	} else {
		color.Red("❌ CDN Not Detected")
	}

	if checkLoadBalancer(headers) {
		color.Green("✅ Load Balancer Detected")
	} else {
		color.Red("❌ Load Balancer Not Detected")
	}

	if checkFirewall(headers) {
		color.Green("✅ Firewall Detected")
	} else {
		color.Red("❌ Firewall Not Detected")
	}

	if checkRateLimiting(headers) {
		color.Green("✅ Rate Limiting Detected")
	} else {
		color.Red("❌ Rate Limiting Not Detected")
	}

	if checkAutoscaling(headers) {
		color.Green("✅ Autoscaling Detected")
	} else {
		color.Red("❌ Autoscaling Not Detected")
	}

	if checkWAF(headers) {
		color.Green("✅ WAF Detected")
	} else {
		color.Red("❌ WAF Not Detected")
	}
	fmt.Println("======================================")

	checkLatency(url)
	checkSSL(url)
	fuzzingCheck(url)
}

func main() {
	// Tampilan ASCII art dan informasi script
	fmt.Println(`


	███▄▄▄▄      ▄██████▄   ▄█     ▄████████  ▄█  ███▄▄▄▄      ▄██████▄  
	███▀▀▀██▄   ███    ███ ███    ███    ███ ███  ███▀▀▀██▄   ███    ███ 
	███   ███   ███    █▀  ███▌   ███    █▀  ███▌ ███   ███   ███    █▀  
	███   ███  ▄███        ███▌   ███        ███▌ ███   ███  ▄███        
	███   ███ ▀▀███ ████▄  ███▌ ▀███████████ ███▌ ███   ███ ▀▀███ ████▄  
	███   ███   ███    ███ ███           ███ ███  ███   ███   ███    ███ 
	███   ███   ███    ███ ███     ▄█    ███ ███  ███   ███   ███    ███ 
	 ▀█   █▀    ████████▀  █▀    ▄████████▀  █▀    ▀█   █▀    ████████▀   v.1.0
	 
coded by: d57 https://github.com/whitehat57
	`)

	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Masukkan URL yang ingin dianalisis (contoh: https://example.com): ")
	url, _ := reader.ReadString('\n')
	url = strings.TrimSpace(url)

	if !strings.HasPrefix(url, "http") {
		fmt.Println("URL harus dimulai dengan http atau https.")
		return
	}

	analyzeWebsite(url)
}
