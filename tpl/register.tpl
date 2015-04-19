<div class="modal-header">
	<button type="button" class="close" data-dismiss="modal"><span aria-hidden="true">&times;</span><span class="sr-only">Schließen</span></button>
</div>
<div class="modal-body">
	<div class="big-text">Trage Deine Organisation ein, um Mitmacher für Deine Veranstaltungen zu finden.</div>
	<form role="form" id="register-upload" class="form-horizontal" action="/upload" method="POST">
		<div class="form-group">
			<div class="col-sm-7">
				<span><input name="name" type="text" id="register-Name" class="form-control" placeholder="Deine Organisation"></span>
				<span><input name="email" type="email" id="register-Email" class="form-control" placeholder="E-Mail-Adresse"></span>
				<span><input name="pwd" type="password" id="register-Pwd" class="form-control" placeholder="Kennwort">
				<input name="pwd2" type="password" id="register-Pwd2" class="form-control" placeholder="Kennwort wiederholen"></span>
			</div>
			<div class="col-sm-4">
				<a id="register-dropzone" class="thumbnail" style="margin: 10px; cursor: pointer">
					<span id="register-spinner" class="fa fa-gear"> </span>
					<img src="/images/thumbnail.png" alt="Bild" id="register-thumbnail" class="img-responsive">
				</a>
				<span id="register-thumbnail-message" class="help-block">Wähle ein Bild im Format jpg, jpeg, png oder gif aus.</span>
				<input type="file" name="file" class="hide">
				<input type="hidden" name="image" id="register-Image">
			</div>
		</div>
		<hr>
		<div class="form-group">
			<div class="col-xs-12" style="margin-left: 10px">
				<span id="event-Category" class="help-block">Wähle eine oder mehrere Kategorien aus:</span>
		{{ range .categories }}
			{{ $id := index $.categoryMap . }}
				<label class="checkbox-inline"><input type="checkbox" name="register-Category" value="{{$id}}"> {{.}} &nbsp;&nbsp;</label>
		{{ end }}
			</div>
		</div>
		<hr>
		<div class="form-group">
			<div class="col-sm-12">
				<textarea name="description" id="register-Descr" class="form-control" placeholder="Beschreibung" rows="5"></textarea>
				<span><input name="website" type="text" id="register-Web" class="form-control" placeholder="Webseite"></span>
			</div>
		</div>
		<hr>
		<p style="margin-left: 12px">
			<span class="help-block">Gib eine Adresse an, um in der Organisatorsuche gefunden zu werden.</span>
		</p>
		<div class="form-group">
			<div class="col-sm-5">
				<input name="street" type="text" id="register-Street" class="form-control" placeholder="Straße">
			</div>
			<div class="col-sm-3">
				<input name="pcode" type="text" id="register-Pcode" class="form-control" placeholder="Postleitzahl">
			</div>
			<div class="col-sm-3">
				<input name="city" type="text" id="register-City" class="form-control" placeholder="Ort">
			</div>
		</div>
		<hr>
		<div class="form-group">
			<div class="col-sm-12" style="margin-left: 5px">
				<label class="checkbox-inline" style="margin-left: 10px"><input type="checkbox" name="agbs" id="register-AGBs" value="Y"> Ich stimme den <a class="highlight" href="/agbs" target="_blank">Allgemeinen Geschäftsbedingungen</a> zu.</label>
			</div>
		</div>
		<div class="form-group">
			<div class="col-sm-4">
				<button type="button" class="btn btn-default" data-dismiss="modal" style="width: 90%">Abbrechen</button>
			</div>
			<div class="col-sm-1">&nbsp;</div>
			<div class="col-sm-7">
				<button id="register-submit" type="submit" class="btn btn-mmr" data-loading-text="Registrieren.." style="width: 90%">Registrieren</button>
			</div>
		</div>
	</form>
</div>
