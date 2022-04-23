package distiot

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
)

type DeviceManager struct {
	token     string // 用户token
	MasterUrl string // master服务器地址
	UserUrl   string // user服务器地址
}

type Device struct {
	ID       int    // 设备ID
	NodeAddr string // 节点服务器地址，只包含IP，不包含头部
	NodePort int    // 节点服务器端口
}

//创建一个新的设备管理器，传入用户token。管理器可以管理多个device
func NewManager(token string) *DeviceManager {
	//未指定URL，则为默认的URL
	masterUrl := "http://master.distiot.ri-co.cn"
	userUrl := "http://user.distiot.ri-co.cn"
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
	params.Set("did", strconv.Itoa(did))
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
		return "", 0, err
	}
	return node.Addr, node.Port, nil
}

type nodeData struct {
	ID   int    `json:"id"`
	Addr string `json:"addr"`
	Port int    `json:"port"`
}
