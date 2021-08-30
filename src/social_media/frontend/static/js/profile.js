// This file handles functions needed for profile related actions

/**
 * Get user profile information
 */
function loadProfile(id) {
    window.location.replace("/profile/" + id)
}

/**
 * change Password
 */
function changePassword() {
    fetch('/profile/password', {
        method: 'PUT',
        body: JSON.stringify({
            "currentPassword": document.getElementById('current-password').value,
            "newPassword": document.getElementById('new-password').value
        })
    }).then(function (response) {
        if (response.ok) {
            document.getElementById("change-password-form").reset();
            toggleCollapsed("password-modal");
            alert("password successfully changed!");
        } else if (response.status === 400) {
            alert("that password is in the top 100 most used passwords of 2021 and is not allowed");
        } else {
            alert("couldn't change password : (")
        }
    });
}

/**
 * change bio
 */
function changeBio() {
    fetch('/profile/bio', {
        method: 'PUT',
        body: JSON.stringify({
            "bio": document.getElementById('new-bio').value
        })
    }).then(function (response) {
        if (response.ok) {
            document.getElementById("change-bio-form").reset();
            toggleCollapsed("bio-modal");
            alert("bio successfully changed!");
        } else {
            alert("something went wrong changing your bio : (")
        }
    }).then(function () {
        loadProfile(localStorage.getItem("username"));
    })
}


/**
 * change profile picture
 */
function changeProfilePicture() {
    fetch('/media/', {
        method: 'POST',
        body: new FormData(document.getElementById("change-pic-form"))
    }).then(function (response) {
        if (!response.ok) {
            alert("could not change profile picture");
        }
        document.getElementById("change-pic-form").reset();
    }).then(function () {
        loadProfile(localStorage.getItem("username"));
    })
}
