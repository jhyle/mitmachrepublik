<div id="events" class="row tiles">
	<div class="col-sm-1 hidden-xs">&nbsp;</div>
	<div class="col-sm-3 col-box">
		{{template "organizer.tpl" .user}}
		<p><a href="/veranstalter/verwaltung/profil" class="btn btn-mmr" style="width: 90%">Profil bearbeiten</a></p>
		<p><a href="/veranstalter/verwaltung/kennwort" class="btn btn-mmr" style="width: 90%">Kennwort ändern</a></p>
		<p><a href="/veranstalter/verwaltung/kennwort" class="btn btn-mmr" style="width: 90%">E-Mail-Adresse ändern</a></p>
		<p><a href="#" id="delete-profile" class="btn btn-mmr" style="width: 90%">Profil löschen</a></p>
	</div>
	<div class="col-sm-7">
		<p class="pull-right"><a href="/veranstalter/verwaltung/veranstaltung" class="btn btn-mmr" style="margin-right: 0">Veranstaltung eintragen</a></p>
		<p class="clearfix"></p>
		{{template "events_list.tpl" .}}
	</div>
	<div class="col-sm-1 hidden-xs">&nbsp;</div>
</div>
