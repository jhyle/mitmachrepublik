<a href="/veranstalter/{{organizerUrl .}}">
	<h3>{{.Addr.Name}}</h3>
	{{if .Image}}<p><img src="/bild/{{.Image}}?width=250" title="{{.Addr.Name}}" class="img-responsive" /></p>{{end}}
	<p>{{.Descr}}</p>
</a>
{{if .Web}}
	<p><a href="{{.Web}}" target="_blank" class="highlight"><span class="fa fa-caret-right"></span> {{.Web}}</a></p>
{{end}}
{{ if not .Addr.IsEmpty }}
	<div><iframe width="250" height="187" src="http://maps.google.de/maps?hl=de&q={{.Addr.Street}}%20{{.Addr.Pcode}}%20{{.Addr.City}}&ie=UTF8&t=&z=14&iwloc=B&output=embed" frameborder="0" scrolling="no" marginheight="0" marginwidth="0"></iframe></div>
{{end}}

