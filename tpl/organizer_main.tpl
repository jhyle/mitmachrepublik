<div class="row tiles">
	<div class="col-xs-1">&nbsp;</div>
	<div class="col-xs-10"><h1>{{if eq .results 1}}Eine Veranstaltung{{else}}{{if eq .results 0}}Keine{{else}}{{.results}}{{end}} Veranstaltungen{{end}} von {{.meta.FB_Title}}</h1></div>
	<div class="col-xs-1">&nbsp;</div>
</div>
<div class="row tiles">
	<div class="col-xs-1">&nbsp;</div>
	<div class="col-xs-3 col-box">
		{{template "organizer_box.tpl" .organizer}}
	</div>
	<div class="col-xs-7">
		{{template "events_list.tpl" .}}
	</div>
	<div class="col-xs-1">&nbsp;</div>
</div>
