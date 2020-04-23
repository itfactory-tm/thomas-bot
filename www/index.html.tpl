<!DOCTYPE html>
<html lang="nl">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>IT Factory: Discord</title>
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
    <form action="/invite" method="POST" id="invite">
        <div class="g-recaptcha" id="g-recaptcha"></div>
    </form>
    <script src="https://www.google.com/recaptcha/api.js?onload=onloadCallback&render=explicit" async defer></script>
</body>
</html>