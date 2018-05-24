<div class="modal-header">
	<button type="button" class="close" data-dismiss="modal"><span aria-hidden="true">&times;</span><span class="sr-only">Schließen</span></button>
</div>
<div class="modal-body">
	<div class="big-text">Trage Deine E-Mail-Adresse ein und wir senden Dir an den ausgewählten Tagen die zu Deiner Suche passenden Veranstaltungen per E-Mail. Du kannst die Zusendung jederzeit über einen Link am Ende des Newsletters abbestellen oder Dich über das Kontaktformular an uns wenden.</div>
	<form role="form" id="email-alert-form" class="form-horizontal" method="POST">
		<div class="form-group">
			<div class="col-sm-12">
				<input name="name" type="text" id="email-alert-Name" class="form-control" placeholder="Dein Name">
				<span><input name="email" type="email" id="email-alert-Email" class="form-control" placeholder="Deine E-Mail-Adresse"></span>
				<input type="hidden" name="query" id="email-alert-Query" value="{{.query}}">
				<input type="hidden" name="place" id="email-alert-Place" value="{{.place}}">
				<input type="hidden" name="targets" id="email-alert-Targets" value="{{.targetIds}}">
				<input type="hidden" name="categories" id="email-alert-Categories" value="{{.categoryIds}}">
				<input type="hidden" name="dates" id="email-alert-Dates" value="{{.dateIds}}">
				<input type="hidden" name="radius" id="email-alert-Radius" value="{{.radius}}">
			</div>
			<div class="col-xs-12">
				<span class="checkbox-inline">Versand am&nbsp;&nbsp;</span>
				<label class="checkbox-inline"><input type="checkbox" name="email-alert-Weekday" value="1"> Mo&nbsp;&nbsp;</label>
				<label class="checkbox-inline"><input type="checkbox" name="email-alert-Weekday" value="2"> Di&nbsp;&nbsp;</label>
				<label class="checkbox-inline"><input type="checkbox" name="email-alert-Weekday" value="3"> Mi&nbsp;&nbsp;</label>
				<label id="email-alert-Weekday" class="checkbox-inline"><input type="checkbox" name="email-alert-Weekday" value="4"> Do&nbsp;&nbsp;</label>
				<label class="checkbox-inline"><input type="checkbox" name="email-alert-Weekday" value="5"> Fr&nbsp;&nbsp;</label>
				<label class="checkbox-inline"><input type="checkbox" name="email-alert-Weekday" value="6"> Sa&nbsp;&nbsp;</label>
				<label class="checkbox-inline"><input type="checkbox" name="email-alert-Weekday" value="0"> So&nbsp;&nbsp;</label>
			</div>
		</div>
		<hr>
		<div class="form-group">
			<div class="col-sm-4">
				<button type="button" class="btn btn-cancel" data-dismiss="modal" style="width: 90%">Abbrechen</button>
			</div>
			<div class="col-sm-1">&nbsp;</div>
			<div class="col-sm-7">
				<button id="email-alert-submit" type="submit" class="btn btn-mmr" data-loading-text="Eintragen.." style="width: 90%">Eintragen</button>
			</div>
		</div>
	</form>
</div>
