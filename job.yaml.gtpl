{{ range $test := . -}}
---
apiVersion: v1
kind: Pod
metadata:
  name: "{{`{{ .Release.Name }}`}}-{{ .TestName }}-test"
  annotations:
    helm.sh/hook: test-success
spec:
  containers:
    - name: test
      image: {{ .Image }}
      command:
{{- range $line := .Command }}
        - {{ . }}
{{- end }}
    volumeMounts:
    - name: test-scripts
      mountPath: "/tests"
      readOnly: true
  volumes:
  - name: test-scripts
    secret:
      secretName: "{{`{{ .Release.Name }}`}}-test-scripts"
{{ end -}}
---
apiVersion: v1
kind: Secret
metadata:
  name: "{{`{{ .Release.Name }}`}}-test-scripts"
type: Opaque
stringData:
{{`{{`}} (tpl (.Files.Glob "test-scripts/*").AsConfig . ) | indent 2 {{`}}`}}
