{{$n := len .events}}
{{if eq $n 0}}
<div class="row-tile">
	<p class="big-text text-center" style="padding-top: 20px">Es wurden keine Veranstaltungen gefunden.</p>
</div>
{{else}}
{{range .events}}
<div class="row-tile">
	{{if not $.user }}
		<a href="/veranstaltung/{{eventUrl .}}">
	{{end}}
	{{if .Image}}
		{{if len .Categories}}{{with index .Categories 0}}<div class="small-icon"><span class="fa fa-{{categoryIcon .}} fa-fw" title="{{categoryTitle .}}"></span></div>{{end}}{{end}}
		<img class="pull-left" style="margin-right: 10px" src="/bild/{{.Image}}?width=220&height=165" title="{{.Title}}">
	{{end}}
	<div class="tile-text">
		{{ if $.user }}
			<p class="pull-right"><a href="#" name="delete-event" title="Löschen" data-target="{{.Id.Hex}}" class="close"><span class="fa fa-times"></span></a></p>
		{{ end }}
		<h3>{{.Title}}</h3>
		<p class="datetime">{{datetimeFormat .Start}} Uhr {{if $.organizerNames}} - {{index $.organizerNames .OrganizerId}}{{end}}</p>
		<p>{{strClip .Descr 100}}</p>
		{{ if not .Addr.IsEmpty }}
			<p class="small-icon pull-left"><span class="fa fa-map-marker fa-fw" title="Ort"></span></p>
			<p class="pull-left place">{{ if .Addr.Name }}{{.Addr.Name}}<br />{{ end }}{{ if .Addr.Street }}{{.Addr.Street}}, {{ end }}{{ if .Addr.Pcode }}{{.Addr.Pcode}} {{ end }}{{citypartName .Addr}}</p>
		{{ end }}
		{{ if $.user }}
			<p class="pull-right"><a href="/veranstalter/verwaltung/veranstaltung/{{.Id.Hex}}" class="btn btn-mmr" style="margin: 0">Bearbeiten</a></p>
		{{ end }}
	</div>
	{{if not $.user }}
		</a>
	{{end}}
</div>
{{end}}
<div class="pages">
	<a href="./0#events"><div class="page">Anfang</div></a>
	{{if gt $.page 0}}<a href="./{{dec $.page}}#events">{{end}}<div class="page">&lt;</div>{{if gt $.page 0}}</a>{{end}}
	{{range $.pages}}
		{{if or (and (ge . (dec (dec $.page))) (le . (inc (inc $.page)))) (or (eq $.page .) (or (le . 1) (ge . (dec $.maxPage))))}}
			<a href="./{{.}}#events"><div class="page {{if eq . $.page}}cur-page{{end}}">{{inc .}}</div></a>
		{{else}}
			{{if or (eq . (dec (dec (dec $.page)))) (eq . (inc (inc (inc $.page))))}}
				<div class="page">..</div>
			{{end}}
		{{end}}
	{{end}}
	{{if lt $.page $.maxPage}}<a href="./{{inc $.page}}#events">{{end}}<div class="page">&gt;</div>{{if lt $.page $.maxPage}}</a>{{end}}
	<a href="./{{$.maxPage}}#events"><div class="page">Ende</div></a>
</div>
{{end}}
