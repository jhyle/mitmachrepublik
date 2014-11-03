<div class="row tiles">
	<div class="col-sm-1 hidden-xs">&nbsp;</div>
	<div class="col-sm-3 col-box">
		{{if .Image}}<p><img src="/bild/{{.Image}}?width=250" class="img-responsive" /></p>{{end}}
		<p>{{.Descr}}</p>
		{{if .Web}}<p><a href="{{.Web}}" target="_blank" class="highlight"><span class="fa fa-caret-right"></span> {{.Web}}</a></p>{{end}}
		<p><a href="#" class="btn btn-mmr" style="width: 90%">Kennwort ändern</a></p>
		<p><a href="/veranstalter/verwaltung/profil" class="btn btn-mmr" style="width: 90%">Profil bearbeiten</a></p>
		<p><a href="#" class="btn btn-mmr" style="width: 90%">Profil löschen</a></p>
	</div>
	<div class="col-sm-7">
		<p class="pull-right"><a href="/veranstalter/verwaltung/veranstaltung" class="btn btn-mmr" style="width: 90%">Veranstaltung eintragen</a></p>
	</div>
	<div class="col-sm-1 hidden-xs">&nbsp;</div>
</div>
