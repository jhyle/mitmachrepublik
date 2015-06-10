{{$n := len .organizers}}
{{if eq $n 0}}
<div class="row-tile">
	<p class="big-text text-center" style="padding-top: 20px">Es wurden keine Veranstalter gefunden.</p>
</div>
{{else}}
{{range .organizers}}
<div class="row-tile">
	<a href="{{.Url}}" title="Veranstaltungen von {{.Name}} ansehen">
	{{if .Image}}
		{{if len .Categories}}{{with index .Categories 0}}<div class="small-icon"><span class="fa fa-{{categoryIcon .}} fa-fw" title="{{categoryTitle .}}"></span></div>{{end}}{{end}}
		<img class="pull-left" style="margin-right: 10px" src="/bild/{{.Image}}?width=220&height=165" alt="Logo {{.Name}}">
	{{end}}
	<div class="tile-text">
		<h3>{{.Name}}</h3>
		<p>{{strClip .PlainDescription 100}}</p>
		{{ if not .Addr.IsEmpty }}
			<p class="small-icon pull-left"><span class="fa fa-map-marker fa-fw" title="Ort"></span></p>
			<p class="pull-left place">{{ if .Addr.Name }}{{.Addr.Name}}<br />{{ end }}{{ if .Addr.Street }}{{.Addr.Street}}, {{ end }}{{ if .Addr.Pcode }}{{.Addr.Pcode}} {{ end }}{{citypartName .Addr}}</p>
		{{ end }}
	</div>
	</a>
</div>
{{end}}
<div class="pages">
	<a href="./0#organizers" title="Zum Anfang der Liste"><div class="page">Anfang</div></a>
	{{if gt $.page 0}}<a href="./{{dec $.page}}#organizers" title="Vorherige Seite">{{end}}<div class="page">&lt;</div>{{if gt $.page 0}}</a>{{end}}
	{{range $.pages}}
		{{if or (and (ge . (dec (dec $.page))) (le . (inc (inc $.page)))) (or (eq $.page .) (or (le . 1) (ge . (dec $.maxPage))))}}
			<a href="./{{.}}#organizers" title="Zu Seite {{inc .}}"><div class="page {{if eq . $.page}}cur-page{{end}}">{{inc .}}</div></a>
		{{else}}
			{{if or (eq . (dec (dec (dec $.page)))) (eq . (inc (inc (inc $.page))))}}
				<div class="page">..</div>
			{{end}}
		{{end}}
	{{end}}
	{{if lt $.page $.maxPage}}<a href="./{{inc $.page}}#organizers" title="Nächste Seite">{{end}}<div class="page">&gt;</div>{{if lt $.page $.maxPage}}</a>{{end}}
	<a href="./{{$.maxPage}}#organizers" title="Zum Ende der Liste"><div class="page">Ende</div></a>
</div>
{{end}}
