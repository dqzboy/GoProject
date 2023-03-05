package main

import (
    "context"
    "fmt"
    "log"
    "net/smtp"
    "os"
    "strconv"
    "strings"
    "time"

    "k8s.io/apimachinery/pkg/api/errors"
    metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
    "k8s.io/client-go/kubernetes"
    "k8s.io/client-go/rest"
    "k8s.io/client-go/tools/clientcmd"
)

const (
    podNamespace = "default" // 设置 Pod 所在的 Namespace
)

func main() {
    // 创建 Kubernetes 的客户端
    kubeconfig := os.Getenv("KUBECONFIG")
    config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
    if err != nil {
        log.Fatalf("Failed to build K8s config: %v", err)
    }
    clientset, err := kubernetes.NewForConfig(config)
    if err != nil {
        log.Fatalf("Failed to create K8s client: %v", err)
    }

    // 创建邮件通知的发送函数
    sendEmail := func(msg string) {
        from := "sender@example.com"
        password := "password"
        to := []string{"recipient@example.com"}

        auth := smtp.PlainAuth("", from, password, "smtp.example.com")

        err := smtp.SendMail("smtp.example.com:587", auth, from, to, []byte(msg))
        if err != nil {
            log.Printf("Failed to send email: %v", err)
        } else {
            log.Println("Email notification sent")
        }
    }

    // 定义发送邮件的时间间隔
    emailInterval := time.Hour

    // 无限循环检查 Pod 状态并发送邮件通知
    for {
        // 获取 Pod 列表
        pods, err := clientset.CoreV1().Pods(podNamespace).List(context.Background(), metav1.ListOptions{})
        if err != nil {
            log.Printf("Failed to list Pods: %v", err)
            continue
        }

        // 统计 Pod 数量和失败的 Pod 数量
        podCount := len(pods.Items)
        failedCount := 0
        var failedPods []string
        for _, pod := range pods.Items {
            if pod.Status.Phase != "Running" {
                failedCount++
                failedPods = append(failedPods, pod.Name)
            }
        }

        // 判断是否有失败的 Pod
        if failedCount > 0 {
            // 发送邮件通知
            subject := fmt.Sprintf("Alert: %d/%d Pods have failed", failedCount, podCount)
            body := fmt.Sprintf("Failed Pods: %s", strings.Join(failedPods, ", "))

            msg := fmt.Sprintf("From: %s\nSubject: %s\n\n%s", "sender@example.com", subject, body)
            sendEmail(msg)

            // 等待指定的时间间隔后再次检查
            time.Sleep(emailInterval)
            continue
        }

        // 如果没有失败的 Pod，则等待 1 分钟后再次检查
        time.Sleep(time.Minute)
    }
}
