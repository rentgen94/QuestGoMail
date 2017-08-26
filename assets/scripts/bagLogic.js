var bagItems = [];
var showBagItems = true;

function showBagItem() {
    if (showBagItems) {
        var settings = {
            "type": "GET",
            "contentType": "application/json; charset=utf-8",
            "xhrFields": {
                withCredentials: true
            },
            "dataType": "text",
            "crossDomain": true,
            "url": "http://localhost:8080/game/look_around/bag"
        };

        $.ajax(settings).done(function(xhr, textStatus, errorThrown) {
            bagItem(xhr);
        }).fail(function (jqXHR, textStatus) {
            alert(jqXHR.status + " " + jqXHR.statusText + ". " + jqXHR.responseText);
            console.log(jqXHR.status + " " + jqXHR.statusText + ". " + jqXHR.responseText);
        });
    } else {
        bagItem("");
    }
}

function bagItem(data) {
    if (data != undefined && data != "") {
        var items = JSON.parse(data);
        if ((items["msg"] != "") || (items["errMsg"] != "")) {
            alert("Msg: " + items["msg"] + "; Error: " + items["errMsg"]);
        }
        if (showBagItems) {
            for (var i = 0; i < items["data"].length; i++) {
                var holder = document.createElement("div");

                var button = document.createElement("button");
                button.className = "buttonBagItem";
                button.innerHTML = "What to do with this???";
                button.onclick = dropBagItem;
                button.setAttribute("bagIt_id", items["data"][i]["id"]);

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

                bagItems.push(holder);

                $('#bagMenu').append(holder);
            }
            showBagItems = false;
        } else {
            for (var i = 0; i < bagItems.length; i++) {
                bagItems[i].parentNode.removeChild(bagItems[i]);

            }
            bagItems.length = 0;
            showBagItems = true;
        }
    } else if (bagItems.length >= 0) {
        for (var i = 0; i < bagItems.length; i++) {
            bagItems[i].parentNode.removeChild(bagItems[i]);

        }
        bagItems.length = 0;
        showBagItems = true;
    }
}

function dropBagItem() {
    var item = this;
    var command = {
        "code": 10,
        "item_key": parseInt(item.getAttribute("bagIt_id")),
        "args": [item.getAttribute("bagIt_id")],
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
        if (bagItems.length >= 0) {
            for (var i = 0; i < bagItems.length; i++) {
                bagItems[i].parentNode.removeChild(bagItems[i]);

            }
            bagItems.length = 0;
            showBagItems = true;
        }
    }).fail(function (jqXHR, textStatus) {
        alert(jqXHR.status + " " + jqXHR.statusText + ". " + jqXHR.responseText);
        console.log(jqXHR.status + " " + jqXHR.statusText + ". " + jqXHR.responseText);
    });
}