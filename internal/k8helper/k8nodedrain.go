package k8helper

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
)

type patchStringValue struct {
	Op    string `json:"op"`
	Path  string `json:"path"`
	Value bool   `json:"value"`
}

// K8NodeDrain drains the node
func K8NodeDrain(nodeList []string) {
	for _, i := range nodeList {
		k8NodeCordon(i)
		k8DeleteNodePods(i)
		fmt.Println("Waiting before proceeding to next node")
		time.Sleep(30 * time.Second)
	}
}

func k8DeleteNodePods(nodeInstance string) {
	clientSet := k8ClientInit()

	pods, err := clientSet.CoreV1().Pods("").List(metav1.ListOptions{
		FieldSelector: "spec.nodeName=" + nodeInstance,
	})

	if err != nil {
		log.Fatal(err)
	}
	for _, i := range pods.Items {
		if i.Namespace == "kube-system" {
			fmt.Printf("Skipping Kube-System Pod - %v\n", i.Name)
			continue
		} else {
			fmt.Printf("Deleting pod: %v from %v\n", i.Name, nodeInstance)
			err := clientSet.CoreV1().Pods(i.Namespace).Delete(i.Name, &metav1.DeleteOptions{})
			if err != nil {
				log.Fatal(err)
			}
		}
	}
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

	fmt.Printf("\n%v marked as unschedulable\n", nodeInstance)

}
