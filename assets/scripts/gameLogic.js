function quit() {
    var url = "http://localhost:8080/game/quit";
    var settings = {
        "type": "POST",
        "contentType": "application/json; charset=utf-8",
        "xhrFields": {
            withCredentials: true
        },
        "crossDomain": true,
        "url": url
    };

    $.ajax(settings).done(function(xhr, textStatus, errorThrown) {
        window.location = "/";
    }).fail(function (jqXHR, textStatus) {
        alert(jqXHR.status + " " + jqXHR.statusText + ". " + jqXHR.responseText);
        console.log(jqXHR.status + " " + jqXHR.statusText + ". " + jqXHR.responseText);
    });
}
