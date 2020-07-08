package k8helper

import (
	"encoding/json"
	"fmt"
	"log"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes"
)

type patchStringValue struct {
	Op    string `json:"op"`
	Path  string `json:"path"`
	Value bool   `json:"value"`
}

// K8NodeDrain drains the node
func K8NodeDrain(nodeList []string) {
	for _, i := range nodeList {
		// k8NodeCordon(i)
		k8NodeEvictPods(i)
	}
}

func k8NodeEvictPods(nodeInstance string) {
	clientSet := k8ClientInit()

	k8GetNodePods(nodeInstance, clientSet)
}

func k8GetNodePods(nodeInstance string, client *kubernetes.Clientset) {

	pods, err := client.CoreV1().Pods("").List(metav1.ListOptions{
		FieldSelector: "spec.nodeName=" + nodeInstance,
	})

	if err != nil {
		log.Fatal(err)
	}
	for _, i := range pods.Items {
		if i.Namespace == "kube-system" {
			continue
		} else {
			//evict
		}
	}
	client.RESTClient().Post().

}

func k8NodeCordon(nodeInstance string) {
	clientSet := k8ClientInit()

	payload := []patchStringValue{{
		Op:    "replace",
		Path:  "/spec/unschedulable",
		Value: true,
	}}
	payloadBytes, _ := json.Marshal(payload)

	_, err := clientSet.
		CoreV1().
		Nodes().
		Patch(nodeInstance, types.JSONPatchType, payloadBytes)

	if err != nil {
		panic(err.Error())
	}

	fmt.Printf("%v marked as unschedulable\n", nodeInstance)

}
