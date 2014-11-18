<div class="row tiles">
	<div class="col-sm-1 hidden-xs">&nbsp;</div>
	<div class="col-sm-3 col-box">
		{{template "organizer.tpl" .}}
	</div>
	<div class="col-sm-7">
		<h3 style="margin-left: 10px">Profil bearbeiten</h3>
		<form role="form" id="profile-upload" class="form-horizontal" action="/upload" method="POST">
			<div class="form-group">
				<div class="col-sm-7">
					<span><input name="name" type="text" id="profile-Name" class="form-control" placeholder="Deine Organisation" value="{{.Addr.Name}}"></span>
				</div>
				<div class="col-sm-4">
					<a id="profile-dropzone" class="thumbnail" style="margin: 10px; cursor: pointer">
						<span id="profile-spinner" class="fa fa-gear"> </span>
						<img src="{{if .Image}}/bild/{{.Image}}?height=200&width=200{{else}}/images/thumbnail.gif{{end}}" alt="Bild" id="profile-thumbnail" class="img-responsive">
					</a>
					<span id="profile-thumbnail-message" class="help-block">Wähle ein Bild im Format jpg, jpeg, png oder gif aus.</span>
					<input type="file" name="file" class="hide">
					<input type="hidden" name="image" id="profile-Image" value="{{.Image}}">
				</div>
			</div>
			<hr>
			<div class="form-group">
				<div class="col-sm-12">
					<textarea name="description" id="profile-Descr" class="form-control" placeholder="Beschreibung">{{.Descr}}</textarea>
					<span><input name="web" type="text" id="profile-Web" class="form-control" placeholder="Webseite" value="{{.Web}}"></span>
				</div>
			</div>
			<hr>
			<div class="form-group">
				<div class="col-sm-5">
					<input name="street" type="text" id="profile-Street" class="form-control" placeholder="Straße" value="{{.Addr.Street}}">
				</div>
				<div class="col-sm-3">
					<input name="postcode" type="text" id="profile-Pcode" class="form-control" placeholder="Postleitzahl" value="{{.Addr.Pcode}}">
				</div>
				<div class="col-sm-3">
					<input name="city" type="text" id="profile-City" class="form-control" placeholder="Ort" value="{{.Addr.City}}">
				</div>
			</div>
			<hr>
			<div class="form-group">
				<div class="col-sm-4">
					<a href="/veranstalter/verwaltung" class="btn btn-default" style="width: 90%">Abbrechen</a>
				</div>
				<div class="col-sm-1">&nbsp;</div>
				<div class="col-sm-7">
					<button type="submit" class="btn btn-mmr" style="width: 90%">Speichern</button>
				</div>
			</div>
		</form>
	</div>
	<div class="col-sm-1 hidden-xs">&nbsp;</div>
</div>
