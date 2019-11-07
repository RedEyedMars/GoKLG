function checkUsername_(){
  var username_val = document.getElementById("username").value;
  if(/^[a-z0-9_-]{3,16}$/igm.test(username_val)){
    document.getElementById("username_status").src = "Success.jpg";
    document.getElementById("username_status").title = 'That username looks good!';
    return true;
  } else {
    if(username_val.length<3){
      document.getElementById("username_status").src = "Pending.jpg";
      document.getElementById("username_status").title = "Usernames must be greater than 3 characters long!";
    } else if(username_val.length>=16){
      document.getElementById("username_status").src = "Fail.jpg";
      document.getElementById("username_status").title = "Usernames must be less than 16 characters long!";
    } else if(/^.*\s.*$/igm.test(username_val)){
      document.getElementById("username_status").src = "Fail.jpg";
      document.getElementById("username_status").title = "Usernames must have no spaces in them!";
    } else if(/^.*[!@#$%^&*+=\(\)\[\]\{\}:;,\.\'\`~<>\/\\].*$/igm.test(username_val)){
      document.getElementById("username_status").src = "Fail.jpg";
      document.getElementById("username_status").title = "Usernames must have no special characters in them!";
    } else {
      document.getElementById("username_status").src = "Fail.jpg";
      document.getElementById("username_status").title = "Your username is not valid!";
    }
    return false;
  }
};
function checkPassword_(){
  return true;
}

function login(username_val){
  location.href = "http://redeyedmars.com/main.html"

  username = document.getElementById("displayusername");
  username.innerHTML = username_val;

};
function logout(){

  location.href = "http://redeyedmars.com/"

};
function signin_() {
    if (checkUsername_()&&checkPassword_()){
      var password = document.getElementById("pass").value;
      var user_val = document.getElementById("username").value;
      conn.send("{attempt_login}"+encrypt_(password+user_val));
    }
};
function signup_() {
    if (checkUsername_()&&checkPassword_()){
      var password = document.getElementById("pass").value;
      var user_val = document.getElementById("username").value;
      conn.send("{attempt_signup}"+user_val+","+encrypt_(password+user_val));
    }
};
function attempt_logout(){
  conn.send("{attempt_logout}");
};

function encrypt_(upwd){
  return forge_sha256(upwd);
};

function appendSigninStatus(val){
  const status = document.getElementById("account_signin_status");
  while (status.firstChild) {
    status.removeChild(status.firstChild);
  }
  var item = document.createElement("div");
  item.innerHTML = createTextLinks_(val);
  status.appendChild(item);
};

commands["login_successful"] = function(msg,user) {

  login(user);
};
commands["login_failed"] = function(msg,user){
  const status = document.getElementById("account_signin_status");
  while (status.firstChild) {
    status.removeChild(status.firstChild);
  }
  var item = document.createElement("div");
  item.innerHTML = createTextLinks_(msg);
  status.appendChild(item);
};
commands["signup_successful"] = function(msg,user){
  login(user);
};
commands["signup_failed"] = function(msg,user){
  appendSigninStatus(msg)
};
commands["logout_successful"] = function(msg,user){
  logout();
};
