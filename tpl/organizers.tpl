{{template "head.tpl" .}}
<form id="events-form" role="form" action="/suche" method="POST">
{{template "banner_search.tpl" .}}
<input type="hidden" name="search" value="organizers" /> 

<div class="row tiles">
	<div class="col-xs-12"><h1>Gemeinschaftliche Organisatoren{{if .place}} in {{.place}}{{end}}</h1></div>
</div>
<div id="organizers" class="row tiles">
	<div class="col-md-3 col-sm-4 col-xs-12 col-box">
		<div class="filter-box">
		<div class="box-headline"><a href="#" onClick="$(this).parent().next().slideToggle();return false">Filter</a></div>
		<div class="box-body">
			<div class="filter-headline">Kategorien</div>
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
	</div>
	<div class="col-md-9 col-sm-8 col-xs-12">
		{{template "organizers_list.tpl" .}}
	</div>
</div>

</form>
{{template "foot.tpl" .}}