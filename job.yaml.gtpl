{{ range $test := . -}}
---
apiVersion: v1
kind: Pod
metadata:
  name: "{{ $test.Name }}-test"
  annotations:
    helm.sh/hook: test-success
spec:
  serviceAccountName: testing-account
  containers:
    - name: test
      image: {{ $test.Image }}
      command:
{{- range $line := $test.Command }}
        - {{ . }}
{{- end }}
      resources:
        requests:
          memory: "1Gi"
          cpu: "500m"
        limits:
          memory: "1Gi"
          cpu: "500m"
      volumeMounts:
        - name: test-scripts
          mountPath: "/tests"
          readOnly: true
  volumes:
    - name: test-scripts
      secret:
        secretName: "helmfile-test-scripts"
  restartPolicy: Never
{{ end -}}
---
apiVersion: v1
kind: Secret
metadata:
  name: "helmfile-test-scripts"
type: Opaque
stringData:
{{`{{`}} (tpl (.Files.Glob "test-scripts/*").AsConfig . ) | indent 2 {{`}}`}}
---
apiVersion: v1
kind: ServiceAccount
metadata:
  name: testing-account
---
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRoleBinding
metadata:
  name: testing-cluster-admin
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: ClusterRole
  name: cluster-admin
subjects:
- kind: ServiceAccount
  name: testing-account
  namespace: {{`{{ .Release.Namespace }}`}}
