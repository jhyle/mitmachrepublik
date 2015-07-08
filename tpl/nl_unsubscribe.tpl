{{template "head.tpl" .}}
{{template "banner_small.tpl" .meta.FB_Title}} 
<div class="row tiles">
	<div class="col-xs-1">&nbsp;</div>
	<div class="col-xs-11">
		<h3>{{if .unsubscribed}}Die  Veranstaltungsbenachrichtigung wurde gelöscht. Du erhältst in Zukunft keine E-Mail mehr.{{else}}Deine Veranstaltungsbenachrichtigung wurde nicht gefunden. Vielleicht hast Du sie schon gelöscht?{{end}}</h3>
	</div>
	<div class="col-xs-1">&nbsp;</div>
</div>
{{template "foot.tpl" .}}