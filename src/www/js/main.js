function open_main_page(){
  conn.send("{view}main_page");
  current_page.style.display = "none";
  current_page = document.getElementById("main_page");
  current_page.style.display = "block";
}
