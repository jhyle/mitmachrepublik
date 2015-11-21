<div class="row tiles">
{{range .events}}
	{{template "tile_event.tpl" .}}
{{end}}
</div>
<div class="row tiles">
	<div class="col-sm-3 col-xs-2">&nbsp;</div>
	<div class="col-sm-6 col-xs-8">
		<a href="/veranstaltungen/{{simpleEventSearchUrl ""}}" class="btn btn-mmr" style="width: 100%">Mehr Veranstaltungen</a>
	</div>
	<div class="col-sm-3 col-xs-2">&nbsp;</div>
</div>
