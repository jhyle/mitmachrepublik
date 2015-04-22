<div class="row tiles">
	<div class="col-xs-1">&nbsp;</div>
	<div class="col-xs-3 col-box">
		{{template "organizer_box.tpl" .}}
	</div>
	<div class="col-xs-7">
		<h3 style="margin-left: 10px">Profil bearbeiten</h3>
		<form role="form" id="profile-upload" class="form-horizontal" action="/upload" method="POST">
			<div class="form-group">
				<div class="col-xs-7">
					<span><input name="name" type="text" id="profile-Name" class="form-control" placeholder="Deine Organisation" value="{{.user.Name}}" autofocus></span>
				</div>
				<div class="col-xs-4">
					<a id="profile-dropzone" class="thumbnail" style="margin: 10px; cursor: pointer">
						<span id="profile-spinner" class="fa fa-gear"> </span>
						<img src="{{if .user.Image}}/bild/{{.user.Image}}?height=200&width=200{{else}}/images/thumbnail.png{{end}}" alt="Bild" id="profile-thumbnail" class="img-responsive">
					</a>
					<span id="profile-thumbnail-message" class="help-block">Wähle ein Bild im Format jpg, jpeg, png oder gif aus.</span>
					<input type="file" name="file" class="hide">
					<input type="hidden" name="image" id="profile-Image" value="{{.user.Image}}">
				</div>
			</div>
			<hr>
			<div class="form-group">
				<div class="col-xs-12" style="margin-left: 10px">
					<span id="event-Category" class="help-block">Wähle eine oder mehrere Kategorien aus:</span>
				{{ range .categories }}
					{{ $id := index $.categoryMap . }}
					<label class="checkbox-inline"><input type="checkbox" name="profile-Category" value="{{$id}}"
					{{ range $.user.Categories }}
						{{ if eq . $id }}checked{{ end }}
					{{ end }}
					> {{.}} &nbsp;&nbsp;</label>
				{{ end }}
				</div>
			</div>
			<hr>
			<div class="form-group">
				<div class="col-xs-12">
					<div id="profile-Descr" placeholder="Beschreibung">{{.user.HtmlDescription}}</div>
					<span><input name="website" type="text" id="profile-Web" class="form-control" placeholder="Webseite" value="{{.user.Web}}"></span>
				</div>
			</div>
			<hr>
			<p style="margin-left: 12px">
				<span class="help-block">Gib eine Adresse an, um in der Veranstaltersuche gefunden zu werden.</span>
			</p>
			<div class="form-group">
				<div class="col-xs-5">
					<input name="street" type="text" id="profile-Street" class="form-control" placeholder="Straße" value="{{.user.Addr.Street}}">
				</div>
				<div class="col-xs-3">
					<input name="pcode" type="text" id="profile-Pcode" class="form-control" placeholder="Postleitzahl" value="{{.user.Addr.Pcode}}">
				</div>
				<div class="col-xs-3">
					<input name="city" type="text" id="profile-City" class="form-control" placeholder="Ort" value="{{.user.Addr.City}}">
				</div>
			</div>
			<hr>
			<div class="form-group">
				<div class="col-xs-4">
					<a href="/veranstalter/verwaltung/0" onClick="history.back(); return false" class="btn btn-default btn-cancel" style="width: 90%">Abbrechen</a>
				</div>
				<div class="col-xs-1">&nbsp;</div>
				<div class="col-xs-7">
					<button type="submit" class="btn btn-mmr" style="width: 90%">Speichern</button>
				</div>
			</div>
		</form>
	</div>
	<div class="col-xs-1">&nbsp;</div>
</div>
