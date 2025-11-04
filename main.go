package main

import (
	"context"
	"fmt"

	"github.com/KanathipP/k8s-lab/kube" // ถ้า module คุณตั้งชื่อแบบนี้
	"k8s.io/client-go/kubernetes/fake"
)

func main() {
	ctx := context.Background()
	client := fake.NewSimpleClientset()

	svc := kube.New(client)

	cm, err := svc.CreateConfigMap(ctx, "default", "from-main", map[string]string{
		"key": "value",
	})
	if err != nil {
		panic(err)
	}

	fmt.Println("created:", cm.Name)
}
