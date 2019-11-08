
function submit_contact_me(){
  conn.send("{contact_me:"+username.innerHTML+"}"+document.getElementById('contact_me_email').value+"::"+document.getElementById('contact_me_message').value)
};


function open_contact_me(){
  conn.send("{view}contact_me");
  current_page.style.display = "none";
  current_page = document.getElementById("contact_me");
  current_page.style.display = "block";
}
