<div class="modal-header">
	<button type="button" class="close" data-dismiss="modal"><span aria-hidden="true">&times;</span><span class="sr-only">Schließen</span></button>
</div>
<div class="modal-body">
	<div class="big-text">Schreibe uns eine Nachricht.</div>
	<form role="form" id="send-mail" class="form-horizontal" method="POST">
		<div class="form-group">
			<div class="col-sm-12">
				<input name="name" type="text" id="send-mail-Name" class="form-control" placeholder="Dein Name">
				<span><input name="email" type="email" id="send-mail-Email" class="form-control" placeholder="Deine E-Mail-Adresse"></span>
				<span><input name="subject" type="text" id="send-mail-Subject" class="form-control" placeholder="Betreff"></span>
				<textarea name="text" id="send-mail-Text" class="form-control" placeholder="Nachricht" rows="5"></textarea>
			</div>
		</div>
		<hr>
		<div class="form-group">
			<div class="col-sm-4">
				<button type="button" class="btn btn-default" data-dismiss="modal" style="width: 90%">Abbrechen</button>
			</div>
			<div class="col-sm-1">&nbsp;</div>
			<div class="col-sm-7">
				<button id="send-mail-submit" type="submit" class="btn btn-mmr" data-loading-text="Senden.." style="width: 90%">Senden</button>
			</div>
		</div>
	</form>
</div>
