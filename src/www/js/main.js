function open_main_page(){
  conn.send("{view}main_page");
  current_page.style.display = "none";
  current_page = document.getElementById("main_page");
  current_page.style.display = "block";
}

function open_account_page(){

    conn.send("{view}account_page");
    current_page.style.display = "none";
    current_page = document.getElementById("account_page");
    current_page.style.display = "block";
}

function change_accounts() {
    if (checkUsername_('change_username','change_username_status')&&checkPassword_('change_pass','change_password_status')){
      var prev_user_val = document.getElementById("username").value;
      var prev_password = document.getElementById("change_old_pass").value;
      var password = document.getElementById("change_pass").value;
      var user_val = document.getElementById("change_username").value;
      conn.send("{attempt_change:"+prev_user_val+"}"+user_val+","+encrypt_(prev_password + prev_user_val)+","+encrypt_(password+user_val));
    }
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
