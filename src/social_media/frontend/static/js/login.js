// This file contains functions needed for login related actions

/**
 * Log user in, alert user on error
 */
 function login() {
    fetch('/login/', {
        method: 'POST',
        body: JSON.stringify({
            "username": document.getElementById('username').value,
            "password": document.getElementById('password').value
        })
    }).then(function (response) {
        if (response.ok) {
            localStorage.setItem("username", document.getElementById("username").value);
            window.location.replace("/");
        } else {
            alert("username or password is wrong  : ( \nyou get 10 login attempts before you're locked out. if you are already locked out, talk to the admin and they will reset your password for you");
        }
    }).catch((error) => {
        console.error("error! something went wrong " + error);
    });
}

/**
 * Roulette of the russian variety
 * 
 * Usually sends you to the youtube video for Rasputin by Boney M
 */
 function roulette() {
    let choices = ["https://www.youtube.com/watch?v=16y1AkoZkmQ", "https://www.youtube.com/watch?v=dQw4w9WgXcQ", "https://www.youtube.com/watch?v=16y1AkoZkmQ", "https://www.youtube.com/watch?v=16y1AkoZkmQ", "https://www.youtube.com/watch?v=16y1AkoZkmQ", "https://www.youtube.com/watch?v=16y1AkoZkmQ"]

    window.location.href = choices[Math.floor(Math.random() * choices.length)];
}