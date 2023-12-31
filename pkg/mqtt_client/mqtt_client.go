package mqtt_client

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/yaxiongwu/webrtc-mqtt-server/pkg/clients"
)

const protocol = "ssl"

//const protocol = "tcp"

//const broker = "k101ceee.ala.cn-hangzhou.emqxsl.cn"
//const broker = "127.0.0.1"
const broker = "www.bxzryd.cn"

//const broker = "broker.emqx.io"
const port = 8883
const topic_connected = "$SYS/brokers/emqx@127.0.0.1/clients/+/connected"
const topic_disconnected = "$SYS/brokers/emqx@127.0.0.1/clients/+/disconnected"
const topic_sub_close_source = "close/source"
const topic_source_reg = "server/reg"
const topic_source_query = "server/query"
const topic_sys = "$SYS/brokers/#"
const username = "wuyaxiong"
const password = "wuyaxiong1982"

type MqttClient struct {
	client                  mqtt.Client
	OnSubscribe             func(string, []byte)
	OnSubscribeConnected    func([]byte)
	OnSubscribeDisconnected func([]byte)
	Topic                   string
	Qos                     int8
	SourceList              clients.SourceList
}

func MqttClientInit() *MqttClient {

	connectAddress := fmt.Sprintf("%s://%s:%d", protocol, broker, port)
	client_id := fmt.Sprintf("go-client1-%d", rand.Int())

	fmt.Println("connect address: ", connectAddress)
	opts := mqtt.NewClientOptions()
	opts.AddBroker(connectAddress)
	opts.SetUsername(username)
	opts.SetPassword(password)
	opts.SetClientID(client_id)
	opts.SetKeepAlive(time.Second * 60)

	// Optional: 设置CA证书
	// opts.SetTLSConfig(loadTLSConfig("caFilePath"))

	client := mqtt.NewClient(opts)
	token := client.Connect()
	if token.WaitTimeout(3*time.Second) && token.Error() != nil {
		log.Fatal(token.Error())
	}

	//go subscribe(client) // 在主函数里, 我们用另起一个 go 协程来订阅消息
	// time.Sleep(time.Second * 1) // 暂停一秒等待 subscribe 完成
	// publish(client)
	log.Println("MqttClientInit ok")

	return &MqttClient{
		client:     client,
		SourceList: clients.SourceList{},
	}
}

func createMqttClient() mqtt.Client {
	connectAddress := fmt.Sprintf("%s://%s:%d", protocol, broker, port)
	client_id := fmt.Sprintf("go-client1-%d", rand.Int())

	fmt.Println("connect address: ", connectAddress)
	opts := mqtt.NewClientOptions()
	opts.AddBroker(connectAddress)
	opts.SetUsername(username)
	opts.SetPassword(password)
	opts.SetClientID(client_id)
	opts.SetKeepAlive(time.Second * 60)

	// Optional: 设置CA证书
	// opts.SetTLSConfig(loadTLSConfig("caFilePath"))

	client := mqtt.NewClient(opts)
	token := client.Connect()
	if token.WaitTimeout(3*time.Second) && token.Error() != nil {
		log.Fatal(token.Error())
	}
	return client
}

// func publish(client mqtt.Client) {
// 	qos := 0
// 	msgCount := 0
// 	for {
// 		payload := fmt.Sprintf("message: %d!", msgCount)
// 		if token := client.Publish(topic, byte(qos), false, payload); token.Wait() && token.Error() != nil {
// 			fmt.Printf("publish failed, topic: %s, payload: %s\n", topic, payload)
// 		} else {
// 			fmt.Printf("publish success, topic: %s, payload: %s\n", topic, payload)
// 		}
// 		msgCount++
// 		time.Sleep(time.Second * 3)
// 	}
// }

// func subscribe(client mqtt.Client) {
// 	qos := 0
// 	client.Subscribe(topic_sys, byte(qos), func(client mqtt.Client, msg mqtt.Message) {
// 		fmt.Printf("Received `%s` from `%s` topic\n", msg.Payload(), msg.Topic())
// 	})
// 	client.Subscribe(topic, byte(qos), func(client mqtt.Client, msg mqtt.Message) {
// 		fmt.Printf("Received2 `%s` from `%s` topic\n", msg.Payload(), msg.Topic())
// 	})
// }

func loadTLSConfig(caFile string) *tls.Config {
	// load tls config
	var tlsConfig tls.Config
	tlsConfig.InsecureSkipVerify = false
	if caFile != "" {
		certpool := x509.NewCertPool()
		ca, err := ioutil.ReadFile(caFile)
		if err != nil {
			log.Fatal(err.Error())
		}
		certpool.AppendCertsFromPEM(ca)
		tlsConfig.RootCAs = certpool
	}
	return &tlsConfig
}

