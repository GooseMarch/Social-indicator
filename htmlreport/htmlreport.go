package htmlreport

import ("bytes"; "html/template"; "lab2_bubble/domain")

const html = `<!DOCTYPE html>
<html lang= "En">
	<head>
		<meta character="UTF-8">
		<title> Demographic analytic report -{{.Country}}</title>
			<style>
			body { font-family: "Times New Roman", Times, serif, "Segoe UI", Roboto, Arial, sans-serif; line-height: 1.6; color: #111010; max-width: 900px; margin: 2px, auto; padding: 20px; background: #826b67c1}
			h1,h2 { color: #554094; padding-bottom: 8px;}
			</style>
	</head>

	<body>
		<h1> Demographic analytic report </h1>

		<div class ="meta-info">
		<p><strong> Country: {{.Country}}</strong></p>
		<p><strong> Mode: {{.Mode}}</strong></p>
		<p><strong> Generator: {{.Generator}}</strong></p>
		<p><strong> Time of access to data: {{.Time}}</strong></p>
		<p><strong> Time of creation file: {{.Reporttime}}</strong></p>
		{{if .Seed}}<p><strong> Seed: {{.Seedvalue}}</strong></p>{{end}}
		</div>

		<h2> main parameters </h2>
		<div class = "parametr-info">
		<p><strong> Count: {{.Count}} </strong></p>
		<p><strong> Minimum: {{.Min}} {{.Whatusing}}</strong></p>
		<p><strong> Maximum: {{.Max}} {{.Whatusing}}</strong></p>
		<p><strong> Average: {{.Average}}{{.Whatusing}}</strong></p>
		</div>

		{{if .Seedresult}}
		<h2>Extra parameters</h2>
		<div class="parametr-info">
			<p><strong>Number of values above average: {{.Extra}}</strong></p>
		</div>
		{{end}}

		<h2> Output data </h2>
		<table>
			<thead>
			<tr><th>№</th><th>Value</th><th>Units</th></tr>
			</thead>
			<tbody>
				{{range .Values}}
				<tr>
					<td> {{.Pointer}} </td>
					<td> {{printf "%.1f" .Value}}</td>
					<td> {{$.Whatusing}}</td>
				</tr>
				{{end}}
			</tbody>
		</table>

		<footers>
			<p> Demographic analytic report. Version 1.0 </p>
		</footers>
	</body>
	</html>`

	func Generatehtml(result domain.Result)(string, error){
		tmpl, errar := template.New("report").Parse(html)
		if errar != nil {
			return "", errar
		}
		var buffer bytes.Buffer

		errar = tmpl.Execute(&buffer, result)
		if errar != nil{
			return "", errar
		}

		return buffer.String(), nil
	}