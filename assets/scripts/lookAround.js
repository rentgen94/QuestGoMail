var roomDescs = [];
var showRoomDesc = true;

function showRoomDescription() {
    if (showRoomDesc) {
        var settings = {
            "type": "GET",
            "contentType": "application/json; charset=utf-8",
            "xhrFields": {
                withCredentials: true
            },
            "dataType": "text",
            "crossDomain": true,
            "url": "http://localhost:8080/game/look_around"
        };

        $.ajax(settings).done(function(xhr, textStatus, errorThrown) {
            roomDesc(xhr);
        }).fail(function (jqXHR, textStatus) {
            alert(jqXHR.status + " " + jqXHR.statusText + ". " + jqXHR.responseText);
            console.log(jqXHR.status + " " + jqXHR.statusText + ". " + jqXHR.responseText);
        });
    } else {
        roomDesc("");
    }
}

function roomDesc(data) {
    if (data != undefined && data != "") {
        var items = JSON.parse(data);
        if ((items["msg"] != "") || (items["errMsg"] != "")) {
            alert("Msg: " + items["msg"] + "; Error: " + items["errMsg"]);
        }
        if (showRoomDesc) {
            var holder = document.createElement("div");

            var button = document.createElement("button");
            button.className = "buttonRoomDesc";
            button.innerHTML = "Interesting place!";
            button.onclick = tookRoomDesc;

            var name = document.createElement("h3");
            name.innerHTML = items["data"]["name"];
            var desc = document.createElement("p");
            desc.innerHTML = items["data"]["description"];

            holder.appendChild(button);
            holder.appendChild(name);
            holder.appendChild(desc);

            roomDescs.push(holder);
            $('#lookMenu').append(holder);
            showRoomDesc = false;
        } else {
            for (var i = 0; i < roomDescs.length; i++) {
                roomDescs[i].parentNode.removeChild(roomDescs[i]);

            }
            roomDescs.length = 0;
            showRoomDesc = true;
        }
    } else if (roomDescs.length > 0) {
        for (var i = 0; i < roomDescs.length; i++) {
            roomDescs[i].parentNode.removeChild(roomDescs[i]);

        }
        roomDescs.length = 0;
        showRoomDesc = true;
    }
}

function tookRoomDesc() {
    var desc = this;
    if (roomDescs.length > 0) {
        for (var i = 0; i < roomDescs.length; i++) {
            roomDescs[i].parentNode.removeChild(roomDescs[i]);

        }
        roomDescs.length = 0;
        showRoomDesc = true;
    }
}