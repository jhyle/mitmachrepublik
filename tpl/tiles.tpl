{{range .events}}
<div class="row tiles">
	<div class="col-sm-1 hidden-xs">&nbsp;</div>
	{{range .}}
		{{template "tile_event.tpl" .}}
	{{end}}
	<div class="col-sm-1 hidden-xs">&nbsp;</div>
</div>
{{end}}
