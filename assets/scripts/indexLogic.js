var labButtons = [];
var show = true;

function showLabyrinths() {
    if (show) {
        var settings = {
            "type": "GET",
            "contentType": "application/json; charset=utf-8",
            "xhrFields": {
                withCredentials: true
            },
            "crossDomain": true,
            "url": "http://localhost:8080/game/labyrinths"
        };

        $.ajax(settings).done(function(xhr, textStatus, errorThrown) {
            labButton(xhr);
        }).fail(function (jqXHR, textStatus) {
            alert(jqXHR.status + " " + jqXHR.statusText + ". " + jqXHR.responseText);
            console.log(jqXHR.status + " " + jqXHR.statusText + ". " + jqXHR.responseText);
        });
    } else {
        labButton();
    }
}

function labButton(data) {
    if (show) {
        var lab = JSON.parse(data);
        for (var i = 0; i < lab.length; i++) {
            var button = document.createElement("button")
            button.className = "button_lab";
            button.innerHTML = lab[i]["name"];
            button.setAttribute("lab_id", lab[i]["id"]);
            button.onclick = setLab;
            labButtons.push(button);
            $('#labMenu').append(button);
        }
        show = false;
    } else {
        for (var i = 0; i < labButtons.length; i++) {
            labButtons[i].parentNode.removeChild(labButtons[i]);

        }
        labButtons.length = 0;
        show = true;
    }
}

function setLab() {
    var button = this;
    $('#show').text("Selected labyrinth: " + this.innerHTML);
    $('#show').attr("lab_id", button.getAttribute("lab_id"));
    show = false;
    labButton("");
}

function startGame() {
    if ($('#show').attr("lab_id") == undefined) {
        alert("Select some labyrinth");
    } else {
        var url = "http://localhost:8080/game/start/" + $('#show').attr("lab_id");
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
            window.location ="/game";
        }).fail(function (jqXHR, textStatus) {
            alert(jqXHR.status + " " + jqXHR.statusText + ". " + jqXHR.responseText);
            console.log(jqXHR.status + " " + jqXHR.statusText + ". " + jqXHR.responseText);
        });
    }
}