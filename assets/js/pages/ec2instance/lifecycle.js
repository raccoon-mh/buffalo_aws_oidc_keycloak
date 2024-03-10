document.addEventListener("DOMContentLoaded", function () {
    var awsTokenjson = localStorage.getItem("awsToken");
    var awsToken = JSON.parse(awsTokenjson);
    document.getElementById("AccessKeyIdCard").textContent = awsToken.AssumeRoleWithWebIdentityResult.Credentials.AccessKeyId
    document.getElementById("Expirationcard").textContent = calculateTimeDifference(awsToken.AssumeRoleWithWebIdentityResult.Credentials.Expiration)
    document.getElementById("AccessKeyId").value = awsToken.AssumeRoleWithWebIdentityResult.Credentials.AccessKeyId
    document.getElementById("SecretAccessKey").value = awsToken.AssumeRoleWithWebIdentityResult.Credentials.SecretAccessKey
    document.getElementById("SessionToken").value = awsToken.AssumeRoleWithWebIdentityResult.Credentials.SessionToken
});

function calculateTimeDifference(targetDate) {
    const targetTime = new Date(targetDate).getTime();
    const currentTime = new Date().getTime();
  
    let difference = targetTime - currentTime;
  
    if (difference <= 0) {
      return "out of date";
    }
  
    const hours = Math.floor(difference / (1000 * 60 * 60));
    difference %= 1000 * 60 * 60;
    const minutes = Math.floor(difference / (1000 * 60));
    difference %= 1000 * 60;
    const seconds = Math.floor(difference / 1000);
  
    let result = '';
    if (hours > 0) {
      result += hours + '시간 ';
    }
    if (minutes > 0) {
      result += minutes + '분 ';
    }
    if (seconds > 0) {
      result += seconds + '초 ';
    }
  
    return result.trim();
}

document.getElementById('fetchDataBtn').addEventListener('click', function() {
    
    var authTokenInput = document.getElementById('authenticity_token');
    var authToken = authTokenInput.value;
    fetchData(authToken);
});


function fetchData(authToken) { 

    var spinner = document.getElementById("spinner");
    spinner.style.display = "block";

    accessKeyId = document.getElementById("AccessKeyId").value
    secretAccessKey = document.getElementById("SecretAccessKey").value
    sessionToken = document.getElementById("SessionToken").value

    var requestData = {
        authenticity_token: authToken,
        accessKeyId: accessKeyId,
        secretAccessKey: secretAccessKey,
        sessionToken: sessionToken,
    };

    $.ajax({
        url: "/ec2/list",
        type: 'POST',
        contentType: 'application/x-www-form-urlencoded',
        data: $.param(requestData),
        success: function(data) {
            console.log(data)
            var tableBody = document.getElementById("vmListTable")
            for (var i = 0; i < data.length; i++) {
                var instance = data[i];
                var row = tableBody.insertRow();

                var cell1 = row.insertCell(0);
                cell1.innerHTML = instance.InstanceID;

                var cell2 = row.insertCell(1);
                cell2.innerHTML = instance.InstanceType;

                var cell3 = row.insertCell(2);
                cell3.innerHTML = "<a href=\'http://"+instance.PublicDNS+"\'>"+instance.PublicDNS+"</a>";

                var cell4 = row.insertCell(3);
                cell4.innerHTML = "<a href='#'>"+instance.State +"</a>";
            }

            spinner.style.display = "none";
        },
        error: function(xhr, status, error) {
            console.error('err:', error);
        }
    });
}