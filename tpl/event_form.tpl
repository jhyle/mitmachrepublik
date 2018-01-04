<div class="row tiles">
	<div class="col-lg-3 col-sm-4 col-xs-5 col-box">
		{{template "organizer_box.tpl" .}}
	</div>
	<div class="col-lg-9 col-sm-8 col-xs-7">
		<h3 style="margin-left: 10px">Veranstaltung {{ if .event.Id }}bearbeiten{{ else }}eintragen{{ end }}</h3>
		<form role="form" id="event-upload" class="form-horizontal" action="/upload" method="POST">
			{{if .event.Id }}<input type="hidden" id="event-Id" value="{{.event.Id.Hex}}">{{ end }}
			<div class="form-group" style="margin-bottom: 0">
				<div class="col-xs-7">
					<span><div id="event-Title-Too-Long"></div><input name="title" type="text" id="event-Title" class="form-control" placeholder="Wie heißt die Veranstaltung?" value="{{.event.Title}}" autofocus></span>
					<span><input name="start" type="text" id="event-Start" class="form-control form-datetime" placeholder="Beginnt" value="{{datetimeFormat .event.Start}}"></span>
					<span><input name="end" type="text" id="event-End" class="form-control form-datetime" placeholder="Endet" value="{{datetimeFormat .event.End}}"></span>
					<label class="checkbox-inline" style="margin-left: 12px; padding-bottom: 9px"><input type="checkbox" name="rsvp" id="event-Rsvp" {{if .event.Rsvp}}checked{{end}}> Anmeldung erforderlich</label>
					<label class="checkbox-inline" style="margin-left: 12px; padding-bottom: 9px"><input type="checkbox" name="facebook" id="event-Facebook" {{if .event.FacebookId}}checked{{end}}> Auf Facebook teilen</label>
					<label class="checkbox-inline" style="margin-left: 12px; padding-bottom: 9px"><input type="checkbox" name="twitter" id="event-Twitter" {{if .event.TwitterId}}checked{{end}}> Auf Twitter teilen</label>
				</div>
				<div class="col-xs-4">
					<span id="event-thumbnail-message" class="help-block" style="text-align: center">Wähle ein Bild im Format jpg, jpeg, png oder gif aus.</span>
					<a id="event-dropzone" class="thumbnail" style="cursor: pointer">
						<span id="event-spinner" class="fa fa-gear"> </span>
						<img src="{{if .event.Image}}/bild/{{.event.Image}}?height=165&width=240{{else}}/images/thumbnail.png{{end}}" alt="Bild" id="event-thumbnail" class="img-responsive">
					</a>
					<input type="file" name="file" class="hide">
					<input type="hidden" name="image" id="event-Image" value="{{.event.Image}}">
				</div>
			</div>
			<div class="form-group">
				<div class="col-xs-12">
					<span class="help-block" style="margin-left: 10px; margin-top: 0;">Hinweis: Auf der Startseite werden nur Veranstaltungen mit Bild angezeigt.</span>
					<input name="credit" type="text" id="event-ImageCredit" class="form-control" placeholder="Bildrechte, falls das Bild nicht von Dir angefertigt wurde" value="{{.event.ImageCredit}}">
					{{if .user.IsAdmin}}<div style="display: inline-block; margin-left: 20px;">
						Veranstalter <select name="organizer" id="event-OrganizerId">
						{{range $id, $name := .organizers}}
						<option value="{{$id.Hex}}" {{if eq $id $.user.Id}}selected{{end}}>{{$name}}</option>
						{{end}}
						</select>
					</div>{{end}}
				</div>
			</div>
			<hr>
			<div class="form-group">
				<div class="col-xs-12" style="margin-left: 10px">
					<span id="event-Recurrency" class="help-block">Wiederholung</span>
					<fieldset>
						<label for="event-Recurrency-None" class="radio-inline"><input type="radio" id="event-Recurrency-None" name="recurrency" value="none" {{if eq .event.Recurrency 0}}checked{{end}}> Keine</label>
						<label for="event-Recurrency-Weekly" class="radio-inline"><input type="radio" id="event-Recurrency-Weekly" name="recurrency" value="weekly" {{if eq .event.Recurrency 1}}checked{{end}}> Wöchentlich</label>
						<label for="event-Recurrency-Monthly" class="radio-inline"><input type="radio" id="event-Recurrency-Monthly" name="recurrency" value="monthly" {{if eq .event.Recurrency 2}}checked{{end}}> Monatlich</label>
					</fieldset>
				</div>
				<div id="event-weekly" class="col-xs-12" style="margin-top: 10px; {{if not (eq .event.Recurrency 1)}}display: none{{end}}">
					<span class="checkbox-inline">jede <select name="weekly-interval" id="event-Recurrency-Weekly-Interval">
						<option value="1" {{if eq .event.Weekly.Interval 1}}selected{{end}}>1. Woche</option>
						<option value="2" {{if eq .event.Weekly.Interval 2}}selected{{end}}>2. Woche</option>
						<option value="3" {{if eq .event.Weekly.Interval 3}}selected{{end}}>3. Woche</option>
						<option value="4" {{if eq .event.Weekly.Interval 4}}selected{{end}}>4. Woche</option>
						<option value="5" {{if eq .event.Weekly.Interval 5}}selected{{end}}>5. Woche</option>
						<option value="6" {{if eq .event.Weekly.Interval 6}}selected{{end}}>6. Woche</option>
						<option value="7" {{if eq .event.Weekly.Interval 7}}selected{{end}}>7. Woche</option>
						<option value="8" {{if eq .event.Weekly.Interval 8}}selected{{end}}>8. Woche</option>
						<option value="9" {{if eq .event.Weekly.Interval 9}}selected{{end}}>9. Woche</option>
						<option value="10" {{if eq .event.Weekly.Interval 10}}selected{{end}}>10. Woche</option>
					</select>&nbsp;&nbsp;am&nbsp;&nbsp;</span>
					<label id="event-Recurrency-Weekly-Weekday" class="checkbox-inline"><input type="checkbox" name="event-Recurrency-Weekly-Weekday" value="1" {{range .event.Weekly.Weekdays}}{{if eq . 1}}checked{{end}}{{end}}> Mo&nbsp;&nbsp;</label>
					<label class="checkbox-inline"><input type="checkbox" name="event-Recurrency-Weekly-Weekday" value="2" {{range .event.Weekly.Weekdays}}{{if eq . 2}}checked{{end}}{{end}}> Di&nbsp;&nbsp;</label>
					<label class="checkbox-inline"><input type="checkbox" name="event-Recurrency-Weekly-Weekday" value="3" {{range .event.Weekly.Weekdays}}{{if eq . 3}}checked{{end}}{{end}}> Mi&nbsp;&nbsp;</label>
					<label class="checkbox-inline"><input type="checkbox" name="event-Recurrency-Weekly-Weekday" value="4" {{range .event.Weekly.Weekdays}}{{if eq . 4}}checked{{end}}{{end}}> Do&nbsp;&nbsp;</label>
					<label class="checkbox-inline"><input type="checkbox" name="event-Recurrency-Weekly-Weekday" value="5" {{range .event.Weekly.Weekdays}}{{if eq . 5}}checked{{end}}{{end}}> Fr&nbsp;&nbsp;</label>
					<label class="checkbox-inline"><input type="checkbox" name="event-Recurrency-Weekly-Weekday" value="6" {{range .event.Weekly.Weekdays}}{{if eq . 6}}checked{{end}}{{end}}> Sa&nbsp;&nbsp;</label>
					<label class="checkbox-inline"><input type="checkbox" name="event-Recurrency-Weekly-Weekday" value="0" {{range .event.Weekly.Weekdays}}{{if eq . 0}}checked{{end}}{{end}}> So&nbsp;&nbsp;</label>
				</div>
				<div id="event-monthly" class="col-xs-12" style="margin-top: 10px; {{if not (eq .event.Recurrency 2)}}display: none{{end}}">
					<span class="checkbox-inline">jeder <select name="monthly-week" id="event-Recurrency-Monthly-Week">
						<option value="0" {{if eq .event.Monthly.Week 0}}selected{{end}}>erste</option>
						<option value="1" {{if eq .event.Monthly.Week 1}}selected{{end}}>zweite</option>
						<option value="2" {{if eq .event.Monthly.Week 2}}selected{{end}}>dritte</option>
						<option value="3" {{if eq .event.Monthly.Week 3}}selected{{end}}>vierte</option>
						<option value="4" {{if eq .event.Monthly.Week 4}}selected{{end}}>letzte</option>
					</select>&nbsp;&nbsp;<select name="monthly-weekday" id="event-Recurrency-Monthly-Weekday">
						<option value="1" {{if eq .event.Monthly.Weekday 1}}selected{{end}}>Montag</option>
						<option value="2" {{if eq .event.Monthly.Weekday 2}}selected{{end}}>Dienstag</option>
						<option value="3" {{if eq .event.Monthly.Weekday 3}}selected{{end}}>Mittwoch</option>
						<option value="4" {{if eq .event.Monthly.Weekday 4}}selected{{end}}>Donnerstag</option>
						<option value="5" {{if eq .event.Monthly.Weekday 5}}selected{{end}}>Freitag</option>
						<option value="6" {{if eq .event.Monthly.Weekday 6}}selected{{end}}>Samstag</option>
						<option value="0" {{if eq .event.Monthly.Weekday 0}}selected{{end}}>Sonntag</option>
					</select>&nbsp;&nbsp;von jedem&nbsp;&nbsp;<select name="monthly-interval" id="event-Recurrency-Monthly-Interval">
						<option value="1" {{if eq .event.Monthly.Interval 1}}selected{{end}}>1. Monat</option>
						<option value="2" {{if eq .event.Monthly.Interval 2}}selected{{end}}>2. Monat</option>
						<option value="3" {{if eq .event.Monthly.Interval 3}}selected{{end}}>3. Monat</option>
						<option value="4" {{if eq .event.Monthly.Interval 4}}selected{{end}}>4. Monat</option>
						<option value="5" {{if eq .event.Monthly.Interval 5}}selected{{end}}>5. Monat</option>
						<option value="6" {{if eq .event.Monthly.Interval 6}}selected{{end}}>6. Monat</option>
						<option value="7" {{if eq .event.Monthly.Interval 7}}selected{{end}}>7. Monat</option>
						<option value="8" {{if eq .event.Monthly.Interval 8}}selected{{end}}>8. Monat</option>
						<option value="9" {{if eq .event.Monthly.Interval 9}}selected{{end}}>9. Monat</option>
						<option value="10" {{if eq .event.Monthly.Interval 10}}selected{{end}}>10. Monat</option>
						<option value="11" {{if eq .event.Monthly.Interval 11}}selected{{end}}>11. Monat</option>
						<option value="12" {{if eq .event.Monthly.Interval 12}}selected{{end}}>12. Monat</option>
					</select>
					</span>
				</div>
			</div>
			<div id="event-recurrencyEnd" class="form-group" style="{{if not (eq .event.Recurrency 1)}}display: none{{end}}">
				<div class="col-xs-7">
					<span><input name="recurrencyEnd" type="text" id="event-RecurrencyEnd" class="form-control form-datetime" placeholder="Wiederholung endet" value="{{datetimeFormat .event.RecurrencyEnd}}"></span>
				</div>
			</div>
			<hr>
			<div class="form-group">
				<div class="col-xs-12" style="margin-left: 10px">
					<span id="event-Target" class="help-block">Wähle eine oder mehrere Zielgruppen aus:</span>
				{{ range .targets }}
					{{ $id := index $.targetMap . }}
					<label class="checkbox-inline"><input type="checkbox" name="event-Target" value="{{$id}}"
					{{ range $.event.Targets }}
						{{ if eq . $id }}checked{{ end }}
					{{ end }}
					> {{.}} &nbsp;&nbsp;</label>
				{{ end }}
				</div>
			</div>
			<div class="form-group">
				<div class="col-xs-12" style="margin-left: 10px">
					<span id="event-Category" class="help-block">Wähle eine oder mehrere Kategorien aus:</span>
				{{ range .categories }}
					{{ $id := index $.categoryMap . }}
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
				<div class="col-xs-12">
					<div id="event-Descr" placeholder="Beschreibung">{{.event.HtmlDescription}}</div>
					<span><input name="website" type="text" id="event-Web" class="form-control" placeholder="Webseite" value="{{.event.Web}}"></span>
				</div>
			</div>
			<hr>
			<div class="form-group">
				<div class="col-xs-12">
					<span><input name="location" type="text" id="event-Name" class="form-control" placeholder="Wo findet die Veranstaltung statt?" value="{{.event.Addr.Name}}"></span>
				</div>
			</div>
			<div class="form-group">
				<div class="col-xs-5">
					<input name="street" type="text" id="event-Street" class="form-control" placeholder="Straße" value="{{.event.Addr.Street}}">
				</div>
				<div class="col-xs-3">
					<input name="pcode" type="text" id="event-Pcode" class="form-control" placeholder="Postleitzahl" value="{{.event.Addr.Pcode}}" maxlength="5">
				</div>
				<div class="col-xs-3">
					<input name="city" type="text" id="event-City" class="form-control" placeholder="Ort" value="{{.event.Addr.City}}">
				</div>
			</div>
			<hr>
			<div class="form-group">
				<div class="col-sm-4 col-xs-5">
					<a href="/veranstalter/verwaltung/0" onClick="history.back();return false" class="btn btn-default btn-cancel" style="width: 90%">Abbrechen</a>
				</div>
				<div class="col-sm-1 hidden-xs">&nbsp;</div>
				<div class="col-xs-7">
					<button type="submit" class="btn btn-mmr" style="width: 90%">Speichern</button>
				</div>
			</div>
		</form>
	</div>
</div>
