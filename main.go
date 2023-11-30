package main

import (
	"flag"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
)

func main() {
	// 解析命令行参数
	listenAddr, proxyAddr, targetURL := parseCommandLineArgs()

	// 验证必需的参数
	if targetURL == "" {
		fmt.Println("Error: Target URL is required.")
		flag.PrintDefaults()
		os.Exit(1)
	}

	// 解析目标URL
	target, err := url.Parse(targetURL)
	if err != nil {
		fmt.Println("Error parsing target URL:", err)
		os.Exit(1)
	}

	proxy, err := url.Parse(proxyAddr)
	if err != nil {
		fmt.Println("Error parsing proxy address:", err)
		os.Exit(1)
	}

	// 创建自定义的 Transport 对象，配置代理
	transport := &http.Transport{
		Proxy: http.ProxyURL(proxy),
	}

	// 创建自定义的 http.Client 对象，使用自定义的 Transport
	client := &http.Client{
		Transport: transport,
	}

	// 设置处理函数，使用自定义的 http.Client 处理所有请求
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// 更新请求的主机头，确保目标服务器知道请求的原始主机
		r.Host = target.Host

		// 打印请求地址和转发地址到日志
		fmt.Printf("Request received: %s %s\n", r.Method, r.URL)
		fmt.Printf("Forwarding to: %s\n", target)

		// 创建新请求，手动设置 RequestURI
		newReq, err := http.NewRequest(r.Method, target.String()+r.RequestURI, r.Body)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error creating new request: %v", err), http.StatusInternalServerError)
			return
		}

		// 复制请求头
		newReq.Header = r.Header

		// 发送新请求到目标服务器
		resp, err := client.Do(newReq)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error forwarding request: %v", err), http.StatusInternalServerError)
			return
		}
		defer resp.Body.Close()

		// 将代理响应返回给客户端
		for key, values := range resp.Header {
			for _, value := range values {
				w.Header().Add(key, value)
			}
		}
		w.WriteHeader(resp.StatusCode)

		// 将代理响应的主体写入响应
		_, _ = io.Copy(w, resp.Body)
	})

	// 启动HTTP服务，监听指定地址和端口
	fmt.Printf("Reverse Proxy is listening on %s\n", listenAddr)
	err = http.ListenAndServe(listenAddr, nil)
	if err != nil {
		fmt.Println("Error starting server:", err)
		os.Exit(1)
	}
}

func parseCommandLineArgs() (string, string, string) {
	// 定义命令行参数
	listenAddr := flag.String("l", "127.0.0.1:8081", "Listen address")
	proxyAddr := flag.String("p", "", "Proxy address (optional)")
	targetURL := flag.String("t", "", "Target URL")
	flag.Parse()

	// 返回解析后的参数
	return *listenAddr, *proxyAddr, *targetURL
}
