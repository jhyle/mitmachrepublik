<div class="row tiles">
	<div class="col-sm-1 hidden-xs">&nbsp;</div>
	<div class="col-sm-3 col-box">
		{{template "organizer.tpl" .user}}
		<p><a href="/veranstalter/verwaltung/profil" class="btn btn-mmr" style="width: 90%">Profil bearbeiten</a></p>
		<p><a href="/veranstalter/verwaltung/kennwort" class="btn btn-mmr" style="width: 90%">Kennwort ändern</a></p>
		<p><a href="/veranstalter/verwaltung/kennwort" class="btn btn-mmr" style="width: 90%">E-Mail-Adresse ändern</a></p>
		<p><a href="#" id="delete-profile" class="btn btn-mmr" style="width: 90%">Profil löschen</a></p>
	</div>
	<div class="col-sm-7">
		<p class="pull-right"><a href="/veranstalter/verwaltung/veranstaltung" class="btn btn-mmr">Veranstaltung eintragen</a></p>
		<p class="clearfix"></p>
		{{ range .events }}
		<div class="row-tile">
			{{ if .Image }}
			<div class="small-icon"><span class="fa fa-futbol-o"></span></div>
			<img class="img-responsive pull-left" style="margin-right: 10px" src="/bild/{{.Image}}?width=220&height=165">
			{{ end }}
			<div class="tile-text">
				<p class="pull-right"><a href="#" name="delete-event" data-target="{{.Id.Hex}}" class="close"><span class="fa fa-times"></span></a></p>
				<h3>{{.Title}}</h3>
				<p class="datetime">{{ dateFormat .Start }}</p>
				<p>{{.Descr}}</p>
				<p class="place-icon pull-left"><span class="fa fa-map-marker"></span></p>
				<p class="pull-left place">{{.Addr.Name}}, {{.Addr.Street}}, {{.Addr.Pcode}} {{.Addr.City}}</p>
				<p class="pull-right"><a href="/veranstalter/verwaltung/veranstaltung/{{.Id.Hex}}" class="btn btn-mmr" style="margin: 0">Bearbeiten</a></p>
			</div>
		</div>
		{{ end }}
	</div>
	<div class="col-sm-1 hidden-xs">&nbsp;</div>
</div>
