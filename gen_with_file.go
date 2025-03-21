package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/skip2/go-qrcode"
)

// IMEI结构体
type IMEI struct {
	IMEINumber string `json:"imei_number"`
	QRCodeURL  string `json:"qr_code_url"`
}

// 生成随机IMEI
func generateIMEI() string {
	rand.Seed(time.Now().UnixNano())

	// IMEI格式: TAC FAC SNR SVN
	// TAC (Type Allocation Code): 6位
	// FAC (Final Assembly Code): 2位
	// SNR (Serial Number): 6位
	// SVN (Software Version Number): 2位
	tac := fmt.Sprintf("%06d", rand.Intn(999999))
	fac := fmt.Sprintf("%02d", rand.Intn(99))
	snr := fmt.Sprintf("%06d", rand.Intn(999999))
	svn := fmt.Sprintf("%02d", rand.Intn(99))

	imei := fmt.Sprintf("%s%s%s%s", tac, fac, snr, svn)
	return imei
}

// 生成二维码并保存为图片
func generateQRCodeImage(imei string) (string, error) {
	// 创建二维码图片
	qrCode, err := qrcode.Encode(imei, qrcode.Medium, 256)
	if err != nil {
		return "", err
	}

	// 获取项目根目录
	projectDir, err := os.Getwd()
	if err != nil {
		return "", err
	}

	// 创建 qrcodes 目录
	qrcodesDir := filepath.Join(projectDir, "qrcodes")
	err = os.MkdirAll(qrcodesDir, 0755)
	if err != nil {
		return "", err
	}

	// 生成唯一文件名
	filename := fmt.Sprintf("qr_%s.png", time.Now().Format("20060102150405"))
	filePath := filepath.Join(qrcodesDir, filename)

	// 保存二维码图片到文件
	err = os.WriteFile(filePath, qrCode, 0644)
	if err != nil {
		return "", err
	}

	// 返回相对于项目目录的路径
	relativePath := filepath.Join("qrcodes", filename)
	return relativePath, nil
}

// 处理HTTP请求
func handleRequest(w http.ResponseWriter, r *http.Request) {
	// 生成随机IMEI
	imei := generateIMEI()

	// 生成二维码图片并获取文件路径
	qrCodePath, err := generateQRCodeImage(imei)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// 设置响应头
	w.Header().Set("Content-Type", "application/json")

	// 创建响应对象
	response := IMEI{
		IMEINumber: imei,
		QRCodeURL:  qrCodePath,
	}

	// 将响应对象编码为JSON并发送
	jsonResponse, err := json.Marshal(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write(jsonResponse)
}

func main() {
	// 设置随机数种子
	rand.Seed(time.Now().UnixNano())

	// 注册路由
	http.HandleFunc("/generate", handleRequest)

	// 启动HTTP服务器
	fmt.Println("Server started on :8080")
	http.ListenAndServe(":8080", nil)
}
