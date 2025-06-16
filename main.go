package main

import (
    "encoding/json"
    "net/http"

    admissionv1 "k8s.io/api/admission/v1"
    corev1       "k8s.io/api/core/v1"
    metav1       "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func validateService(w http.ResponseWriter, r *http.Request) {
    var review admissionv1.AdmissionReview
    if err := json.NewDecoder(r.Body).Decode(&review); err != nil {
        http.Error(w, "invalid request", http.StatusBadRequest)
        return
    }

    req := review.Request
    resp := admissionv1.AdmissionResponse{
        UID:     req.UID,
        Allowed: true,
    }

    if req.Kind.Kind == "Service" && req.Namespace == "dor" {
        var svc corev1.Service
        _ = json.Unmarshal(req.Object.Raw, &svc)
        if svc.Spec.Type == corev1.ServiceTypeLoadBalancer {
            resp.Allowed = false
            resp.Result = &metav1.Status{
                Message: "LoadBalancer services are forbidden in this namespace",
            }
        }
    }

    review.Response  = &resp
    review.APIVersion = "admission.k8s.io/v1"
    review.Kind       = "AdmissionReview"

    w.Header().Set("Content-Type", "application/json")
    _ = json.NewEncoder(w).Encode(review)
}

func main() {
    http.HandleFunc("/validate", validateService)
    if err := http.ListenAndServeTLS(
        ":8443",
        "/certs/tls.crt",
        "/certs/tls.key",
        nil,
    ); err != nil {
        panic(err)
    }
}

