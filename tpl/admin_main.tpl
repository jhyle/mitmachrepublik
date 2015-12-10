<div id="events" class="row tiles">
	<div class="col-lg-3 col-sm-4 col-xs-5 col-box-admin">
		{{template "organizer_box.tpl" .}}
		<p></p>
		<p><a href="/veranstalter/verwaltung/profil" class="btn btn-mmr" style="width: 90%">Beschreibung ändern</a></p>
		<p><a href="/veranstalter/verwaltung/kennwort" class="btn btn-mmr" style="width: 90%">Kennwort ändern</a></p>
		<p><a href="/veranstalter/verwaltung/kennwort" class="btn btn-mmr" style="width: 90%">E-Mail-Adresse ändern</a></p>
		<p><a href="#" id="delete-profile" class="btn btn-mmr" style="width: 90%">Profil löschen</a></p>
	</div>
	<div class="col-lg-9 col-sm-8 col-xs-7">
		<form id="adminsearch" class="form-inline pull-left" role="form" action="0?">
			<input type="text" name="query" class="form-control search-field" placeholder="Veranstaltungen" value="{{.query}}" autocomplete="off">
			<button type="submit" class="btn btn-mmr btn-search"><span class="fa fa-search"></span></button>
		</form>
		<p class="pull-right"><a href="/veranstalter/verwaltung/veranstaltung" class="btn btn-mmr">Veranstaltung eintragen</a></p>
		<p class="clearfix"></p>
		{{if not .user.Approved}}
		<div class="row-tile">
			<p class="big-text text-center" style="padding-top: 20px">Deine E-Mail-Adresse wurde noch nicht bestätigt, daher werden Deine Veranstaltungen nicht veröffentlicht. <a href="#" id="send-double-opt-in" class="btn btn-mmr" data-loading-text="Senden..">Bestätigungs-E-Mail erneut versenden</a></p>
		</div>
		{{end}}
		{{template "events_list.tpl" .}}
	</div>
</div>
