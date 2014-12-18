<div class="row tiles">
	<div class="col-sm-1 hidden-xs">&nbsp;</div>
	<div class="col-sm-3 col-box">
		{{template "organizer.tpl" .user}}
	</div>
	<div class="col-sm-7" style="padding-top: 15px">
		<h3 style="margin-left: 10px">Veranstaltung {{ if .event.Id }}bearbeiten{{ else }}eintragen{{ end }}</h3>
		<form role="form" id="event-upload" class="form-horizontal" action="/upload" method="POST">
			{{if .event.Id }}<input type="hidden" id="event-Id" value="{{.event.Id.Hex}}">{{ end }}
			<div class="form-group">
				<div class="col-sm-7">
					<span><input name="title" type="text" id="event-Title" class="form-control" placeholder="Was willst Du machen?" value="{{.event.Title}}" maxlength="40"></span>
					<span><input name="start" type="datetime-local" id="event-Start" class="form-control form-datetime" placeholder="Fängt an" value="{{dateFormat .event.Start}}"></span>
					<span><input name="end" type="datetime-local" id="event-End" class="form-control form-datetime" placeholder="Endet" value="{{dateFormat .event.End}}"></span>
				</div>
				<div class="col-sm-4">
					<a id="event-dropzone" class="thumbnail" style="margin: 10px; cursor: pointer">
						<span id="event-spinner" class="fa fa-gear"> </span>
						<img src="{{if .event.Image}}/bild/{{.event.Image}}?height=200&width=200{{else}}/images/thumbnail.gif{{end}}" alt="Bild" id="event-thumbnail" class="img-responsive">
					</a>
					<span id="event-thumbnail-message" class="help-block">Wähle ein Bild im Format jpg, jpeg, png oder gif aus.</span>
					<input type="file" name="file" class="hide">
					<input type="hidden" name="image" id="event-Image" value="{{.event.Image}}">
				</div>
			</div>
			<hr>
			<div class="form-group">
				<div class="col-sm-12" style="margin-left: 10px">
					<span id="event-Category" class="help-block">Wähle eine oder mehrere Kategorien aus:</span>
				{{ range .categories }}
					{{ $id := index $.categoryIds . }}
					<label class="checkbox-inline"><input type="checkbox" name="event-Category" value="{{$id}}"
					{{ range $.event.Categories }}
						{{ if eq . $id }}checked{{ end }}
					{{ end }}
					> {{.}} &nbsp;&nbsp;</label>
				{{ end }}
				</div>
			</div>
			<hr>
			<div class="form-group">
				<div class="col-sm-12">
					<textarea name="description" id="event-Descr" class="form-control" placeholder="Beschreibung">{{.event.Descr}}</textarea>
					<span><input name="website" type="text" id="event-Web" class="form-control" placeholder="Webseite" value="{{.event.Web}}"></span>
				</div>
			</div>
			<hr>
			<div class="form-group">
				<div class="col-sm-12">
					<span><input name="location" type="text" id="event-Name" class="form-control" placeholder="Wo findet die Veranstaltung statt?" value="{{.event.Addr.Name}}"></span>
				</div>
			</div>
			<div class="form-group">
				<div class="col-sm-5">
					<input name="street" type="text" id="event-Street" class="form-control" placeholder="Straße" value="{{.event.Addr.Street}}">
				</div>
				<div class="col-sm-3">
					<input name="pcode" type="text" id="event-Pcode" class="form-control" placeholder="Postleitzahl" value="{{.event.Addr.Pcode}}" maxlength="5">
				</div>
				<div class="col-sm-3">
					<input name="city" type="text" id="event-City" class="form-control" placeholder="Ort" value="{{.event.Addr.City}}">
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
