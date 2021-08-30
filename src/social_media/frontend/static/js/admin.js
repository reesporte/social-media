function createUser() {
    fetch('/admin/create', {
        method: 'POST',
        body: JSON.stringify({
            "username" : document.getElementById("username-create").value,
            "admin" : document.getElementById("admin-create").checked
        })
    }).then(function (response){
        if (response.ok) {
            return response.json()
        } else {
            throw "could not create user. status: " + response.status
        }
    }).then(function (json) {
        alert("user created! " + json['pw'])
    }).catch(function (error) {
        alert(error)
    })
    document.getElementById('create-user-form').reset()
    toggleCollapsed('create-user-modal')
}

function deleteUser() {
    fetch('/admin/delete', {
        method: 'DELETE',
        body: JSON.stringify({
            "username" : document.getElementById("username-delete").value,
        })
    }).then(function (response){
        if (response.ok) {
            alert("user deleted! : ) ")
        } else {
            alert("could not delete user. status: " + response.status)
        }
    })   
    document.getElementById('delete-user-form').reset()
    toggleCollapsed('delete-user-modal')
}

function resetUserPassword() {
    fetch('/admin/reset', {
        method: 'PUT',
        body: JSON.stringify({
            "username" : document.getElementById("username-reset").value,
        })
    }).then(function (response){
        if (response.ok) {
            return response.json()
        } else {
            throw "could not reset user password. status: " + response.status
        }
    }).then(function (json) {
        alert("user password reset! " + json['pw'])
    }).catch(function (error) {
        alert(error)
    })
    document.getElementById('reset-user-pass-form').reset()
    toggleCollapsed('reset-user-pass-modal')
}
