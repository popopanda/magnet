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
	k8NodeCordon(nodeList)
}

// func k8NodeEvictPods() {

// }

func k8NodeCordon(nodeList []string) {
	clientSet := k8ClientInit()

	payload := []patchStringValue{{
		Op:    "replace",
		Path:  "/spec/unschedulable",
		Value: true,
	}}
	payloadBytes, _ := json.Marshal(payload)

	for _, i := range nodeList {
		_, err := clientSet.CoreV1().Nodes().Patch(i, types.JSONPatchType, payloadBytes)

		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf("%v marked as unschedulable", i)
	}

}
