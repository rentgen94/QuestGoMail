function objectifyForm(formArray) {//serialize data function

    var returnArray = {};
    for (var i = 0; i < formArray.length; i++){
        returnArray[formArray[i]['name']] = formArray[i]['value'];
    }
    return returnArray;
}

$("#nameform").on('submit', login);

function login() {
    var credentials = $("#nameform").serializeArray();
    credentials = objectifyForm(credentials);
    credentials = JSON.stringify(credentials);
    console.log(credentials);

    var settings = {
        "type": "POST",
        "contentType": "application/json; charset=utf-8",
        "xhrFields": {
            withCredentials: true
        },
        "dataType": "json",
        "crossDomain": true,
        "url": "http://localhost:8080/player/login",
        "data": credentials
    }

    $.ajax(settings).done(function(response) {
        var in30Minutes = 1 / (24 * 2); // было 12(5 minutes)
        Cookies.set('player', $("#nameform #name").val(), {
           expires: in30Minutes
        });
        //console.log("Cookie: " + xhr.getResponseHeader("Set-Cookie"));
        window.location = "/";
    }).fail(function (jqXHR, textStatus) {
        alert(jqXHR.status + " " + jqXHR.statusText + ". " + jqXHR.responseText);
        console.log(jqXHR.status + " " + jqXHR.statusText + ". " + jqXHR.responseText);
    });
}

$("#regform").on('submit', register);

function register() {
    var credentials = $("#regform").serializeArray();
    credentials = objectifyForm(credentials);
    credentials = JSON.stringify(credentials);
    console.log(credentials);

    var settings = {
        "type": "POST",
        "contentType": "application/json; charset=utf-8",
        "xhrFields": {
            withCredentials: true
        },
        "dataType": "json",
        "crossDomain": true,
        "url": "http://localhost:8080/player/register",
        "data": credentials
    }

    $.ajax(settings).done(function(response) {
        window.location = "/player/login";
    }).fail(function (jqXHR, textStatus) {
        alert(jqXHR.status + " " + jqXHR.statusText + ". " + jqXHR.responseText);
        console.log(jqXHR.status + " " + jqXHR.statusText + ". " + jqXHR.responseText);
    });
}