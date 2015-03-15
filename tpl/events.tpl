{{template "head.tpl"}}
<form id="events-form" role="form" action="/suche" method="POST"> 
{{template "banner_search.tpl" .}}

<!-- div class="row tiles">
	{{ $len := len .events }}
	<div class="col-sm-1 hidden-xs">&nbsp;</div>
	<div class="col-sm-10"><h1>{{if eq $len 1}}Eine Veranstaltung{{else}}{{$len}} Veranstaltungen{{end}} in Berlin</h1></div>
	<div class="col-sm-1 hidden-xs">&nbsp;</div>
</div-->
<div id="events" class="row tiles">
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
		>  Am Wochenende</label>
	</div>
	<div class="col-sm-7">
		{{template "events_list.tpl" .}}
	</div>
</div>

</form>
{{template "foot.tpl"}}