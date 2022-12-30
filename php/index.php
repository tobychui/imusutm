<!DOCTYPE html>
<html>
	<head>
		<meta charset="UTF-8">
		<meta name="description" content="Uptime monitor for imuslab">
		<meta name="keywords" content="uptime, monitor, automatic">
		<meta name="author" content="tobychui">
		<title>Uptime Monitor | imuslab</title>
		<link rel="icon" type="image/png" href="favicon.png" />

		<!-- Responsive -->
		<meta name="apple-mobile-web-app-capable" content="yes" />
		<meta name="viewport" content="user-scalable=no, width=device-width, initial-scale=0.9, maximum-scale=1"/>
		<?php //include_once("analyticstracking.php") ?>

		<!-- HTML Meta Tags -->
		<title>Uptime Monitor | imuslab</title>
		<meta name="description" content="imuslab Uptime monitor">

		<!-- Facebook Meta Tags -->
		<meta property="og:url" content="https://imuslab.com/">
		<meta property="og:type" content="website">
		<meta property="og:title" content="imuslab">
		<meta property="og:description" content="imuslab Uptime monitor">
		<meta property="og:image" content="">

		<!-- Twitter Meta Tags -->
		<meta name="twitter:card" content="summary_large_image">
		<meta property="twitter:domain" content="imuslab.com">
		<meta property="twitter:url" content="https://imuslab.com/">
		<meta name="twitter:title" content="imuslab">
		<meta name="twitter:description" content="imuslab Uptime monitor">
		<meta name="twitter:image" content="">

		<!-- css and js -->
		<link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/fomantic-ui/2.9.0/semantic.min.css">
		<script src="https://code.jquery.com/jquery-3.6.3.min.js" integrity="sha256-pvPw+upLPUjgMXY0G+8O0xUf+/Im1MZjXxxgOcBQBXU=" crossorigin="anonymous"></script>
		<script src="https://cdnjs.cloudflare.com/ajax/libs/fomantic-ui/2.9.0/semantic.min.js"></script>
		<link href="https://unpkg.com/aos@2.3.1/dist/aos.css" rel="stylesheet">
		<script src="https://unpkg.com/aos@2.3.1/dist/aos.js"></script>

		<!-- fonts -->
		<link rel="preconnect" href="https://fonts.googleapis.com">
		<link rel="preconnect" href="https://fonts.gstatic.com" crossorigin>
		<link href="https://fonts.googleapis.com/css2?family=Noto+Sans+TC:wght@300;400;500;700;900&display=swap" rel="stylesheet">

		<!-- css override -->
		<style>
			body{
				padding:0px !important;
                background-color: #f5f5f5;
			}

			h1, h2, h3, h4, h5, p, span, div, a { 
				font-family: 'Noto Sans TC', sans-serif;
                color: #333333;
			}

            #utm{
                background-color: white;
                border-radius: 1em;
            }


            .domain{
                margin-bottom: 1em;
                position: relative;
            }

            .statusDot{
                height: 1.8em;
                border-radius: 0.4em;
                width: 0.4em;
                background-color: #e8e8e8;
                display:inline-block;
                cursor: pointer;
                margin-left: 0.1em;
            }
            
            .online.statusDot{
                background-color: #3bd671;
            }
            .error.statusDot{
                background-color: #f29030;
            }
            .offline.statusDot{
                background-color: #df484a;
            }
            .padding.statusDot{
                cursor: auto;
            }
            
		</style>
	</head>
	<body>
        <br>
		<div class="ui container">
            <div id="utm" class="ui basic segment">
                <div class="ui basic segment">
                    <h4 class="ui header">
                        <i class="red remove icon"></i>
                        <div class="content">
                            Uptime Monitoring service is currently unavailable
                          <div class="sub header">This might cause by an error in cluster communication within the host servers. Please wait for administrator to resolve the issue.</div>
                        </div>
                    </h4> 
                </div>
               
            </div>
            <div class="ui divider"></div>
            <a href="//imuslab.com"><img class="ui tiny right floated image" src="img/logo.png"></a>
            <small>Uptime service provided by imuslab. We do not ensure the data shown is accurate in any means and shall be used as reference only.</small>
            <br><br><br>
		</div>
		<script>
			AOS.init();
            let records = JSON.parse(`<?php echo file_get_contents("http://localhost:8089/")?>`);
			
			//For every 5 minutes
            /*
			setInterval(function(){
				$.get("api.php?token={{your_token}}", function(data){
					console.log("Status Updated");
					records = data;
					renderRecords();
				});
			}, (300 * 1000));
            */
			
            function renderRecords(){
                $("#utm").html("");
                for (let [key, value] of Object.entries(records)) {
                    renderUptimeData(key, value);
                }
            }

            function format_time(s) {
                const date = new Date(s * 1e3);
                return(date.toLocaleString());
            }


            function renderUptimeData(key, value){
                if (value.length == 0){
                    return
                }

                let id = value[0].ID;
                let name = value[0].Name;
                let url = value[0].URL;
                let protocol = value[0].Protocol;

                //Generate the status dot
                let statusDotList = ``;
                for(var i = 0; i < (288 - value.length); i++){
                    //Padding
                    statusDotList += `<div class="padding statusDot"></div>`
                }

                let ontimeRate = 0;
                for (var i = 0; i < value.length; i++){
                    //Render status to html
                    let thisStatus = value[i];
                    let dotType = "";
                    if (thisStatus.Online){
                        if (thisStatus.StatusCode < 200 || thisStatus.StatusCode >= 300){
                            dotType = "error";
                        }else{
                            dotType = "online";
                        }
                        ontimeRate++;
                    }else{
                        dotType = "offline";
                    }

                    let datetime = format_time(thisStatus.Timestamp);
                    statusDotList += `<div title="${datetime}" class="${dotType} statusDot"></div>`
                }

                ontimeRate = ontimeRate / value.length * 100;
                let ontimeColor = "#df484a"
                if (ontimeRate > 0.8){
                    ontimeColor = "#3bd671";
                }else if(ontimeRate > 0.5) {
                    ontimeColor = "#f29030";
                }
                //Check of online status now
                let currentOnlineStatus = "Unknown";
                let onlineStatusCss = ``;
                if (value[value.length - 1].Online){
                    currentOnlineStatus = `<i class="circle icon"></i> Online`;
                    onlineStatusCss = `color: #3bd671;`;
                }else{
                    currentOnlineStatus = `<i class="circle icon"></i> Offline`;
                    onlineStatusCss = `color: #df484a;`;
                }

                //Generate the html
                $("#utm").append(`<div class="ui basic segment statusbar">
                    <div class="domain">
                        <div style="position: absolute; top: 0; right: 0.4em;">
                            <p class="onlineStatus" style="display: inline-block; font-size: 1.3em; padding-right: 0.5em; padding-left: 0.3em; ${onlineStatusCss}">${currentOnlineStatus}</p>
                        </div>
                        <div>
                            <h3 class="ui header" style="margin-bottom: 0.2em;">${name}</h3>
                            <a href="${url}" target="_blank">${url}</a> | <span style="color: ${ontimeColor};">${(ontimeRate).toFixed(2)}%<span>
                        </div>
						<div class="ui basic label protocol" style="position: absolute; bottom: 0; right: 0.2em; margin-bottom: -0.6em;">
							proto: ${protocol}
						</div>
                    </div>
                    <div class="status" style="marign-top: 1em;">
                        ${statusDotList}
                    </div>
                    <div class="ui divider"></div>
                </div>`);
            }

            renderRecords(records);
			
		</script>
	</body>
</html>