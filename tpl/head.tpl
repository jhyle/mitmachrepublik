<!DOCTYPE html><html lang="de">
<head>
	<meta charset="utf-8">
	<meta http-equiv="X-UA-Compatible" content="IE=edge">
	<meta name="viewport" content="width=1170">
	<title>{{.meta.Title}}</title>
	<meta name="og:title" content="{{.meta.FB_Title}}">
	<meta name="og:site_name" content="Mitmach-Republik">
	<meta name="og:image" content="{{.meta.FB_Image}}">
	<meta name="og:description" content="{{.meta.FB_Descr}}">
	<link rel="shortcut icon" href="/favicon.ico" type="image/x-icon">
	<link rel="icon" href="/favicon.ico" type="image/x-icon">
	<link href="http://fonts.googleapis.com/css?family=Open+Sans:400,300,700,600" rel="stylesheet" type="text/css">
	<link href="/css/styles-1.css" rel="stylesheet">
	<link href="//maxcdn.bootstrapcdn.com/font-awesome/4.2.0/css/font-awesome.min.css" rel="stylesheet">
	<!-- HTML5 Shim and Respond.js IE8 support of HTML5 elements and media queries -->
	<!-- WARNING: Respond.js doesn't work if you view the page via file:// -->
	<!--[if lt IE 9]>
		<script src="https://oss.maxcdn.com/html5shiv/3.7.2/html5shiv.min.js"></script>
		<script src="https://oss.maxcdn.com/respond/1.4.2/respond.min.js"></script>
	<![endif]-->
</head>
<body>
<script>
  (function(i,s,o,g,r,a,m){i['GoogleAnalyticsObject']=r;i[r]=i[r]||function(){
  (i[r].q=i[r].q||[]).push(arguments)},i[r].l=1*new Date();a=s.createElement(o),
  m=s.getElementsByTagName(o)[0];a.async=1;a.src=g;m.parentNode.insertBefore(a,m)
  })(window,document,'script','//www.google-analytics.com/analytics.js','ga');
  ga('create', '{{.ga_code}}', 'auto');
  ga('send', 'pageview');
</script>
<div id="fb-root"></div>
<script>(function(d, s, id) {
  var js, fjs = d.getElementsByTagName(s)[0];
  if (d.getElementById(id)) return;
  js = d.createElement(s); js.id = id;
  js.src = "//connect.facebook.net/de_DE/sdk.js#xfbml=1&version=v2.3";
  fjs.parentNode.insertBefore(js, fjs);
}(document, 'script', 'facebook-jssdk'));</script>
	<div class="modal fade" id="login" tabindex="-2" role="dialog" aria-labelledby="login-dialog" aria-hidden="true">
		<div class="modal-dialog">
			<div class="modal-content">
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
								<button name="register" type="button" data-dismiss="modal" data-toggle="modal" data-target="#register" class="btn btn-mmr" style="margin-top: 18px; width: 90%">Kostenlos registrieren</button>
							</div>
						</div>
					</form>
				</div>
			</div>
		</div>
	</div>
	<div class="modal fade" id="register" tabindex="-3" role="dialog" aria-labelledby="login-dialog" aria-hidden="true">
		<div class="modal-dialog">
			<div class="modal-content">
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
			</div>
		</div>
	</div>
	<div class="modal fade" id="registered" tabindex="-4" role="dialog" aria-labelledby="login-dialog" aria-hidden="true">
		<div class="modal-dialog">
			<div class="modal-content">
				<div class="modal-header">
					<button type="button" class="close" data-dismiss="modal"><span aria-hidden="true">&times;</span><span class="sr-only">Schließen</span></button>
				</div>
				<div class="modal-body">
					<div class="big-text">
						Danke für Deine Registrierung. Zur Überprüfung Deiner E-Mail-Adresse haben wir Dir eine E-Mail mit einem Aktivierungslink zugesendet.
						Von Dir eingegebene Veranstaltungen werden erst nach der Aktivierung Deines Profils sichtbar.<br />
						Bitte hab Verständnis dafür, dass wir noch in der Testphase sind. Schreib uns eine Nachricht, wenn etwas nicht so funktioniert wie erwartet.
					</div>
					<div class="form-group">
						<div class="col-sm-4">
							<a href="/" class="btn btn-default" style="width: 90%">Schließen</a>
						</div>
						<div class="col-sm-1">&nbsp;</div>
						<div class="col-sm-7">
							<a href="/veranstalter/verwaltung/0" class="btn btn-mmr" style="width: 90%">Veranstaltung eintragen</a>
						</div>
					</div>
					<div class="clearfix"></div>
				</div>
			</div>
		</div>
	</div>
	<div class="modal fade" id="mail" tabindex="-5" role="dialog" aria-labelledby="email-dialog" aria-hidden="true">
		<div class="modal-dialog">
			<div class="modal-content">
				<div class="modal-header">
					<button type="button" class="close" data-dismiss="modal"><span aria-hidden="true">&times;</span><span class="sr-only">Schließen</span></button>
				</div>
				<div class="modal-body">
					<div class="big-text">Schreib uns eine Nachricht.</div>
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
			</div>
		</div>
	</div>
{{if .event}}
	<div class="modal fade" id="share" tabindex="-5" role="dialog" aria-labelledby="email-dialog" aria-hidden="true">
		<div class="modal-dialog">
			<div class="modal-content">
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

die Veranstaltung {{.event.Title}} in {{citypartName .event.Addr}} finde ich interessant, schau doch mal rein: http://{{$.hostname}}{{eventUrl .event | encodePath}}.

Liebe Grüße! 
								</textarea>
							</div>
						</div>
						<hr>
						<div class="form-group">
							<div class="col-sm-4">
								<button type="button" class="btn btn-default" data-dismiss="modal" style="width: 90%">Abbrechen</button>
							</div>
							<div class="col-sm-1">&nbsp;</div>
							<div class="col-sm-7">
								<button id="send-event-submit" type="submit" class="btn btn-mmr" data-loading-text="Senden.." style="width: 90%">Senden</button>
							</div>
						</div>
					</form>
				</div>
			</div>
		</div>
	</div>
{{end}}
	<div class="container">
		<div class="row">
			<div class="col-xs-1">&nbsp;</div>
			<div class="col-xs-3 col-head"><a href="/" title="MitmachRepublik"><img src="/images/mitmachrepublik.png" style="max-width:80%" alt="MitmachRepublik"/></a></div>
			<div class="col-xs-1">&nbsp;</div>
			<div class="col-xs-3 col-head"><span id="head-organizer">Du bist ein Organisator?</span><br /><a id="head-events" href="#" data-toggle="modal" data-target="#login" class="highlight"><span class="fa fa-caret-right"></span> Trage Deine Veranstaltungen ein.</a></div>
			<div class="col-xs-1">&nbsp;</div>
			<div class="col-xs-2 col-head"><a id="head-login" href="#" data-toggle="modal" data-target="#login"><span class="fa fa-user highlight"></span> Anmelden</a><br /><!-- a href="#">Über uns</a--><a href="https://www.facebook.com/mitmachrepublik" target="_blank"><span class="fa fa-facebook"></span></a> <a href="https://twitter.com/mitmachrepublik" target="_blank"><span class="fa fa-twitter" ></span></a> | <a href="#" class="highlight" data-toggle="modal" data-target="#mail">Kontakt</a></div>
			<div class="col-xs-1">&nbsp;</div>
		</div>