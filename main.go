package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"net/netip"

	"github.com/mdlayher/arp"
	"github.com/mdlayher/ethernet"
)

func main() {
	// 定义命令行参数
	ifaceName := flag.String("i", "", "网络接口名 (如 eth0)")
	macStr := flag.String("m", "", "源MAC地址 (如 00:11:22:33:44:55)")
	ipStr := flag.String("ip", "", "源IP地址 (如 192.168.1.100)")

	flag.Parse()

	// 验证必需参数
	if *ifaceName == "" || *macStr == "" || *ipStr == "" {
		fmt.Println("错误：缺少必需参数")
		fmt.Println("用法示例:")
		flag.PrintDefaults()
		log.Fatalln("退出：参数不完整")
	}
	// 解析MAC地址
	srcMAC, err := net.ParseMAC(*macStr)
	if err != nil {
		log.Fatalf("MAC地址解析失败: %v", err)
	}

	// 解析IP地址
	srcIP, err := netip.ParseAddr(*ipStr)
	if err != nil {
		log.Fatalf("IP地址解析失败: %v", err)
	}

	// 获取网络接口
	iface, err := net.InterfaceByName(*ifaceName)
	if err != nil {
		log.Fatalf("无法获取接口: %v", err)
	}

	// 创建ARP客户端
	client, err := arp.Dial(iface)
	if err != nil {
		log.Fatalf("创建ARP客户端失败: %v", err)
	}
	defer client.Close()

	// 构造GARP请求
	packet := &arp.Packet{
		HardwareType:       1, // 以太网
		Operation:          arp.OperationRequest,
		SenderHardwareAddr: srcMAC,
		SenderIP:           srcIP,
		TargetHardwareAddr: ethernet.Broadcast, // 广播地址
		TargetIP:           srcIP,              // 目标IP与源IP相同
	}

	// 发送GARP包
	if err := client.WriteTo(packet, ethernet.Broadcast); err != nil {
		log.Fatalf("发送失败: %v", err)
	}

	log.Printf("GARP包已通过接口 %s 发送! 源IP: %s, 源MAC: %s\n", *ifaceName, srcIP, srcMAC)
}
