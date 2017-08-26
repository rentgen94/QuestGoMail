var doors = [];
var showDoors = true;

function showDoorsItems() {
    if (showDoors) {
        var settings = {
            "type": "GET",
            "contentType": "application/json; charset=utf-8",
            "xhrFields": {
                withCredentials: true
            },
            "dataType": "text",
            "crossDomain": true,
            "url": "http://localhost:8080/game/look_around/entities/doors"
        };

        $.ajax(settings).done(function(xhr, textStatus, errorThrown) {
            doorItem(xhr);
        }).fail(function (jqXHR, textStatus) {
            alert(jqXHR.status + " " + jqXHR.statusText + ". " + jqXHR.responseText);
            console.log(jqXHR.status + " " + jqXHR.statusText + ". " + jqXHR.responseText);
        });
    } else {
        doorItem("");
    }
}

function doorItem(data) {
    if (data != undefined && data != "") {
        var items = JSON.parse(data);
        if ((items["msg"] != "") || (items["errMsg"] != "")) {
            alert("Msg: " + items["msg"] + "; Error: " + items["errMsg"]);
        }
        if (showDoors) {
            for (var i = 0; i < items["data"].length; i++) {
                var holder = document.createElement("div");

                var button = document.createElement("button");
                button.className = "buttonDoorItem";
                button.innerHTML = "Interesting door";
                button.onclick = openDoor;
                button.setAttribute("door_id", items["data"][i]["id"]);

                var name = document.createElement("h3");
                name.innerHTML = items["data"][i]["name"];

                holder.appendChild(button);
                holder.appendChild(name);

                doors.push(holder);

                $('#doorMenu').append(holder);
            }
            showDoors = false;
        } else {
            for (var i = 0; i < doors.length; i++) {
                doors[i].parentNode.removeChild(doors[i]);
            }
            doors.length = 0;
            showDoors = true;
        }
    } else if (doors.length >= 0) {
        for (var i = 0; i < doors.length; i++) {
            doors[i].parentNode.removeChild(doors[i]);

        }
        doors.length = 0;
        showDoors = true;
    }
}

function openDoor() {
    var item = this;
    alert("Take item with id: " + item.getAttribute("door_id"));
    var command = {
        "code": 9,
        "item_key": parseInt(item.getAttribute("door_id")),
        "args": [item.getAttribute("door_id")],
        "items": []
    };
    var settings = {
        "type": "POST",
        "contentType": "application/json; charset=utf-8",
        "xhrFields": {
            withCredentials: true
        },
        "dataType": "text",
        "crossDomain": true,
        "url": "http://localhost:8080/game/command",
        "data": JSON.stringify(command)

    };

    $.ajax(settings).done(function(xhr, textStatus, errorThrown) {
        var items = JSON.parse(xhr);
        if ((items["msg"] != "") || (items["errMsg"] != "")) {
            alert("Msg: " + items["msg"] + "; Error: " + items["errMsg"]);
        }
        if (doors.length >= 0) {
            for (var i = 0; i < doors.length; i++) {
                doors[i].parentNode.removeChild(doors[i]);

            }
            doors.length = 0;
            showDoors = true;
        }
    }).fail(function (jqXHR, textStatus) {
        alert(jqXHR.status + " " + jqXHR.statusText + ". " + jqXHR.responseText);
        console.log(jqXHR.status + " " + jqXHR.statusText + ". " + jqXHR.responseText);
    });
}