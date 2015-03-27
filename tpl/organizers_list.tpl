{{$n := len .organizers}}
{{if eq $n 0}}
<div class="row-tile">
	<p class="big-text text-center" style="padding-top: 20px">Es wurden keine Veranstalter gefunden.</p>
</div>
{{else}}
{{range .organizers}}
<div class="row-tile">
	<a href="/veranstalter/{{organizerUrl .}}">
	{{if .Image}}
		<div class="small-icon"><span class="fa fa-{{if len .Categories}}{{with index .Categories 0}}{{categoryIcon .}}{{end}}{{end}} fa-fw"></span></div>
		<img class="pull-left" style="margin-right: 10px" src="/bild/{{.Image}}?width=220&height=165">
	{{end}}
	<div class="tile-text" {{if .Image}}style="margin-left: 230px"{{end}}>
		<h3>{{.Addr.Name}}</h3>
		<p>{{strClip .Descr 100}}</p>
		{{ if not .Addr.IsEmpty }}
			<p class="small-icon pull-left"><span class="fa fa-map-marker"></span></p>
			<p class="pull-left place">{{ if .Addr.Name }}{{.Addr.Name}}<br />{{ end }}{{ if .Addr.Street }}{{.Addr.Street}}, {{ end }}{{ if .Addr.Pcode }}{{.Addr.Pcode}} {{ end }}{{.Addr.City}}</p>
		{{ end }}
	</div>
	</a>
</div>
{{end}}
<div class="pages">
	<a href="./0#organizers"><div class="page">Anfang</div></a>
	{{if gt $.page 0}}<a href="./{{dec $.page}}#organizers">{{end}}<div class="page">&lt;</div>{{if gt $.page 0}}</a>{{end}}
	{{range $.pages}}
		{{if or (and (ge . (dec (dec $.page))) (le . (inc (inc $.page)))) (or (eq $.page .) (or (le . 1) (ge . (dec $.maxPage))))}}
			<a href="./{{.}}#organizers"><div class="page {{if eq . $.page}}cur-page{{end}}">{{inc .}}</div></a>
		{{else}}
			{{if or (eq . (dec (dec (dec $.page)))) (eq . (inc (inc (inc $.page))))}}
				<div class="page">..</div>
			{{end}}
		{{end}}
	{{end}}
	{{if lt $.page $.maxPage}}<a href="./{{inc $.page}}#organizers">{{end}}<div class="page">&gt;</div>{{if lt $.page $.maxPage}}</a>{{end}}
	<a href="./{{$.maxPage}}#organizers"><div class="page">Ende</div></a>
</div>
{{end}}
