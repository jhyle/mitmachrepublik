<!DOCTYPE html><html lang="en">
<head>
	<meta charset="utf-8">
	<meta http-equiv="X-UA-Compatible" content="IE=edge">
	<meta name="viewport" content="width=device-width, initial-scale=1">
	<title>{{.title}}</title>
	<link href="http://fonts.googleapis.com/css?family=Open+Sans:400,300,700,600" rel="stylesheet" type="text/css">
	<link href="/css/bootstrap.min.css" rel="stylesheet">
	<link href="/css/bootstrap-datetimepicker.min.css" rel="stylesheet" type="text/css">
	<link href="/css/mmr.css" rel="stylesheet" type="text/css">
	<link href="//maxcdn.bootstrapcdn.com/font-awesome/4.2.0/css/font-awesome.min.css" rel="stylesheet">
	<!-- HTML5 Shim and Respond.js IE8 support of HTML5 elements and media queries -->
	<!-- WARNING: Respond.js doesn't work if you view the page via file:// -->
	<!--[if lt IE 9]>
		<script src="https://oss.maxcdn.com/html5shiv/3.7.2/html5shiv.min.js"></script>
		<script src="https://oss.maxcdn.com/respond/1.4.2/respond.min.js"></script>
	<![endif]-->
</head>
<body>
	<div class="modal" id="login" tabindex="-2" role="dialog" aria-labelledby="login-dialog" aria-hidden="true">
		<div class="modal-dialog">
			<div class="modal-content">
				<div class="modal-header">
					<button type="button" class="close" data-dismiss="modal"><span aria-hidden="true">&times;</span><span class="sr-only">Schließen</span></button>
				</div>
				<div class="modal-body">
					<form role="form" id="login-form" class="form-horizontal">
						<div class="form-group">
							<div class="col-sm-7" style="border-right: 1px solid #ccc; padding-right: 0">
								<div class="big-text">Ich bin bereits als Veranstalter registriert.</div>
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
	<div class="modal" id="register" tabindex="-3" role="dialog" aria-labelledby="login-dialog" aria-hidden="true">
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
									<img src="/images/thumbnail.gif" alt="Bild" id="register-thumbnail" class="img-responsive">
								</a>
								<span id="register-thumbnail-message" class="help-block">Wähle ein Bild im Format jpg, jpeg, png oder gif aus.</span>
								<input type="file" name="file" class="hide">
								<input type="hidden" name="image" id="register-Image">
							</div>
						</div>
						<hr>
						<div class="form-group">
							<div class="col-sm-12">
								<textarea name="description" id="register-Descr" class="form-control" placeholder="Beschreibung"></textarea>
								<span><input name="website" type="text" id="register-Web" class="form-control" placeholder="Webseite"></span>
							</div>
						</div>
						<hr>
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
	<div class="modal" id="registered" tabindex="-4" role="dialog" aria-labelledby="login-dialog" aria-hidden="true">
		<div class="modal-dialog">
			<div class="modal-content">
				<div class="modal-header">
					<button type="button" class="close" data-dismiss="modal"><span aria-hidden="true">&times;</span><span class="sr-only">Schließen</span></button>
				</div>
				<div class="modal-body">
					<div class="big-text">
						Danke für Deine Registrierung. Zur Überprüfung Deiner E-Mail-Adresse haben wir Dir eine E-Mail mit einem Aktivierungslink zugesendet.
						Von Dir eingegebene Veranstaltungen werden erst nach der Aktivierung Deiner Anmeldung sichtbar.
					</div>
					<div class="form-group">
						<div class="col-sm-4">
							<a href="/" class="btn btn-default" style="width: 90%">Schließen</a>
						</div>
						<div class="col-sm-1">&nbsp;</div>
						<div class="col-sm-7">
							<a href="/veranstalter/verwaltung/veranstaltung" class="btn btn-mmr" style="width: 90%">Veranstaltung eintragen</a>
						</div>
					</div>
					<div class="clearfix"></div>
				</div>
			</div>
		</div>
	</div>
	<div class="container">
		<div class="row">
			<div class="col-xs-1">&nbsp;</div>
			<div class="col-xs-3 col-head"><a href="/" title="MitmachRepublik"><img src="/images/mitmachrepublik.gif" style="max-width:80%" alt="MitmachRepublik"/></a></div>
			<div class="col-xs-1">&nbsp;</div>
			<div class="col-xs-3 col-head"><span id="head-organizer">Du bist Veranstalter?</span><br /><a id="head-events" href="#" data-toggle="modal" data-target="#login" class="highlight"><span class="fa fa-caret-right"></span> Trage Deine Veranstaltungen ein.</a></div>
			<div class="col-xs-1">&nbsp;</div>
			<div class="col-xs-2 col-head"><a id="head-login" href="#" data-toggle="modal" data-target="#login"><span class="fa fa-user highlight"></span> Anmelden</a><br /><a href="#">Über uns</a> | <a href="#">Kontakt</a></div>
			<div class="col-xs-1">&nbsp;</div>
		</div>