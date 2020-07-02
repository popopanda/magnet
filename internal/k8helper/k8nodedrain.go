package k8helper

func K8NodeDrain(nodeList []string) {
	clientSet := k8ClientInit()
	test, err := clientSet.CoreV1().Nodes().Delete()
}
