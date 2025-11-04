package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/KanathipP/k8s-lab/internal/config"
	"github.com/KanathipP/k8s-lab/internal/kube"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

func main() {
	cfg := config.Load()

	restCfg, err := loadRestConfig()
	if err != nil {
		log.Fatalf("cannot create K8s config: %v", err)
	}

	kubeSvc, err := kube.NewFromConfig(restCfg)
	if err != nil {
		log.Fatalf("cannot create kube service: %v", err)
	}

	mux := http.NewServeMux()

	// /healthz
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusOK, map[string]any{"ok": true})
	})

	// /pods?component=&instance=
	mux.HandleFunc("/pods", func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
		defer cancel()

		component := r.URL.Query().Get("component")
		instance := r.URL.Query().Get("instance")

		pods, err := kubeSvc.ListPods(ctx, cfg.Namespace, cfg.LabelSelector, component, instance)
		if err != nil {
			log.Printf("ListPods error: %v", err)
			http.Error(w, "internal error", http.StatusInternalServerError)
			return
		}

		writeJSON(w, http.StatusOK, map[string]any{"items": pods})
	})

	addr := ":8080"
	log.Printf("listening on %s", addr)
	if err := http.ListenAndServe(addr, corsMiddleware(mux, cfg.CORSAllowOrigins)); err != nil {
		log.Fatal(err)
	}
}

// ---------- helpers ----------

func loadRestConfig() (*rest.Config, error) {
	// พยายาม in-cluster ก่อน
	if cfg, err := rest.InClusterConfig(); err == nil {
		return cfg, nil
	}
	// ไม่ได้ก็ลองใช้ kubeconfig บนเครื่อง dev
	kubeconfig := clientcmd.RecommendedHomeFile // ~/.kube/config
	if env := os.Getenv("KUBECONFIG"); env != "" {
		kubeconfig = env
	}
	return clientcmd.BuildConfigFromFlags("", kubeconfig)
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

// corsMiddleware แบบง่าย ๆ (ไม่ครบเท่า FastAPI แต่ใช้เล่นได้)
func corsMiddleware(next http.Handler, origins []string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")
		if origin != "" && (len(origins) == 0 || origins[0] == "*" || contains(origins, origin)) {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			w.Header().Set("Vary", "Origin")
		}
		if r.Method == http.MethodOptions {
			w.Header().Set("Access-Control-Allow-Methods", "GET,OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func contains(list []string, v string) bool {
	for _, x := range list {
		if x == v {
			return true
		}
	}
	return false
}
