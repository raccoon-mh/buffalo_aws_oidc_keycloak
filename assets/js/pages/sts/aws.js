document.addEventListener("DOMContentLoaded", function () {
    getAWSSTStokenfromStorage()
});

document.getElementById('fetchDataBtn').addEventListener('click', function() {
    var authTokenInput = document.getElementById('authenticity_token');
    var rolearnInput = document.getElementById('selectedRolearn');
    var authToken = authTokenInput.value;
    var rolearn = rolearnInput.value;
    fetchDataAndStore(authToken, rolearn);
});

document.getElementsByName('rolearn').forEach(function(arntag) {
    arntag.addEventListener("click", function(event) {
        document.getElementById('selectedRolearn').value = this.textContent
    });
});

function getAWSSTStokenfromStorage(){
    var awsTokenjson = localStorage.getItem("awsToken");
    var awsToken = JSON.parse(awsTokenjson);
    document.getElementById("Audience").value  = awsToken.AssumeRoleWithWebIdentityResult.Audience
    document.getElementById("AssumedRoleId").value = awsToken.AssumeRoleWithWebIdentityResult.AssumedRoleUser.AssumedRoleId
    document.getElementById("Arn").value = awsToken.AssumeRoleWithWebIdentityResult.AssumedRoleUser.Arn
    document.getElementById("Provider").value = awsToken.AssumeRoleWithWebIdentityResult.Provider
    document.getElementById("AccessKeyId").value = awsToken.AssumeRoleWithWebIdentityResult.Credentials.AccessKeyId
    document.getElementById("SecretAccessKey").value = awsToken.AssumeRoleWithWebIdentityResult.Credentials.SecretAccessKey
    document.getElementById("SessionToken").value = awsToken.AssumeRoleWithWebIdentityResult.Credentials.SessionToken
    document.getElementById("Expiration").value = awsToken.AssumeRoleWithWebIdentityResult.Credentials.Expiration
    document.getElementById("RequestId").value = awsToken.ResponseMetadata.RequestId
}

function fetchDataAndStore(authToken, rolearn) { 

    var spinner = document.getElementById("spinner");
    spinner.style.display = "block";

    var requestData = {
        authenticity_token: authToken,
        rolearn: rolearn
    };

    $.ajax({
        url: "/sts/aws/token",
        type: 'POST',
        contentType: 'application/x-www-form-urlencoded',
        data: $.param(requestData),
        success: function(data) {
            console.log(data)
            localStorage.setItem('awsToken', JSON.stringify(data));
            getAWSSTStokenfromStorage()
            spinner.style.display = "none";
        },
        error: function(xhr, status, error) {
            console.error('err:', error);
        }
    });
}



