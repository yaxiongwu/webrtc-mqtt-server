package mqtt_client

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

//const protocol = "ssl"
const protocol = "tcp"

//const broker = "k101ceee.ala.cn-hangzhou.emqxsl.cn"
const broker = "127.0.0.1"

//const broker = "broker.emqx.io"
const port = 1883
const topic_connected = "$SYS/brokers/emqx@127.0.0.1/clients/+/connected"
const topic_disconnected = "$SYS/brokers/emqx@127.0.0.1/clients/+/disconnected"
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
		client: client,
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
	// client.Subscribe(topic, byte(qos), func(client mqtt.Client, msg mqtt.Message) {
	// 	fmt.Printf("Received2 `%s` from `%s` topic\n", msg.Payload(), msg.Topic())
	// })
}
