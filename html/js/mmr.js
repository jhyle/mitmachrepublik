var PcodePattern = /^[0-9]{5}$/;
var DateTimePattern = /^(\d{2})\.(\d{2})\.(\d{4}) (\d{2})\:(\d{2})$/;
var EmailPattern = /^([a-zA-Z0-9_.+-])+\@(([a-zA-Z0-9-])+\.)+([a-zA-Z0-9]{2,4})+$/;
var WebPattern = /(http|ftp|https):\/\/[\w-]+(\.[\w-]+)+([\w.,@?^=%&amp;:\/~+#-]*[\w@?^=%&amp;\/~+#-])?/;

function zeroFill( number, width )
{
	width -= number.toString().length;
	if ( width > 0 )
	{
		return new Array( width + (/\./.test( number ) ? 2 : 1) ).join( '0' ) + number;
	}
	return number + ""; // always return a string
}

function str2date(s)
{
	DateTimePattern.exec(s);
	var year = RegExp.$3;
	var month = RegExp.$2;
	var day = RegExp.$1;
	var hour = RegExp.$4;
	var minute = RegExp.$5;
	return new Date(parseInt(year), parseInt(month) - 1, parseInt(day), parseInt(hour), parseInt(minute), 0);
}

function initProfileForm(id)
{
	$('#' + id + '-Descr').summernote({
		height: 300,
		lang: 'de-DE'
	});
	
	$("#" + id + "-dropzone").click(function() {
		$(this).parent().find("input").click();
	});

	$("#" + id + "-Name").popover({content: "Bitte gib den Namen Deiner Organisation ein.", trigger: "manual", placement: "right"});
	$("#" + id + "-Name").focus(function () { $("#" + id + "-Name").popover('hide'); });
	if ($("#" + id + "-Email").length) {
		$("#" + id + "-Email").popover({content: "Bitte gib eine gültige E-Mail-Adresse ein.", trigger: "manual", placement: "right"});
		$("#" + id + "-Email").focus(function () { $("#" + id + "-Email").popover('hide'); });
	}
	if ($("#" + id + "-Pwd").length) {
		$("#" + id + "-Pwd").popover({content: "Bitte gib ein Kennwort ein.", trigger: "manual", placement: "right"});
		$("#" + id + "-Pwd").focus(function () { $("#" + id + "-Pwd").popover('hide'); });
		$("#" + id + "-Pwd2").popover({content: "Kennwort und Kennwortwiederholung stimmen nicht überein.", trigger: "manual", placement: "right"});
		$("#" + id + "-Pwd2").focus(function () { $("#" + id + "-Pwd2").popover('hide'); });
	}
	$("#" + id + "-Web").popover({content: "Bitte gib eine gültige Web-Adresse (mit http://) ein.", trigger: "manual", placement: "auto"});
	$("#" + id + "-Web").focus(function () { $("#" + id + "-Web").popover('hide'); });
	$("#" + id + "-Pcode").popover({content: "Bitte gib eine vollständige Postleitzahl ein.", trigger: "manual", placement: "auto top"});	
	$("#" + id + "-Pcode").focus(function () { $("#" + id + "-Pcode").popover('hide'); });
}

function initEmailAndPwdForm(id)
{
	$("#" + id + "-Email").popover({content: "Bitte gib eine gültige E-Mail-Adresse ein.", trigger: "manual", placement: "auto right"});
	$("#" + id + "-Email").focus(function () { $("#" + id + "-Email").popover('hide'); });
	$("#" + id + "-Pwd2").popover({content: "Kennwort und Kennwortwiederholung stimmen nicht überein.", trigger: "manual", placement: "auto right"});
	$("#" + id + "-Pwd2").focus(function () { $("#" + id + "-Pwd2").popover('hide'); });
}

function initEventForm(id)
{
	$('#' + id + '-Descr').summernote({
		height: 300,
		lang: 'de-DE'
	});
	
	$("#" + id + "-dropzone").click(function() {
		$(this).parent().find("input").click();
	});

	$("#" + id + "-Title").popover({content: "Bitte gib der Veranstaltung einen Namen.", trigger: "manual", placement: "auto right"});
	$("#" + id + "-Title").focus(function () { $("#" + id + "-Title").popover('hide'); });
	$("#" + id + "-Title-Too-Long").popover({content: "Der Veranstaltungsname ist zu lang. Bitte kürze ihn auf 40 Zeichen.", trigger: "manual", placement: "auto right"});
	$("#" + id + "-Title-Too-Long").focus(function () { $("#" + id + "-Title-Too-Long").popover('hide'); });
	$("#" + id + "-Start").popover({content: "Bitte gib den Veranstaltungsbeginn an.", trigger: "manual", placement: "auto right"});
	$("#" + id + "-Start").focus(function () { $("#" + id + "-Start").popover('hide'); });
	$("#" + id + "-End").popover({content: "Bitte gib ein gültiges Ende an.", trigger: "manual", placement: "auto right"});
	$("#" + id + "-End").focus(function () { $("#" + id + "-End").popover('hide'); });
	$("#" + id + "-Target").popover({content: "Bitte wähle mindestens eine Zielgruppe aus.", trigger: "manual", placement: "auto right"});
	$("#" + id + "-Target").focus(function () { $("#" + id + "-Target").popover('hide'); });
	$("#" + id + "-Recurrency-Weekly-Weekday").popover({content: "Bitte wähle mindestens einen Wochentag aus.", trigger: "manual", placement: "auto right"});
	$("#" + id + "-Recurrency-Weekly-Weekday").focus(function () { $("#" + id + "--Recurrency-Weekly-Weekday").popover('hide'); });
	$("#" + id + "-Category").popover({content: "Bitte wähle mindestens eine Kategorie aus.", trigger: "manual", placement: "auto right"});
	$("#" + id + "-Category").focus(function () { $("#" + id + "-Category").popover('hide'); });

	$("#" + id + "-Web").popover({content: "Bitte gib eine gültige Web-Adresse (mit http://) ein.", trigger: "manual", placement: "auto right"});
	$("#" + id + "-Web").focus(function () { $("#" + id + "-Web").popover('hide'); });
	$("#" + id + "-Pcode").popover({content: "Bitte gib eine vollständige Postleitzahl ein.", trigger: "manual", placement: "auto top"});	
	$("#" + id + "-Pcode").focus(function () { $("#" + id + "-Pcode").popover('hide'); });
}

function initEmailAlertForm(id)
{
	$("#" + id + "-Email").popover({content: "Bitte gib eine gültige E-Mail-Adresse ein.", trigger: "manual", placement: "auto"});
	$("#" + id + "-Email").focus(function () { $("#" + id + "-Email").popover('hide'); });
	$("#" + id + "-Weekday").popover({content: "Bitte wähle mindestens einen Wochentag aus.", trigger: "manual", placement: "auto"});
	$("#" + id + "-Weekday").focus(function () { $("#" + id + "-Weekday").popover('hide'); });
}

function initSendMailForm(id)
{
	$("#" + id + "-Email").popover({content: "Bitte gib eine gültige E-Mail-Adresse ein.", trigger: "manual", placement: "auto"});
	$("#" + id + "-Email").focus(function () { $("#" + id + "-Email").popover('hide'); });
	$("#" + id + "-Subject").popover({content: "Bitte gib der Nachricht einen Betreff.", trigger: "manual", placement: "auto"});
	$("#" + id + "-Subject").focus(function () { $("#" + id + "-Subject").popover('hide'); });
}

function validate(ok, id) {

	if (!ok) {
		$(id).parent().addClass("has-error");
		$(id).popover("show");
	} else {
		$(id).parent().removeClass("has-error");
		$(id).popover("hide");
	}
	
	return ok;
}

function validateProfileForm(id)
{
	var ok = validate($("#" + id + "-Name").val().trim().length > 0, "#" + id +"-Name");
	if ($("#" + id +"-Email").length) {
		ok &= validate(EmailPattern.test($("#" + id +"-Email").val()), "#" + id +"-Email");
	}
	if ($("#" + id + "-Pwd").length) {
		if (validate($("#" + id +"-Pwd").val().trim().length > 0, "#" + id +"-Pwd")) {
			ok &= validate($("#" + id +"-Pwd").val() == $("#" + id +"-Pwd2").val(), "#" + id +"-Pwd2");
		} else {
			ok = false;
		}
	}
	ok &= validate($("#" + id +"-Web").val().trim().length == 0 || WebPattern.test($("#" + id +"-Web").val()), "#" + id +"-Web");
	ok &= validate($("#" + id +"-Pcode").val().trim().length == 0 || PcodePattern.test($("#" + id +"-Pcode").val()), "#" + id +"-Pcode");
	return ok;
}

function validateEmailAndPwdForm(id)
{
	var ok = validate(EmailPattern.test($("#" + id +"-Email").val()), "#" + id +"-Email");
	if ($("#" + id +"-Pwd").val().trim().length > 0) {
		ok &= validate($("#" + id +"-Pwd").val() == $("#" + id +"-Pwd2").val(), "#" + id +"-Pwd2");
	} else {
		validate(true, "#" + id +"-Pwd2");
	}
	return ok;
}

function validateEventForm(id)
{
	var ok = validate($("#" + id + "-Title").val().trim().length > 0, "#" + id +"-Title");
	var ok = validate($("#" + id + "-Title").val().trim().length <= 40, "#" + id +"-Title-Too-Long");
	ok &= validate(DateTimePattern.test($("#" + id + "-Start").val()), "#" + id + "-Start");
	ok &= validate($("#" + id +"-End").val().trim().length == 0 || DateTimePattern.test($("#" + id + "-End").val()), "#" + id + "-End");
	ok &= validate($("#" + id +"-Web").val().trim().length == 0 || WebPattern.test($("#" + id +"-Web").val()), "#" + id +"-Web");
	ok &= validate($("#" + id +"-Pcode").val().trim().length == 0 || PcodePattern.test($("#" + id +"-Pcode").val()), "#" + id +"-Pcode");
	ok &= validate(!$("#" + id + "-Recurrency-Weekly").is(":checked") || $("input[name=" + id + "-Recurrency-Weekly-Weekday]:checked").map(function () {return this.value;}).get().length > 0, "#" + id + "-Recurrency-Weekly-Weekday");
	ok &= validate($("input[name=" + id + "-Target]:checked").map(function () {return this.value;}).get().length > 0, "#" + id + "-Target");
	ok &= validate($("input[name=" + id + "-Category]:checked").map(function () {return this.value;}).get().length > 0, "#" + id + "-Category");
	return ok;
}

function validateEmailAlertForm(id)
{
	var ok = validate(EmailPattern.test($("#" + id +"-Email").val()), "#" + id +"-Email");
	ok &= validate($("input[name=" + id + "-Weekday]:checked").map(function () {return this.value;}).get().length > 0, "#" + id + "-Weekday");
	return ok;
}

function validateSendMailForm(id)
{
	var ok = validate(EmailPattern.test($("#" + id +"-Email").val()), "#" + id +"-Email");
	ok &= validate($("#" + id + "-Subject").val().trim().length > 0, "#" + id +"-Subject");
	return ok;
}

function gatherProfileForm(id)
{
	var data = {};
	var user_fields = ["Name", "Email", "Pwd", "Image", "ImageCredit", "Web"];
	for (var i = 0, len = user_fields.length; i < len; i++) {
		if ($("#" + id + "-" + user_fields[i]).length) {
			data[user_fields[i]] = $("#" + id + "-" + user_fields[i]).val();
		}
	}

	data["Descr"] = $("#" + id + "-Descr").code();
	
	data["Addr"] = {}
	var addr_fields = ["Street", "Pcode", "City"];
	for (var i = 0, len = addr_fields.length; i < len; i++) {
		if ($("#" + id + "-" + addr_fields[i]).length) {
			data["Addr"][addr_fields[i]] = $("#" + id + "-" + addr_fields[i]).val();
		}
	}
	data["Categories"] = new Array();
	var categories = $("input[name=" + id + "-Category]:checked").map(function () {return this.value;}).get();
	for (i = 0; i < categories.length; i++) {
		data["Categories"][i] = parseInt(categories[i]);
	}
	return data;
}

function gatherEventForm(id)
{
	var data = {};
	var event_fields = ["Id", "Title", "Start", "End", "RecurrencyEnd", "Image", "ImageCredit", "OrganizerId", "Web"];
	for (var i = 0, len = event_fields.length; i < len; i++) {
		if ($("#" + id + "-" + event_fields[i]).length) {
			data[event_fields[i]] = $("#" + id + "-" + event_fields[i]).val();
		}
	}
	
	data["Descr"] = $("#" + id + "-Descr").code();
	
	if (data["Start"].length > 0) {
		data["Start"] = str2date(data["Start"]);
	} else {
		delete data["Start"];
	}
	if (data["End"].length > 0) {
		data["End"] = str2date(data["End"]);
	} else {
		delete data["End"];
	}
	if (data["RecurrencyEnd"].length > 0) {
		data["RecurrencyEnd"] = str2date(data["RecurrencyEnd"]);
	} else {
		delete data["RecurrencyEnd"];
	}

	data["Rsvp"] = $("#" + id + "-Rsvp").is(':checked');	
	
	var recurrency = $("input[type=radio][name=recurrency]:checked").val();
	if (recurrency == "weekly") {
		data["Recurrency"] = 1;
		data["Weekly"] = {
			"Interval": parseInt($("#" + id + "-Recurrency-Weekly-Interval").val()),
			"Weekdays": []
		};
		var weekdays = $("input[name=" + id + "-Recurrency-Weekly-Weekday]:checked").map(function () {return this.value;}).get();
		for (i = 0; i < weekdays.length; i++) {
			data["Weekly"]["Weekdays"][i] = parseInt(weekdays[i]);
		}
	} else if (recurrency == "monthly") {
		data["Recurrency"] = 2;
		data["Monthly"] = {
			"Interval": parseInt($("#" + id + "-Recurrency-Monthly-Interval").val()),
			"Week": parseInt($("#" + id + "-Recurrency-Monthly-Week").val()),
			"Weekday": parseInt($("#" + id + "-Recurrency-Monthly-Weekday").val())
		};
	}
	
	data["Addr"] = {}
	var addr_fields = ["Name", "Street", "Pcode", "City"];
	for (var i = 0, len = addr_fields.length; i < len; i++) {
		if ($("#" + id + "-" + addr_fields[i]).length) {
			data["Addr"][addr_fields[i]] = $("#" + id + "-" + addr_fields[i]).val();
		}
	}
	
	data["Targets"] = new Array();
	var targets = $("input[name=" + id + "-Target]:checked").map(function () {return this.value;}).get();
	for (i = 0; i < targets.length; i++) {
		data["Targets"][i] = parseInt(targets[i]);
	}
	data["Categories"] = new Array();
	var categories = $("input[name=" + id + "-Category]:checked").map(function () {return this.value;}).get();
	for (i = 0; i < categories.length; i++) {
		data["Categories"][i] = parseInt(categories[i]);
	}
	
	return data;
}

function gatherEmailAlertForm(id) {
	
	var data = {};
	var alert_fields = ["Name", "Email", "Query", "Place", "Targets", "Categories", "Dates", "Radius"];
	for (var i = 0, len = alert_fields.length; i < len; i++) {
		if ($("#" + id + "-" + alert_fields[i]).length) {
			data[alert_fields[i]] = $("#" + id + "-" + alert_fields[i]).val();
		}
	}

	if (data["Targets"].length > 0) {
		data["Targets"] = data["Targets"].split(",");
		for (i = 0; i < data["Targets"].length; i++) {
			data["Targets"][i] = parseInt(data["Targets"][i]);
		}
	}
	if (data["Categories"].length > 0) {
		data["Categories"] = data["Categories"].split(",");
		for (i = 0; i < data["Categories"].length; i++) {
			data["Categories"][i] = parseInt(data["Categories"][i]);
		}
	}
	if (data["Dates"].length > 0) {
		data["Dates"] = data["Dates"].split(",");
		for (i = 0; i < data["Dates"].length; i++) {
			data["Dates"][i] = parseInt(data["Dates"][i]);
		}
	}
	if (data["Dates"].length > 0) {
		data["Radius"] = parseInt(data["Radius"]);
	}

	data["Weekdays"] = new Array();
	var weekdays = $("input[name=" + id + "-Weekday]:checked").map(function () {return this.value;}).get();
	for (i = 0; i < weekdays.length; i++) {
		data["Weekdays"][i] = parseInt(weekdays[i]);
	}
	
	return data;
}

function activateSpinner(id)
{
	$("#" + id + "-spinner").addClass("fa-spin");
	$("#" + id + "-spinner").css("visibility", "visible");
}

function deactivateSpinner(id)
{
	$("#" + id + "-spinner").css("visibility", "hidden");
	$("#" + id + "-spinner").removeClass("fa-spin");
}

function setErrorMessage(id, msg)
{
	$("#" + id + "-message").parent().addClass("has-error");
	$("#" + id + "-message").text(msg);
}

function removeErrorMessage(id)
{
	$("#" + id + "-message").parent().removeClass("has-error");
	$("#" + id + "-message").text("");
}

function gatherSearchForm()
{
	var data = {};
	if ($("input[name=query]").length) {
		data["query"] = $("input[name=query]").val().trim();
	}
	if ($("input[name=place]").length) {
		data["place"] = $("input[name=place]").val().trim();
	}
	
	var target = "";
	if ($("select[name=target]").length) {
		target = $("select[name=target]").val();
	} else {
		var targets = $("input[name=target]:checked").map(function () {return this.value;}).get();
		for (i = 0; i < targets.length; i++) {
			if (i > 0) target += ",";
			target += targets[i];
		}
	}
	if (target == "") target = "0";
	data["target"] = target;
	
	var category = "";
	if ($("select[name=category]").length) {
		category = $("select[name=category]").val();
	} else {
		var categories = $("input[name=category]:checked").map(function () {return this.value;}).get();
		for (i = 0; i < categories.length; i++) {
			if (i > 0) category += ",";
			category += categories[i];
		}
	}
	if (category == "") category = "0";
	data["category"] = category;
	
	var date = "";
	if ($("select[name=date]").length) {
		date = $("select[name=date]").val();
	} else {
		var dates = $("input[name=date]:checked").map(function () {return this.value;}).get();
		for (i = 0; i < dates.length; i++) {
			if (i > 0) date += ",";
			date += dates[i];
		}
	}
	if (date == "") date = "0";
	data["date"] = date;

	return data;
}

function updateEventCount()
{
	var data = gatherSearchForm();
	$.ajax({cache: false, url : "/eventcount/" + (data["query"] ? data["query"] : "") + "/" + data["place"] + "/" + data["date"] + "/" + data["target"]+ "/" + data["category"], type: "GET", dataType: "json",
		success: function(data) {
			$("button[value=events]").html(data + (data != 1 ? " Veranstaltungen" : " Veranstaltung"));
		}
	});
}

function updateOrganizerCount()
{
	var data = gatherSearchForm();
	$.ajax({cache: false, url : "/organizercount/" + data["place"] + "/" + data["category"], type: "GET", dataType: "json",
		success: function(data) {
			$("button[value=organizers]").html(data + (data != 1 ? " Organisatoren" : " Organisator"));
		}
	});
}

$(function() {

	initEmailAndPwdForm("password");

	$("#password").submit(function(e) {
		
		e.preventDefault();
		if (!validateEmailAndPwdForm("password")) return;
		$("#password-submit").button('loading');
		var data = {"Email": $("#password-Email").val(), "Pwd": $("#password-Pwd").val()};
		
		$.ajax({cache: false, url : "/password", type: "POST", data : JSON.stringify(data),
			success: function() {
				window.location.href = "/veranstalter/verwaltung/0";
			},
			error: function() {
				$("#password-submit").button('reset');
				alert("Es gab ein Problem in der Kommunikation mit dem Server. Bitte versuche es später noch einmal.");
			}
		});
	});
	
	initProfileForm("profile");

	$("#profile-upload").fileupload({
		
		dataType: "json",
		dropZone: $("#profile-dropzone"),
		
		add: function (e, data) {
			activateSpinner("profile");
			data.submit();
		},
		
		done: function(e, data) {
			deactivateSpinner("profile");
			removeErrorMessage("profile-thumbnail");
			$("#profile-thumbnail").attr("src", "/bild/" + data.result + "?height=200&width=200");
			$("#profile-Image").attr("value", data.result);
		},
		
		fail: function(e, data) {
			deactivateSpinner("profile");
			setErrorMessage("profile-thumbnail", "Beim Hochladen des Bildes ist ein Fehler aufgetreten.");
		}
	});
	
	$("#profile-upload").submit(function(e) {
		
		e.preventDefault();
		if (!validateProfileForm("profile")) return;
		var data = gatherProfileForm("profile");
		
		$.ajax({cache: false, url : "/profile", type: "POST", data : JSON.stringify(data),
			success: function() {
				window.location.href = "/veranstalter/verwaltung/0";
			},
			error: function() {
				alert("Es gab ein Problem in der Kommunikation mit dem Server. Bitte versuche es später noch einmal.");
			}
		});
	});
	
	initEventForm("event");
	
	$("#event-upload").fileupload({
		
		dataType: "json",
		dropZone: $("#event-dropzone"),
		
		add: function (e, data) {
			activateSpinner("event");
			data.submit();
		},
		
		done: function(e, data) {
			deactivateSpinner("event");
			removeErrorMessage("event-thumbnail");
			$("#event-thumbnail").attr("src", "/bild/" + data.result + "?height=200&width=200");
			$("#event-Image").attr("value", data.result);
		},
		
		fail: function(e, data) {
			deactivateSpinner("event");
			setErrorMessage("event-thumbnail", "Beim Hochladen des Bildes ist ein Fehler aufgetreten.");
		}
	});
	
	$("#event-upload").submit(function(e) {
		
		e.preventDefault();
		if (!validateEventForm("event")) return;
		var data = gatherEventForm("event");
		
		$.ajax({cache: false, url : "/event", type: "POST", data : JSON.stringify(data),
			success: function() {
				window.location.href = "/veranstalter/verwaltung/0";
			},
			error: function() {
				alert("Es gab ein Problem in der Kommunikation mit dem Server. Bitte versuche es später noch einmal.");
			}
		});
	});

	$("#login").on('shown.bs.modal', function () {
		$("#login-form").submit(function(e) {
			e.preventDefault();
			var data = {"Email": $("#login-Email").val(), "Pwd": $("#login-Pwd").val()};
			$.ajax({cache: false, url : "/login", type: "POST", dataType : "json", data : JSON.stringify(data),
				success: function(sessionid) {
					$("#login").modal("hide");
					$.cookie("SESSIONID", sessionid, {path: '/'});
					window.location.href = "/veranstalter/verwaltung/0";
				},
				error : function(result) {
					if (result.status == 404) {
						alert("Diese Anmeldedaten sind uns nicht bekannt.");
					} else {
						alert("Es gab ein Problem in der Kommunikation mit dem Server. Bitte versuche es später noch einmal.");
					}
				}
			});
		});
		
	    $("#login input").first().focus();
	});
	
	$("#register").on('shown.bs.modal', function () {
		initProfileForm("register");

		$("#register-upload").fileupload({
			
			dataType: "json",
			dropZone: $("#register-dropzone"),
			
			add: function (e, data) {
				activateSpinner("register");
				data.submit();
			},
			
			done: function(e, data) {
				deactivateSpinner("register");
				removeErrorMessage("register-thumbnail");
				$("#register-thumbnail").attr("src", "/bild/" + data.result + "?height=200&width=200");
				$("#register-Image").attr("value", data.result);
			},
			
			fail: function(e, data) {
				deactivateSpinner("register");
				setErrorMessage("register-thumbnail", "Beim Hochladen des Bildes ist ein Fehler aufgetreten.");
			}
		});
		
		$("#register-upload").submit(function(e) {
			
			e.preventDefault();
			if (!validateProfileForm("register")) return;
			if (!$("#register-AGBs").is(':checked')) {
				alert("Bitte stimme den Allgemeinen Geschäftsbedingungen zu.");
				return;
			}
			var data = gatherProfileForm("register");
			data["agbs"] = $("#register-AGBs").is(':checked')
			$("#register-submit").button('loading');

			$.ajax({cache: false, url : "/register", type: "POST", dataType : "json", data : JSON.stringify(data),
				success: function(sessionid) {
					$("#register-submit").button('reset');
					$("#register").modal("hide");
					$.cookie("SESSIONID", sessionid, {path: '/'});
					$("#registered").load("/dialog/registered").modal("show");
				},
				error: function(result) {
					if (result.status == 409) {
						alert("Die E-Mail-Adresse ist schon registriert. Bitte wähle eine andere.");
					} else if (result.status == 400) {
						alert("Bitte stimme den Allgemeinen Geschäftsbedingungen zu.");
					} else {
						alert("Es gab ein Problem in der Kommunikation mit dem Server. Bitte versuche es später noch einmal.");
					}
					$("#register-submit").button('reset');
				}
			});
		});
		
	    $("#register input").first().focus();
	});
	
	$("#mail").on('shown.bs.modal', function () {
		initSendMailForm("send-mail");
		
		$("#send-mail").submit(function (e) {
			e.preventDefault();
			if (!validateSendMailForm("send-mail")) return;
			$("#send-mail-submit").button('loading');
			var data = {"Name": $("#send-mail-Name").val(), "Email": $("#send-mail-Email").val(), "Subject": $("#send-mail-Subject").val(), "Text": $("#send-mail-Text").val()};
			$.ajax({cache: false, url : "/sendcontactmail", type: "POST", data : JSON.stringify(data),
				success: function() {
					$("#send-mail-submit").button('reset');
					alert("Deine Nachricht wurde verschickt.")
					$("#mail").modal("hide");
				},
				error: function() {
					$("#send-mail-submit").button('reset');
					alert("Es gab ein Problem in der Kommunikation mit dem Server. Bitte versuche es später noch einmal.");
				}
			});
		});
		
	    $("#mail input").first().focus();
	});
	
	$("#share").on('shown.bs.modal', function () {
		initSendMailForm("send-event");
		
		$("#send-event").submit(function (e) {
			e.preventDefault();
			if (!validateSendMailForm("send-event")) return;
			$("#send-event-submit").button('loading');
			var data = {"Name": $("#send-event-Name").val(), "Email": $("#send-event-Email").val(), "Subject": $("#send-event-Subject").val(), "Text": $("#send-event-Text").val()};
			$.ajax({cache: false, url : "/sendeventmail", type: "POST", data : JSON.stringify(data),
				success: function() {
					$("#send-event-submit").button('reset');
					alert("Deine Nachricht wurde verschickt.")
					$("#share").modal("hide");
				},
				error: function() {
					$("#send-event-submit").button('reset');
					alert("Es gab ein Problem in der Kommunikation mit dem Server. Bitte versuche es später noch einmal.");
				}
			});
		});
		
	    $("#share input").first().focus();
	});
	
	$("#email-alert").on('shown.bs.modal', function () {
		initEmailAlertForm("email-alert");
		
		$("#email-alert-form").submit(function (e) {
			e.preventDefault();
			if (!validateEmailAlertForm("email-alert")) return;
			$("#email-alert-submit").button('loading');
			var data = gatherEmailAlertForm("email-alert");
			$.ajax({cache: false, url : "/emailalert", type: "POST", data : JSON.stringify(data),
				success: function() {
					$("#email-alert-submit").button('reset');
					alert("Die Benachrichtigung wurde angelegt. Du bekommst jetzt regelmäßig Veranstaltungen für Deine Suche zugesendet.")
					$("#email-alert").modal("hide");
				},
				error: function() {
					$("#email-alert-submit").button('reset');
					alert("Es gab ein Problem in der Kommunikation mit dem Server. Bitte versuche es später noch einmal.");
				}
			});
		});
		
	    $("#email-alert input").first().focus();
	});
	
	$("#delete-profile").click(function(e) {
		e.preventDefault();
		if (confirm("Dein Profil und Deine Veranstaltungen werden unwiederbringlich gelöscht.")) {
			$.ajax({cache: false, url : "/unregister", type: "POST",
				success: function(sessionid) {
					$.removeCookie("SESSIONID", {path: '/'});
					window.location.href = "/";
				},
				error: function() {
					alert("Es gab ein Problem in der Kommunikation mit dem Server. Bitte versuche es später noch einmal.");
				}
			});
		}
	});
	
	$("#send-double-opt-in").click(function(e) {
		e.preventDefault();
		$("#send-double-opt-in").button('loading');
		$.ajax({cache: false, url : "/sendcheckmail", type: "POST",
			success: function() {
				$("#send-double-opt-in").button('reset');
				alert("Die E-Mail wurde versendet. Bitte überprüfe Dein Postfach und klicke auf den Link in der Mail.")
			},
			error: function() {
				$("#send-double-opt-in").button('reset');
				alert("Es gab ein Problem in der Kommunikation mit dem Server. Bitte versuche es später noch einmal.");
			}
		});
	});
	
	$("a[name=delete-event]").click(function(e) {
		e.preventDefault();
		if (confirm("Die Veranstaltung wird unwiederbringlich gelöscht.")) {
			$.ajax({cache: false, url : "/event/" + $(this).data("target"), type: "DELETE",
				success: function() {
					window.location.reload();
				},
				error: function() {
					alert("Es gab ein Problem in der Kommunikation mit dem Server. Bitte versuche es später noch einmal.");
				}
			});
		}
	});

	$("input[name=fulltextsearch]").typeahead({
		source: function(query, process) {
			$.ajax({cache: false, url : "/typeahead/" + query, type: "GET", dataType: "json",
				success: function(data) {
					process(data);
				}
			});
		},
		afterSelect: function(item) {
			$('#fulltextsearch').submit();
		}
	});

	$("input[name=query]").typeahead({
		source: function(query, process) {
			$.ajax({cache: false, url : "/typeahead/" + query, type: "GET", dataType: "json",
				success: function(data) {
					process(data);
				}
			});
		},
		afterSelect: function(item) {
			$('#adminsearch').submit();
		}
	});

	$("input[name=query]").change(function() {
		updateEventCount();
		updateOrganizerCount();
	});

	$("input[name=place]").typeahead({ source: function(query, process) {
		$.ajax({cache: false, url : "/location/" + query, type: "GET", dataType: "json",
			success: function(data) {
				process(data);
			}
		});
	}});
	
	$("input[name=place]").change(function() {
		updateEventCount();
		updateOrganizerCount();
	});
	
	$("select[name=target]").change(function() {
		updateEventCount();
	});
	
	$("select[name=category]").change(function() {
		updateEventCount();
		updateOrganizerCount();
	});
	
	$("select[name=date]").change(function() {
		updateEventCount();
	});
	
	$("input[name=target]").change(function() {
		$("#events-form").submit();
	});
	
	$("input[name=category]").change(function() {
		$("#events-form").submit();
	});
	
	$("input[name=date]").change(function() {
		$("#events-form").submit();
	});

	var now = new Date();
	$(".form-datetime").datetimepicker({
		format: "dd.mm.yyyy hh:ii",
		autoclose: true,
		fontAwesome: true,
		language: "de",
		pickerPosition: "bottom-right",
		todayHighlight: true,
		startDate: zeroFill(now.getDate(), 2) + "." + zeroFill(now.getMonth() + 1, 2) + "." + now.getFullYear() + " " + zeroFill(now.getHours(), 2) + ":" + zeroFill(now.getMinutes(), 2)
	});
	
	$("input[type=radio][name=recurrency]").change(function () {
		if (this.value == "weekly") {
			$("#event-monthly").hide();
			$("#event-weekly").show();
			$("#event-recurrencyEnd").show();
		} else if (this.value == "monthly") {
			$("#event-weekly").hide();
			$("#event-monthly").show();
			$("#event-recurrencyEnd").show();
		} else {
			$("#event-weekly").hide();
			$("#event-monthly").hide();
			$("#event-recurrencyEnd").hide();
		}
	});

	$("a[data-href]").each(function() {
        var $this = $(this);
        $this.attr("href", $this.attr("data-href"));
    });

	var hash = window.location.hash;
	if (hash.substring(1) == "login") {
		$.removeCookie("SESSIONID", {path: '/'});
		$("#login .modal-content").load("/dialog/login");
		$("#login").modal("show");
	}

	if ($.cookie("SESSIONID")) {
		$("#head-organizer").html("Du bist angemeldet.");
		$("#head-events").html("<span class='fa fa-caret-right'></span> Meine Veranstaltungen");
		$("#head-events").attr("data-toggle", "");
		$("#head-events").attr("data-target", "");
		$("#head-events").attr("href", "/veranstalter/verwaltung/0");
		$("#head-events").attr("title", "Verwalte Deine Veranstaltungen und Profileinstellungen.");
		$("#head-login").html("<span class='fa fa-user highlight'></span> Abmelden");
		$("#head-login").attr("data-toggle", "");
		$("#head-login").attr("data-target", "");
		$("#head-login").attr("href", "#");
		$("#head-login").attr("title", "Abmelden");
		$("#head-login").click(function() {
			$.ajax({cache: false, url : "/logout", type: "POST",
				success: function() {
					$.removeCookie("SESSIONID", {path: '/'});
					window.location.href = "/";
				},
				error: function() {
					$.removeCookie("SESSIONID", {path: '/'});
					window.location.href = "/";
				}
			});
			return false;
		});
	}
});
