package services

import (
	"ConfigManager/kube"
	"fmt"
	"gopkg.in/yaml.v3"
	coreV1 "k8s.io/api/core/v1"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"log"
)

type ConfigMap struct {
	ApiVersion string            `yaml:"apiVersion"`
	Kind       string            `yaml:"kind"`
	Metadata   ObjectMeta        `yaml:"metadata"`
	Data       map[string]string `yaml:"data"`
}

type ObjectMeta struct {
	Name                       string            `yaml:"name,omitempty" protobuf:"bytes,1,opt,name=name"`
	GenerateName               string            `yaml:"generateName,omitempty" protobuf:"bytes,2,opt,name=generateName"`
	Namespace                  string            `yaml:"namespace,omitempty" protobuf:"bytes,3,opt,name=namespace"`
	UID                        types.UID         `yaml:"uid,omitempty" protobuf:"bytes,5,opt,name=uid,casttype=k8s.io/kubernetes/pkg/types.UID"`
	ResourceVersion            string            `yaml:"resourceVersion,omitempty" protobuf:"bytes,6,opt,name=resourceVersion"`
	Generation                 int64             `yaml:"generation,omitempty" protobuf:"varint,7,opt,name=generation"`
	DeletionGracePeriodSeconds *int64            `yaml:"deletionGracePeriodSeconds,omitempty" protobuf:"varint,10,opt,name=deletionGracePeriodSeconds"`
	Labels                     map[string]string `yaml:"labels,omitempty" protobuf:"bytes,11,rep,name=labels"`
	Annotations                map[string]string `yaml:"annotations,omitempty" protobuf:"bytes,12,rep,name=annotations"`
	Finalizers                 []string          `yaml:"finalizers,omitempty" patchStrategy:"merge" protobuf:"bytes,14,rep,name=finalizers"`
	ClusterName                string            `yaml:"clusterName,omitempty" protobuf:"bytes,15,opt,name=clusterName"`
}

var (
	dataKey   string
	dataValue string
)

func (c ConfigMap) GetConfigMap(configName string) (info coreV1.ConfigMap, err error) {

	configmapInfo, err := kube.KubeCliSet.CoreV1().ConfigMaps(c.Metadata.Namespace).Get(kube.KubeCtx, configName, metaV1.GetOptions{})
	if err != nil {
		return
	}

	return *configmapInfo, nil

}

func (c ConfigMap) CreateConfigMap(configMapData string) (info string, err error) {

	fmt.Println(configMapData)
	err = yaml.Unmarshal([]byte(configMapData), &c)
	if err != nil {
		log.Println("Yaml 配置解析失败 ->", err)
		return
	}

	//var (
	//	dataKey   string
	//	dataValue string
	//)

	for k, v := range c.Data {
		dataKey = k
		dataValue = v
	}

	// TODO: 判断namespce是否存在

	var configMap = &coreV1.ConfigMap{
		ObjectMeta: metaV1.ObjectMeta{
			Name:                       c.Metadata.Name,
			GenerateName:               c.Metadata.GenerateName,
			Namespace:                  c.Metadata.Namespace,
			UID:                        c.Metadata.UID,
			ResourceVersion:            c.Metadata.ResourceVersion,
			Generation:                 c.Metadata.Generation,
			DeletionGracePeriodSeconds: c.Metadata.DeletionGracePeriodSeconds,
			Labels:                     c.Metadata.Labels,
			Annotations:                c.Metadata.Annotations,
			Finalizers:                 c.Metadata.Finalizers,
			ClusterName:                c.Metadata.ClusterName,
		},
		Data: map[string]string{
			dataKey: dataValue,
		},
	}

	_, err = kube.KubeCliSet.CoreV1().ConfigMaps(c.Metadata.Namespace).Create(kube.KubeCtx, configMap, metaV1.CreateOptions{})
	if err != nil {
		return "创建失败", err
	}
	return "创建成功", nil

}

func (c ConfigMap) UpdateConfigMap(configMapData string) (info string, err error) {

	fmt.Println(configMapData)
	err = yaml.Unmarshal([]byte(configMapData), &c)
	if err != nil {
		log.Println("配置解析失败 ->", err)
		return
	}

	for k, v := range c.Data {
		dataKey = k
		dataValue = v
	}

	configmapInfo, err := c.GetConfigMap(c.Metadata.Name)
	if err != nil {
		return "ConfigMap 不存在", err
	}

	configmapInfo.Data = c.Data
	_, err = kube.KubeCliSet.CoreV1().ConfigMaps(c.Metadata.Namespace).Update(kube.KubeCtx, &configmapInfo, metaV1.UpdateOptions{})
	if err != nil {
		log.Println("configMap 更新失败 ->", err)
		return "", err
	}

	return "", nil

}

func (c ConfigMap) DeleteConfigMap(configMapName, configMapNamespace string) (err error) {

	err = kube.KubeCliSet.CoreV1().ConfigMaps(configMapNamespace).Delete(kube.KubeCtx, configMapName, metaV1.DeleteOptions{})
	if err != nil {
		log.Println(err)
		return err
	}
	return nil

}
