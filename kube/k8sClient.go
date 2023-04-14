package kube

import (
	"context"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"log"
	"os"
)

var (
	KubeCliSet   *kubernetes.Clientset
	DyKubeCliSet dynamic.Interface
	KubeCtx      context.Context
)

func init() {
	dir, _ := os.Getwd()
	config, err := clientcmd.BuildConfigFromFlags("", dir+"/kube/"+"kubeconfig")
	if err != nil {
		log.Fatal(err)
		return
	}

	clientSet, err := kubernetes.NewForConfig(config)
	if clientSet == nil {
		return
	}

	dynamicClient, err := dynamic.NewForConfig(config)
	if err != nil {
		log.Fatalf("Error creating Kubernetes dynamic client: %v", err)
	}

	KubeCliSet = clientSet
	DyKubeCliSet = dynamicClient
	KubeCtx = context.TODO()

}
