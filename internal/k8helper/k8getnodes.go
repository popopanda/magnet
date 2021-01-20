package k8helper

import (
	"fmt"
	"log"
	"regexp"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

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

		combinedList = append(combinedList, fmt.Sprintf("%v:%v - %v", i.Name, i.Status.NodeInfo.KubeletVersion, parseEc2ID(i.Spec.ProviderID)))
	}

	return nodeNameList, nodeIDList, combinedList
}

func parseEc2ID(nodeSpec string) string {
	s := regexp.MustCompile("/").Split(nodeSpec, -1)
	return s[4]
}

// K8GetNodesListByGroups Test
func K8GetNodesListByGroups(nodegroup string) ([]string, []string) {

	groupNodeNameList := []string{}
	groupNodeIDList := []string{}

	clientSet := k8ClientInit()

	nodes, err := clientSet.CoreV1().Nodes().List(metav1.ListOptions{})
	if err != nil {
		panic(err.Error())
	}

	for _, i := range nodes.Items {
		for _, label := range i.Labels {
			if label == nodegroup {
				groupNodeIDList = append(groupNodeIDList, parseEc2ID(i.Spec.ProviderID))
				groupNodeNameList = append(groupNodeNameList, i.Name)
			}
		}
	}

	if len(groupNodeNameList) <= 0 {
		log.Fatal("Error: No nodes found in ", nodegroup)
	}

	return groupNodeNameList, groupNodeIDList
}
