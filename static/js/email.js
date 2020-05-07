// Urls
var select_url = window.location.origin + "/select/";
var page_url = window.location.origin + "/page/";
var item_url = window.location.origin + "/item/";
var reply_rul = window.location.origin + "/reply/";

var item_html = window.location.origin + "/static/html/itemView.html";

// control footer btn
var controlA_str = '<a class="footerBtn btn btn-outline-secondary" style="margin: 0px 10px;" role="button">';

// Select by id items
var msgContainer = $("#msgContainer");
var controlFooter = $("#controlFooter");

// ##################################### Generate Boxes List
var uSpanStr = '<span/>';
var boxDivStr = '<div class="boxDiv" style="text-align:center"></div>';
var liStr = '<li></li>';

function addrsToString(addrs){
    var str;
    for(var addr of addrs){
        ;
    }
}

function genBoxs(){
    // var boxs;
    $ul = $("#boxUl")
    for(var idx in boxs){
        // icon
        //$icon = $(iconStr).attr("data-id", key);
        //if(map[key]==true) {$icon.css("color", "#339533");}
        //$(usernameSpan_str).text(key);
        var box = boxs[idx];
        $boxDiv = $(boxDivStr).append($(uSpanStr).text(box));
        $boxDiv.attr("data-idx", idx);
        $li = $(liStr).append($boxDiv);
        $ul.append($li);
    }
}

function selectBtn($this){
    $(".boxDiv").removeClass("boxDiv-selected");
    $this.addClass("boxDiv-selected");
}

// ##################################### Generate page items

var msgDiv_str          = '<div class="row msgDiv"/>';     

var subDiv_str       = '<div class="col-3"/>';
var subSpan_str    = '<span style="display:block;margin-right: 15px;font-weight: 700;overflow: hidden;font-size:20px;"/>';

var fromDiv_str       = '<div class="col-2"/>';
var fromSpan_str     = '<span style="display:block;color: gray;font-size:15px"/>';

var toDiv_str       = '<div class="col-2" />'
var toSpan_str     = '<span style="display:block;color: gray;font-size:15px"/>';

var bodyDiv_str      = '<div class="col-3"/>';
var bodySpan_str        = '<span style="display:block;color: gray;font-size:15px"/>';

var dateDiv_str      = '<div class="col-2" style="text-align:right"/>';
var dateSpan_str        = '<span style="display:block;color: gray;font-size:small"/>';

function genPageItem(data, idx){
    var subject = data["subject"];
    var from = data["from"];//[0]["Address"];
    var to = data["to"];//[0]["Address"];
    var date = data["date"];
    var body = data["body"].substring(0, 40);
    ///var trimmedString = string.substring(0, length);
    var $subDiv = $(subDiv_str).append($(subSpan_str).text(subject));
    var $fromDiv = $(fromDiv_str).append($(fromSpan_str).text(from));
    var $toDiv = $(toDiv_str).append($(toSpan_str).text(to));
    var $bodyDiv = $(bodyDiv_str).append($(bodySpan_str).text(body));
    var $dateDiv = $(dateDiv_str).append($(dateSpan_str).text(date));

    var $msgDiv = $(msgDiv_str).append($subDiv, $fromDiv, $toDiv, $bodyDiv, $dateDiv);
    $msgDiv.attr("data-idx", idx);
    return $msgDiv;
}

function genPageFooter(page, hasnext){
    var pagenum = parseInt(page);
    $prev = $(controlA_str).append($("<strong>").text("<<Prev"));
    $next = $(controlA_str).append($("<strong>").text("Next>>"));
    if(pagenum==1){
        $prev.addClass("disabled");
        //$prev.prop('disabled', true);
    }else{
        $prev.addClass("activeBtn");
        $prev.attr("data-method", "page");
        $prev.attr("data-param", pagenum-1);
    }
    if(!hasnext){
        $next.addClass("disabled");
    }else{
        $next.addClass("activeBtn");
        $next.attr("data-method", "page");
        $next.attr("data-param", pagenum+1);
    }
    controlFooter.children().remove();
    controlFooter.append($prev, $next);
}

function genPage(data){
    var msglist = data["msglist"];
    var hasnext = data["hasnext"];
    var pagenum = data["page"];
    var hasprev = (pagenum != 1);
    msgContainer.children().remove();
    var i = 0;
    for (i=0;i < msglist.length;i++){
        var msg = msglist[i];
        var $msgDiv = genPageItem(msg, i);
        msgContainer.append($msgDiv);
    }
    genPageFooter(pagenum, hasnext);
}

//################################## Generate Items


// ids:     itmeSub, itemDate, itemToName, 
//          itemToAddr, itemFromName, itemFromAddr, itemBody
function genItem(data){
    var subject = data["subject"];
    var fromName = data["from"];//[0]["Name"];
    var fromAddr = data["from"];//[0]["Address"];
    var toName = data["to"];//[0]["Name"];
    var toAddr = data["to"];//[0]["Address"];
    var date = data["date"];
    var body = data["body"];
    console.log(data);
    msgContainer.load(item_html, function(){
        $("#itemSub").text(subject);
        $("#itemDate").text(date);
        $("#itemToName").text(toName);
        $("#itemToAddr").text(toAddr);
        $("#itemFromName").text(fromName);
        $("#itemFromAddr").text(fromAddr);

        var textBody = '<span class="col-11" id="itemBody" style="white-space: pre-line"></span>';
        var htmlBody = '<div class="col-11" id="itemBody">';
        var $div=$('<div>').html(body);
        // if there are any children it is html
        if($div.children().length){
            $("#itemBodyContainer").append($(htmlBody).html(body));
        }else{
            //body = body.replace(/(?:\r\n|\r|\n)/g, '<br>');
             $("#itemBodyContainer").append($(textBody).text(body));
        }
    });
}

function page(page){
    var url = page_url+page.toString();
    $.getJSON(url, function(data){
        genPage(data);
    });
}

function select(idx){
    var url = select_url+idx.toString();
    $.getJSON(url, function(data){
        genPage(data);
    });
}

function item(idx){
    var url = item_url + idx.toString();
    $.getJSON(url, function(data){
        genItem(data);
    });
}

$(function(){
    genBoxs();
    $(".boxDiv").click(function(){
        var idx = $(this).attr("data-idx");
        select(idx);
        selectBtn($(this));
    });
    msgContainer.on('click', 'div.msgDiv', function() {
        // $(".msgDiv").click(function(){
        var idx = $(this).attr("data-idx");
        item(idx);
    });
    $("body").on('click', 'a.footerBtn', function() {
        // $(".msgDiv").click(function(){
        var method = $(this).attr("data-method");
        var param = $(this).attr("data-param");
        if(method == "page"){
            page(param);
        }
    });
});