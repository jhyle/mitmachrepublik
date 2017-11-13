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
				<a class="highlight" style="margin-left: 12px" href="#" data-dismiss="modal" data-toggle="modal" data-target="#send-password" data-remote="/dialog/password"><span class="fa fa-caret-right"></span> Kennwort vergessen?</a>
				<button name="login" id="login-submit" type="submit" class="btn btn-mmr" style="width: 90%" data-loading-text="Anmelden..">Anmelden</button>
			</div>
			<div class="col-sm-5">
			<div class="big-text">Ich bin neu hier und suche Mitmacher für meine gemeinschaftlichen Projekte und Veranstaltungen.<br><br></div>						
				<button name="register" type="button" data-dismiss="modal" data-toggle="modal" data-target="#register" data-remote="/dialog/register" class="btn btn-mmr" style="margin-top: 49px; width: 90%">Kostenlos registrieren</button>
			</div>
		</div>
	</form>
</div>
