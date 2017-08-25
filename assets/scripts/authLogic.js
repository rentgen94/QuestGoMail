function objectifyForm(formArray) {//serialize data function

    var returnArray = {};
    for (var i = 0; i < formArray.length; i++){
        returnArray[formArray[i]['name']] = formArray[i]['value'];
    }
    return returnArray;
}

function login() {
    var credentials = $("#nameform").serializeArray();
    credentials = objectifyForm(credentials);
    console.log(credentials);

    var settings = {
        "method": "POST",
        "contentType": "application/json; charset=utf-8",
        "dataType": "jsonp",
        "crossDomain": true,
        "url": "http://localhost:8080/player/login",
        "data": credentials,
    }

    $.ajax(settings).done(function(response) {
        alert("Wow");
//            var in30Minutes = 1 / (24 * 2); // было 12(5 minutes)
//            Cookies.set('user', $("#nameform #user").val(), {
//                expires: in30Minutes
//            });
//            Cookies.set('token', response, {
//                expires: in30Minutes
//            });
        window.location = "/";
    }).fail(function (jqXHR, textStatus) {
        alert(jqXHR.status + " " + jqXHR.statusText + ". " + jqXHR.responseText);
        console.log(jqXHR.status + " " + jqXHR.statusText + ". " + jqXHR.responseText);
    });
}