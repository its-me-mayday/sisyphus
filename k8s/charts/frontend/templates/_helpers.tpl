{{/*
Nome base dell'app
*/}}
{{- define "sisyphus-frontend.name" -}}
{{- .Chart.Name }}
{{- end }}

{{/*
Fullname — release + name
*/}}
{{- define "sisyphus-frontend.fullname" -}}
{{- printf "%s-%s" .Release.Name .Chart.Name | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Labels standard
*/}}
{{- define "sisyphus-frontend.labels" -}}
helm.sh/chart: {{ .Chart.Name }}-{{ .Chart.Version }}
app.kubernetes.io/name: {{ include "sisyphus-frontend.name" . }}
app.kubernetes.io/instance: {{ .Release.Name }}
app.kubernetes.io/version: {{ .Chart.AppVersion | quote }}
app.kubernetes.io/managed-by: {{ .Release.Service }}
{{- end }}

{{/*
Selector labels
*/}}
{{- define "sisyphus-frontend.selectorLabels" -}}
app.kubernetes.io/name: {{ include "sisyphus-frontend.name" . }}
app.kubernetes.io/instance: {{ .Release.Name }}
{{- end }}