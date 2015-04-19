<div class="modal-header">
	<button type="button" class="close" data-dismiss="modal"><span aria-hidden="true">&times;</span><span class="sr-only">Schließen</span></button>
</div>
<div class="modal-body">
	<form role="form" id="login-form" class="form-horizontal">
		<div class="form-group">
			<div class="col-sm-7" style="border-right: 1px solid #ccc; padding-right: 0">
				<div class="big-text">Ich bin bereits als Organisator registriert.</div>
				<input name="email" id="login-Email" type="email" class="form-control" placeholder="E-Mail-Adresse">
				<input name="password" id="login-Pwd" type="password" class="form-control" placeholder="Kennwort">
				<button name="login" type="submit" class="btn btn-mmr" style="width: 90%">Anmelden</button>
			</div>
			<div class="col-sm-5">
			<div class="big-text">Ich bin neu hier und suche Mitmacher für meine nichtkommerziellen und gemeinschaftlichen Veranstaltungen.</div>						
				<button name="register" type="button" data-dismiss="modal" data-toggle="modal" data-target="#register" data-remote="/dialog/register" class="btn btn-mmr" style="margin-top: 18px; width: 90%">Kostenlos registrieren</button>
			</div>
		</div>
	</form>
</div>
