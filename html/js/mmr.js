var PcodePattern = /^[0-9]{5}$/;
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
					window.location.href = "/veranstalter/verwaltung";
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
					window.location.href = "/veranstalter/verwaltung";
				} else {
					$("#password-submit").button('reset');
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
				window.location.href = "/veranstalter/verwaltung";
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

	if ($.cookie("SESSIONID")) {
		$("#head-events").html("<span class='fa fa-caret-right'></span> Meine Veranstaltungen");
		$("#head-events").attr("data-toggle", "");
		$("#head-events").attr("data-target", "");
		$("#head-events").attr("href", "/veranstalter/verwaltung");
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
