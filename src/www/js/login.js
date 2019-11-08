function checkUsername_(usernameNameId,usernameStatusId){
  var username_val = document.getElementById(usernameNameId).value;
  if(/^[a-z0-9_-]{3,16}$/igm.test(username_val)){
    document.getElementById(usernameStatusId).src = "Success.png";
    document.getElementById(usernameStatusId).title = 'That username looks good!';
    return true;
  } else {
    if(username_val.length<3){
      document.getElementById(usernameStatusId).src = "Pending.png";
      document.getElementById(usernameStatusId).title = "Usernames must be greater than 3 characters long!";
    } else if(username_val.length>=16){
      document.getElementById(usernameStatusId).src = "Fail.png";
      document.getElementById(usernameStatusId).title = "Usernames must be less than 16 characters long!";
    } else if(/^.*\s.*$/igm.test(username_val)){
      document.getElementById(usernameStatusId).src = "Fail.png";
      document.getElementById(usernameStatusId).title = "Usernames must have no spaces in them!";
    } else if(/^.*[!@#$%^&*+=\(\)\[\]\{\}:;,\.\'\`~<>\/\\].*$/igm.test(username_val)){
      document.getElementById(usernameStatusId).src = "Fail.png";
      document.getElementById(usernameStatusId).title = "Usernames must have no special characters in them!";
    } else {
      document.getElementById(usernameStatusId).src = "Fail.png";
      document.getElementById(usernameStatusId).title = "Your username is not valid!";
    }
    return false;
  }
};
function checkPassword_(passwordNameId,passwordStatusId){
  return true;
}

function login(username_val){
  document.getElementById("popup").style.display = "none";
  document.getElementById("site_div").style.display = "block";

  conn.send("{view}main_page");
  current_page = document.getElementById("main_page");
  current_page.style.display = "block";

  username = document.getElementById("displayusername");
  username.innerHTML = username_val;

}
function logout(){
  const status = document.getElementById("account_signin_status");
  while (status.firstChild) {
    status.removeChild(status.firstChild);
  }
  document.getElementById("popup").style.display = "block";
  document.getElementById("site_div").style.display = "none";
  current_page.style.display = "none";
}
function signin_() {
    if (checkUsername_('username','username_status')&&checkPassword_('pass','password_status')){
      var password = document.getElementById("pass").value;
      var user_val = document.getElementById("username").value;
      conn.send("{attempt_login:"+user_val+"}"+encrypt_(password+user_val));
    }
};
function signup_() {
    if (checkUsername_('username','username_status')&&checkPassword_('pass','password_status')){
      var password = document.getElementById("pass").value;
      var user_val = document.getElementById("username").value;
      conn.send("{attempt_signup:"+user_val+"}"+encrypt_(password+user_val));
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
