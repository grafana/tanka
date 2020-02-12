package main

import (
	"io/ioutil"
	"strings"
	"testing"

	"github.com/fatih/color"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	red = color.New(color.FgRed).SprintFunc()
	grn = color.New(color.FgGreen).SprintFunc()
	mgB = color.New(color.FgMagenta, color.Bold).SprintFunc()
	blB = color.New(color.FgBlue, color.Bold).SprintFunc()
	non = color.New().SprintFunc()
)

func TestColordiff(t *testing.T) {
	data := `diff -u -N /tmp/LIVE-155518783/apps.v1.Deployment.default.grafana /tmp/MERGED-942280082/apps.v1.Deployment.default.grafana
--- /tmp/LIVE-155518783/apps.v1.Deployment.default.grafana      2020-02-12 20:56:42.472375264 +0100
+++ /tmp/MERGED-942280082/apps.v1.Deployment.default.grafana    2020-02-12 20:56:42.475708598 +0100
@@ -6,7 +6,7 @@
     kubectl.kubernetes.io/last-applied-configuration: |
       {"apiVersion":"apps/v1","kind":"Deployment","metadata":{"annotations":{},"name":"grafana","namespace":"default"},"spec":{"minReadySeconds":10,"replicas":1,"revisionHistoryLimit":10,"selector":{"matchLabels":{"name":"grafana"}},"template":{"metadata":{"labels":{"name":"grafana"}},"spec":{"containers":[{"image":"grafana/grafana","imagePullPolicy":"IfNotPresent","name":"grafana"}]}}}}
   creationTimestamp: "2020-02-11T21:31:40Z"
-  generation: 1
+  generation: 2
   name: grafana
   namespace: default
   resourceVersion: "3041"
@@ -15,7 +15,7 @@
 spec:
   minReadySeconds: 10
   progressDeadlineSeconds: 600
-  replicas: 1
+  replicas: 2
   revisionHistoryLimit: 10
   selector:
     matchLabels:
@@ -32,9 +32,9 @@
         name: grafana
     spec:
       containers:
-      - image: grafana/grafana
+      - image: grfana/grafana
         imagePullPolicy: IfNotPresent
-        name: grafana
+        name: grafna
         resources: {}`

	want := strings.Join([]string{
		blB(`diff -u -N /tmp/LIVE-155518783/apps.v1.Deployment.default.grafana /tmp/MERGED-942280082/apps.v1.Deployment.default.grafana`),
		red(`--- /tmp/LIVE-155518783/apps.v1.Deployment.default.grafana      2020-02-12 20:56:42.472375264 +0100`),
		grn(`+++ /tmp/MERGED-942280082/apps.v1.Deployment.default.grafana    2020-02-12 20:56:42.475708598 +0100`),
		mgB(`@@ -6,7 +6,7 @@`),
		non(`     kubectl.kubernetes.io/last-applied-configuration: |`),
		non(`       {"apiVersion":"apps/v1","kind":"Deployment","metadata":{"annotations":{},"name":"grafana","namespace":"default"},"spec":{"minReadySeconds":10,"replicas":1,"revisionHistoryLimit":10,"selector":{"matchLabels":{"name":"grafana"}},"template":{"metadata":{"labels":{"name":"grafana"}},"spec":{"containers":[{"image":"grafana/grafana","imagePullPolicy":"IfNotPresent","name":"grafana"}]}}}}`),
		non(`   creationTimestamp: "2020-02-11T21:31:40Z"`),
		red(`-  generation: 1`),
		grn(`+  generation: 2`),
		non(`   name: grafana`),
		non(`   namespace: default`),
		non(`   resourceVersion: "3041"`),
		mgB(`@@ -15,7 +15,7 @@`),
		non(` spec:`),
		non(`   minReadySeconds: 10`),
		non(`   progressDeadlineSeconds: 600`),
		red(`-  replicas: 1`),
		grn(`+  replicas: 2`),
		non(`   revisionHistoryLimit: 10`),
		non(`   selector:`),
		non(`     matchLabels:`),
		mgB(`@@ -32,9 +32,9 @@`),
		non(`         name: grafana`),
		non(`     spec:`),
		non(`       containers:`),
		red(`-      - image: grafana/grafana`),
		grn(`+      - image: grfana/grafana`),
		non(`         imagePullPolicy: IfNotPresent`),
		red(`-        name: grafana`),
		grn(`+        name: grafna`),
		non(`         resources: {}`),
		"", // newline
	}, "\n")

	r := colordiff(data)
	got, err := ioutil.ReadAll(r)
	require.NoError(t, err)

	assert.Equal(t, want, string(got))
}
