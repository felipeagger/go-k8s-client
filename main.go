package main

import (
	"fmt"
	"log"

	batchv1 "k8s.io/api/batch/v1"
	core "k8s.io/api/core/v1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

const namespace = "default"

func main() {
	
	//From kube config file use
	//clientset, err := getClient("./kube_config.yml")
	//if err != nil {
	//	panic(err)
	//}
	//
	
	//From load the in-cluster config
	config, err := rest.InClusterConfig()
	if err != nil {
		panic(err.Error())
	}
	
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	getNodes(clientset)
	getPods(clientset)

	createPod(clientset, "pod-created-by-go")

	createJob(clientset, "job-created-by-go")

	deletePod(clientset, "pod-created-by-go")

}

func getClient(kubeconfig string) (*kubernetes.Clientset, error) {

	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		return nil, err
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	return clientset, nil
}

func createJob(clientset *kubernetes.Clientset, jobName string) bool {

	jobsClient := clientset.BatchV1().Jobs(namespace)
	newJob := getJobSpec(jobName)

	_, err := jobsClient.Create(newJob)
	if err != nil {
		log.Fatalln("\nFailed to create K8s job.")
		fmt.Println(err)
		return false
	}

	fmt.Println("\nCreated job with Success!")
	return true
}

func createPod(clientset *kubernetes.Clientset, PodName string) bool {

	newPod := getPodSpec(PodName)

	newPod, err := clientset.CoreV1().Pods(newPod.Namespace).Create(newPod)
	if err != nil {
		fmt.Println(err)
		return false
	}

	return true
}

func getJobSpec(jobName string) *batchv1.Job {
	var backOffLimit int32 = 0

	command := []string{"/bin/sh", "-c", "apk add go && go version"}

	return &batchv1.Job{
		ObjectMeta: metav1.ObjectMeta{
			Name:      jobName,
			Namespace: namespace,
		},
		Spec: batchv1.JobSpec{
			Template: v1.PodTemplateSpec{
				Spec: v1.PodSpec{
					Containers: []v1.Container{
						{
							Name:    jobName,
							Image:   "alpine",
							Command: command,
						},
					},
					RestartPolicy: v1.RestartPolicyNever,
				},
			},
			BackoffLimit: &backOffLimit,
		},
	}

}

func getPodSpec(PodName string) *core.Pod {
	return &core.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      PodName,
			Namespace: namespace,
			Labels: map[string]string{
				"app": "demo",
			},
		},
		Spec: core.PodSpec{
			Containers: []core.Container{
				{
					Name:            "alpine",
					Image:           "alpine",
					ImagePullPolicy: core.PullIfNotPresent,
					Command: []string{
						"sleep",
						"3600",
					},
				},
			},
		},
	}
}

func deletePod(clientset *kubernetes.Clientset, podname string) bool {

	err := clientset.CoreV1().Pods(namespace).Delete(podname, &metav1.DeleteOptions{})
	if err != nil {
		fmt.Println(err)
		return false
	}

	return true
}

func getPods(clientset *kubernetes.Clientset) {

	pods, err := clientset.CoreV1().Pods(namespace).List(metav1.ListOptions{})
	if err != nil {
		panic(err.Error())
	}

	fmt.Println("Listing Pods of Namespace " + namespace)

	for i := range pods.Items {
		fmt.Printf("\n PodName: %s - Namespace: %s - PodIP: %s - Status: %s",
			pods.Items[i].ObjectMeta.Name,
			pods.Items[i].ObjectMeta.Namespace,
			pods.Items[i].Status.PodIP,
			pods.Items[i].Status.Phase)
	}

}

func getNodes(clientset *kubernetes.Clientset) {

	nodes, err := clientset.CoreV1().Nodes().List(metav1.ListOptions{})
	if err != nil {
		panic(err.Error())
	}

	fmt.Println("\nListing Nodes of Cluster")

	for i := range nodes.Items {
		fmt.Printf("\n NodeName: %s - PodCIDR: %s - KubeletVersion: %s\nOSImage: %s - ContainerRuntime: %s - Architecture: %s\n\n",
			nodes.Items[i].ObjectMeta.Name,
			nodes.Items[i].Spec.PodCIDR,
			nodes.Items[i].Status.NodeInfo.KubeletVersion,
			nodes.Items[i].Status.NodeInfo.OSImage,
			nodes.Items[i].Status.NodeInfo.ContainerRuntimeVersion,
			nodes.Items[i].Status.NodeInfo.Architecture)
	}

}
