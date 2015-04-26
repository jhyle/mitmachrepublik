<!DOCTYPE html><html lang="de">
<head>
	<meta charset="utf-8">
	<meta http-equiv="X-UA-Compatible" content="IE=edge">
	<meta name="viewport" content="width=1170">
	<title>{{.meta.Title}}</title>
	<meta name="description" content="{{.meta.FB_Descr}}">
	<meta name="og:title" content="{{.meta.FB_Title}}">
	<meta name="og:site_name" content="Mitmach-Republik">
	<meta name="og:image" content="{{.meta.FB_Image}}">
	<meta name="og:description" content="{{.meta.FB_Descr}}">
	<link rel="shortcut icon" href="/favicon.ico" type="image/x-icon">
	<link rel="icon" href="/favicon.ico" type="image/x-icon">
	<link href="http://fonts.googleapis.com/css?family=Open+Sans:400,300,700,600" rel="stylesheet" type="text/css">
	<link href="/css/styles-2.css" rel="stylesheet">
	<link href="//maxcdn.bootstrapcdn.com/font-awesome/4.2.0/css/font-awesome.min.css" rel="stylesheet">
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
			</div>
		</div>
	</div>
	<div class="modal fade" id="register" tabindex="-3" role="dialog" aria-labelledby="login-dialog" aria-hidden="true">
		<div class="modal-dialog">
			<div class="modal-content">
			</div>
		</div>
	</div>
	<div class="modal fade" id="registered" tabindex="-4" role="dialog" aria-labelledby="login-dialog" aria-hidden="true">
	</div>
	<div class="modal fade" id="mail" tabindex="-5" role="dialog" aria-labelledby="email-dialog" aria-hidden="true">
		<div class="modal-dialog">
			<div class="modal-content">
			</div>
		</div>
	</div>
{{if .event}}
	<div class="modal fade" id="share" tabindex="-5" role="dialog" aria-labelledby="email-dialog" aria-hidden="true">
		<div class="modal-dialog">
			<div class="modal-content">
			</div>
		</div>
	</div>
{{end}}
	<div class="container">
		<div class="row">
			<div class="col-xs-1">&nbsp;</div>
			<div class="col-xs-3 col-head"><a href="/" title="Mitmach-Republik"><img src="/images/mitmachrepublik.png" style="max-width:80%" alt="Mitmach-Republik"/></a></div>
			<div class="col-xs-1">&nbsp;</div>
			<div class="col-xs-3 col-head"><span id="head-organizer">Du bist ein Organisator?</span><br /><a id="head-events" href="/dialog/login" rel="nofollow" data-toggle="modal" data-target="#login" class="highlight"><span class="fa fa-caret-right"></span> Trage Deine Veranstaltungen ein.</a></div>
			<div class="col-xs-1">&nbsp;</div>
			<div class="col-xs-2 col-head"><a id="head-login" href="/dialog/login" rel="nofollow" data-toggle="modal" data-target="#login"><span class="fa fa-user highlight"></span> Anmelden</a><br /><!-- a href="#">Über uns</a--><a class="highlight" title="Like uns auf Facebook." href="https://www.facebook.com/mitmachrepublik" target="_blank"><span class="fa fa-facebook"></span></a> <a class="highlight" title="Folge uns auf twitter." href="https://twitter.com/mitmachrepublik" target="_blank"><span class="fa fa-twitter"></span></a> | <a href="/dialog/contact" rel="nofollow" title="Schreibe uns eine Nachricht." data-toggle="modal" data-target="#mail">Kontakt</a></div>
			<div class="col-xs-1">&nbsp;</div>
		</div>