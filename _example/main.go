package main

import "fmt"

func main() {
	//首先，初始化一个设备管理器
	man := NewManager("d703fcc1-655e-4a4f-bdb1-5fecd89b07cb")
	//手动设置master和user服务器，用于本地测试，正式上线后不需要设置
	man.MasterUrl = "http://localhost:8090/master"
	man.UserUrl = "http://localhost:8091/user"

	//通过设备管理器创建设备，可以创建多个，参数为设备的ID
	d1, err := man.NewDevice(14)
	if err != nil {
		fmt.Println("创建设备出错", err.Error())
		return
	}

	//上传采集的数据，无论设备数据类型为何，统一使用字符串上传
	err = d1.UploadDataHttp("10.002")
	if err != nil {
		fmt.Println("上传数据出错 ", err.Error())
	}
}
