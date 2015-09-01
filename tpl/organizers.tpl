{{template "head.tpl" .}}
<form id="events-form" role="form" action="/suche" method="POST">
{{template "banner_search.tpl" .}}
<input type="hidden" name="search" value="organizers" /> 

<div class="row tiles">
	<div class="col-xs-1">&nbsp;</div>
	<div class="col-xs-11"><h1>Organisatoren{{if .place}} in {{.place}}{{end}}{{if and (gt .results 0) (gt .maxPage 0)}} - Seite {{inc .page}} von {{inc .maxPage}}{{end}}</h1></div>
</div>
<div id="organizers" class="row tiles">
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
	</div>
	<div class="col-xs-7">
		{{template "organizers_list.tpl" .}}
	</div>
</div>

</form>
{{template "foot.tpl" .}}