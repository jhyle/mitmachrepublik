{{template "head.tpl" .}}
<form id="events-form" role="form" action="/suche" method="POST">
{{template "banner_search.tpl" .}}
<input type="hidden" name="search" value="organizers" /> 

<div class="row tiles">
	{{ $len := len .organizers }}
	<div class="col-sm-1 hidden-xs">&nbsp;</div>
	<div class="col-sm-10"><h1>{{if eq $len 0}}Keine{{else}}{{if eq $len 1}}Ein{{else}}{{$len}}{{end}}{{end}} Veranstalter in {{.place}} gefunden</h1></div>
	<div class="col-sm-1 hidden-xs">&nbsp;</div>
</div>
<div id="organizers" class="row tiles">
	<div class="col-sm-1 hidden-xs">&nbsp;</div>
	<div class="col-sm-3 col-box">
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
	<div class="col-sm-7">
		{{template "organizers_list.tpl" .}}
	</div>
</div>

</form>
{{template "foot.tpl" .}}