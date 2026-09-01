{{- define "mychart.fullname" -}}
{{ .Release.Name }}-{{ .Chart.Name }}
{{- end -}}

{{- define "mychart.labels" -}}
app: {{ .Release.Name }}
chart: {{ .Chart.Name }}-{{ .Chart.Version }}
{{- end -}}
