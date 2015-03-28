{{range .events}}
<div class="row tiles">
	<div class="col-xs-1 hidden-xs">&nbsp;</div>
	{{range .}}
		{{template "tile_event.tpl" .}}
	{{end}}
	<div class="col-xs-1 hidden-xs">&nbsp;</div>
</div>
{{end}}
