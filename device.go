package distiot

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
)

/* 设备管理器
需要先调用NewManager初始化管理器，使用管理器生成设备
设备id和用户token需要从平台获取
*/
type DeviceManager struct {
	token     string // 用户token
	MasterUrl string // master服务器地址
	UserUrl   string // user服务器地址
}

/* 单个设备
用以上传，读取数据等操作
*/
type Device struct {
	ID       int    // 设备ID
	token    string //用户token
	NodeAddr string // 节点服务器地址，只包含IP，不包含头部
	NodePort int    // 节点服务器端口
}

//创建一个新的设备管理器，传入用户token。管理器可以管理多个device
func NewManager(token string) *DeviceManager {
	//未指定URL，则为默认的URL
	masterUrl := "http://master.distiot.ri-co.cn/master"
	userUrl := "http://user.distiot.ri-co.cn/user"
	return &DeviceManager{
		token:     token,
		MasterUrl: masterUrl,
		UserUrl:   userUrl,
	}
}

//创建一个新设备，并获取该设备的node信息
func (m *DeviceManager) NewDevice(did int) (*Device, error) {
	var device Device
	device.ID = did
	device.token = m.token
	addr, port, err := m.getNode(did)
	if err != nil {
		return nil, err
	}
	device.NodeAddr = addr
	device.NodePort = port
	return &device, nil
}

//获取设备的node节点信息，返回节点addr，端口
func (m *DeviceManager) getNode(did int) (string, int, error) {
	//向master获取node信息
	Url, err := url.Parse(m.MasterUrl + "/getNode")
	if err != nil {
		return "", 0, err
	}
	params := url.Values{}
	params.Set("token", m.token)
	params.Set("id", strconv.Itoa(did))
	Url.RawQuery = params.Encode()
	res, err := http.Get(Url.String())
	if err != nil {
		return "", 0, err
	}
	//解析node信息
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", 0, err
	}
	var node nodeData
	err = json.Unmarshal(body, &node)
	if err != nil {
		err2 := errors.New(string(body))
		return "", 0, err2
	}
	return node.Addr, node.Port, nil
}

type nodeData struct {
	ID   int    `json:"id"`
	Addr string `json:"addr"`
	Port int    `json:"port"`
}

/* Http上传数据 一次只能上传一条数据，非异步请求，传入字符串格式的数据
请确保正确初始化后再上传
*/
func (d *Device) UploadDataHttp(data string) error {
	//使用高效的拼接方式，先初始化请求地址
	var strs bytes.Buffer
	strs.WriteString("http://")
	strs.WriteString(d.NodeAddr)
	strs.WriteString(":")
	strs.WriteString(strconv.Itoa(d.NodePort))
	strs.WriteString("/node/dataWriteSingle")

	//设置GET请求参数
	params := url.Values{}
	params.Set("token", d.token)
	params.Set("did", strconv.Itoa(d.ID))
	params.Set("data", data)

	Url, err := url.Parse(strs.String())
	if err != nil {
		return err
	}
	Url.RawQuery = params.Encode()
	res, err := http.Get(Url.String())
	if err != nil {
		return err
	}
	defer res.Body.Close()
	if res.Status != "200 OK" {
		//解析node信息
		body, _ := ioutil.ReadAll(res.Body)
		return errors.New(string(body))
	}
	return nil
}
