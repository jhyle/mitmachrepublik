<div class="row tiles">
	<div class="col-xs-1">&nbsp;</div>
	<div class="col-xs-3 col-box">
		{{template "organizer_box.tpl" .}}
	</div>
	<div class="col-xs-7">
		<h3 style="margin-left: 10px">E-Mail-Adresse oder Kennwort ändern</h3>
		<form role="form" id="password" class="form-horizontal" action="/password" method="POST">
			<div class="form-group">
				<div class="col-xs-12">
					<span><input name="email" type="email" id="password-Email" class="form-control" placeholder="E-Mail-Adresse" value="{{.Email}}"></span>
					<span class="help-block" style="margin: 10px; width: 90%">Bei Änderung der E-Mail-Adresse werden Ihre Veranstaltungen auf nicht sichtbar geschaltet, bis Sie die neue E-Mail-Adresse bestätigt haben.</span>
					<span><input name="pwd" type="password" id="password-Pwd" class="form-control" placeholder="Neues Kennwort">
					<input name="pwd2" type="password" id="password-Pwd2" class="form-control" placeholder="Neues Kennwort wiederholen"></span>
				</div>
			</div>
			<hr>
			<div class="form-group">
				<div class="col-xs-4">
					<a href="/veranstalter/verwaltung/0" onClick="history.back(); return false" class="btn btn-default" style="width: 90%">Abbrechen</a>
				</div>
				<div class="col-xs-1">&nbsp;</div>
				<div class="col-xs-7">
					<button id="password-submit" type="submit" class="btn btn-mmr" style="width: 90%" data-loading-text="Speichern..">Speichern</button>
				</div>
			</div>
		</form>
	</div>
	<div class="col-xs-1">&nbsp;</div>
</div>
