// main.go
package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	admissionv1 "k8s.io/api/admission/v1"
	corev1 "k8s.io/api/core/v1"
)

func admitServices(w http.ResponseWriter, r *http.Request) {
	body, _ := ioutil.ReadAll(r.Body)

	var review admissionv1.AdmissionReview
	json.Unmarshal(body, &review)

	req := review.Request
	resp := admissionv1.AdmissionResponse{
		UID:     req.UID,
		Allowed: true,
	}

	if req.Kind.Kind == "Service" &&
		req.Namespace == "target-namespace" {

		var svc corev1.Service
		json.Unmarshal(req.Object.Raw, &svc)

		if svc.Spec.Type == corev1.ServiceTypeLoadBalancer {
			resp.Allowed = false
			resp.Result = &metav1.Status{
				Message: "LoadBalancer services are forbidden in this namespace",
			}
		}
	}

	review.Response = &resp
	respBytes, _ := json.Marshal(review)
	w.Header().Set("Content-Type", "application/json")
	w.Write(respBytes)
}

func main() {
	http.HandleFunc("/validate", admitServices)
	server := &http.Server{
		Addr:      ":8443",
		TLSConfig: getTLSConfig("/certs/tls.crt", "/certs/tls.key"),
	}
	server.ListenAndServeTLS("", "")
}

