package kube

import (
	"context"
	"testing"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/fake"
)

func TestListPods_FilterByComponentAndInstance(t *testing.T) {
	ctx := context.Background()

	client := fake.NewSimpleClientset(
		&corev1.Pod{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "server-0",
				Namespace: "flwr",
				Labels: map[string]string{
					"app.kubernetes.io/component": "serverapp",
					"app.kubernetes.io/instance":  "0",
				},
			},
			Spec: corev1.PodSpec{
				Containers: []corev1.Container{{Name: "main"}},
			},
			Status: corev1.PodStatus{Phase: corev1.PodRunning},
		},
		&corev1.Pod{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "client-0",
				Namespace: "flwr",
				Labels: map[string]string{
					"app.kubernetes.io/component": "clientapp",
					"app.kubernetes.io/instance":  "0",
				},
			},
			Spec: corev1.PodSpec{
				Containers: []corev1.Container{{Name: "main"}},
			},
			Status: corev1.PodStatus{Phase: corev1.PodRunning},
		},
	)

	svc := NewFromClientSet(client)

	pods, err := svc.ListPods(ctx, "flwr", "", "clientapp", "0")
	if err != nil {
		t.Fatalf("ListPods error %v", err)
	}
	if len(pods) != 1 {
		t.Fatalf("expected 1 pod, got %d", len(pods))
	}
	if pods[0].Name != "client-0" {
		t.Fatalf("expected client-0, got %s", pods[0].Name)
	}
	if pods[0].Component != "clientapp" || pods[0].Instance != "0" {
		t.Fatalf("unexpected labels: %+v", pods[0])
	}
}
