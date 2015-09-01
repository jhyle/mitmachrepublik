<div class="row tiles">
	<div class="col-xs-1">&nbsp;</div>
	<div class="col-xs-11"><h1>Veranstaltungen von {{.organizer.Name}} {{if and (gt .results 0) (gt .maxPage 0)}} - Seite {{inc .page}} von {{inc .maxPage}}{{end}}</h1></div>
</div>
<div class="row tiles">
	<div class="col-xs-1">&nbsp;</div>
	<div class="col-xs-3 col-box">
		{{template "organizer_box.tpl" .}}
	</div>
	<div class="col-xs-7">
		{{template "events_list.tpl" .}}
	</div>
	<div class="col-xs-1">&nbsp;</div>
</div>
