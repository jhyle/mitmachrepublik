<div id="events" class="row tiles">
	<div class="col-xs-1">&nbsp;</div>
	<div class="col-xs-3 col-box">
		{{template "organizer_box.tpl" .}}
		<p><a href="/veranstalter/verwaltung/veranstaltung" class="btn btn-mmr" style="width: 90%">Veranstaltung eintragen</a></p>
		<p><a href="/veranstalter/verwaltung/profil" class="btn btn-mmr" style="width: 90%">Beschreibung ändern</a></p>
		<p><a href="/veranstalter/verwaltung/kennwort" class="btn btn-mmr" style="width: 90%">Kennwort ändern</a></p>
		<p><a href="/veranstalter/verwaltung/kennwort" class="btn btn-mmr" style="width: 90%">E-Mail-Adresse ändern</a></p>
		<p><a href="#" id="delete-profile" class="btn btn-mmr" style="width: 90%">Profil löschen</a></p>
	</div>
	<div class="col-xs-7">
		<form class="form-inline pull-left" role="form" action="0?">
			<input type="text" name="query" class="form-control" style="margin-left: 0" placeholder="Veranstaltungen" value="{{.query}}">
			von
			<select name="location" style="margin-left: 10px; width: 230px">
				<option value="">allen<option>
				{{range $.locations}}
				<option value="{{.}}"{{if eq . $.location}} selected{{end}}>{{.}}</option>
				{{end}}
			</select>
			<button type="submit" class="btn btn-mmr" style="padding-left: 20px; padding-right: 20px">Suchen</button>
		</form>
		<p class="pull-right"></p>
		<p class="clearfix"></p>
		{{if not .user.Approved}}
		<div class="row-tile">
			<p class="big-text text-center" style="padding-top: 20px">Deine E-Mail-Adresse wurde noch nicht bestätigt, daher werden Deine Veranstaltungen nicht veröffentlicht. <a href="#" id="send-double-opt-in" class="btn btn-mmr" data-loading-text="Senden..">Bestätigungs-E-Mail erneut versenden</a></p>
		</div>
		{{end}}
		{{template "events_list.tpl" .}}
	</div>
	<div class="col-xs-1">&nbsp;</div>
</div>
