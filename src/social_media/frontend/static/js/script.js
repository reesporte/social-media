/**
 * Display / don't display the item with id `id`
 * @param {string} id 
 */
function toggleCollapsed(id, replyingTo="") {
    if (document.getElementById("replyingTo")) {
        document.getElementById("replyingTo").value = replyingTo   
    }
    
    var element = document.getElementById(id);
    var style = element.style.display;
    if (style == "block") {
        document.getElementById(id).style.display = "none";
    } else {
        document.getElementById(id).style.display = "block";
    }
}

/**
 * call /logout endpoint
 */
function logout() {
    window.location.replace("/logout/");
}


/**
 * go to /profile endpoint
 */
 function profile() {
    window.location.replace("/profile/");
}