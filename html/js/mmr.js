var PcodePattern = /^[0-9]{5}$/;
var DateTimePattern = /^(\d{2})\.(\d{2})\.(\d{4}) (\d{2}).(\d{2})$/;
var EmailPattern = /^([a-zA-Z0-9_.+-])+\@(([a-zA-Z0-9-])+\.)+([a-zA-Z0-9]{2,4})+$/;
var WebPattern = /(http|ftp|https):\/\/[\w-]+(\.[\w-]+)+([\w.,@?^=%&amp;:\/~+#-]*[\w@?^=%&amp;\/~+#-])?/;

function initProfileForm(id)
{
	$("#" + id + "-dropzone").click(function() {
		$(this).parent().find("input").click();
	});

	$("#" + id + "-Name").popover({content: "Bitte gib den Namen Deiner Organisation ein.", trigger: "manual", placement: "auto right"});
	if ($("#" + id + "-Email").length) {
		$("#" + id + "-Email").popover({content: "Bitte gib eine gültige E-Mail-Adresse ein.", trigger: "manual", placement: "auto right"});
	}
	if ($("#" + id + "-Pwd").length) {
		$("#" + id + "-Pwd").popover({content: "Bitte gib ein Kennwort ein.", trigger: "manual", placement: "auto right"});
		$("#" + id + "-Pwd2").popover({content: "Kennwort und Kennwortwiederholung stimmen nicht überein.", trigger: "manual", placement: "auto right"});
	}
	$("#" + id + "-Web").popover({content: "Bitte gib eine gültige Web-Adresse (URL) ein.", trigger: "manual", placement: "auto right"});
	$("#" + id + "-Pcode").popover({content: "Bitte gib eine vollständige Postleitzahl ein.", trigger: "manual", placement: "auto top"});	
}

function initEmailAndPwdForm(id)
{
	$("#" + id + "-Email").popover({content: "Bitte gib eine gültige E-Mail-Adresse ein.", trigger: "manual", placement: "auto right"});
	$("#" + id + "-Pwd2").popover({content: "Kennwort und Kennwortwiederholung stimmen nicht überein.", trigger: "manual", placement: "auto right"});
}

function initEventForm(id)
{
	$("#" + id + "-dropzone").click(function() {
		$(this).parent().find("input").click();
	});

	$("#" + id + "-Title").popover({content: "Bitte gib der Veranstaltung einen Titel.", trigger: "manual", placement: "auto right"});
	$("#" + id + "-Start").popover({content: "Bitte gib den Veranstaltungsbeginn an.", trigger: "manual", placement: "auto right"});
	$("#" + id + "-End").popover({content: "Bitte gib ein gültiges Ende an.", trigger: "manual", placement: "auto right"});
	$("#" + id + "-Category").popover({content: "Bitte wähle mindestens eine Kategorie aus.", trigger: "manual", placement: "auto right"});
	$("#" + id + "-Web").popover({content: "Bitte gib eine gültige Web-Adresse (URL) ein.", trigger: "manual", placement: "auto right"});
	$("#" + id + "-Pcode").popover({content: "Bitte gib eine vollständige Postleitzahl ein.", trigger: "manual", placement: "auto top"});	
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
	ok = validate(EmailPattern.test($("#" + id +"-Email").val()), "#" + id +"-Email");
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
	ok &= validate(DateTimePattern.test($("#" + id + "-Start").val()), "#" + id + "-Start");
	ok &= validate($("#" + id +"-End").val().trim().length == 0 || DateTimePattern.test($("#" + id + "-End").val()), "#" + id + "-End");
	ok &= validate($("#" + id +"-Web").val().trim().length == 0 || WebPattern.test($("#" + id +"-Web").val()), "#" + id +"-Web");
	ok &= validate($("#" + id +"-Pcode").val().trim().length == 0 || PcodePattern.test($("#" + id +"-Pcode").val()), "#" + id +"-Pcode");
	ok &= validate($("input[name=" + id + "-Category]:checked").map(function () {return this.value;}).get().length > 0, "#" + id + "-Category");
	return ok;
}

function gatherProfileForm(id)
{
	var data = {};
	var user_fields = ["Email", "Pwd", "Image", "Descr", "Web"];
	for (var i = 0, len = user_fields.length; i < len; i++) {
		if ($("#" + id + "-" + user_fields[i]).length) {
			data[user_fields[i]] = $("#" + id + "-" + user_fields[i]).val();
		}
	}
	data["Addr"] = {}
	var addr_fields = ["Name", "Street", "Pcode", "City"];
	for (var i = 0, len = addr_fields.length; i < len; i++) {
		if ($("#" + id + "-" + addr_fields[i]).length) {
			data["Addr"][addr_fields[i]] = $("#" + id + "-" + addr_fields[i]).val();
		}
	}
	data["Categories"] = new Array();
	categories = $("input[name=" + id + "-Category]:checked").map(function () {return this.value;}).get();
	for (i = 0; i < categories.length; i++) {
		data["Categories"][i] = parseInt(categories[i]);
	}
	return data;
}

function gatherEventForm(id)
{
	var data = {};
	var event_fields = ["Id", "Title", "Start", "End", "Image", "Descr", "Web"];
	for (var i = 0, len = event_fields.length; i < len; i++) {
		if ($("#" + id + "-" + event_fields[i]).length) {
			data[event_fields[i]] = $("#" + id + "-" + event_fields[i]).val();
		}
	}
	
	if (data["Start"].length > 0) {
		DateTimePattern.exec(data["Start"]);
		data["Start"] = new Date(RegExp.$2 + " " + RegExp.$1 + ", " + RegExp.$3 + " " + RegExp.$4 + ":" + RegExp.$5 + ":00");
	} else {
		delete data["Start"];
	}
	if (data["End"].length > 0) {
		DateTimePattern.exec(data["End"]);
		data["End"] = new Date(RegExp.$2 + " " + RegExp.$1 + ", " + RegExp.$3 + " " + RegExp.$4 + ":" + RegExp.$5 + ":00");
	} else {
		delete data["End"];
	}

	data["Addr"] = {}
	var addr_fields = ["Name", "Street", "Pcode", "City"];
	for (var i = 0, len = addr_fields.length; i < len; i++) {
		if ($("#" + id + "-" + addr_fields[i]).length) {
			data["Addr"][addr_fields[i]] = $("#" + id + "-" + addr_fields[i]).val();
		}
	}
	
	data["Categories"] = new Array();
	categories = $("input[name=" + id + "-Category]:checked").map(function () {return this.value;}).get();
	for (i = 0; i < categories.length; i++) {
		data["Categories"][i] = parseInt(categories[i]);
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
	place = $("input[name=place]").val().trim();
	if (place.length == 0) place = "Berlin";
	data["place"] = place;
	
	var category = "";
	if ($("select[name=category]").length) {
		category = $("select[name=category]").val();
	} else {
		categories = $("input[name=category]:checked").map(function () {return this.value;}).get();
		for (i = 0; i < categories.length; i++) {
			if (i > 0) category += ",";
			category += categories[i];
		}
	}
	if (category == "") category = "0";
	data["category"] = category;
	
	var date = "";
	dates = $("input[name=date]:checked").map(function () {return this.value;}).get();
	for (i = 0; i < dates.length; i++) {
		if (i > 0) date += ",";
		date += dates[i];
	}
	if (date == "") date = "0";
	data["date"] = date;

	return data;
}

function updateEventCount()
{
	data = gatherSearchForm();
	$.ajax({cache: false, url : "/eventcount/" + data["place"] + "/" + data["date"] + "/" + data["category"], type: "GET", dataType: "json",
		success: function(data) {
			$("button[value=events]").html(data + (data != 1 ? " Veranstaltungen" : " Veranstaltung"));
		}
	});
}

function updateOrganizerCount()
{
	data = gatherSearchForm();
	$.ajax({cache: false, url : "/organizercount/" + data["place"] + "/" + data["category"], type: "GET", dataType: "json",
		success: function(data) {
			$("button[value=organizers]").html(data + " Veranstalter");
		}
	});
}

$(function() {

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
		var data = gatherProfileForm("register");
		$("#register-submit").button('loading');

		$.ajax({cache: false, url : "/register", type: "POST", dataType : "json", data : JSON.stringify(data),
			success: function(sessionid) {
				$("#register-submit").button('reset');
				$("#register").modal("hide");
				$.cookie("SESSIONID", sessionid, {path: '/'});
				$("#registered").modal("show");
			},
			error : function(result) {
				if (result.status == 409) {
					alert("Die E-Mail-Adresse ist schon registriert. Bitte wähle eine andere.");
				} else {
					alert("Es gab ein Problem in der Kommunikation mit dem Server. Bitte versuche es später noch einmal.");
				}
				$("#register-submit").button('reset');
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
		
		$.ajax({cache: false, url : "/profile", type: "POST", dataType : "json", data : JSON.stringify(data),
			error : function(result) {
				if (result.status == 200) {
					window.location.href = "/veranstalter/verwaltung/0";
				} else {
					alert("Es gab ein Problem in der Kommunikation mit dem Server. Bitte versuche es später noch einmal.");
				}
			}
		});
	});
	
	initEmailAndPwdForm("password");

	$("#password").submit(function(e) {
		
		e.preventDefault();
		if (!validateEmailAndPwdForm("password")) return;
		$("#password-submit").button('loading');
		var data = {"Email": $("#password-Email").val(), "Pwd": $("#password-Pwd").val()};
		
		$.ajax({cache: false, url : "/password", type: "POST", dataType : "json", data : JSON.stringify(data),
			error : function(result) {
				if (result.status == 200) {
					window.location.href = "/veranstalter/verwaltung/0";
				} else {
					$("#password-submit").button('reset');
					alert("Es gab ein Problem in der Kommunikation mit dem Server. Bitte versuche es später noch einmal.");
				}
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
		
		$.ajax({cache: false, url : "/event", type: "POST", dataType : "json", data : JSON.stringify(data),
			error : function(result) {
				if (result.status == 200 || result.status == 201) {
					window.location.href = "/veranstalter/verwaltung/0";
				} else {
					alert("Es gab ein Problem in der Kommunikation mit dem Server. Bitte versuche es später noch einmal.");
				}
			}
		});
	});

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
	
	$("#delete-profile").click(function(e) {
		e.preventDefault();
		if (confirm("Dein Profil und Deine Veranstaltungen werden unwiederbringlich gelöscht.")) {
			$.ajax({cache: false, url : "/unregister", type: "POST",
				success: function(sessionid) {
					$.removeCookie("SESSIONID", {path: '/'});
					window.location.href = "/";
				},
				error : function(result) {
					alert("Es gab ein Problem in der Kommunikation mit dem Server. Bitte versuche es später noch einmal.");
				}
			});
		}
	});
	
	$("a[name=delete-event]").click(function(e) {
		e.preventDefault();
		if (confirm("Die Veranstaltung wird unwiederbringlich gelöscht.")) {
			$.ajax({cache: false, url : "/event/" + $(this).data("target"), type: "DELETE",
				success: function() {
					window.location.reload();
				},
				error : function() {
					alert("Es gab ein Problem in der Kommunikation mit dem Server. Bitte versuche es später noch einmal.");
				}
			});
		}
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
	
	$("select[name=category]").change(function() {
		updateEventCount();
		updateOrganizerCount();
	});
	
	$("#events-form").submit(function() {
		if (!$("input[name=place]").val()) {
			$("input[name=place]").val($("input[name=place]").attr("placeholder"))
		}
	});
	
	$("input[name=category]").change(function() {
		$("#events-form").submit();
	});
	
	$("input[name=date]").change(function() {
		$("#events-form").submit();
	});

	$(".form-datetime").datetimepicker({
		format: "dd.mm.yyyy hh:ii",
		autoclose: true,
		language: "de",
		pickerPosition: "bottom-right"
	});
	
	var hash = window.location.hash;
	if (hash.substring(1) == "login") {
		$.removeCookie("SESSIONID", {path: '/'});
		$("#login").modal("show");
	}

	if ($.cookie("SESSIONID")) {
		$("#head-organizer").html("Du bist angemeldet.");
		$("#head-events").html("<span class='fa fa-caret-right'></span> Meine Veranstaltungen");
		$("#head-events").attr("data-toggle", "");
		$("#head-events").attr("data-target", "");
		$("#head-events").attr("href", "/veranstalter/verwaltung/0");
		$("#head-login").html("<span class='fa fa-user highlight'></span> Abmelden");
		$("#head-login").attr("data-toggle", "");
		$("#head-login").attr("data-target", "");
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
