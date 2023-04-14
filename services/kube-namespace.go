package services

import (
	"ConfigManager/kube"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"log"
)

func ListNamespace() (info interface{}, err error) {

	resp, err := kube.KubeCliSet.CoreV1().Namespaces().List(kube.KubeCtx, metaV1.ListOptions{})
	if err != nil {
		log.Println("获取 NameSpace 失败")
		return nil, err
	}

	var nsList []string
	for _, v := range resp.Items {
		nsList = append(nsList, v.Name)
	}

	return nsList, nil

}
