package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"time"

	"github.com/skip2/go-qrcode"
)

// 生成随机IMEI
func generateIMEI() string {
	rand.Seed(time.Now().UnixNano())

	// IMEI格式: TAC FAC SNR SVN
	tac := fmt.Sprintf("%06d", rand.Intn(999999))
	fac := fmt.Sprintf("%02d", rand.Intn(99))
	snr := fmt.Sprintf("%06d", rand.Intn(999999))
	svn := fmt.Sprintf("%02d", rand.Intn(99))

	imei := fmt.Sprintf("%s%s%s%s", tac, fac, snr, svn)
	return imei
}

// 处理HTTP请求
func handleRequest(w http.ResponseWriter, r *http.Request) {
	// 生成随机IMEI
	imei := generateIMEI()

	// 生成二维码
	qrCode, err := qrcode.Encode(imei, qrcode.Medium, 256)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// 设置响应头为PNG图像
	w.Header().Set("Content-Type", "image/png")
	w.Header().Set("Content-Disposition", "inline; filename=\"qrcode.png\"")

	// 将二维码直接写入响应体
	w.Write(qrCode)
}

func main() {
	// 注册路由
	http.HandleFunc("/generate", handleRequest)

	// 启动HTTP服务器
	fmt.Println("Server started on :8080")
	http.ListenAndServe(":8080", nil)
}
