package main

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"log"
	"strings"
)

func main() {
	data := `
apiVersion: v1
kind: ConfigMap
metadata:
  name: test
  namespace: test
data:
  micro.config.prod.yaml: |-
    service:
      # 消息服务hedwig
      message_microservice:
        name: MESSAGE_MICROSERVICE
        package: message
        version: 1
        ip: js-design-staging-hedwig
        port: 50002

      # 大数据服务
      dataCenter_microservice:
        name: DATACENTER_MICROSERVICE
        package: dataCenter
        version: 1
        ip: js-design-staging-hermes
        port: 50056

      # 用户服务wood
      user_microservice:
        name: USER_MICROSERVICE
        package: user
        version: 1
        ip: js-design-staging-wood
        port: 50054
    `
	var yamlData interface{}
	if err := yaml.Unmarshal([]byte(data), &yamlData); err != nil {
		log.Println("1 ->", err)
	}
	yamlBytes, err := yaml.Marshal(yamlData)
	if err != nil {
		log.Println("2 ->", err)
		return
	}
	// 将 \n 转义成 \\n
	yamlStr := string(yamlBytes)
	yamlStr = strings.ReplaceAll(yamlStr, "\n", "\\n")
	fmt.Println(yamlStr)
}
