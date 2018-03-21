package apidoc

const Tmpl_oapi2 = `swagger: "2.0"
info:
  title: HashiCorp Vault API
  description: |
    The Vault HTTP API gives you full access to Vault via HTTP. Every aspect of Vault can be controlled via this API. The Vault CLI uses the HTTP API to access Vault. You can read more about the Vault API at: https://www.vaultproject.io/api/index.html
  version: {{ .Version }}
  contact:
    name: HashiCorp Vault
    url: https://www.vaultproject.io
  license:
    name: Mozilla Public License 2.0
    url: https://www.mozilla.org/en-US/MPL/2.0

paths:{{ range .PathList }}{{ $path := . }}
  {{ .Pattern }}:{{ range $method, $el := .Methods }}
    {{ lower $method }}:{{ if $el.Summary }}
      summary: {{ $el.Summary }}{{ end }}{{ if $el.Description }}
      description: {{ $el.Description }}{{ end }}
      produces:
        - application/json
      tags:
        - {{ $path.Prefix }}
	  {{- if (or $el.PathFields $el.BodyFields) }}
      parameters:{{ range $el.PathFields }}
        - name: {{ .Name }}
          description: {{ .Description }}
          in: path
          type: {{ .Type }}
          required: true{{ end -}}
      {{ end -}}
      {{ if $el.BodyFields }}
        - name: Data
          in: body
          schema:
            type: object
            properties:{{ range $el.BodyFields }}
              {{ .Name }}:
                description: {{ .Description }}
                type: {{ .Type }}
                {{- if (eq .Type "array") }}
                items:
                  type: {{ .SubType }}
                {{- end }}
            {{- end }}
      {{-  end }}
      responses:
	  {{- range $el.Responses }}
        {{ .Code }}:
          description: {{ .Description }}
		  {{- if .Example }}
          examples:
            application/json:
              {{  .Example }}
		  {{- end }}
	  {{- end }}
    {{ end }}
{{- end }}
`
