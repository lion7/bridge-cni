package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"

	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

func main() {
	// check args
	args := os.Args
	if len(args) != 2 {
		fmt.Println("Usage: bridge-cni [output-conf-path]")
		os.Exit(1)
	}

	// lookup the node name
	nodeName, found := os.LookupEnv("NODE_NAME")
	if !found {
		fmt.Printf("Env var NODE_NAME is not defined\n")
		os.Exit(1)
	}

	// creates the in-cluster config
	config, err := rest.InClusterConfig()
	if err != nil {
		log.Panic(err.Error())
	}
	// creates the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Panic(err.Error())
	}

	// get the node object
	node, err := clientset.CoreV1().Nodes().Get(context.TODO(), nodeName, metav1.GetOptions{})
	if errors.IsNotFound(err) {
		fmt.Printf("Node %s not found\n", nodeName)
	} else if statusError, isStatus := err.(*errors.StatusError); isStatus {
		fmt.Printf("Error getting node %s: %v\n", nodeName, statusError.ErrStatus.Message)
		os.Exit(1)
	} else if err != nil {
		panic(err.Error())
	}

	// extract the pod CIDR
	podCidr := node.Spec.PodCIDR
	if len(podCidr) == 0 {
		log.Panicf("PodCIDR is empty for node %s", nodeName)
	}

	// generate the CNI config
	conf := NetConfList{
		CNIVersion: "1.0.0",
		Name:       "cbr0",
		Plugins: []*PluginConf{
			{
				Type:             "bridge",
				IsDefaultGateway: true,
				IPAM: IPAM{
					Type:   "host-local",
					Subnet: podCidr,
				},
			},
		},
	}
	confJson, err := json.MarshalIndent(conf, "", "  ")
	if err != nil {
		panic(err.Error())
	}

	// open a writer for the desired output path
	outputPath := args[1]
	var output io.Writer
	if outputPath == "-" {
		output = os.Stdout
	} else {
		var err error
		output, err = os.OpenFile(outputPath, os.O_RDWR|os.O_CREATE, 0o644)
		if err != nil {
			panic(err.Error())
		}
	}

	// write the json to the desired output
	_, err = output.Write(confJson)
	if err != nil {
		panic(err.Error())
	}
}
