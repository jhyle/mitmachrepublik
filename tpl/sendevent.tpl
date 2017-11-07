<div class="modal-header">
	<button type="button" class="close" data-dismiss="modal"><span aria-hidden="true">&times;</span><span class="sr-only">Schließen</span></button>
</div>
<div class="modal-body">
	<div class="big-text">Sende die Veranstaltung per Mail an einen Freund.</div>
	<form role="form" id="send-event" class="form-horizontal" method="POST">
		<div class="form-group">
			<div class="col-sm-12">
				<input name="name" type="text" id="send-event-Name" class="form-control" placeholder="Name des Empfängers">
				<span><input name="email" type="email" id="send-event-Email" class="form-control" placeholder="E-Mail-Adresse des Empfängers"></span>
				<span><input name="subject" type="text" id="send-event-Subject" class="form-control" placeholder="Betreff" value="Veranstaltung {{.event.Title}} auf mitmach-republik.de"></span>
				<textarea name="text" id="send-event-Text" class="form-control" placeholder="Nachricht" rows="5">Hallo,

die Veranstaltung {{.event.Title}} in {{citypartName .event.Addr}} finde ich interessant, schau doch mal rein: http://{{$.hostname}}{{.event.Url | encodePath}}?from={{$.from}}.

Liebe Grüße! 
				</textarea>
			</div>
		</div>
		<hr>
		<div class="form-group">
			<div class="col-sm-4">
				<button type="button" class="btn btn-cancel" data-dismiss="modal" style="width: 90%">Abbrechen</button>
			</div>
			<div class="col-sm-1">&nbsp;</div>
			<div class="col-sm-7">
				<button id="send-event-submit" type="submit" class="btn btn-mmr" data-loading-text="Senden.." style="width: 90%">Senden</button>
			</div>
		</div>
	</form>
</div>
