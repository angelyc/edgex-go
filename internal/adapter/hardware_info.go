package adapter

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"os/exec"
	"regexp"
	"sort"
	"strings"
)

type PersonInfo struct {
	Name    string
	age     int32
	Sex     bool
	Hobbies []string
}

func InitGatewayID() {
	hardware, err := getHardwareInfo()
	if err != nil {
		fmt.Println("ERROR: failed to get hardware info")
	}
	// 计算
	GatewayId = generateGatewayId(hardware)
	hardware.Id = GatewayId
	// 保存hardware.json
	saveHardwareInfo(hardware)
}

func getHardwareInfo() (GatewayInfo, error) {
	var hardware GatewayInfo
	fileJson, err := os.Open(HardwareInfoFile)
	defer fileJson.Close()
	// 读取文件中保存的信息
	if err == nil {
		decoder := json.NewDecoder(fileJson)
		err = decoder.Decode(&hardware)
	}
	// 文件不存在 || 读文件失败，获取系统信息
	if err != nil {
		// 获取物理信息
		hardware.ProcessorID, err = processorId()
		if err == nil {
			hardware.NetIfs = getNetIfs()
			if err != nil {
				fmt.Println("ERROR: failed get Net interfaces")
			}
		} else {
			fmt.Println("ERROR: failed get processid")
		}
	}
	return hardware, err
}

func getNetIfs() (netIfs []NetIf) {
	nis, err := net.Interfaces()
	if err != nil {
		fmt.Println("ERROR: failed get Net interfaces")
		return
	}
	for _, ni := range nis {
		var ipv4 = ""
		addresses, _ := ni.Addrs()
		for _, v := range addresses {
			if ipv4 != "" {
				ipv4 += ", "
			}
			ipv4 += v.String()
		}
		netIf := NetIf{Name: ni.Name,
			DisplayName: ni.Name,
			MacAddr:     ni.HardwareAddr.String(),
			Ipv4Addr:    ipv4,
			Ipv6Addr:    "",
			IsUp:        netIfIsUp(ni)}
		netIfs = append(netIfs, netIf)
	}
	return
}

func netIfIsUp(p net.Interface) bool {
	if p.Flags%2 != 0 {
		return true
	} else {
		return false
	}

}
func generateGatewayId(hardware GatewayInfo) string {
	sort.Slice(hardware.NetIfs, func(i, j int) bool {
		stri := hardware.NetIfs[i].MacAddr
		strj := hardware.NetIfs[j].MacAddr
		ret := strings.Compare(stri, strj)
		if ret < 0 {
			return false
		}
		return true
	})

	var macs string
	for _, netIf := range hardware.NetIfs {
		macs += "#" + netIf.MacAddr
	}
	gateway := macs + hardware.ProcessorID
	fmt.Println("INFO: gateway id is:" + gateway)
	return getMd5String(gateway, false, false)
}

//生成32位md5字串
func getMd5String(s string, upper bool, half bool) string {
	h := md5.New()
	h.Write([]byte(s))
	result := hex.EncodeToString(h.Sum(nil))
	if upper == true {
		result = strings.ToUpper(result)
	}
	if half == true {
		result = result[8:24]
	}
	return result
}
func saveHardwareInfo(info GatewayInfo) {
	// 创建文件
	filePtr, err := os.Create(HardwareInfoFile)
	if err != nil {
		fmt.Println("Create file failed", err.Error())
		return
	}
	defer filePtr.Close()

	// 创建Json编码器
	encoder := json.NewEncoder(filePtr)

	err = encoder.Encode(info)
	if err != nil {
		fmt.Println("Encoder failed", err.Error())

	} else {
		fmt.Println("Encoder success")
	}

}
func processorId() (string, error) {
	cmd := exec.Command("wmic", "cpu", "get", "ProcessorID")
	out, err := cmd.CombinedOutput()
	if err != nil {
		return "", err
	}
	str := string(out)
	//匹配一个或多个空白符的正则表达式
	reg := regexp.MustCompile("\\s+")
	str = reg.ReplaceAllString(str, "")
	//正则处理后结果为：ProcessorIdBFEBFBFF000906E9，截取ProcessorId后面的值
	return str[11:], nil
}
func ReadFile() GatewayInfo {
	var hardware GatewayInfo
	filePtr, err := os.Open(HardwareInfoFile)
	if err != nil {
		fmt.Println("Open file failed [Err:%s]", err.Error())
		return hardware
	}
	defer filePtr.Close()

	// 创建json解码器
	decoder := json.NewDecoder(filePtr)
	err = decoder.Decode(&hardware)
	if err != nil {
		fmt.Println("Decoder failed", err.Error())
	}
	return hardware
}

func WriteFile() {
	personInfo := []PersonInfo{{"David", 30, true, []string{"跑步", "读书", "看电影"}}, {"Lee", 27, false, []string{"工作", "读书", "看电影"}}}

	// 创建文件
	filePtr, err := os.Create("person_info.json")
	if err != nil {
		fmt.Println("Create file failed", err.Error())
		return
	}
	defer filePtr.Close()

	// 创建Json编码器
	encoder := json.NewEncoder(filePtr)

	err = encoder.Encode(personInfo)
	if err != nil {
		fmt.Println("Encoder failed", err.Error())

	} else {
		fmt.Println("Encoder success")
	}

	// 带JSON缩进格式写文件
	data, err := json.MarshalIndent(personInfo, "", "  ")
	if err != nil {
		fmt.Println("Encoder failed", err.Error())

	} else {
		fmt.Println("Encoder success")
	}

	filePtr.Write(data)
}
