package k8helper

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"regexp"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

func k8ClientInit() *kubernetes.Clientset {
	var kubeconfig *string
	if home := homeDir(); home != "" {
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}
	flag.Parse()

	// use the current context in kubeconfig
	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		panic(err.Error())
	}

	// create the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	return clientset
}

func homeDir() string {
	if h := os.Getenv("HOME"); h != "" {
		return h
	}
	return os.Getenv("USERPROFILE") // windows
}

// K8GetNodesList Obtain node from current kubernetes context
func K8GetNodesList() ([]string, []string, []string) {

	nodeNameList := []string{}
	nodeIDList := []string{}
	combinedList := []string{}

	clientSet := k8ClientInit()

	nodes, err := clientSet.CoreV1().Nodes().List(metav1.ListOptions{})
	if err != nil {
		panic(err.Error())
	}
	for _, i := range nodes.Items {
		nodeIDList = append(nodeIDList, parseEc2ID(i.Spec.ProviderID))
		nodeNameList = append(nodeNameList, i.Name)

		combinedList = append(combinedList, fmt.Sprintf("%v - %v", i.Name, parseEc2ID(i.Spec.ProviderID)))
	}

	return nodeNameList, nodeIDList, combinedList
}

func parseEc2ID(nodeSpec string) string {
	s := regexp.MustCompile("/").Split(nodeSpec, -1)
	return s[4]
}
