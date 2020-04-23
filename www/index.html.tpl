<!DOCTYPE html>
<html lang="nl">
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
          var onloadCallback = function() {
            var verifyCallback = function() {
                document.getElementById("invite").submit();
            };
            grecaptcha.render('g-recaptcha', {
              'sitekey' : '{{ .RecaptchaKey }}',
              'callback' : verifyCallback,
            });
          };
    </script>
</head>
<body>
    <h1>Nice to meet you!</h1>
    <div class="container">
        <div class="speech-bubble">
            <p>Je bent nog 1 stap weg van de ITFactory Discord! Thomas Bot is de enige robot die binnen mag. Wil je daarom even bevestigen dat jij geen collega robot bent?</p>
            <form action="/invite" method="POST" id="invite">
                <div class="g-recaptcha" id="g-recaptcha"></div>
            </form>
        </div>
        <p><img src="/static/thomasbot.png" alt="Thomas Bot"></p>
    </div>

    <script src="https://www.google.com/recaptcha/api.js?onload=onloadCallback&render=explicit" async defer></script>
</body>
</html>