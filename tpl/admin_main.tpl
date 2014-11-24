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
		{{ $n := len .events }}
		{{ if eq $n 0 }}
		<div class="row-tile">
			<p class="big-text text-center" style="padding-top: 20px">Du hast keine Veranstaltungen eingetragen.</p>
		</div>
		{{ else }}
		{{ range .events }}
		<div class="row-tile">
			{{ if .Image }}
			<div class="small-icon"><span class="fa fa-futbol-o"></span></div>
			<img class="img-responsive pull-left" style="margin-right: 10px" src="/bild/{{.Image}}?width=220&height=165">
			{{ end }}
			<div class="tile-text">
				<p class="pull-right"><a href="#" name="delete-event" title="Löschen" data-target="{{.Id.Hex}}" class="close"><span class="fa fa-times"></span></a></p>
				<h3>{{.Title}}</h3>
				<p class="datetime">{{ dateFormat .Start }}</p>
				<p>{{.Descr}}</p>
				{{ if not .Addr.IsEmpty }}
					<p class="place-icon pull-left"><span class="fa fa-map-marker"></span></p>
					<p class="pull-left place">{{ if .Addr.Name }}{{.Addr.Name}}, {{ end }}{{ if .Addr.Street }}{{.Addr.Street}}, {{ end }}{{ if .Addr.Pcode }}{{.Addr.Pcode}} {{ end }}{{.Addr.City}}</p>
				{{ end }}
				<p class="pull-right"><a href="/veranstalter/verwaltung/veranstaltung/{{.Id.Hex}}" class="btn btn-mmr" style="margin: 0">Bearbeiten</a></p>
			</div>
		</div>
		{{ end }}
		{{ end }}
	</div>
	<div class="col-sm-1 hidden-xs">&nbsp;</div>
</div>
