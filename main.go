package main

import (
	"encoding/json"
	// "flag"
	// "fmt"
	//"net"
	"net/http"
	// _ "net/http/pprof"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"github.com/unrolled/secure"
	mqtt "github.com/yaxiongwu/webrtc-mqtt-server/pkg/mqtt_client"
)

const CONFIG_FILE = "config.toml"

type Config struct {
	logV int //`mapstructure:"v"`

}

var (
	conf       = Config{}
	mqttClient *mqtt.MqttClient
	logger     *log.Logger
)

type disconnectInfo struct {
	Username string `json:"username"` //必须为大写开头才能解析！
	// Ts              int    `json:"ts"`
	// Sockport        int    `json:"sockport"`
	Reason string `json:"reason"`
	// Protocol        string `json:"protocol"`
	// Proto_ver       int8   `json:"proto_ver"`
	// Proto_name      string `json:"proto_name"`
	Ipaddress string `json:"ipaddress"`
	// Disconnected_at int    `json:"disconnected_at"`
	// Donnected_at    int    `json:"connected_at"`
	Clientid string `json:"clientid"`
}
type connectInfo struct {
	Username string `json:"username"` //必须为大写开头才能解析！
	// Ts              int    `json:"ts"`
	// Sockport        int    `json:"sockport"`
	Reason string `json:"reason"`

	// Protocol        string `json:"protocol"`
	// Proto_ver       int8   `json:"proto_ver"`
	// Proto_name      string `json:"proto_name"`
	Ipaddress string `json:"ipaddress"`
	//Keepalive       int `json:"keepalive"`
	//Expiry_interval int `json:"expiry_interval"`

	// Disconnected_at int    `json:"disconnected_at"`
	// Donnected_at    int    `json:"connected_at"`
	Clientid string `json:"clientid"`
	//Clean_start bool   `json:"clean_start"`
}

func loadConfig() bool {
	_, err := os.Stat(CONFIG_FILE)
	if err != nil {
		return false
	}

	viper.SetConfigFile(CONFIG_FILE)
	viper.SetConfigType("toml")

	err = viper.ReadInConfig()
	if err != nil {
		logger.Printf("config file read failed,%v", err)
		return false
	}
	err = viper.GetViper().Unmarshal(&conf)
	if err != nil {
		logger.Printf("sfu config file loaded failed,%v", err)
		return false
	}
	conf.logV = viper.GetInt("log.v")
	logger.Printf("Config file loaded,%v", conf)
	return true
}

func ginRun() {
	r := gin.Default()

	r.GET("/sourcesList", func(c *gin.Context) {
		logger.Printf("c.Query(sourceType),%v ", c.Query("sourceType"))
		//list, _ := json.Marshal(getSourceList(s, rtc.SourceType(rtc.SourceType_value[c.Query("sourceType")])))
		c.JSON(200, gin.H{
			//"list": string(list), //[]byte会自动转换成base64传输
			"message": "pong",
		})
	})

	r.LoadHTMLGlob("web/*")

	r.GET("/vedio", func(c *gin.Context) {
		logger.Printf("c.Query(id):%v", c.Query("id"))
		c.HTML(http.StatusOK, "index.html", gin.H{
			"id": c.Query("id"),
		})
	})

	r.GET("/mqtt_source", func(c *gin.Context) {
		c.HTML(http.StatusOK, "mqtt_source.html", gin.H{})
	})

	r.GET("/mqtt_client", func(c *gin.Context) {
		c.HTML(http.StatusOK, "mqtt_client.html", gin.H{})
	})

	r.GET("/index/:destination", func(context *gin.Context) {
		context.HTML(http.StatusOK, "index_phone.html", gin.H{
			"destination": context.Param("destination"),
		})
	})
	r.GET("/index_pc/:destination", func(context *gin.Context) {
		context.HTML(http.StatusOK, "index_pc.html", gin.H{
			"destination": context.Param("destination"),
		})
	})
	r.Use(LoadTls())
	r.RunTLS(":8080", "certs/bxzryd.pem", "certs/bxzryd.key")
	//r.Run(":8080")
}

func LoadTls() gin.HandlerFunc {
	return func(c *gin.Context) {
		middleware := secure.New(secure.Options{
			SSLRedirect: true,
			SSLHost:     "localhost:8080",
		})
		err := middleware.Process(c.Writer, c.Request)
		if err != nil {
			logger.Printf("error:%v", err)
			return
		}
		c.Next()
	}
}

func onSubscribeConnect(payload []byte) {

	//logger.Printf("Connect Received `%s`\n", payload)
	info := connectInfo{}

	err := json.Unmarshal(payload, &info)
	//解析失败会报错，如json字符串格式不对，缺"号，缺}等。
	if err != nil {
		logger.Println(err)
	}

	logger.Printf("connect Clientid `%v`\n", info.Clientid)
	// switch expression {
	// case condition:

	// }

}

func onSubscribeDisconnect(payload []byte) {
	info := disconnectInfo{}
	//logger.Printf("Disconnect payload `%s`\n", payload)

	err := json.Unmarshal(payload, &info)
	//解析失败会报错，如json字符串格式不对，缺"号，缺}等。
	if err != nil {
		logger.Println(err)
	}

	logger.Printf("Disconnect Clientid `%v`,reason: %v\n", info.Clientid, info.Reason)
	// switch expression {
	// case condition:

	// }

}

func main() {
	//	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
	logger = log.New(os.Stderr, "", log.Ldate|log.Ltime|log.Lshortfile)
	loadConfig()
	mqttClient = mqtt.MqttClientInit()

	mqttClient.OnSubscribeConnected = onSubscribeConnect
	mqttClient.OnSubscribeDisconnected = onSubscribeDisconnect
	go mqttClient.Subscribe()
	//go mqttClient.Subscribe("$SYS/#", 1)

	mqttClient.Publish("wuyaxiong", "paload", 1)

	logger.Println("MqttClientInit ok")
	ginRun()
}
