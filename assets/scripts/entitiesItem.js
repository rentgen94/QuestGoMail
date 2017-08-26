var roomItems = [];
var showRoomItems = true;

function showRoomItem() {
    if (showRoomItems) {
        var settings = {
            "type": "GET",
            "contentType": "application/json; charset=utf-8",
            "xhrFields": {
                withCredentials: true
            },
            "dataType": "text",
            "crossDomain": true,
            "url": "http://localhost:8080/game/look_around/entities/items"
        };

        $.ajax(settings).done(function(xhr, textStatus, errorThrown) {
             roomItem(xhr);
        }).fail(function (jqXHR, textStatus) {
            alert(jqXHR.status + " " + jqXHR.statusText + ". " + jqXHR.responseText);
            console.log(jqXHR.status + " " + jqXHR.statusText + ". " + jqXHR.responseText);
        });
    } else {
         roomItem("");
    }
}

function roomItem(data) {
    if (data != undefined && data != "") {
        var items = JSON.parse(data);
        if ((items["msg"] != "") || (items["errMsg"] != "")) {
            alert("Msg: " + items["msg"] + "; Error: " + items["errMsg"]);
        }
        if (showRoomItems) {
            for (var i = 0; i < items["data"].length; i++) {
                var holder = document.createElement("div");

                var button = document.createElement("button");
                button.className = "buttonRoomItem";
                button.innerHTML = "Take this awesome item!";
                button.onclick = tookRoomItem;
                button.setAttribute("item_id", items["data"][i]["id"]);

                var name = document.createElement("h3");
                name.innerHTML = items["data"][i]["name"];
                var desc = document.createElement("p");
                desc.innerHTML = items["data"][i]["description"];
                var size = document.createElement("p");
                size.innerHTML = items["data"][i]["size"];

                holder.appendChild(button);
                holder.appendChild(name);
                holder.appendChild(desc);
                holder.appendChild(size);

                roomItems.push(holder);

                $('#itemMenu').append(holder);
            }
            showRoomItems = false;
        } else {
            for (var i = 0; i < roomItems.length; i++) {
                roomItems[i].parentNode.removeChild(roomItems[i]);

            }
            roomItems.length = 0;
            showRoomItems = true;
        }
    } else if (roomItems.length >= 0) {
        for (var i = 0; i < roomItems.length; i++) {
            roomItems[i].parentNode.removeChild(roomItems[i]);

        }
        roomItems.length = 0;
        showRoomItems = true;
    }
}

function tookRoomItem() {
    var item = this;
    alert("Take item with id: " + item.getAttribute("item_id"));
    var command = {
            "code": 9,
            "item_key": parseInt(item.getAttribute("item_id")),
            "args": [item.getAttribute("item_id")],
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
        if (roomItems.length >= 0) {
            for (var i = 0; i < roomItems.length; i++) {
                roomItems[i].parentNode.removeChild(roomItems[i]);

            }
            roomItems.length = 0;
            showRoomItems = true;
        }
    }).fail(function (jqXHR, textStatus) {
        alert(jqXHR.status + " " + jqXHR.statusText + ". " + jqXHR.responseText);
        console.log(jqXHR.status + " " + jqXHR.statusText + ". " + jqXHR.responseText);
    });
}