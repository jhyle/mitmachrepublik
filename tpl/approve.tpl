{{template "head.tpl" .}}
{{template "banner_small.tpl" .meta.FB_Title}} 
<div class="row tiles">
	<div class="col-xs-1">&nbsp;</div>
	<div class="col-xs-10">
		<p>{{if .approved}}Vielen Dank für Deine Bestätigung! Dein Profil ist jetzt aktiv.{{else}}Deine ID wurde nicht gefunden. Bitte registriere Dich noch einmal oder wende Dich an den Support.{{end}}
	</div>
	<div class="col-xs-1">&nbsp;</div>
</div>
{{template "foot.tpl" .}}