{{template "head.tpl"}}
{{$title := "Registrierung bestätigen"}}
{{template "banner_small.tpl"  $title}} 
<div class="row tiles">
	<div class="col-sm-1 hidden-xs">&nbsp;</div>
	<div class="col-sm-10">
		<p>{{if .approved}}Vielen Dank für Deine Bestätigung! Dein Profil ist jetzt aktiv.{{else}}Deine ID wurde nicht gefunden. Bitte registriere Dich noch einmal oder wende Dich an den Support.{{end}}
	</div>
	<div class="col-sm-1 hidden-xs">&nbsp;</div>
</div>
{{template "foot.tpl"}}