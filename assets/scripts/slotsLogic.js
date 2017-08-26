var slots = [];
var showSlots = true;

function showSlotItems() {
    if (showSlots) {
        var settings = {
            "type": "GET",
            "contentType": "application/json; charset=utf-8",
            "xhrFields": {
                withCredentials: true
            },
            "dataType": "text",
            "crossDomain": true,
            "url": "http://localhost:8080/game/look_around/entities/slots"
        };

        $.ajax(settings).done(function(xhr, textStatus, errorThrown) {
            slotItem(xhr);
        }).fail(function (jqXHR, textStatus) {
            alert(jqXHR.status + " " + jqXHR.statusText + ". " + jqXHR.responseText);
            console.log(jqXHR.status + " " + jqXHR.statusText + ". " + jqXHR.responseText);
        });
    } else {
        slotItem("");
    }
}

function slotItem(data) {
    if (data != undefined && data != "") {
        var items = JSON.parse(data);
        if ((items["msg"] != "") || (items["errMsg"] != "")) {
            alert("Msg: " + items["msg"] + "; Error: " + items["errMsg"]);
        }
        if (showSlots) {
            for (var i = 0; i < items["data"].length; i++) {
                var holder = document.createElement("div");

                var button = document.createElement("button");
                button.className = "buttonSlotItem";
                button.innerHTML = "Hmm, looks like we need to put something";
                button.onclick = putSlotItem;
                button.setAttribute("slot_id", items["data"][i]["id"]);

                var name = document.createElement("h3");
                name.innerHTML = items["data"][i]["name"];
                var capacity = document.createElement("p");
                capacity.innerHTML = "Capacity: " + items["data"][i]["capacity"];
                var contains = document.createElement("p");
                contains.innerHTML = "Contains: " + items["data"][i]["contains"];

                holder.appendChild(button);
                holder.appendChild(name);
                holder.appendChild(capacity);
                holder.appendChild(contains);

                slots.push(holder);

                $('#slotsMenu').append(holder);
            }
            showSlots = false;
        } else {
            for (var i = 0; i < slots.length; i++) {
                slots[i].parentNode.removeChild(slots[i]);

            }
            slots.length = 0;
            showSlots = true;
        }
    } else if (slots.length >= 0) {
        for (var i = 0; i < slots.length; i++) {
            slots[i].parentNode.removeChild(slots[i]);

        }
        slots.length = 0;
        showSlots = true;
    }
}

function putSlotItem() {
    var item = this;
    alert("Take item with id: " + item.getAttribute("slot_id"));
    var command = {
        "code": 9,
        "item_key": parseInt(item.getAttribute("slot_id")),
        "args": [item.getAttribute("slot_id")],
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
        if (slots.length >= 0) {
            for (var i = 0; i < slots.length; i++) {
                slots[i].parentNode.removeChild(slots[i]);

            }
            slots.length = 0;
            showSlots = true;
        }
    }).fail(function (jqXHR, textStatus) {
        alert(jqXHR.status + " " + jqXHR.statusText + ". " + jqXHR.responseText);
        console.log(jqXHR.status + " " + jqXHR.statusText + ". " + jqXHR.responseText);
    });
}