func (m *MqttClient) Publish(topic string, payload string, qos int8) {
	if token := m.client.Publish(topic, byte(qos), false, payload); token.Wait() && token.Error() != nil {
		fmt.Printf("publish failed, topic: %s, payload: %s\n", topic, payload)
	} else {
		fmt.Printf("publish success, topic: %s, payload: %s\n", topic, payload)
	}
}

// func (m *MqttClient) subscribe(client mqtt.Client) {
// 	//qos := 0
// 	client.Subscribe(m.Topic, byte(m.Qos), func(client mqtt.Client, msg mqtt.Message) {
// 		fmt.Printf("Received `%s` from `%s` topic\n", msg.Payload(), msg.Topic())
// 		m.OnSubscribe(msg.Topic(), msg.Payload())
// 	})
// 	// client.Subscribe(topic, byte(qos), func(client mqtt.Client, msg mqtt.Message) {
// 	// 	fmt.Printf("Received2 `%s` from `%s` topic\n", msg.Payload(), msg.Topic())
// 	// })
// }

func (m *MqttClient) Subscribe() {
	qos := 2
	m.client.Subscribe(topic_connected, byte(qos), func(client mqtt.Client, msg mqtt.Message) {
		//fmt.Printf("Received `%s` from `%s` topic\n", msg.Payload(), msg.Topic())
		m.OnSubscribeConnected(msg.Payload())
	})
	m.client.Subscribe(topic_disconnected, byte(qos), func(client mqtt.Client, msg mqtt.Message) {
		//fmt.Printf("Received `%s` from `%s` topic\n", msg.Payload(), msg.Topic())
		m.OnSubscribeDisconnected(msg.Payload())
	})
	m.client.Subscribe(topic_sub_close_source, byte(qos), func(client mqtt.Client, msg mqtt.Message) {
		fmt.Printf("Received `%s` from `%s` topic\n", msg.Payload(), msg.Topic())
		fmt.Printf("before,list: %v\n", m.SourceList.SList)
		sourceId := msg.Payload()
		for index, value := range m.SourceList.SList {
			if value.Id == string(sourceId) {
				//m.SourceList.SList[index]=nil
				m.SourceList.SList = append(m.SourceList.SList[:index], m.SourceList.SList[index+1:]...)
			}
		}
		fmt.Printf("after,list: %v\n", m.SourceList.SList)
		m.OnSubscribeDisconnected(msg.Payload())
	})
	m.client.Subscribe(topic_source_reg, byte(qos), func(client mqtt.Client, msg mqtt.Message) {
		fmt.Printf("topic_source_reg `%s` from `%s` topic\n", msg.Payload(), msg.Topic())
		sourceInfo := clients.Source{}
		//logger.Printf("Disconnect payload `%s`\n", payload)
		err := json.Unmarshal(msg.Payload(), &sourceInfo)
		//解析失败会报错，如json字符串格式不对，缺"号，缺}等。
		if err != nil {
			fmt.Println(err)
		}
		fmt.Printf("source info:%v\n", sourceInfo)
		m.SourceList.SList = append(m.SourceList.SList, sourceInfo)
		//m.OnSubscribeDisconnected(msg.Payload())
	})

	m.client.Subscribe(topic_source_query, byte(qos), func(client mqtt.Client, msg mqtt.Message) {
		fmt.Printf("topic_source_query `%s` from `%s` topic\n", msg.Payload(), msg.Topic())
		clientQuery := clients.ClientQuerySource{}

		err := json.Unmarshal(msg.Payload(), &clientQuery)
		//解析失败会报错，如json字符串格式不对，缺"号，缺}等。
		if err != nil {
			fmt.Println(err)
		}
		fmt.Printf("clientQuery%v\n", clientQuery)

		//把list的sourceId放到一个列表中
		// len :=len( m.SourceList.SList)
		// clientsIdSlice := make([]string, len)
		// for e := m.SourceList.SList.Front(); e != nil; e = e.Next() {
		// 	fmt.Print(e.Value)
		// 	clientsIdSlice = append(clientsIdSlice, e.Value.Id)
		// }
		//bytes, _ := json.Marshal(m.SourceList.SList.Front().Value)
		bytes, _ := json.Marshal(m.SourceList.SList)
		if token := m.client.Publish("sourceList/"+clientQuery.Id, byte(qos), false, bytes); token.Wait() && token.Error() != nil {
			fmt.Printf("publish failed, topic: %s, payload: %s\n", "sourceList/"+clientQuery.Id, bytes)
		}

	})

	// client.Subscribe(topic, byte(qos), func(client mqtt.Client, msg mqtt.Message) {
	// 	fmt.Printf("Received2 `%s` from `%s` topic\n", msg.Payload(), msg.Topic())
	// })
}
