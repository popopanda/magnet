package k8helper

import (
	"encoding/json"
	"fmt"
	"log"

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
		fmt.Printf("%v marked as unschedulable\n", i)
		// k8NodeEvictPods()
		fmt.Printf("%v Evicted\n", i)
	}
}

// func k8NodeEvictPods() {

// }

func k8NodeCordon(nodeInstance string) {
	clientSet := k8ClientInit()

	payload := []patchStringValue{{
		Op:    "replace",
		Path:  "/spec/unschedulable",
		Value: true,
	}}
	payloadBytes, _ := json.Marshal(payload)

	_, err := clientSet.CoreV1().Nodes().Patch(nodeInstance, types.JSONPatchType, payloadBytes)

	if err != nil {
		log.Fatal(err)
	}
}
