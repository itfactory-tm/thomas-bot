<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>IT Factory: Discord</title>

     <link rel="apple-touch-icon" sizes="180x180" href="/static/apple-touch-icon.png">
     <link rel="icon" type="image/png" sizes="32x32" href="/static/favicon-32x32.png">
     <link rel="icon" type="image/png" sizes="16x16" href="/static/favicon-16x16.png">
     <link rel="manifest" href="/static/site.webmanifest">
     <link rel="mask-icon" href="safari-pinned-tab.svg" color="#5bbad5">
     <meta name="msapplication-TileColor" content="#2b5797">
     <meta name="theme-color" content="#ffffff">

     <link rel="stylesheet" href="/static/bot.css">
    <script type="text/javascript">
          var verifyCallback = function() {
            document.getElementById("invite").submit();
          };
    </script>
</head>
<body>
    <h1>Nice to meet you!</h1>
    <div class="container">
        <div class="speech-bubble">
            <p>You're just one stem removed fron joining the ITFactory Discord! Thomas Bot is the only robot who may enter. Can you confirm you're not a fellow robot?</p>
            <form action="/invite" method="POST" id="invite">
                <div class="h-captcha" data-sitekey="{{.HCaptchaSiteKey}}" data-callback="verifyCallback"></div>
            </form>
        </div>
        <p><img src="/static/thomasbot.png" alt="Thomas Bot"></p>
    </div>

    <script src="https://hcaptcha.com/1/api.js" async defer></script>
</body>
</html>