window.WS=""
function FindEnemy(){
    if (localStorage.getItem("EnemyId")==null){
        var settings = {
            "url": "http://work.zxcvbnm.online:80/find/"+localStorage.getItem("UID")+"?token="+localStorage.getItem("Token"),
            "method": "GET",
            "timeout": 0,
        };
        $.ajax(settings).done(function (response) {
            if(response["Finded"]){
                $('#text_info').text("Сделайте выбор:");
                localStorage.setItem("EnemyId",response["Playingwith"]);
                $('#enemy-info').text("Вы играете против: "+response["Playingwith"]).show(1200);
                $('#choose').show(1200);
                console.log(localStorage.getItem("EnemyId"))
            }
        })
    }

}

function WinOrNot(){
    if (localStorage.getItem("chosen")!==null){
        var settings = {
            "url": "http://work.zxcvbnm.online:80/check/"+localStorage.getItem("UID")+"?token="+localStorage.getItem("Token"),
            "method": "GET",
            "timeout": 0,
        };
        $.ajax(settings).done(function (response) {

            if(!response["Wait"]){
                if (response["Windata"] === "Вы выиграли!"){
                    $('#text_info').text("Вы выиграли!!!");
                    $('#gameend').attr('src',"/src/win-cat.gif").show(1000);
                    localStorage.removeItem("chosen")
                }

                if (response["Windata"] === "Вы проиграли!"){
                    $('#text_info').text("Вы проиграли!!!");
                    $('#gameend').attr('src',"/src/lose-cat.gif").show(1000);
                    localStorage.removeItem("chosen")
                }

                if (response["Windata"] === "Ничья"){
                    $('#text_info').text("Ничья!!!");
                    $('#gameend').attr('src',"/src/draw.gif").show(1000);
                    localStorage.removeItem("chosen")
                }
            }
        })
    }
}

$( window ).on( "load",function() {
    window.WS = new WebSocket("ws://work.zxcvbnm.online:80/ws");


    window.WS.onmessage = function(event) {
        var response=JSON.parse(event.data);
        if(response["Finded"]){
            $('#text_info').text("Сделайте выбор:");
            localStorage.setItem("EnemyId",response["Playingwith"]);
            $('#enemy-info').text("Вы играете против: "+response["Playingwith"]).show(1200);
            $('#choose').show(1200);
            console.log(localStorage.getItem("EnemyId"))
        }
        if(response["Wait"]===false){
            if (response["Windata"] === "Вы выиграли!"){
                $('#text_info').text("Вы выиграли!!!");
                $('#gameend').attr('src',"/src/win-cat.gif").show(1000);
                localStorage.removeItem("chosen")
            }

            if (response["Windata"] === "Вы проиграли!"){
                $('#text_info').text("Вы проиграли!!!");
                $('#gameend').attr('src',"/src/lose-cat.gif").show(1000);
                localStorage.removeItem("chosen")
            }

            if (response["Windata"] === "Ничья"){
                $('#text_info').text("Ничья!!!");
                $('#gameend').attr('src',"/src/draw.gif").show(1000);
                localStorage.removeItem("chosen")
            }
        }
        if (response["Token"]!=null){
            if (localStorage.getItem("UID")!=null) {
                localStorage.removeItem("UID");
            }
            if (localStorage.getItem("Token")!=null) {
                localStorage.removeItem("Token");
            }
            if (localStorage.getItem("chosen")!=null) {
                localStorage.removeItem("chosen");
            }
            console.log(response["Token"])
            localStorage.setItem("UID",response["UserUid"]);
            localStorage.setItem("Token",response["Token"]);
            if (localStorage.getItem("EnemyId")!=null) {
                localStorage.removeItem("EnemyId");
            }
            console.log(localStorage.getItem("EnemyId"))
            $('#userid').text("UID: "+response["UserUid"]);
            //alert(response["UserUid"]);
        }
        //alert("Получены данные " + event.data);
    };

    window.WS.onopen = function() {
        window.WS.send('{"type":"new"}');
    };
    // var settings = {
    //     "url": "http://work.zxcvbnm.online:80/new",
    //     "method": "GET",
    //     "timeout": 0,
    // };
    $('#rock').on('click',function(event){
        if (localStorage.getItem("chosen")!=null){
            return
        }
        localStorage.setItem("chosen","ok")
        window.WS.send('{"type":"choose","uid":"'+localStorage.getItem("UID")+'","chose":"1"}');
        $('#text_info').text("Ожидаем противника...");
        $('#choose').hide(500);


    });
    $('#scissors').on('click',function(event){
        if (localStorage.getItem("chosen")!=null){
            return
        }
        localStorage.setItem("chosen","ok")
        window.WS.send('{"type":"choose","uid":"'+localStorage.getItem("UID")+'","chose":"2"}');
        $('#text_info').text("Ожидаем противника...");
        $('#choose').hide(500);
    });
    $('#paper').on('click',function(event){
        if (localStorage.getItem("chosen")!=null){
            return
        }
        localStorage.setItem("chosen","ok")
        window.WS.send('{"type":"choose","uid":"'+localStorage.getItem("UID")+'","chose":"3"}');
        $('#text_info').text("Ожидаем противника...");
        $('#choose').hide(500);
    });

    $('#gameend').on('click',function(event){
        window.WS.send('{"type":"restart","uid":"'+localStorage.getItem("UID")+'","chose":"1"}');
        if (localStorage.getItem("chosen")!=null) {
            localStorage.removeItem("chosen");
        }
        if (localStorage.getItem("EnemyId")!=null) {
            localStorage.removeItem("EnemyId");
        }
        $('#text_info').text("Ищем противника...");
        $('#gameend').hide(500);
        $('#enemy-info').hide(500);

    });


    // $.ajax(settings).done(function (response) {
    //
    // })


    //setInterval(FindEnemy, 500);
    //setInterval(WinOrNot, 500);
})