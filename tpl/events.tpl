{{template "head.tpl" .}}
<form id="events-form" role="form" action="/suche" method="POST"> 
{{template "banner_search.tpl" .}}

<div class="row tiles">
	<div class="col-xs-1">&nbsp;</div>
	<div class="col-xs-11"><h1>{{if eq .results 1}}Eine Veranstaltung{{else}}{{if eq .results 0}}Keine{{else}}{{.results}}{{end}} Veranstaltungen{{end}}{{if .place}} in {{.place}}{{end}} gefunden{{if gt .results 0}} - Seite {{inc .page}} von {{inc .maxPage}}{{end}}</h1></div>
</div>
<div id="events" class="row tiles">
	<div class="col-xs-1">&nbsp;</div>
	<div class="col-xs-3 col-box">
		<h3>Filter</h3>
		<h5>Kategorien</h5>
		<hr>
		{{ range .categories }}
			{{ $id := index $.categoryMap . }}
			<label class="checkbox"><input type="checkbox" name="category" value="{{$id}}"
			{{ range $.categoryIds }}
				{{ if eq $id . }} checked {{ end }}
			{{ end }}
			>  {{.}}</label>					
		{{ end }}
		<h5>Datum</h5>
		<hr>
		<label class="checkbox"><input type="checkbox" name="date" value="1"
		{{ range $.dates }}
			{{ if eq "heute" . }} checked {{ end }}
		{{ end }}
		>  Heute</label>
		<label class="checkbox"><input type="checkbox" name="date" value="2"
		{{ range $.dates }}
			{{ if eq "morgen" . }} checked {{ end }}
		{{ end }}
		>  Morgen</label>
		<label class="checkbox"><input type="checkbox" name="date" value="3"
		{{ range $.dates }}
			{{ if eq "wochenende" . }} checked {{ end }}
		{{ end }}
		>  Nächstes Wochenende</label>
	</div>
	<div class="col-xs-7">
		{{template "events_list.tpl" .}}
	</div>
</div>

</form>
{{template "foot.tpl" .}}