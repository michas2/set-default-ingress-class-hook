package main

import (
	"encoding/json"
	adm "k8s.io/api/admission/v1beta1"
	ext "k8s.io/api/extensions/v1beta1"
	mach "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/klog"
	"strings"
)

type Patch struct {
	Op    string      `json:"op"`
	Path  string      `json:"path"`
	Value interface{} `json:"value"`
}

const (
	key string = "kubernetes.io/ingress.class"
)

func patch( annotations map[string]string, name string) []byte {
	var patches []Patch
	if len(annotations) == 0 {
		klog.V(0).Infof("New annotation to %v.\n", name)
		patchNew := Patch{"add", "/metadata/annotations", map[string]string{}}
		patches = append(patches, patchNew)
	}
	if _, ok := annotations[key]; !ok {
		klog.V(0).Infof("Added class %v to %v.\n", Class, name)
		patchAdd := Patch{"add", "/metadata/annotations/" + strings.Replace(key, "/", "~1", -1), Class}
		patches = append(patches, patchAdd)
	} else {
		klog.V(0).Info("No need to add class.\n")
	}
	bytes, _ := json.Marshal(patches)
	return bytes
}

func addDefaultIngressClass(ar adm.AdmissionReview ) *adm.AdmissionResponse {
	ingressResource := mach.GroupVersionResource{Group: "extensions", Version: "v1beta1", Resource: "ingresses"}

	if ar.Request.Resource != ingressResource {
		klog.Errorf("expect resource to be %s", ingressResource)
		return nil
	}

	raw := ar.Request.Object.Raw
	ingress := ext.Ingress{}
	deserializer := codecs.UniversalDeserializer()
	if _, _, err := deserializer.Decode(raw, nil, &ingress); err != nil {
		klog.Error(err)
		return toAdmissionResponse(err)
	}

	reviewResponse := adm.AdmissionResponse{}
	reviewResponse.Allowed = true
	reviewResponse.Patch = patch(ingress.ObjectMeta.Annotations, ingress.ObjectMeta.Name)

	pt := adm.PatchTypeJSONPatch
	reviewResponse.PatchType = &pt

	return &reviewResponse
}
