<!DOCTYPE html><html lang="de">
<head>
	<meta http-equiv="X-UA-Compatible" content="IE=edge">
	<meta charset="utf-8">
	<meta name="viewport" content="width=device-width, initial-scale=0.75">
	<title>{{.meta.Title}}</title>
	<meta name="description" content="{{.meta.Descr}}">
	<meta property="og:title" content="{{.meta.FB_Title}}">
	<meta property="og:site_name" content="Mitmach-Republik">
	<meta property="og:image" content="{{.meta.FB_Image}}">
	<meta property="og:description" content="{{.meta.FB_Descr}}">
{{if .noindex}}
	<meta name="robots" content="noindex, follow">
{{end}}
{{if and (.event) (not .user)}}
	<meta property="og:url" content="http://{{.hostname}}{{.event.Url}}">
	<link rel="canonical" href="http://{{.hostname}}{{.event.Url}}">
{{end}}
	<link rel="shortcut icon" href="/favicon.ico" type="image/x-icon">
	<link rel="icon" href="/favicon.ico" type="image/x-icon">
	{{if .meta.RSS}}
		<link rel="alternate" type="application/rss+xml" title="RSS" href="?fmt=RSS">
	{{end}}
	<link href="http://fonts.googleapis.com/css?family=Open+Sans:400,300,700,600" rel="stylesheet" type="text/css">
	<link href="/css/styles-17.css" rel="stylesheet">
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
{{if and (.event) (not .user)}}
<script>(function(d, s, id) {
  var js, fjs = d.getElementsByTagName(s)[0];
  if (d.getElementById(id)) return;
  js = d.createElement(s); js.id = id;
  js.src = "//connect.facebook.net/de_DE/sdk.js#xfbml=1&version=v2.3";
  fjs.parentNode.insertBefore(js, fjs);
}(document, 'script', 'facebook-jssdk'));</script>
{{end}}
	<div class="modal fade" id="login" tabindex="-2" role="dialog" aria-labelledby="login-dialog" aria-hidden="true">
		<div class="modal-dialog">
			<div class="modal-content">
			</div>
		</div>
	</div>
	<div class="modal fade" id="register" tabindex="-3" role="dialog" aria-labelledby="register-dialog" aria-hidden="true">
		<div class="modal-dialog">
			<div class="modal-content">
			</div>
		</div>
	</div>
	<div class="modal fade" id="registered" tabindex="-4" role="dialog" aria-labelledby="registered-dialog" aria-hidden="true">
	</div>
	<div class="modal fade" id="mail" tabindex="-5" role="dialog" aria-labelledby="email-dialog" aria-hidden="true">
		<div class="modal-dialog">
			<div class="modal-content">
			</div>
		</div>
	</div>
	<div class="modal fade" id="email-alert" tabindex="-5" role="dialog" aria-labelledby="email-alert-dialog" aria-hidden="true">
		<div class="modal-dialog">
			<div class="modal-content">
			</div>
		</div>
	</div>
	<div class="modal fade" id="send-password" tabindex="-7" role="dialog" aria-labelledby="password-dialog" aria-hidden="true">
		<div class="modal-dialog">
			<div class="modal-content">
			</div>
		</div>
	</div>
{{if .event}}
	<div class="modal fade" id="share" tabindex="-6" role="dialog" aria-labelledby="email-dialog" aria-hidden="true">
		<div class="modal-dialog">
			<div class="modal-content">
			</div>
		</div>
	</div>
{{end}}
	<div class="container">
		<div class="row">
			<div class="col-lg-3 col-sm-3 col-xs-6 col-logo">
				<span class="logo-helper"></span><a href="/" title="Zur Startseite"><img src="/images/mitmachrepublik.png" style="width: 100%" alt="Logo"></a>
			</div>
			<div class="col-lg-2 visible-lg-block col-md-1 visible-md-block">&nbsp;</div>
			<div class="col-lg-3 col-sm-4 col-xs-6 col-head">
				<span id="head-organizer">Du bist ein Organisator?</span><br>
				<a id="head-events" href="javascript:void(0)" data-href="/dialog/login" title="Melde Dich an, um Deine Veranstaltungen einzutragen." rel="nofollow" data-toggle="modal" data-target="#login" class="highlight"><span class="fa fa-caret-right"></span> Finde kostenlos Mitmacher</a>
			</div>
			<div class="col-lg-1 visible-lg-block">&nbsp;</div>
			<div class="col-lg-3 col-md-4 col-sm-5 col-xs-12 col-head">
				<a class="highlight" title="Like uns auf Facebook." href="https://www.facebook.com/mitmachrepublik" target="_blank"><span class="fa fa-facebook"></span></a>
				<a class="highlight" title="Folge uns auf twitter." href="https://twitter.com/mitmachrepublik" target="_blank"><span class="fa fa-twitter"></span></a>
				| <a href="javascript:void(0)" data-href="/dialog/contact" rel="nofollow" title="Schreibe uns eine Nachricht." data-toggle="modal" data-target="#mail">Kontakt</a>
				| <a id="head-login" href="javascript:void(0)" data-href="/dialog/login" rel="nofollow" data-toggle="modal" data-target="#login" title="Melde Dich an, um Deine Veranstaltungen einzutragen.">Login</a>
				| <a href="/wir-ueber-uns" title="Wer wir sind und warum wir das machen.">Wir über uns</a><br>
				<form id="fulltextsearch" role="form" action="/suche" method="POST" class="form-inline"><input class="form-control form-search" style="width: 83%" name="fulltextsearch" placeholder="Veranstaltungen suchen" autocomplete="off"> <button type="submit" class="btn btn-mmr form-search" style="margin-left: 5px" name="search" title="Veranstaltungen suchen" value="query"><span class="fa fa-search"></span></button></form>
			</div>
		</div>
