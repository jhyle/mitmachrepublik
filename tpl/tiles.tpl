{{range .events}}
<div class="row tiles">
	<div class="col-xs-1 hidden-xs">&nbsp;</div>
	{{range .}}
		{{template "tile_event.tpl" .}}
	{{end}}
	<div class="col-xs-1 hidden-xs">&nbsp;</div>
</div>
{{end}}
<div class="row tiles">
	<div class="col-xs-1 hidden-xs">&nbsp;</div>
	<div class="col-xs-2">&nbsp;</div>
	<div class="col-xs-6">
		<a href="/veranstaltungen/{{simpleEventSearchUrl ""}}" class="btn btn-mmr" style="width: 100%">Mehr Veranstaltungen</a>
	</div>
	<div class="col-xs-2">&nbsp;</div>
	<div class="col-xs-1 hidden-xs">&nbsp;</div>
</div>
