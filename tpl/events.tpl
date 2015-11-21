{{template "head.tpl" .}}
<form id="events-form" role="form" action="/suche" method="POST"> 
{{template "banner_search.tpl" .}}

<div class="row tiles">
	<div class="col-xs-12"><h1>{{if and (eq 1 (len .categoryIds)) (ne 0 (index .categoryIds 0))}}{{ range .categories }}{{if eq (index $.categoryMap .) (index $.categoryIds 0)}}{{.}}{{end}}{{end}}{{else}}Veranstaltungen{{end}}{{if eq 1 (len .targetIds)}}{{ range .targets }}{{if eq (index $.targetMap .) (index $.targetIds 0)}} für {{.}}{{end}}{{end}}{{end}}{{.dateNames}}{{if .place}} in {{.place}}{{end}}{{if and (gt .results 0) (gt .maxPage 0)}} - Seite {{inc .page}} von {{inc .maxPage}}{{end}}</h1></div>
</div>
<div class="row tiles">
	<div class="col-md-3 col-sm-4 col-xs-12 col-box">
		<div class="filter-box">
		<h3>Filter</h3>
		<h5>Datum</h5>
		<hr>
		{{ range .dates }}
			{{ $id := . }}
			{{ if gt $id 0 }}
				<label class="checkbox"><input type="checkbox" name="date" value="{{$id}}"
				{{ range $.dateIds }}
					{{ if eq $id . }} checked {{ end }}
				{{ end }}
				>  {{ index $.dateMap $id }}</label>
			{{ end }}
		{{ end }}
		<h5>Zielgruppen</h5>
		<hr>
		{{ range .targets }}
			{{ $id := index $.targetMap . }}
			<label class="checkbox"><input type="checkbox" name="target" value="{{$id}}"
			{{ range $.targetIds }}
				{{ if eq $id . }} checked {{ end }}
			{{ end }}
			>  {{.}}</label>					
		{{ end }}
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
		</div>
	</div>
	<div class="col-md-9 col-sm-8 col-xs-12">
		{{template "events_list.tpl" .}}
	</div>
</div>

</form>
{{template "foot.tpl" .}}