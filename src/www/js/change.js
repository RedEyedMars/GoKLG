
function open_account_page(){

    conn.send("{view}account_page");
    current_page.style.display = "none";
    current_page = document.getElementById("account_page");
    current_page.style.display = "block";
    document.getElementById("change_old_pwd_div").style.display = "block";
    document.getElementById("change_user_div").style.display = "none";
    document.getElementById("change_pwd_div").style.display = "none";
}

function change_profile() {
    if (heckPassword_('change_old_pass','change_old_password_status')){
      var prev_user_val = document.getElementById("username").value;
      var prev_password = document.getElementById("change_old_pass").value;
      var password = document.getElementById("change_pass").value;
      var user_val = document.getElementById("change_username").value;
      conn.send("{request_change_permission:"+prev_user_val+"}"+encrypt_(prev_password + prev_user_val));
    }
};

function change_user_name() {
    if (checkUsername_('change_username','change_username_status')&&checkPassword_('change_pass','change_password_status')){
      var prev_user_val = document.getElementById("username").value;
      var prev_password = document.getElementById("change_old_pass").value;
      var password = document.getElementById("change_pass").value;
      var user_val = document.getElementById("change_username").value;
      conn.send("{change_username:"+prev_user_val+"}"+user_val+","+encrypt_(prev_password + prev_user_val));
    }
};

function change_password() {
    if (checkUsername_('change_username','change_username_status')&&checkPassword_('change_pass','change_password_status')){
      var prev_user_val = document.getElementById("username").value;
      var prev_password = document.getElementById("change_old_pass").value;
      var password = document.getElementById("change_pass").value;
      var user_val = document.getElementById("change_username").value;
      conn.send("{change_password:"+prev_user_val+"}"+encrypt_(password + user_val)+","+encrypt_(prev_password + prev_user_val));
    }
};

commands["change_request_granted"] = function(msg,user){
  document.getElementById("change_old_pwd_div").style.display = "none";
  document.getElementById("change_user_div").style.display = "block";
  document.getElementById("change_pwd_div").style.display = "block";
};
commands["change_request_denied"] = function(msg,user){
    const status = document.getElementById("account_change_status");
    while (status.firstChild) {
      status.removeChild(status.firstChild);
    }
    var item = document.createElement("div");
    item.innerHTML = createTextLinks_(msg);
    status.appendChild(item);
};

commands["change_successful"] = function(msg,user){

  username = document.getElementById("displayusername");
  username.innerHTML = user;

  document.getElementById("change_pass").value = "";
  document.getElementById("change_old_pass").value = "";
  document.getElementById("change_username").value = "";
};
commands["change_failed"] = function(msg,user){
  const status = document.getElementById("account_change_status");
  while (status.firstChild) {
    status.removeChild(status.firstChild);
  }
  var item = document.createElement("div");
  item.innerHTML = createTextLinks_(msg);
  status.appendChild(item);
};
