var baseUrl = 'http://localhost:1025';
var api_prefix = '/caro/api/v1';

var user = {};
document.getElementById('btn-login').addEventListener('click', doLogin);

function checkLogin() {
    if (!user?.name) {
        showLogin();
    } else {
        showLobby();
    }

    function showLogin() {
        document.getElementById('screen-login').classList.remove('hidden');
        document.getElementById('screen-lobby').classList.add('hidden');
    }

    function showLobby() {
        document.getElementById('screen-login').classList.add('hidden');
        document.getElementById('screen-lobby').classList.remove('hidden');
        document.getElementById('user-name').innerText = user.name;
    }
}

function doLogin() {
    var userName = document.getElementById('input-user-name').value;
    user = {
        name: userName
    }
    checkLogin();
}

checkLogin();


