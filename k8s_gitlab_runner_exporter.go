package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

type Runner struct {
	Name string
}

var (
	listenAddress = getEnv("LISTEN_PORT", ":9191")              // Порт для получения метрик
	namespace     = getEnv("NAMESPACE_RUNNER", "default")       // Неймспейс с gitlab раннером
	label         = getEnv("LABEL_RUNNER", "app=gitlab-runner") // Лейбл по которому определяется под раннера
)

func main() {
	http.HandleFunc("/runners", getRunners)
	http.HandleFunc("/metrics", getMetrics)
	err := http.ListenAndServe(listenAddress, nil)
	if err != nil {
		log.Fatal("ListenAndServer: ", err)
	}
}

func getEnv(key, fallback string) string {
	value := os.Getenv(key)
	if len(value) == 0 {
		return fallback
	}
	return value
}

func getRunners(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	config, err := rest.InClusterConfig()
	if err != nil {
		panic(err.Error())
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	pods, err := clientset.CoreV1().Pods(namespace).List(context.TODO(), metav1.ListOptions{LabelSelector: label})
	if err != nil {
		panic(err.Error())
	}
	var a []Runner
	for _, v := range pods.Items {
		a = append(a, Runner{Name: v.Name})
	}
	j, err := json.Marshal(a)
	if err != nil {
		fmt.Printf("Error: %s", err.Error())
	}
	fmt.Fprintln(w, string(j))
}

func getMetrics(w http.ResponseWriter, r *http.Request) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", "http://"+getIP(r.URL.Query().Get("runner"))+":9252/metrics", nil)
	if err != nil {
		fmt.Println("Error Collecting JSON from API: ", err)
	}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error Collecting JSON from API: ", err)
	}
	if resp.StatusCode != 200 {
		fmt.Println("Error Collecting JSON from API: ", resp.Status)
	}
	body, error := ioutil.ReadAll(resp.Body)
	if error != nil {
		fmt.Println(error)
	}
	fmt.Fprintln(w, string(body))
	resp.Body.Close()
}

func getIP(name string) string {
	var ip string
	config, err := rest.InClusterConfig()
	if err != nil {
		panic(err.Error())
	}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}
	pods, err := clientset.CoreV1().Pods(namespace).List(context.TODO(), metav1.ListOptions{LabelSelector: label})
	if err != nil {
		panic(err.Error())
	}

	for _, v := range pods.Items {
		if v.Name == name {
			ip = v.Status.PodIP
		}
	}
	return ip
}
