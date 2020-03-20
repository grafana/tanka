package client

import (
	"os"
	"strings"

	"github.com/grafana/tanka/pkg/kubernetes/manifest"
)

// Order in which install different kinds of Kubernetes objects.
// Inspired by https://github.com/helm/helm/blob/8c84a0bc0376650bc3d7334eef0c46356c22fa36/pkg/releaseutil/kind_sorter.go
var kindOrder = []string{
	"Namespace",
	"NetworkPolicy",
	"ResourceQuota",
	"LimitRange",
	"PodSecurityPolicy",
	"PodDisruptionBudget",
	"ServiceAccount",
	"Secret",
	"ConfigMap",
	"StorageClass",
	"PersistentVolume",
	"PersistentVolumeClaim",
	"CustomResourceDefinition",
	"ClusterRole",
	"ClusterRoleList",
	"ClusterRoleBinding",
	"ClusterRoleBindingList",
	"Role",
	"RoleList",
	"RoleBinding",
	"RoleBindingList",
	"Service",
	"DaemonSet",
	"Pod",
	"ReplicationController",
	"ReplicaSet",
	"Deployment",
	"HorizontalPodAutoscaler",
	"StatefulSet",
	"Job",
	"CronJob",
	"Ingress",
	"APIService",
}

// Apply applies the given yaml to the cluster
func (k Kubectl) Apply(data manifest.List, opts ApplyOpts) error {
	// Manifests have already been pre-sorted during the Reconcile phase.
	return k.apply(data, opts)
}

func (k Kubectl) apply(data manifest.List, opts ApplyOpts) error {
	argv := []string{"-f", "-"}
	if opts.Force {
		argv = append(argv, "--force")
	}

	if !opts.Validate {
		argv = append(argv, "--validate=false")
	}

	cmd := k.ctl("apply", argv...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	cmd.Stdin = strings.NewReader(data.String())

	return cmd.Run()
}
