"use strict";

function shortenLink(link, callback) {

    console.log("submit url:" + link);

    $("#txt").ajaxStart(function () {
        $("#submitBtn").addClass("am-disabled");
    });
    $("#txt").ajaxComplete(function () {
        $("#submitBtn").removeClass("am-disabled");
    });
    $.ajax({
        url: "/add/",
        type: "post",
        data: {rawLink: link},
        dataType: "json",
        success: function (resp) {
            if (resp.success) {
                document.getElementById("shortLink").innerText = window.location.href.substr(0, window.location.href.length - window.location.pathname.length) + resp.data.link;
            } else {
                alert(resp.msg);
            }
        },
        error: function (XMLHttpRequest, textStatus, errorThrown) {
            alert(XMLHttpRequest.status);
            alert(textStatus);
        }
    });
}

var submitBtn = document.getElementById("submitBtn");
submitBtn.onclick = function (e) {
    var rawLink = document.getElementById("link").value;

    shortenLink(rawLink, function (shortLink) {
        document.getElementById("shortLink").innerText = shortLink;
    });
};