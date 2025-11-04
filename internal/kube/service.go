package kube

import (
	"context"
	"fmt"
	"strings"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

type Service struct {
	cs kubernetes.Interface
}

func NewFromClientSet(cs kubernetes.Interface) *Service {
	return &Service{cs}
}

func NewFromConfig(cfg *rest.Config) (*Service, error) {
	cs, err := kubernetes.NewForConfig(cfg)
	if err != nil {
		return nil, fmt.Errorf("create clientset : %w", err)
	}

	return &Service{cs}, nil
}

type PodSummary struct {
	Name       string            `json:"name"`
	Namespace  string            `json:"namespace"`
	NodeName   string            `json:"nodename"`
	Phase      string            `json:"phase"`
	Labels     map[string]string `json:"labels"`
	Component  string            `json:"component"`
	Instance   string            `json:"instance"`
	Containers []string          `json:"containers"`
}

func (s *Service) ListPods(
	ctx context.Context,
	namespace string,
	baseSelector string,
	component string,
	instance string,
) ([]PodSummary, error) {
	var parts []string
	if baseSelector != "" {
		parts = append(parts, baseSelector)
	}

	if component != "" {
		parts = append(parts, fmt.Sprintf("app.kubernetes.io/component=%s", component))
	}

	if instance != "" {
		parts = append(parts, fmt.Sprintf("app.kubernetes.io/instance=%s", instance))
	}
	selector := strings.Join(parts, ",")

	pods, err := s.cs.CoreV1().Pods(namespace).List(ctx, metav1.ListOptions{LabelSelector: selector})
	if err != nil {
		return nil, fmt.Errorf("list pods: %w", err)
	}

	var out []PodSummary
	for _, p := range pods.Items {
		lbls := p.Labels
		if lbls == nil {
			lbls = map[string]string{}
		}
		var containers []string
		for _, c := range p.Spec.Containers {
			containers = append(containers, c.Name)
		}
		out = append(out, PodSummary{
			Name:       p.Name,
			Namespace:  p.Namespace,
			NodeName:   p.Spec.NodeName,
			Phase:      string(p.Status.Phase),
			Labels:     lbls,
			Component:  lbls["app.kubernetes.io/component"],
			Instance:   lbls["app.kubernetes.io/instance"],
			Containers: containers,
		})
	}
	return out, nil
}
