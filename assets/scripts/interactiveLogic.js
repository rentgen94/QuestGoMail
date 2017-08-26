var interactives = [];
var showInters = true;

function showIntersItems() {
    if (showInters) {
        var settings = {
            "type": "GET",
            "contentType": "application/json; charset=utf-8",
            "xhrFields": {
                withCredentials: true
            },
            "dataType": "text",
            "crossDomain": true,
            "url": "http://localhost:8080/game/look_around/entities/interactives"
        };

        $.ajax(settings).done(function(xhr, textStatus, errorThrown) {
            interItem(xhr);
        }).fail(function (jqXHR, textStatus) {
            alert(jqXHR.status + " " + jqXHR.statusText + ". " + jqXHR.responseText);
            console.log(jqXHR.status + " " + jqXHR.statusText + ". " + jqXHR.responseText);
        });
    } else {
        interItem("");
    }
}

function interItem(data) {
    if (data != undefined && data != "") {
        var items = JSON.parse(data);
        if ((items["msg"] != "") || (items["errMsg"] != "")) {
            alert("Msg: " + items["msg"] + "; Error: " + items["errMsg"]);
        }
        if (showInters) {
            for (var i = 0; i < items["data"].length; i++) {
                var holder = document.createElement("div");

                var button = document.createElement("button");
                button.className = "buttonInterItem";
                button.innerHTML = "Hmm, with this item can do something...";
                button.onclick = interactWithItem;
                button.setAttribute("inter_id", items["data"][i]["id"]);

                var name = document.createElement("h3");
                name.innerHTML = items["data"][i]["name"];
                var desc = document.createElement("p");
                desc.innerHTML = items["data"][i]["description"];

                holder.appendChild(button);
                holder.appendChild(name);
                holder.appendChild(desc);

                interactives.push(holder);

                $('#intMenu').append(holder);
            }
            showInters = false;
        } else {
            for (var i = 0; i < interactives.length; i++) {
                interactives[i].parentNode.removeChild(interactives[i]);

            }
            interactives.length = 0;
            showInters = true;
        }
    } else if (interactives.length >= 0) {
        for (var i = 0; i < interactives.length; i++) {
            interactives[i].parentNode.removeChild(interactives[i]);
        }
        interactives.length = 0;
        showInters = true;
    }
}

function interactWithItem() {
    var item = this;
    alert("Take item with id: " + item.getAttribute("inter_id"));
    var command = {
        "code": 9,
        "item_key": parseInt(item.getAttribute("inter_id")),
        "args": [item.getAttribute("inter_id")],
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
        if (interactives.length >= 0) {
            for (var i = 0; i < interactives.length; i++) {
                interactives[i].parentNode.removeChild(interactives[i]);

            }
            interactives.length = 0;
            showInters = true;
        }
    }).fail(function (jqXHR, textStatus) {
        alert(jqXHR.status + " " + jqXHR.statusText + ". " + jqXHR.responseText);
        console.log(jqXHR.status + " " + jqXHR.statusText + ". " + jqXHR.responseText);
    });
}